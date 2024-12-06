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
	_services "github.com/okdp/okdp-server/api/openapi/v3/_api/services"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/services"
)

type IServiceController struct {
	service *services.Service
}

func ServiceController() *IServiceController {
	service, err := services.NewService()
	if err != nil {
		return nil
	}
	return &IServiceController{
		service: service,
	}
}

func (r IServiceController) ListServices(c *gin.Context, kadInstanceID string, params _services.ListServicesParams) {
	services, err := r.service.List(kadInstanceID, params.Catalog)
	if err != nil {
		log.Error("Unable to list services, details: %+v", err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, services)
}
