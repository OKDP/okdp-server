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

package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/controllers"
	"github.com/okdp/okdp-server/internal/security"
	"github.com/okdp/okdp-server/internal/security/authc"
	"github.com/okdp/okdp-server/internal/security/authz"
)

func NewOKDPServer(config *config.ApplicationConfig) *http.Server {

	gin.SetMode(config.Server.Mode)
	r := &controllers.Router{gin.New()}                               //nolint:all
	apiV1 := &controllers.Group{r.Group(constants.OkdpServerBaseUrl)} //nolint:all

	r.Use(log.Logger()...)
	r.Use(gin.Recovery())
	
	// Apply http security (cors, headers, etc) on the root path (/) and the groups (/api/v1)
	r.Use(security.HttpSecurity(config.Security)...)
	// https://github.com/gin-gonic/gin/issues/3546
	apiV1.Use(security.HttpSecurity(config.Security)...)

	// Authentication
	apiV1.Use(authc.Authenticator(config.Security.AuthN)...)
	// Authorization
	apiV1.Use(authz.Authorizer(config.Security.AuthZ))

	// Register Controllers
	apiV1.RegisterControllers()
	r.RegisterSwaggerApiDoc()

	server := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", config.Server.ListenAddress, config.Server.Port),
	}

	return server
}
