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

	"github.com/okdp/okdp-server/internal/common/constants"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

func (c KubeClient) ListGitRepos(namespaces ...string) ([]*model.GitRepository, *model.ServerResponse) {

	ctx := context.Background()

	var kustomizations []*kustomizev1.Kustomization
	kustomizations, err := c.ListKutomizations(ctx, namespaces...)
	if err != nil {
		return nil, err
	}

	kustomizations = utils.Filter2(kustomizations, func(k kustomizev1.Kustomization) bool {
		return k.Spec.SourceRef.Kind == constants.GitRepository
	})

	gitRepos := make([]*model.GitRepository, 0, len(kustomizations))

	for _, k := range kustomizations {
		namespace := utils.DefaultIfEmpty(k.Spec.SourceRef.Namespace, k.Namespace)
		fluxRepo, err := c.GetGitRepository(ctx, k.Spec.SourceRef.Name, namespace)
		if err != nil {
			return nil, err
		}
		gitRepo := &model.GitRepository{
			RepoURL:   fluxRepo.Spec.URL,
			Ref:       BranchOrTag(fluxRepo.Spec.Reference.Branch, fluxRepo.Spec.Reference.Tag),
			Name:      k.Name,
			Namespace: k.Namespace,
			Path:      k.Spec.Path,
		}
		gitRepo.Credentials.SecretRef = fluxRepo.Spec.SecretRef.Name
		gitRepos = append(gitRepos, gitRepo)
	}

	return gitRepos, nil
}

func BranchOrTag(branch, tag string) string {
	if branch != "" {
		return "refs/heads/" + branch
	}
	return "refs/tags/" + tag
}
