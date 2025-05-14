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

package git

import (
	"path/filepath"

	"github.com/okdp/okdp-server/internal/integrations/git/client"
	k8sclient "github.com/okdp/okdp-server/internal/integrations/k8s/client"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

type Repository struct {
	*k8sclient.KubeClients
}

func NewRepository() *Repository {
	return &Repository{
		k8sclient.GetClients(),
	}
}

func (r Repository) ListGitRepos(clusterID string, namespace string) ([]*model.GitRepository, *model.ServerResponse) {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return kubeClient.ListGitRepos(namespace)
}

func (r Repository) GetGitRepo(clusterID string, namespace string, kustomizationName string) (*model.GitRepository, *model.ServerResponse) {
	gitRepos, err := r.ListGitRepos(clusterID, namespace)
	if err != nil {
		return nil, err
	}

	for _, repo := range gitRepos {
		if repo.Name == kustomizationName {
			return repo, nil
		}
	}
	return nil, model.RepoNotFoundError(clusterID, namespace, kustomizationName)
}

func (r Repository) GetContents(clusterID string, namespace string, kustomizationName string) ([]*model.GitContent, *model.ServerResponse) {
	repo, err := r.GetGitRepo(clusterID, namespace, kustomizationName)
	if err != nil {
		return nil, err
	}
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	auth, err := kubeClient.GetAuthMethod(repo.Credentials.SecretRef, namespace)
	if err != nil {
		return nil, err
	}
	return client.DoReadContent(repo, auth)
}

func (r Repository) Write(clusterID string, namespace string, kustomizationName string, content string, commitOpts *model.GitCommit, path string) *model.ServerResponse {
	repo, err := r.GetGitRepo(clusterID, namespace, kustomizationName)
	if err != nil {
		return err
	}
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return err
	}
	auth, err := kubeClient.GetAuthMethod(repo.Credentials.SecretRef, namespace)
	if err != nil {
		return err
	}
	return client.DoPushContent(repo, auth, content, commitOpts, utils.PathOrFallback(path, filepath.Join(repo.Path, path)))
}

func (r Repository) DeleteFile(clusterID string, namespace string, kustomizationName string, commitOpts *model.GitCommit, path string) *model.ServerResponse {
	repo, err := r.GetGitRepo(clusterID, namespace, kustomizationName)
	if err != nil {
		return err
	}
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return err
	}
	auth, err := kubeClient.GetAuthMethod(repo.Credentials.SecretRef, namespace)
	if err != nil {
		return err
	}
	return client.DoDeleteFile(repo, auth, commitOpts, path)
}
