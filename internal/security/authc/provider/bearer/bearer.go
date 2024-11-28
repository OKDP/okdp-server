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

package bearer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/model/auth"
	"golang.org/x/net/context"
)

var (
	rolesAttributePath  string
	groupsAttributePath string
)

type BearerProvider struct {
	context.Context
	*oidc.IDTokenVerifier
}

func NewProvider(bearerConf config.BearerAuth) BearerProvider {
	ctx := context.Background()
	config := &oidc.Config{
		SkipClientIDCheck:          true,
		InsecureSkipSignatureCheck: bearerConf.SkipSignatureCheck,
		SkipIssuerCheck:            bearerConf.SkipIssuerCheck,
	}
	verifier := oidc.NewVerifier(bearerConf.IssuerUri, oidc.NewRemoteKeySet(ctx, bearerConf.JwksURL), config)

	rolesAttributePath = bearerConf.RolesAttributePath
	groupsAttributePath = bearerConf.GroupsAttributePath

	return BearerProvider{ctx, verifier}
}

// Auth returns a middleware which authenticates the user with a bearer authentication 
// and propagates the user info (roles/groups) into the autorization provider.
func (p *BearerProvider) Auth() []gin.HandlerFunc {
	return []gin.HandlerFunc{p.authenticate()}
}

func (p *BearerProvider) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userInfo model.UserInfo
		)

		authorization := c.Request.Header.Get("Authorization")
		accessToken := strings.TrimPrefix(authorization, "Bearer ")
		err := p.verifyAccessToken(accessToken)
		if err != nil {
			log.Warn("Failed to verify access Token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to verify access Token: " + err.Error()})
			return
		}

		userInfo, err = p.getUserInfo(accessToken)
		if err != nil {
			log.Warn("Unable to get user roles/groups from access token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unable to get user roles/groups from access token: " + err.Error()})
			return
		}
		log.Debug("Successfully authenticated user : %s", userInfo.AsJsonString())
		c.Set(constants.OAuth2UserInfo, &userInfo)
		c.Next()
	}
}

func (p *BearerProvider) verifyAccessToken(accessToken string) error {
	if len(strings.TrimSpace(accessToken)) == 0 {
		return fmt.Errorf("no access token found in the Authorization header")
	}
	_, err := p.Verify(p.Context, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (p *BearerProvider) getUserInfo(accessToken string) (model.UserInfo, error) {
	token := &model.Token {AccessToken: accessToken}
	return token.GetUserInfo(rolesAttributePath, groupsAttributePath)
}
