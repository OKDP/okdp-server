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

package kad

import (
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/kad/client"
	"github.com/okdp/okdp-server/internal/model"
)

type ComponentReleaseClient struct {
	KAD *client.KadClients
}

func NewComponentReleaseClient() *ComponentReleaseClient {
	return &ComponentReleaseClient{
		KAD: client.GetClients(),
	}
}

func (c ComponentReleaseClient) Get(kadInstanceId string, name string, catalog *string) (*model.ComponentRelease, error) {
	kadClient := c.KAD.ID(kadInstanceId)
	if kadClient == nil {
		return nil, client.InvalidInstanceError(kadInstanceId)
	}
	req := kadClient.NewRequest(constants.ComponentReleaseURL)
	if catalog != nil {
		req = req.SetQueryParam("catalog", *catalog)
	}
	return client.DoGet[model.ComponentRelease](req)
}

func (c ComponentReleaseClient) List(kadInstanceId string, catalog *string) (*model.ComponentReleases, error) {
	kadClient := c.KAD.ID(kadInstanceId)
	if kadClient == nil {
		return nil, client.InvalidInstanceError(kadInstanceId)
	}
	req := kadClient.NewRequest(constants.ComponentReleaseURL)
	if catalog != nil {
		req = req.SetQueryParam("catalog", *catalog)
	}
	return client.DoGet[model.ComponentReleases](req)
}
