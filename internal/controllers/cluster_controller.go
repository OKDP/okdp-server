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
)

type IClusterController struct {
	clusterService *services.ClusterService
}

func ClusterController() *IClusterController {
	return &IClusterController{
		clusterService: services.NewClusterService(),
	}
}

func (r IClusterController) ListClusters(c *gin.Context) {
	clusters := r.clusterService.ListClusters()
	c.JSON(http.StatusOK, clusters)
}

func (r IClusterController) GetCluster(c *gin.Context, clusterID string) {
	cluster, err := r.clusterService.GetCluster(clusterID)
	if err != nil {
		log.Error("%+v", clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, cluster)
}

func (r IClusterController) ListNamespaces(c *gin.Context, clusterID string) {
	namespaces, err := r.clusterService.ListNamespaces(clusterID)
	if err != nil {
		log.Error("%+v", clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, namespaces)
}

func (r IClusterController) GetNamespace(c *gin.Context, clusterID string, namespace string) {
	ns, err := r.clusterService.GetNamespaceByName(clusterID, namespace)
	if err != nil {
		log.Error("%+v", clusterID, err)
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, ns)
}

func (r IClusterController) CreateNamespace(c *gin.Context, clusterID string) {
	var namespace model.Namespace

	if err := c.ShouldBindJSON(&namespace); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.AbortWithStatusJSON(resp.Status, resp)
		return
	}

	response := r.clusterService.CreateNamespace(clusterID, &namespace)

	c.JSON(response.Status, response)

}

func (r IClusterController) UpdateNamespace(c *gin.Context, clusterID string) {
	var namespace model.Namespace

	if err := c.ShouldBindJSON(&namespace); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.AbortWithStatusJSON(resp.Status, resp)
		return
	}

	response := r.clusterService.UpdateNamespace(clusterID, &namespace)

	c.JSON(response.Status, response)
}

func (r IClusterController) DeleteNamespace(c *gin.Context, clusterID string, namespace string) {
	response := r.clusterService.DeleteNamespace(clusterID, namespace)
	c.JSON(response.Status, response)
}
