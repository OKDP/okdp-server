/*
 *    Copyright 2024 okdp.io
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package authz

import (
	"errors"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/security/authc/model"
	"github.com/okdp/okdp-server/internal/utils"
	"github.com/okdp/okdp-server/internal/logging"
)

type Enforcer struct {
	*casbin.Enforcer
}

func newEnforcer(authZConf config.AuthZ) (*Enforcer, error) {
	file := authZConf.File
	log.Info("Loading casbin configuration files, Model file: %s, Policy file: %s", file.ModelPath, file.PolicyPath)
	e, err := casbin.NewEnforcer(file.ModelPath, file.PolicyPath)
	return &Enforcer{e}, err
}

// Authorizer returns a middleware that will authorize the user to access resources based on the policy and model configuration.
func Authorizer(authzConfig config.AuthZ) gin.HandlerFunc {
	e, err := newEnforcer(authzConfig)
	if err != nil {
		log.Panic("Unable to get enforcer instance: %s", err)
	}
	return e.authorize()
}

func (e *Enforcer) authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			allowed bool = false
			err     error
		)
		userInfo, ok := c.Get(constants.OAuth2UserInfo)
		if !ok {
			log.Warn("Unable to authorize user, no user informtaion found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unable to authorize user, no user informtaion found in context"})
			return
		}
		email := userInfo.(*model.UserInfo).Email
		sub := userInfo.(*model.UserInfo).Subject

		rSub := userInfo.(*model.UserInfo).Roles
		rObj := c.Request.URL.Path
		rAct := c.Request.Method

		roles := utils.Map(rSub, func(s string) string {
			return constants.CasbinRolePrefix + s
		})
		e.Enforce(roles, rObj, rAct)

		// Check the role is allowed to access the path with the action
		for _, role := range roles {
			allowed1, err1 := e.Enforce(role, rObj, rAct)
			allowed = allowed || allowed1
			err = errors.Join(err, err1)
			if allowed {
				break
			}
		}

		if err != nil {
			log.Warn("Unable to authorize user (%s/%s): %s",email, sub, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		if !allowed {
			log.Warn("User (%s/%s) not allowed to execute the action", email, sub)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized action"})
			return
		}

		c.Next()
	}
}

