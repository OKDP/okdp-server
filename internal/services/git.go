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

package services

import (
	"github.com/okdp/okdp-server/internal/integrations/git"
	"github.com/okdp/okdp-server/internal/model"
)

type GitRepoService struct {
	git *git.Repository
}

func NewGitRepoService() *GitRepoService {
	return &GitRepoService{
		git: git.NewRepository(),
	}
}

func (s GitRepoService) ListGitRepos(clusterID string, namespace string) ([]*model.GitRepository, *model.ServerResponse) {
	return s.git.ListGitRepos(clusterID, namespace)
}

func (s GitRepoService) GetGitRepo(clusterID string, namespace string, kustomizationName string) (*model.GitRepository, *model.ServerResponse) {
	return s.git.GetGitRepo(clusterID, namespace, kustomizationName)
}

func (s GitRepoService) ListReleases(clusterID string, namespace string, kustomizationName string) ([]*model.ReleaseInfo, *model.ServerResponse) {
	contents, err := s.git.GetContents(clusterID, namespace, kustomizationName)
	if err != nil {
		return nil, err
	}
	return s.toReleaseInfo(contents)
}

func (s GitRepoService) GetRelease(clusterID string, namespace string, kustomizationName string, releaseName string) (*model.Release, *model.ServerResponse) {
	contents, err := s.git.GetContents(clusterID, namespace, kustomizationName)
	if err != nil {
		return nil, err
	}

	for _, content := range contents {
		release, err := content.ToRelease()
		if err != nil {
			return nil, err
		}

		if release.ObjectMeta.Name == releaseName {
			return release, nil
		}
	}

	return nil, model.KuboCDGitReleaseNotFoundError(clusterID, namespace, kustomizationName, releaseName)
}

func (s GitRepoService) CreateGitRelease(clusterID string, namespace string, kustomizationName string, release *model.Release, commitOpts *model.GitCommit) (*model.Release, *model.ServerResponse) {
	releaseInfo, er := s.getReleaseInfo(clusterID, namespace, kustomizationName, release.Namespace, release.Name)
	if er != nil {
		return nil, er
	}
	if releaseInfo == nil {
		content, err := release.SanitizeMetadata().SanitizeStatus().ToYAML()
		if err != nil {
			return nil, model.
				NewServerResponse(model.OkdpServerResponse).
				UnprocessableEntity("Unable to convert KuboCD release %s/%s into yaml", release.Namespace, release.Name)
		}
		return release, s.git.Write(clusterID, namespace, kustomizationName, content, commitOpts, release.Name+".yaml")
	}

	return nil, model.NewServerResponse(model.OkdpServerResponse).ConflictError("Release '%s' already exists in the git repo %s (%s)", release.Name, releaseInfo.Git.Path, releaseInfo.Git.URL)
}

func (s GitRepoService) UpdateGitRelease(clusterID string, namespace string, kustomizationName string, release *model.Release, commitOpts *model.GitCommit) (*model.Release, *model.ServerResponse) {
	releaseInfo, er := s.getReleaseInfo(clusterID, namespace, kustomizationName, release.Namespace, release.Name)
	if er != nil {
		return nil, er
	}
	if releaseInfo == nil {
		return nil, model.NewServerResponse(model.OkdpServerResponse).NotFoundError("Release '%s' does not exist in the git repo", release.Name)
	}

	content, err := release.SanitizeMetadata().SanitizeStatus().ToYAML()
	if err != nil {
		return nil, model.
			NewServerResponse(model.OkdpServerResponse).
			UnprocessableEntity("Unable to convert KuboCD release %s/%s into yaml", release.Namespace, release.Name)
	}

	return release, s.git.Write(clusterID, namespace, kustomizationName, content, commitOpts, releaseInfo.Git.Path)
}

func (s GitRepoService) DeleteGitRelease(clusterID string, namespace string, kustomizationName string, releaseName string, commitOpts *model.GitCommit) *model.ServerResponse {
	releaseInfo, er := s.getReleaseInfo(clusterID, namespace, kustomizationName, "default", releaseName)
	if er != nil {
		return er
	}
	if *releaseInfo == (model.ReleaseInfo{}) {
		return model.NewServerResponse(model.OkdpServerResponse).NotFoundError("Release '%s' does not exist in the git repo", releaseName)
	}

	return s.git.DeleteFile(clusterID, namespace, kustomizationName, commitOpts, releaseInfo.Git.Path)
}

func (s GitRepoService) toReleaseInfo(contents []*model.GitContent) ([]*model.ReleaseInfo, *model.ServerResponse) {
	releasesInfo := []*model.ReleaseInfo{}
	for _, content := range contents {
		release, err := content.ToRelease()
		if err != nil {
			return nil, err
		}

		name := release.ObjectMeta.Name
		namespace := release.ObjectMeta.Namespace
		releaseInfo := &model.ReleaseInfo{
			Name:        name,
			Namespace:   &namespace,
			Description: &release.Spec.Description,
		}
		releaseInfo.Package.Repository = release.Spec.Package.Repository
		releaseInfo.Package.Tag = release.Spec.Package.Tag
		releaseInfo.Git.Path = content.Path
		releaseInfo.Git.URL = content.URL

		releasesInfo = append(releasesInfo, releaseInfo)
	}
	return releasesInfo, nil
}

func (s GitRepoService) getReleaseInfo(clusterID, namespace, kustomizationName, releaseNamespace, releaseName string) (*model.ReleaseInfo, *model.ServerResponse) {
	releases, err := s.ListReleases(clusterID, namespace, kustomizationName)
	if err != nil {
		return nil, err
	}

	for _, rel := range releases {
		if rel.Name == releaseName && rel.Namespace != nil && *rel.Namespace == releaseNamespace {
			return rel, nil
		}
	}

	return nil, nil
}
