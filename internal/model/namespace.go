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

	corev1 "k8s.io/api/core/v1"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type Namespace _api.Namespace

func (ns Namespace) ToProject() *Project {
	var displayName, description, environment *string

	if ns.Metadata.Annotations != nil {
		annotations := *ns.Metadata.Annotations
		if val, ok := annotations["okdp.io/display-name"]; ok {
			displayName = &val
		}
		if val, ok := annotations["okdp.io/description"]; ok {
			description = &val
		}
		if val, ok := annotations["okdp.io/environment"]; ok {
			environment = &val
		}
	}

	var status *_api.ProjectStatus
	if ns.Status != nil && ns.Status.Phase != nil {
		s := _api.ProjectStatus(*ns.Status.Phase)
		status = &s
	}

	return &Project{
		Name:              ns.Metadata.Name,
		CreationTimestamp: ns.Metadata.CreationTimestamp,
		DisplayName:       displayName,
		Description:       description,
		Environment:       environment,
		Status:            status,
	}
}

func ToNamespace(ns corev1.Namespace) *Namespace {
	meta := ns.ObjectMeta

	m := &Namespace{
		ApiVersion: "v1",
		Kind:       _api.NamespaceKind("Namespace"),
		Metadata: struct {
			Annotations       *map[string]string `json:"annotations,omitempty"`
			CreationTimestamp *time.Time         `json:"creationTimestamp,omitempty"`
			Labels            *map[string]string `json:"labels,omitempty"`
			Name              string             `json:"name"`
		}{
			Name: meta.Name,
		},
	}

	if len(meta.Annotations) > 0 {
		annotations := make(map[string]string, len(meta.Annotations))
		for k, v := range meta.Annotations {
			annotations[k] = v
		}
		m.Metadata.Annotations = &annotations
	}

	if len(meta.Labels) > 0 {
		labels := make(map[string]string, len(meta.Labels))
		for k, v := range meta.Labels {
			labels[k] = v
		}
		m.Metadata.Labels = &labels
	}

	if !meta.CreationTimestamp.IsZero() {
		t := meta.CreationTimestamp.Time
		m.Metadata.CreationTimestamp = &t
	}

	if ns.Status.Phase != "" {
		phase := _api.NamespaceStatusPhase(ns.Status.Phase)
		m.Status = &struct {
			Phase *_api.NamespaceStatusPhase `json:"phase,omitempty"`
		}{
			Phase: &phase,
		}
	}

	return m
}
