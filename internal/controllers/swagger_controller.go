package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
)


func Swagger(c *gin.Context) {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic("Error loading swagger spec: "+err.Error())
	}
	c.JSON(http.StatusOK, swagger)
}


