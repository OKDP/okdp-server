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
	"github.com/okdp/okdp-server/internal/utils"
)

type ComponentRelease _api.ComponentRelease
type ComponentReleases []*ComponentRelease
type FlatComponent _api.FlatComponent
type FlatComponents []*FlatComponent

func (r *ComponentReleases) Flatten() *FlatComponents {
	var flatComponents = make(FlatComponents, 0, 100)
	for _, c := range *r {
		extC := &FlatComponent{
			ComponentName:        c.Spec.Component.Ref.Name,
			ComponentVersion:     c.Spec.Component.Ref.Version,
			PackageName:          c.Spec.HelmReleaseName,
			PackageVersion:       c.Spec.Component.Source.Version,
			ComponentReleaseName: c.Spec.Name,
			Enabled:              c.Spec.Enabled,
			Suspended:            c.Spec.Component.Suspended,
			Protected:            c.Spec.Component.Protected,
			Catalogs:             utils.ArrayNullToEmpty(c.Status.Catalogs),
		}
		flatComponents = append(flatComponents, extC)
	}
	return &flatComponents
}

func (f *FlatComponents) AddTemplateReleaseInfo(tri *TemplateReleaseInfo) *FlatComponents {
	for _, c := range *f {
		summary, found := tri.GetByComponentReleaseName(c.ComponentReleaseName)
		if found {
			c.TemplateName = summary.TemplateName
			c.TemplateVersion = summary.TemplateVersion
			c.TemplateReleaseName = summary.TemplateReleaseName
			c.Catalogs = append(c.Catalogs, summary.Catalogs...)
		}
	}
	return f
}

func (f *FlatComponents) ConvertToService() *Services {
	var services Services

	componentsByServiceName := make(map[string][]FlatComponent)
	for _, c := range *f {
		var serviceName string
		if c.TemplateReleaseName != "" {
			serviceName = c.TemplateReleaseName
		} else {
			serviceName = c.ComponentReleaseName
		}
		componentsByServiceName[serviceName] = append(componentsByServiceName[serviceName], *c)
	}

	for serviceName, components := range componentsByServiceName {
		services = append(services, &Service{
			Name:           serviceName,
			IsComposition:  len(components) > 1,
			FlatComponents: components,
		})
	}

	return &services

}
