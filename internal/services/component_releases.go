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
 )
 
 
 type ComponentReleaseService struct {
	 componentRelease *kad.ComponentReleaseClient
 }
 
 func NewComponentReleaseService() *ComponentReleaseService {
	 return &ComponentReleaseService{
		componentRelease: kad.NewComponentReleaseClient(),
	 }
 }
 
 func (s ComponentReleaseService) Get(kadInstanceId string, name string, catalog *string) (*model.ComponentRelease, error) {
	 return s.componentRelease.Get(kadInstanceId, name, catalog)
 }
 
 func (s ComponentReleaseService) List(kadInstanceId string, catalog *string) (*model.ComponentReleases, error) {
	 return s.componentRelease.List(kadInstanceId, catalog)
 }
 