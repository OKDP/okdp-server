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
	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"

	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

type LocalRepo struct {
	*git.Repository
	*git.Worktree
}

func NewLocalRepo(url string, ref string, auth transport.AuthMethod) (*LocalRepo, *model.ServerResponse) {

	fs := memfs.New()
	repo, er := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
		Auth:          auth,
		Depth:         1,
		Progress:      nil,
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(ref),
	})

	if er != nil {
		return nil, model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity(er.Error())
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Unable to access local git repo filesystem %s, %s, details: %s", url, ref, err.Error())
	}

	return &LocalRepo{
		repo, worktree,
	}, nil
}

func DoReadContent(repo *model.GitRepository, auth transport.AuthMethod) ([]*model.GitContent, *model.ServerResponse) {
	contents := []*model.GitContent{}

	localrepo, err := NewLocalRepo(repo.RepoURL, repo.Ref, auth)
	if err != nil {
		return nil, err
	}

	fs := localrepo.Filesystem

	filePaths, err := ListFiles(fs, repo.Path)
	if err != nil {
		return nil, err
	}
	for _, filePath := range filePaths {
		if !utils.IsYaml(filePath) {
			log.Warn("Ignoring file '%s' in folder '%s': not a YAML file", filePath, repo.Path)
			continue
		}
		content, err := ReadContent(fs, filePath)
		if err != nil {
			return nil, err
		}
		content.URL = repo.RepoURL
		contents = append(contents, content)
	}

	return contents, nil

}

func DoPushContent(repo *model.GitRepository, auth transport.AuthMethod, content string, commitOpts *model.GitCommit, path string) *model.ServerResponse {
	localrepo, er := NewLocalRepo(repo.RepoURL, repo.Ref, auth)
	if er != nil {
		return er
	}

	err := localrepo.CommitFile(content, path, commitOpts)

	if err != nil {
		return model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Unable to commit file %s (%s/%s) => %s (%s), details: %s", path, repo.Namespace, repo.Name, repo.RepoURL, repo.Ref, err.Error())
	}

	err = localrepo.SafePush(auth)

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Unable to push %s (%s/%s) => %s (%s), details: %s", path, repo.Namespace, repo.Name, repo.RepoURL, repo.Ref, err.Error())
	}

	return nil
}

func DoDeleteFile(repo *model.GitRepository, auth transport.AuthMethod, commitOpts *model.GitCommit, path string) *model.ServerResponse {
	localrepo, er := NewLocalRepo(repo.RepoURL, repo.Ref, auth)
	if er != nil {
		return er
	}

	err := localrepo.DeleteFile(path, commitOpts)

	if err != nil {
		return model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Unable to delete file %s (%s/%s) => %s (%s), details: %s", path, repo.Namespace, repo.Name, repo.RepoURL, repo.Ref, err.Error())
	}

	err = localrepo.SafePush(auth)

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Unable to push %s (%s/%s) => %s (%s), details: %s", path, repo.Namespace, repo.Name, repo.RepoURL, repo.Ref, err.Error())
	}

	return model.NewServerResponse(model.OkdpServerResponse).Deleted("Release name %s successfully deleted", path)
}
