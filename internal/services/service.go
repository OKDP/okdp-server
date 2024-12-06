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
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/errors"
)

type Service struct {
	componentReleaseService *ComponentReleaseService
	templateReleaseService  *TemplateReleaseService
}

func NewService() (*Service, error) {
	return &Service{
		componentReleaseService: NewComponentReleaseService(),
		templateReleaseService:  NewTemplateReleaseService(),
	}, nil
}

func (s Service) List(kadInstanceId string, catalog *string) (*model.Services, *errors.ServerError) {

	componentReleases, err := s.componentReleaseService.List(kadInstanceId, catalog)
	if err != nil {
		return nil, err
	}
	temmplateReleases, err := s.templateReleaseService.List(kadInstanceId, catalog)
	if err != nil {
		return nil, err
	}

	tri := temmplateReleases.GroupTemplateReleaseInfoByComponentRelease()

	return componentReleases.Flatten().AddTemplateReleaseInfo(tri).ConvertToService(), nil
}


