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

package controllers

import (
	"github.com/gin-gonic/gin"
	_catalog "github.com/okdp/okdp-server/api/openapi/v3/_api/catalogs"
	_componentrelease "github.com/okdp/okdp-server/api/openapi/v3/_api/componentreleases"
	_component "github.com/okdp/okdp-server/api/openapi/v3/_api/components"
	_kad "github.com/okdp/okdp-server/api/openapi/v3/_api/kad"
	_service "github.com/okdp/okdp-server/api/openapi/v3/_api/services"
	_templaterelease "github.com/okdp/okdp-server/api/openapi/v3/_api/templatereleases"
	_user "github.com/okdp/okdp-server/api/openapi/v3/_api/users"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
)

type Router struct {
	*gin.Engine
}

type Group struct {
	*gin.RouterGroup
}

func (g *Group) RegisterControllers() {
	_user.RegisterHandlers(g, UserProfileController())
	_catalog.RegisterHandlers(g, CatalogController())

	_componentrelease.RegisterHandlers(g, ComponentReleaseController())
	_templaterelease.RegisterHandlers(g, TemplateReleaseController())
	_component.RegisterHandlers(g, ComponentController())
	_service.RegisterHandlers(g, ServiceController())
	_kad.RegisterHandlers(g, KadController())
}

func (r *Router) RegisterSwaggerAPIDoc(swaggerConf config.Swagger) {
	r.GET(constants.SwaggerAPIDocsURI, Swagger(swaggerConf))
}

func (r *Router) RegisterHealth() {
	r.GET(constants.HealthzURI, Healthz)
	r.GET(constants.ReadinessURI, Readiness)
}
