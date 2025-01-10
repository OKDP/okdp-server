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
	"net/http"

	"github.com/gin-gonic/gin"
	_component "github.com/okdp/okdp-server/api/openapi/v3/_api/components"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/services"
)

type IComponentController struct {
	componentService *services.ComponentService
}

func ComponentController() *IComponentController {
	return &IComponentController{
		componentService: services.NewComponentService(),
	}
}

func (r IComponentController) ListComponents(c *gin.Context, kadInstanceID string, params _component.ListComponentsParams) {
	components, err := r.componentService.List(kadInstanceID, params.Catalog)
	if err != nil {
		log.Error("Unable to list Components on kad instance %s, details: %+v", kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, components)
}

func (r IComponentController) GetComponentsByName(c *gin.Context, kadInstanceID string, name string, params _component.GetComponentsByNameParams) {
	// Workaround - Kad should return a components by name
	log.Info("Get component by name '%s' and catalog '%s' on kad instanceId '%s'", name, params.Catalog, kadInstanceID)
	components, err := r.componentService.List(kadInstanceID, params.Catalog)
	if err != nil {
		log.Error("Unable to list Components on kad instance %s, details: %+v", kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}
	var filtered model.Components = []*model.Component{}
	for _, component := range *components {
		if component.Spec.Name == name {
			filtered = append(filtered, component)
		}
	}
	c.JSON(http.StatusOK, filtered)
}
