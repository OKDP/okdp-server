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
