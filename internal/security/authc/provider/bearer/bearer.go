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
	"github.com/okdp/okdp-server/internal/common/constants"
	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/model"
	authc "github.com/okdp/okdp-server/internal/security/authc/model"
	"golang.org/x/net/context"
)

var (
	rolesAttributePath  string
	groupsAttributePath string
)

type Provider struct {
	context.Context
	*oidc.IDTokenVerifier
}

func NewProvider(bearerConf config.BearerAuth) Provider {
	ctx := context.Background()
	config := &oidc.Config{
		SkipClientIDCheck:          true,
		InsecureSkipSignatureCheck: bearerConf.SkipSignatureCheck,
		SkipIssuerCheck:            bearerConf.SkipIssuerCheck,
	}
	verifier := oidc.NewVerifier(bearerConf.IssuerURI, oidc.NewRemoteKeySet(ctx, bearerConf.JwksURL), config)

	rolesAttributePath = bearerConf.RolesAttributePath
	groupsAttributePath = bearerConf.GroupsAttributePath

	return Provider{ctx, verifier}
}

// Auth returns a middleware which authenticates the user with a bearer authentication
// and propagates the user info (roles/groups) into the autorization provider.
func (p *Provider) Auth() []gin.HandlerFunc {
	return []gin.HandlerFunc{p.authenticate()}
}

func (p *Provider) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userInfo authc.UserInfo
		)

		authorization := c.Request.Header.Get("Authorization")
		accessToken := strings.TrimPrefix(authorization, "Bearer ")
		err := p.verifyAccessToken(accessToken)
		if err != nil {
			log.Warn("Failed to verify access Token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.
				NewServerResponse(model.OkdpServerResponse).GenericError(http.StatusUnauthorized, "Failed to verify access Token: "+err.Error()))
			return
		}

		userInfo, err = p.getUserInfo(accessToken)
		if err != nil {
			log.Warn("Unable to get user roles/groups from access token: %w", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.
				NewServerResponse(model.OkdpServerResponse).GenericError(http.StatusUnauthorized, "Unable to get user roles/groups from access token: "+err.Error()))
			return
		}
		log.Debug("Successfully authenticated user : %s", userInfo.AsJSONString())
		c.Set(constants.OAuth2UserInfo, &userInfo)
		c.Next()
	}
}

func (p *Provider) verifyAccessToken(accessToken string) error {
	if len(strings.TrimSpace(accessToken)) == 0 {
		return fmt.Errorf("no access token found in the Authorization header")
	}
	_, err := p.Verify(p.Context, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) getUserInfo(accessToken string) (authc.UserInfo, error) {
	token := &authc.Token{AccessToken: accessToken}
	return token.GetUserInfo(rolesAttributePath, groupsAttributePath)
}
