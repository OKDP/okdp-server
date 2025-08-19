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
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/services"
	"github.com/okdp/okdp-server/internal/utils"
)

type IProjectController struct {
	clusterService *services.ClusterService
}

func ProjectController() *IProjectController {
	return &IProjectController{
		clusterService: services.NewClusterService(),
	}
}

func (r IProjectController) ListProjects(c *gin.Context, clusterID string) {
	namespaces, err := r.clusterService.ListNamespaces(clusterID)
	if err != nil {
		log.Error("%+v", clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}

	projects := utils.Map(namespaces, func(ns *model.Namespace) *model.Project {
		return ns.ToProject()
	})

	c.JSON(http.StatusOK, projects)
}

func (r IProjectController) GetProject(c *gin.Context, clusterID string, projectName string) {
	ns, err := r.clusterService.GetNamespaceByName(clusterID, projectName)
	if err != nil {
		log.Error("%+v", clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, ns.ToProject())
}

func (r IProjectController) CreateProject(c *gin.Context, clusterID string) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.AbortWithStatusJSON(resp.Status, resp)
		return
	}
	response := r.clusterService.CreateNamespace(clusterID, project.ToNamespace())
	c.JSON(response.Status, response)

}

func (r IProjectController) UpdateProject(c *gin.Context, clusterID string) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.AbortWithStatusJSON(resp.Status, resp)
		return
	}
	response := r.clusterService.UpdateNamespace(clusterID, project.ToNamespace())
	c.JSON(response.Status, response)

}

func (r IProjectController) DeleteProject(c *gin.Context, clusterID string, projectName string) {
	response := r.clusterService.DeleteNamespace(clusterID, projectName)
	c.JSON(response.Status, response)
}
