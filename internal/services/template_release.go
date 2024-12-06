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
	 "github.com/okdp/okdp-server/internal/kad"
	 "github.com/okdp/okdp-server/internal/model"
	 "github.com/okdp/okdp-server/internal/errors"
 )
 
 type TemplateReleaseService struct {
	 templateRelease *kad.TemplateReleaseClient
 }
 
 func NewTemplateReleaseService() *TemplateReleaseService {
	 return &TemplateReleaseService{
		templateRelease: kad.NewTemplateReleaseClient(),
	 }
 }
 
 func (s TemplateReleaseService) Get(kadInstanceId string, name string, catalog *string) (*model.TemplateRelease, *errors.ServerError) {
	 return s.templateRelease.Get(kadInstanceId, name, catalog)
 }
 
 func (s TemplateReleaseService) List(kadInstanceId string, catalog *string) (*model.TemplateReleases, *errors.ServerError) {
	 return s.templateRelease.List(kadInstanceId, catalog)
 }
 