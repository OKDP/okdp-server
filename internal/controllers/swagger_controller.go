/*
 *    Copyright 2026 The OKDP Authors.
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
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
	"github.com/okdp/okdp-server/internal/config"
)

func Swagger(swaggerConf config.Swagger) gin.HandlerFunc {
	return func(c *gin.Context) {
		swagger, err := api.GetSwagger()
		if err != nil {
			panic("Error loading swagger spec: " + err.Error())
		}
		if swagger.Components.SecuritySchemes == nil {
			swagger.Components.SecuritySchemes = make(map[string]*openapi3.SecuritySchemeRef)
		}
		for key, value := range swaggerConf.SecuritySchemes {
			if swagger.Components.SecuritySchemes[key] == nil {
				swagger.Components.SecuritySchemes[key] = &openapi3.SecuritySchemeRef{}
			}
			swagger.Components.SecuritySchemes[key].Value = value
		}
		swagger.Security = swaggerConf.Security
		c.JSON(http.StatusOK, swagger)
	}
}
