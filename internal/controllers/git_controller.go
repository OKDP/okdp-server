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

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/services"
)

type IGitRepoController struct {
	gitRepoService *services.GitRepoService
}

func GitRepoController() *IGitRepoController {
	return &IGitRepoController{
		gitRepoService: services.NewGitRepoService(),
	}
}

func (r IGitRepoController) ListGitRepos(c *gin.Context, clusterID string, namespace string) {
	gitRepos, err := r.gitRepoService.ListGitRepos(clusterID, namespace)
	if err != nil {
		log.Error("Unable to list Git repos on namespace '%s' with cluster ID '%s', details: %+v", namespace, clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, gitRepos)
}

func (r IGitRepoController) GetGitRepo(c *gin.Context, clusterID string, namespace string, kustomizationName string) {
	gitRepo, err := r.gitRepoService.GetGitRepo(clusterID, namespace, kustomizationName)
	if err != nil {
		log.Error("Unable to find Git repo '%s' on namespace '%s' with cluster id '%s', details: %+v", kustomizationName, namespace, clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, gitRepo)
}

func (r IGitRepoController) ListGitReleases(c *gin.Context, clusterID string, namespace string, kustomizationName string) {
	releasesInfo, err := r.gitRepoService.ListReleases(clusterID, namespace, kustomizationName)
	if err != nil {
		log.Error("Unable to get releases from Git repo '%s/%s' on cluster ID '%s', details: %+v", namespace, kustomizationName, clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, releasesInfo)
}

func (r IGitRepoController) GetGitRelease(c *gin.Context, clusterID string, namespace string, kustomizationName string, releaseName string) {
	release, err := r.gitRepoService.GetRelease(clusterID, namespace, kustomizationName, releaseName)
	if err != nil {
		log.Error("Unable to get release from Git repo '%s/%s' on cluster ID '%s', details: %+v", namespace, kustomizationName, clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, release)
}

func (r IGitRepoController) CreateGitRelease(c *gin.Context, clusterID string, namespace string, kustomizationName string) {
	var release model.Release
	userInfo, err := GetUserInfo(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
	}

	if err := c.ShouldBindJSON(&release); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.AbortWithStatusJSON(resp.Status, resp)
		return
	}

	msg := fmt.Sprintf("Create new KuboCD release %s/%s on cluster id %s", namespace, release.Name, clusterID)
	commitOpts := model.NewGitCommitOptions(msg).
		Author(userInfo.Name).
		Email(userInfo.Email)
	resp, err := r.gitRepoService.CreateGitRelease(clusterID, namespace, kustomizationName, &release, commitOpts)

	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, resp)

}

func (r IGitRepoController) UpdateGitRelease(c *gin.Context, clusterID string, namespace string, kustomizationName string) {
	var release model.Release
	userInfo, err := GetUserInfo(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
	}

	if err := c.ShouldBindJSON(&release); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.AbortWithStatusJSON(resp.Status, resp)
		return
	}

	msg := fmt.Sprintf("Update KuboCD release %s/%s on cluster id %s", namespace, release.Name, clusterID)
	commitOpts := model.NewGitCommitOptions(msg).
		Author(userInfo.Name).
		Email(userInfo.Email)
	resp, err := r.gitRepoService.UpdateGitRelease(clusterID, namespace, kustomizationName, &release, commitOpts)

	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r IGitRepoController) DeleteGitRelease(c *gin.Context, clusterID string, namespace string, kustomizationName string, releaseName string) {

	userInfo, err := GetUserInfo(c)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
		return
	}

	msg := fmt.Sprintf("Delete KuboCD release %s/%s on cluster id %s", namespace, releaseName, clusterID)
	commitOpts := model.NewGitCommitOptions(msg).
		Author(userInfo.Name).
		Email(userInfo.Email)

	resp := r.gitRepoService.DeleteGitRelease(clusterID, namespace, kustomizationName, releaseName, commitOpts)

	c.JSON(http.StatusOK, resp)

}
