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

package client

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/okdp/okdp-server/internal/model"
)

func (c KubeClient) ListNamespaces(ctx context.Context) ([]string, *model.ServerResponse) {
	var namespaceList corev1.NamespaceList
	err := c.List(ctx, &namespaceList)
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to list Kubernetes namespaces on clusterId '%s', details: '%s'", c.clusterID, err.Error())
	}

	exclude := map[string]bool{
		"local-path-storage": true,
	}

	namespaces := []string{}
	for _, ns := range namespaceList.Items {
		if exclude[ns.Name] || strings.HasPrefix(ns.Name, "kube-") {
			continue
		}
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}
