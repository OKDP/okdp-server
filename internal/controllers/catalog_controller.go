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
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/services"
)

type ICatalogController struct {
	catalogService *services.CatalogService
}

func CatalogController() *ICatalogController {
	return &ICatalogController{
		catalogService: services.NewCatalogService(),
	}
}

func (r ICatalogController) ListCatalogs(c *gin.Context, kadInstanceID string) {
	catalogs, err := r.catalogService.List(kadInstanceID)
	if err != nil {
		log.Error("Unable to list Catalogs on kad instance %s, details: %+v", kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, catalogs)
}

func (r ICatalogController) GetCatalog(c *gin.Context, kadInstanceID string, name string) {
	catalog, err := r.catalogService.Get(kadInstanceID, name)
	if err != nil {
		log.Error("Unable to get Catalog info '%s' on kad instance %s, details: %+v", name, kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}
	if catalog == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Component with id %s not found", name))
		return
	}
	c.JSON(http.StatusOK, catalog)

}
