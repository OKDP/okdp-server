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
	"io"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/okdp/okdp-server/internal/model"
)

func (r LocalRepo) SafePush(auth transport.AuthMethod) error {

	// err := r.Worktree.Pull(&git.PullOptions{
	// 	RemoteName:   "origin",
	// 	Auth:         auth,
	// })
	// if err != nil && err != git.NoErrAlreadyUpToDate {
	// 	return err
	// }

	err := r.Push(&git.PushOptions{
		Auth: auth,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil

}

func (r LocalRepo) CommitFile(content string, path string, commitOpts *model.GitCommit) error {

	fs := r.Worktree.Filesystem

	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(content)); err != nil {
		return err
	}

	_, err = r.Worktree.Add(path)
	if err != nil {
		return err
	}
	_, err = r.Worktree.Commit(commitOpts.Message, &commitOpts.CommitOptions)

	return err
}

func (r LocalRepo) DeleteFile(path string, commitOpts *model.GitCommit) error {

	fs := r.Worktree.Filesystem

	if err := fs.Remove(path); err != nil {
		return err
	}
	_, err := r.Worktree.Remove(path)
	if err != nil {
		return err
	}

	_, err = r.Worktree.Commit(commitOpts.Message, &commitOpts.CommitOptions)

	return err
}

func ListFiles(fs billy.Filesystem, folder string) ([]string, *model.ServerResponse) {
	var fileNames []string

	var walk func(string) *model.ServerResponse
	walk = func(currentPath string) *model.ServerResponse {
		stat, err := fs.Stat(currentPath)
		if err != nil || !stat.IsDir() {
			return nil // Skip if not a directory or doesn't exist
		}

		dirEntries, err := fs.ReadDir(currentPath)
		if err != nil {
			return model.
				NewServerResponse(model.GitRepoResponse).
				UnprocessableEntity("Failed to read git folder %s: %s", currentPath, err.Error())
		}

		for _, entry := range dirEntries {
			fullPath := filepath.Join(currentPath, entry.Name())
			if entry.IsDir() {
				if walkErr := walk(fullPath); walkErr != nil {
					return walkErr
				}
			} else {
				fileNames = append(fileNames, fullPath)
			}
		}
		return nil
	}

	if err := walk(folder); err != nil {
		return nil, err
	}

	return fileNames, nil
}

func ReadContent(fs billy.Filesystem, filePath string) (*model.GitContent, *model.ServerResponse) {
	file, err := fs.Open(filePath)
	if err != nil {
		return nil, model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Failed to open file path %s", filePath, err.Error())
	}

	content, err := io.ReadAll(file)

	if err != nil {
		return nil, model.
			NewServerResponse(model.GitRepoResponse).
			UnprocessableEntity("Failed to read file %s", filePath, err.Error())
	}

	return &model.GitContent{
		Content: &content,
		Path:    filePath,
	}, nil
}
