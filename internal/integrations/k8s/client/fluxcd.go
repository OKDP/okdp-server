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

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1b2 "github.com/fluxcd/source-controller/api/v1beta2"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

func (c KubeClient) ListKutomizations(ctx context.Context, namespaces ...string) ([]*kustomizev1.Kustomization, *model.ServerResponse) {
	var kustomizationList kustomizev1.KustomizationList
	if err := c.List(ctx, &kustomizationList); err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to list Kustomizations '%s' ", err.Error())
	}

	filtered := utils.Filter(kustomizationList.Items, func(k kustomizev1.Kustomization) bool {
		return len(namespaces) == 0 || utils.Contains(namespaces, k.Namespace)
	})

	return filtered, nil

}

func (c KubeClient) ListGitRepositories(ctx context.Context, namespaces ...string) ([]*sourcev1.GitRepository, *model.ServerResponse) {

	var repos sourcev1.GitRepositoryList
	if err := c.List(ctx, &repos); err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to list Git Repositories '%s'", err.Error())
	}

	filtered := utils.Filter(repos.Items, func(k sourcev1.GitRepository) bool {
		return utils.Contains(namespaces, k.Namespace)
	})

	return filtered, nil
}

func (c KubeClient) GetGitRepository(ctx context.Context, name string, namespace string) (*sourcev1.GitRepository, *model.ServerResponse) {
	repoKey := k8s.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	var repo sourcev1.GitRepository
	err := c.Get(ctx, repoKey, &repo)
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get fluxcd Git Repository '%s' in namespace '%s', details: '%s'", name, namespace, err.Error())
	}

	return &repo, nil
}

func (c KubeClient) ListOCIRepositories(ctx context.Context, namespaces ...string) ([]*sourcev1b2.OCIRepository, *model.ServerResponse) {

	var repos sourcev1b2.OCIRepositoryList
	if err := c.List(ctx, &repos); err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to list OCI Repositories '%s'", err.Error())
	}

	filtered := utils.Filter(repos.Items, func(k sourcev1b2.OCIRepository) bool {
		return utils.Contains(namespaces, k.Namespace)
	})

	return filtered, nil
}

func (c KubeClient) GetOCIRepository(ctx context.Context, name string, namespace string) (*sourcev1b2.OCIRepository, *model.ServerResponse) {
	repoKey := k8s.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	var repo sourcev1b2.OCIRepository
	err := c.Get(ctx, repoKey, &repo)
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get fluxcd OCI Repository '%s' in namespace '%s', details: '%s'", name, namespace, err.Error())
	}

	return &repo, nil
}
