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
	errr "errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/common/constants"
	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/model"
	authc "github.com/okdp/okdp-server/internal/security/authc/model"
	"github.com/okdp/okdp-server/internal/utils"
)

type Enforcer struct {
	*casbin.Enforcer
}

func newEnforcer(authZConf config.AuthZ) (*Enforcer, error) {
	var e *casbin.Enforcer
	var err error
	authzProvider := strings.ToLower(authZConf.Provider)
	switch authzProvider {
	case "inline":
		modelFilePath, policyFilePath, er := writeFilesToTmp(authZConf.InLine.Model, authZConf.InLine.Policy)
		if er != nil {
			return nil, er
		}
		log.Info("Loading casbin configuration files, Model file: %s, Policy file: %s", modelFilePath, policyFilePath)
		e, err = casbin.NewEnforcer(modelFilePath, policyFilePath)
	case "file":
		file := authZConf.File
		log.Info("Loading casbin configuration files, Model file: %s, Policy file: %s", file.ModelPath, file.PolicyPath)
		e, err = casbin.NewEnforcer(file.ModelPath, file.PolicyPath)
	case "database":
		return nil, fmt.Errorf("database provider option not implemented")
	default:
		return nil, fmt.Errorf("provider option '%s' not recognized, valid ones: inline or file", authZConf.Provider)
	}
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
			allowed = false
			err     error
		)
		userInfo, ok := c.Get(constants.OAuth2UserInfo)
		if !ok {
			log.Warn("Unable to authorize user, no user informtaion found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.
				NewServerResponse(model.OkdpServerResponse).GenericError(http.StatusUnauthorized, "Unable to authorize user, no user informtaion found in context"))
			return
		}
		email := userInfo.(*authc.UserInfo).Email
		sub := userInfo.(*authc.UserInfo).Subject

		rSub := userInfo.(*authc.UserInfo).Roles
		rObj := c.Request.URL.Path
		rAct := c.Request.Method

		roles := utils.Map(rSub, func(s string) string {
			return constants.CasbinRolePrefix + s
		})

		// Check the role is allowed to access the path with the action
		for _, role := range roles {
			allowed1, err1 := e.Enforce(role, rObj, rAct)
			allowed = allowed || allowed1
			err = errr.Join(err, err1)
			if allowed {
				break
			}
		}

		if err != nil {
			log.Warn("Unable to authorize user (%s/%s): %s", email, sub, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.
				NewServerResponse(model.OkdpServerResponse).GenericError(http.StatusUnauthorized, err.Error()))
			return
		}
		if !allowed {
			log.Warn("User (%s/%s) not allowed to execute the action", email, sub)
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.
				NewServerResponse(model.OkdpServerResponse).GenericError(http.StatusUnauthorized, "Unauthorized action"))
			return
		}

		c.Next()
	}

}

func writeFilesToTmp(modelStr, policyStr string) (modelFilePath, policyFilePath string, err error) {
	modelFilePath = filepath.Join("/tmp", "authz-model.conf")
	policyFilePath = filepath.Join("/tmp", "authz-policy.csv")

	// Write model string to the model.conf file
	err = os.WriteFile(modelFilePath, []byte(modelStr), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write model file: %v", err)
	}

	// Write policy string to the policy.csv file
	err = os.WriteFile(policyFilePath, []byte(policyStr), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write policy file: %v", err)
	}

	return modelFilePath, policyFilePath, nil
}
