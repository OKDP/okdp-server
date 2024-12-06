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

package authc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/errors"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/security/authc/provider/basic"
	"github.com/okdp/okdp-server/internal/security/authc/provider/bearer"
	"github.com/okdp/okdp-server/internal/security/authc/provider/oidc"
)

func Authenticator(authNConfig config.AuthN) []gin.HandlerFunc {
	var handlers = []gin.HandlerFunc{}
	log.Info("Loading authentication providers: ", authNConfig.Provider)
	for _, provider := range authNConfig.Provider {
		switch provider {
		case "basic":
			p, err := basic.NewProvider(config.GetAppConfig().Security.AuthN.Basic)
			if err != nil {
				log.Panic("Unable to get a basic auth provider: %w", err)
			}
			handlers = append(handlers, p.Auth()...)
		case "openid":
			p, err := oidc.NewProvider(config.GetAppConfig().Security.AuthN.OpenID)
			if err != nil {
				log.Panic("Unable to get an OIDC provider: %w", err)
			}
			handlers = append(handlers, p.Auth()...)
		case "bearer":
			p := bearer.NewProvider(config.GetAppConfig().Security.AuthN.Bearer)
			handlers = append(handlers, p.Auth()...)
		default:
			log.Panic("Unknown authentication provider: %s", provider)
		}
	}
	// Ensure the user was authenticated by any of the autentication provides
	handlers = append(handlers, ensureUserAuthenticated())
	return handlers
}

// Ensure the user was authenticated by any of the autentication providers
func ensureUserAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, found := c.Get(constants.OAuth2UserInfo)
		if !found {
			log.Warn("Failed to authenticate user")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.OfType(errors.OkdpServer).GenericError(http.StatusUnauthorized, "Authentication failed"))
			return
		}
	}
}
