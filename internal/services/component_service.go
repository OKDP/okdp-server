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

package services

import (
	kad "github.com/okdp/okdp-server/internal/kad/webserver/component"
	model "github.com/okdp/okdp-server/internal/model/kad"
	response "github.com/okdp/okdp-server/internal/kad/response"
	"github.com/okdp/okdp-server/internal/utils"
)

type IComponentService interface {
	List() []string
	Get(string) string
}

type ComponentService struct {
	kadComponent *kad.KadComponent
}

func NewComponentService() (*ComponentService, error) {
	kadComponent, err := kad.NewKadComponent()
	if err != nil {
		return nil, err
	}
	return &ComponentService{
		kadComponent: kadComponent,
	}, nil
}

func (s ComponentService) Get(id string, includeRawSpec bool) ([]response.KadResponseWrapper[model.Component], error) {
	resp, err := s.kadComponent.Get(id)
	return withRawSpec(resp, includeRawSpec), err
}

func (s ComponentService) List(includeRawSpec bool) ([]response.KadResponseWrapper[model.Component], error) {
	resp, err := s.kadComponent.List()
	return withRawSpec(resp, includeRawSpec), err
}

func withRawSpec (components []response.KadResponseWrapper[model.Component], includeRawSpec bool) []response.KadResponseWrapper[model.Component]{
	if !includeRawSpec {
		components = utils.Map(components, 
		  func (c response.KadResponseWrapper[model.Component]) response.KadResponseWrapper[model.Component]{ 
			c.Raw = nil; return c 
		})
	  }
	return components
}

