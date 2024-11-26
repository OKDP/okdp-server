package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
	"github.com/okdp/okdp-server/internal/services"
)

type IDeploymentController struct {
	deploymentService *services.DeploymentService
}

func DeploymentController() *IDeploymentController {
	deploymentService, err := services.NewDeploymentService()
	if err != nil {
		return nil
	}
	return &IDeploymentController{
		deploymentService: deploymentService,
	}
}

func (r IDeploymentController) ListDeployments(c *gin.Context, spaceid string, compositionid string) {
	c.JSON(http.StatusOK, r.deploymentService.List())
}

func (r IDeploymentController) GetDeployments(c *gin.Context, spaceid string, compositionid string, deploymentid string) {
	deployment := r.deploymentService.Get(deploymentid)
	if deployment == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Deployment with id %s not found", deploymentid))
		return
	}
	c.JSON(http.StatusOK, deployment)

}

func (r IDeploymentController) CreateDeployment(c *gin.Context, spaceid string, compositionid string) {
	var deployment api.Deployment
	err := c.BindJSON(&deployment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.deploymentService.Create(deployment)
	c.JSON(http.StatusCreated, deployment)
}

func (r IDeploymentController) UpdateDeployment(c *gin.Context, spaceid string, compositionid string, deploymentid string) {
	var deployment api.Deployment
	err := c.BindJSON(&deployment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	updatedDeployment := r.deploymentService.Update(deploymentid, deployment)
	if updatedDeployment == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Deployment with id %s not found", deploymentid))
		return
	}
	c.JSON(http.StatusOK, deployment)
}

func (r IDeploymentController) DeleteDeployment(c *gin.Context, spaceid string, compositionid string, deploymentid string) {
	deployment := r.deploymentService.Delete(deploymentid)
	if deployment == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Deployment with id %s not found", deploymentid))
		return
	}
	c.JSON(http.StatusOK, deployment)
}
