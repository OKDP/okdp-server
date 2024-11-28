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
	"github.com/okdp/okdp-server/internal/kad"
	model "github.com/okdp/okdp-server/internal/model/kad"
	response "github.com/okdp/okdp-server/internal/kad/response"
)

type IKadComponent interface {
	List() []string
	Get(string) string
}

type KadComponent struct {
	kadClient *kad.KadClient
}

func NewKadComponent() (*KadComponent, error) {
	return &KadComponent{kad.GetClient()}, nil
}

func (k *KadComponent) Get(id string) ([]response.KadResponseWrapper[model.Component], error) {
	// kad implementation for GET by id??????
	resp, err := k.kadClient.Get("1")
	if err != nil {
		return nil, err
	}
	return model.ToComponents(resp)
}

func (k *KadComponent) List() ([]response.KadResponseWrapper[model.Component], error) {
	resp, err := k.kadClient.List()
	if err != nil {
		return nil, err
	}
	return model.ToComponents(resp)
}
