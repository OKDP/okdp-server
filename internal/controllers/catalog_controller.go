/*
 *    Copyright 2025 okdp.io
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
	log "github.com/okdp/okdp-server/internal/common/logging"
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

func (r ICatalogController) ListCatalogs(c *gin.Context) {
	catalogs := r.catalogService.ListCatalogs()
	c.JSON(http.StatusOK, catalogs)
}

func (r ICatalogController) GetCatalog(c *gin.Context, catalogID string) {
	catalog, err := r.catalogService.GetCatalog(catalogID)
	if err != nil {
		log.Error("Unable to find the Catalog with ID '%s', details: %+v", catalogID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, catalog)
}

func (r ICatalogController) ListPackages(c *gin.Context, catalogID string) {
	packages, err := r.catalogService.GetPackages(catalogID)
	if err != nil {
		log.Error("Unable to find the packages with Catalog ID '%s', details: %+v", catalogID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, packages)
}

func (r ICatalogController) GetPackage(c *gin.Context, catalogID string, name string) {
	result, err := r.catalogService.GetPackage(catalogID, name)
	if err != nil {
		log.Error("Unable to find the package '%s' with Catalog ID '%s', details: %+v", name, catalogID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (r ICatalogController) GetPackageVersions(c *gin.Context, catalogID string, name string) {
	r.GetPackage(c, catalogID, name)
}

func (r ICatalogController) GetPackageDefinition(c *gin.Context, catalogID string, name string, version string) {
	definition, err := r.catalogService.GetPackageDefinition(catalogID, name, version)
	if err != nil {
		log.Error("Unable to find the package definition for package '%s:%s' with Catalog ID '%s', details: %+v", name, version, catalogID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, definition)
}

func (r ICatalogController) GetPackageSchema(c *gin.Context, catalogID string, name string, version string) {
	definition, err := r.catalogService.GetPackageDefinition(catalogID, name, version)
	if err != nil {
		log.Error("Unable to find the package definition for package '%s:%s' with Catalog ID '%s', details: %+v", name, version, catalogID, err)
		c.JSON(err.Status, err)
		return
	}
	schema, ok := definition["schema"]
	if !ok {
		c.JSON(http.StatusOK, struct{}{})
		return
	}
	c.JSON(http.StatusOK, schema)
}
