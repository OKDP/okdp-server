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
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
	"sigs.k8s.io/yaml"
)

type GitRepository _api.GitRepository
type GitCommit struct {
	Message string
	git.CommitOptions
}
type GitContent struct {
	Content *[]byte
	Path    string
	URL     string
}

func (c GitContent) ToRelease() (*Release, *ServerResponse) {
	var release Release
	if err := yaml.Unmarshal(*c.Content, &release); err != nil {
		return nil, NewServerResponse(GitRepoResponse).
			UnprocessableEntity("Failed to convert yaml file content '%s' into KuboCD Release: %v", c.Path, err)
	}
	return &release, nil
}

func RepoNotFoundError(clusterID string, namespace string, fluxrepo string) *ServerResponse {
	return NewServerResponse(OkdpServerResponse).
		NotFoundError("The git repo %s not found on namespace %s with cluster id %s.", fluxrepo, namespace, clusterID)
}

func KuboCDGitReleaseNotFoundError(clusterID string, namespace string, fluxrepo string, releaseName string) *ServerResponse {
	return NewServerResponse(OkdpServerResponse).
		NotFoundError("Unable to find KuboCD release '%s' in the git repo referenced by fluxcd repo %s/%s (clusterID: %s).", releaseName, namespace, fluxrepo, clusterID)
}

func NewGitCommitOptions(message string) *GitCommit {
	commit := &GitCommit{
		Message: "[okdp-server] " + message,
		CommitOptions: git.CommitOptions{
			Author: &object.Signature{},
		},
	}
	commit.CommitOptions.Author.When = time.Now()
	return commit
}

func (c *GitCommit) Author(name string) *GitCommit {
	c.CommitOptions.Author.Name = name
	return c
}

func (c *GitCommit) Email(email string) *GitCommit {
	c.CommitOptions.Author.Email = email
	return c
}
