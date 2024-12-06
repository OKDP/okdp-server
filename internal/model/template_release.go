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

package model

import (
	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type TemplateRelease _api.TemplateRelease
type TemplateReleases []TemplateRelease

type TemplateReleaseSummary struct {
	TemplateName        string
	TemplateVersion     string
	TemplateReleaseName string
	Catalogs            []string
}

type TemplateReleaseInfo struct {
	mapping map[string]TemplateReleaseSummary
}

func (r *TemplateReleases) GroupTemplateReleaseInfoByComponentRelease() *TemplateReleaseInfo {
	m := &TemplateReleaseInfo{
		mapping: make(map[string]TemplateReleaseSummary),
	}

	for _, t := range *r {
		for _, c := range t.Status.Children {
			m.mapping[c] = TemplateReleaseSummary{
				TemplateName:        t.Spec.Template.Ref.Name,
				TemplateVersion:     t.Spec.Template.Ref.Version,
				TemplateReleaseName: t.Spec.Name,
				Catalogs:            t.Status.Catalogs,
			}
		}
	}

	return m
}

func (r *TemplateReleaseInfo) GetByComponentReleaseName (name string) (*TemplateReleaseSummary, bool) {
    s, found := r.mapping[name]
	if !found {
        return nil, false
	}
	return &s, true
}

