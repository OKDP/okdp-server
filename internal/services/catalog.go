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

package services

import (
	"github.com/okdp/okdp-server/internal/integrations/registry"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/servererrors"
)

type CatalogService struct {
	catalog *registry.RepoCatalog
}

func NewCatalogService() *CatalogService {
	return &CatalogService{
		catalog: registry.NewRepoCatalog(),
	}
}

func (s CatalogService) ListCatalogs() []*model.Catalog {
	return s.catalog.ListCatalogs()
}

func (s CatalogService) GetCatalog(catalogID string) (*model.Catalog, *servererrors.ServerError) {
	return s.catalog.GetCatalog(catalogID)
}

func (s CatalogService) GetPackages(catalogID string) ([]*model.Package, *servererrors.ServerError) {
	return s.catalog.GetPackages(catalogID)
}

func (s CatalogService) GetPackageByName(catalogID string, name string) (*model.Package, *servererrors.ServerError) {
	return s.catalog.GetPackageByName(catalogID, name)
}

func (s CatalogService) GetPackageDefinition(catalogID string, name string, version string) (map[string]interface{}, *servererrors.ServerError) {
	return s.catalog.GetPackageDefinition(catalogID, name, version)
}

// func (s CatalogService) GetDefinition(name string) (*model.Catalog, *errors.ServerError) {
// 	return s.catalog.GetDefinition(name)
// }

// func (s CatalogService) GetTags(name string) (*model.Catalog, *errors.ServerError) {
// 	return s.catalog.GetTags(name)
// }
