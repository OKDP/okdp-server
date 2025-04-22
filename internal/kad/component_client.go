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
	"github.com/okdp/okdp-server/internal/servererrors"
)

type ComponentClient struct {
	KAD *client.KadClients
}

func NewComponentClient() *ComponentClient {
	return &ComponentClient{
		KAD: client.GetClients(),
	}
}

func (c ComponentClient) Get(kadInstanceID string, name string, catalog *string) (*model.Component, *servererrors.ServerError) {
	kadClient, err := c.KAD.ID(kadInstanceID)
	if err != nil {
		return nil, err
	}
	req := kadClient.NewRequest(constants.ComponentURL + "/" + name)
	if catalog != nil {
		req = req.SetQueryParam("catalog", *catalog)
	}
	return client.DoGet[model.Component](req)
}

func (c ComponentClient) List(kadInstanceID string, catalog *string) (*model.Components, *servererrors.ServerError) {
	kadClient, err := c.KAD.ID(kadInstanceID)
	if err != nil {
		return nil, err
	}
	req := kadClient.NewRequest(constants.ComponentURL)
	if catalog != nil {
		req = req.SetQueryParam("catalog", *catalog)
	}
	return client.DoGet[model.Components](req)
}
