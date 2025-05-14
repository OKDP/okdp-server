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

package oci

import (
	"github.com/okdp/okdp-server/internal/integrations/oci/client"
	"github.com/okdp/okdp-server/internal/model"
)

type RepoCatalog struct {
	r *client.RepositoryClients
}

func NewRepoCatalog() *RepoCatalog {
	return &RepoCatalog{
		r: client.GetClients(),
	}
}

func (r RepoCatalog) ListCatalogs() []*model.Catalog {
	return client.ListCatalogs()
}

func (r RepoCatalog) GetCatalog(catalogID string) (*model.Catalog, *model.ServerResponse) {
	return client.GetCatalog(catalogID)
}

func (r RepoCatalog) GetPackages(catalogID string) ([]*model.Package, *model.ServerResponse) {
	return client.GetPackages(catalogID)
}

func (r RepoCatalog) GetPackage(catalogID string, name string) (*model.Package, *model.ServerResponse) {
	return client.GetPackage(catalogID, name)
}

func (r RepoCatalog) GetPackageDefinition(catalogID string, name string, version string) (map[string]interface{}, *model.ServerResponse) {
	return client.GetPackageDefinition(catalogID, name, version)
}
