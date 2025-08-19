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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

func (c KubeClient) ListNamespaces(ctx context.Context) ([]*model.Namespace, *model.ServerResponse) {
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

	namespaces := []*model.Namespace{}
	for _, ns := range namespaceList.Items {
		if exclude[ns.Name] || strings.HasPrefix(ns.Name, "kube-") {
			continue
		}
		namespaces = append(namespaces, model.ToNamespace(ns))
	}

	return namespaces, nil
}

func (c KubeClient) GetNamespaceByName(ctx context.Context, clusterID string, name string) (*model.Namespace, *model.ServerResponse) {
	var ns corev1.Namespace
	key := ctrlclient.ObjectKey{Name: name}

	err := c.Get(ctx, key, &ns)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, model.
				NewServerResponse(model.K8sClusterResponse).
				NotFoundError("Namespace '%s' not found on clusterId '%s'", name, clusterID)
		}
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get namespace '%s' on clusterId '%s', details: '%s'", name, clusterID, err.Error())
	}

	return model.ToNamespace(ns), nil
}

func (c KubeClient) CreateNamespace(ctx context.Context, namespace *model.Namespace) *model.ServerResponse {
	name := namespace.Metadata.Name
	if name == "" {
		return model.NewServerResponse(model.K8sClusterResponse).
			BadRequest("Namespace name must be provided")
	}

	var existing corev1.Namespace
	err := c.Get(ctx, ctrlclient.ObjectKey{Name: name}, &existing)
	if err == nil {
		return model.NewServerResponse(model.K8sClusterResponse).
			ConflictError("Namespace '%s' already exists '%s'", name)
	} else if !apierrors.IsNotFound(err) {
		return model.NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to check for namespace '%s': %v", name, err)
	}

	newNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      nil,
			Annotations: nil,
		},
	}

	if namespace.Metadata.Labels != nil {
		newNS.Labels = *namespace.Metadata.Labels
	}
	if namespace.Metadata.Annotations != nil {
		newNS.Annotations = *namespace.Metadata.Annotations
	}

	if err := c.Create(ctx, newNS); err != nil {
		return model.NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to create namespace '%s': %v", name, err)
	}

	return model.NewServerResponse(model.K8sClusterResponse).
		Created("Namespace '%s' created successfully", name)
}

func (c KubeClient) UpdateNamespace(ctx context.Context, namespace *model.Namespace) *model.ServerResponse {
	name := namespace.Metadata.Name
	if name == "" {
		return model.NewServerResponse(model.K8sClusterResponse).
			BadRequest("Namespace name must be provided")
	}

	var existing corev1.Namespace
	err := c.Get(ctx, ctrlclient.ObjectKey{Name: name}, &existing)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return model.NewServerResponse(model.K8sClusterResponse).
				NotFoundError("Namespace '%s' not found", name)
		}
		return model.NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get namespace '%s': %v", name, err)
	}

	utils.MergeLabels(&existing, namespace.Metadata.Labels)
	utils.MergeAnnotations(&existing, namespace.Metadata.Annotations)

	if err := c.Update(ctx, &existing); err != nil {
		return model.NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to update namespace '%s': %v", name, err)
	}

	return model.NewServerResponse(model.K8sClusterResponse).
		Updated("Namespace '%s' updated successfully", name)
}

// TODO: delete all resources inside the namespace + patch those resources to remove finalizers
func (c KubeClient) DeleteNamespace(ctx context.Context, namespace string) *model.ServerResponse {
	if namespace == "" {
		return model.NewServerResponse(model.K8sClusterResponse).
			BadRequest("Namespace name must be provided")
	}

	var existing corev1.Namespace
	err := c.Get(ctx, ctrlclient.ObjectKey{Name: namespace}, &existing)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return model.NewServerResponse(model.K8sClusterResponse).
				NotFoundError("Namespace '%s' not found", namespace)
		}
		return model.NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get namespace '%s': %v", namespace, err)
	}

	if err := c.Delete(ctx, &existing); err != nil {
		return model.NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to delete namespace '%s': %v", namespace, err)
	}

	return model.NewServerResponse(model.K8sClusterResponse).
		Deleted("Namespace '%s' deleted successfully", namespace)
}
