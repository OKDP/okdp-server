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

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_component "github.com/okdp/okdp-server/api/openapi/v3/_api/componentreleases"
	"github.com/okdp/okdp-server/internal/constants"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/model"
	auth "github.com/okdp/okdp-server/internal/security/authc/model"
	"github.com/okdp/okdp-server/internal/services"
)

type IComponentReleaseController struct {
	componentReleaseService *services.ComponentReleaseService
}

func ComponentReleaseController() *IComponentReleaseController {
	return &IComponentReleaseController{
		componentReleaseService: services.NewComponentReleaseService(),
	}
}

func (r IComponentReleaseController) ListComponentReleases(c *gin.Context, kadInstanceID string, params _component.ListComponentReleasesParams) {
	components, err := r.componentReleaseService.List(kadInstanceID, params.Catalog)
	if err != nil {
		log.Error("Unable to list Component Releases on kad instance %s, details: %+v", kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, components)
}

func (r IComponentReleaseController) GetComponentRelease(c *gin.Context, kadInstanceID string, name string, params _component.GetComponentReleaseParams) {
	component, err := r.componentReleaseService.Get(kadInstanceID, name, params.Catalog)
	if err != nil {
		log.Error("Unable to get Component Release '%s' on kad instance %s, details: %+v", name, kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}
	if component == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Component with id %s not found", name))
		return
	}
	c.JSON(http.StatusOK, component)

}

func (r IComponentReleaseController) CreateOrUpdateComponentRelease(c *gin.Context, kadInstanceID string, name string) {

	var componentReleaseRequest model.ComponentReleaseRequest
	if err := c.ShouldBindJSON(&componentReleaseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maybeUserInfo, _ := c.Get(constants.OAuth2UserInfo)
	userInfo, _ := maybeUserInfo.(*auth.UserInfo)
	commitData := map[string]string{
		"commit-message":  componentReleaseRequest.Comment,
		"committer-name":  userInfo.Name,
		"committer-email": userInfo.Email,
	}

	gitCommitReponse, err := r.componentReleaseService.CreateOrUpdateComponentRelease(kadInstanceID, name, componentReleaseRequest, commitData)
	if err != nil {
		log.Error("Unable to create or update component release '%s' in Git folder '%s' on kad instance %s, details: %+v", name, componentReleaseRequest.GitRepoFolder, kadInstanceID, err)
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gitCommitReponse)

}
