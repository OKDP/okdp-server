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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_component "github.com/okdp/okdp-server/api/openapi/v3/_api/components"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/services"
)

type IComponentController struct {
	componentService *services.ComponentService
}

func ComponentController() *IComponentController {
	componentService, err := services.NewComponentService()
	if err != nil {
		return nil
	}
	return &IComponentController{
		componentService: componentService,
	}
}

func (r IComponentController) ListComponents(c *gin.Context, params _component.ListComponentsParams) {
	components, err := r.componentService.List(params.IncludeRawSpec)
	if err != nil {
		log.Error("Unable to get data from kad: %s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, components)
}

func (r IComponentController) GetComponents(c *gin.Context, componentid string, params _component.GetComponentsParams) {
	components, err := r.componentService.Get(componentid, params.IncludeRawSpec)
	if err != nil {
		log.Error("Unable to get data from kad: %s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(components) == 0 {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Component with id %s not found", componentid))
		return
	}
	c.JSON(http.StatusOK, components)

}



