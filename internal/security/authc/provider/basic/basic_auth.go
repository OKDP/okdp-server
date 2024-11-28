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

package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/model/auth"
)

type BasicProvider struct {
	accounts gin.Accounts
	roles    map[string][]string
}

func NewProvider(basicUsers []config.BasicAuth) (*BasicProvider, error) {
	accounts := make(gin.Accounts)
	userRoles := make(map[string][]string)
	for _, user := range basicUsers {
		accounts[user.Login] = user.Password
		userRoles[user.Login] = user.Roles
	}
	return &BasicProvider{accounts: accounts, roles: userRoles}, nil
}


// Auth returns a middleware which authenticates the user with a basic authentication 
// and returns a second middleware which propagates the user info (roles) into the autorization provider.
func (p *BasicProvider) Auth() []gin.HandlerFunc {
	return []gin.HandlerFunc{p.authenticate(), p.setUserInfo()}
}

// Authenticate User
func (p *BasicProvider) authenticate() gin.HandlerFunc {
	return gin.BasicAuth(p.accounts)
}

// Propagate userInfo to authorization
func (p *BasicProvider) setUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		login, ok := c.Get(gin.AuthUserKey)
		if ok {
			c.Set(constants.OAuth2UserInfo, &model.UserInfo{Login: login.(string), Roles: p.getUserRoles(login.(string))})
		}
	}
}

func (p *BasicProvider) getUserRoles(login string) []string {
	roles, found := p.roles[login]
	if found {
		return roles
	}
	return []string{}
}

