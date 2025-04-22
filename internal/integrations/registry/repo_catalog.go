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

package registry

import (
	"github.com/okdp/okdp-server/internal/integrations/registry/client"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/servererrors"
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

func (r RepoCatalog) GetCatalog(catalogID string) (*model.Catalog, *servererrors.ServerError) {
	return client.GetCatalog(catalogID)
}

func (r RepoCatalog) GetPackages(catalogID string) ([]*model.Package, *servererrors.ServerError) {
	return client.GetPackages(catalogID)
}

func (r RepoCatalog) GetPackageByName(catalogID string, name string) (*model.Package, *servererrors.ServerError) {
	return client.GetPackageByName(catalogID, name)
}

func (r RepoCatalog) GetPackageDefinition(catalogID string, name string, version string) (map[string]interface{}, *servererrors.ServerError) {
	return client.GetPackageDefinition(catalogID, name, version)
}
