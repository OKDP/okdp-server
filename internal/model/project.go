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

package model

import (
	"time"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type Project _api.Project

func (p Project) ToNamespace() *Namespace {
	labels := map[string]string{
		"okdp.io/project": "true",
	}
	annotations := map[string]string{}

	if p.Description != nil {
		annotations["okdp.io/description"] = *p.Description
	}
	if p.DisplayName != nil {
		annotations["okdp.io/display-name"] = *p.DisplayName
	}
	if p.Environment != nil {
		annotations["okdp.io/environment"] = *p.Environment
	}

	return &Namespace{
		ApiVersion: "v1",
		Kind:       "Namespace",
		Metadata: struct {
			Annotations       *map[string]string `json:"annotations,omitempty"`
			CreationTimestamp *time.Time         `json:"creationTimestamp,omitempty"`
			Labels            *map[string]string `json:"labels,omitempty"`
			Name              string             `json:"name"`
		}{
			Name:              p.Name,
			Labels:            &labels,
			Annotations:       &annotations,
			CreationTimestamp: p.CreationTimestamp,
		},
	}
}
