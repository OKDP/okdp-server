package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_api "github.com/okdp/okdp-server/api/openapi/v3/_api/k8s"
	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/services"
	"github.com/okdp/okdp-server/internal/utils"
)

type IKuboCDController struct {
	k8sService *services.KuboCDService
}

func KuboCDController() *IKuboCDController {
	return &IKuboCDController{
		k8sService: services.NewKuboCDService(),
	}
}

func (r IKuboCDController) ListK8sReleases(c *gin.Context, clusterID string, namespace string) {
	releasesInfo, err := r.k8sService.ListReleases(clusterID, namespace)
	if err != nil {
		log.Error("Unable to get releases from Kubernetes cluster '%s' on namespace '%s', details: %+v", clusterID, namespace, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, releasesInfo)
}

func (r IKuboCDController) GetK8sRelease(c *gin.Context, clusterID string, namespace string, releaseName string) {
	release, err := r.k8sService.GetRelease(clusterID, namespace, releaseName)
	if err != nil {
		log.Error("Unable to get release from Kubernetes cluster  '%s' on namespace '%s', details: %+v", clusterID, namespace, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, release)
}

func (r IKuboCDController) GetK8sReleaseStatus(c *gin.Context, clusterID string, namespace string, releaseName string) {
	release, err := r.k8sService.GetReleaseStatus(clusterID, namespace, releaseName)
	if err != nil {
		log.Error("Unable to get release status from Kubernetes cluster  '%s' on namespace '%s', details: %+v", clusterID, namespace, err)
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, release)
}

func (r IKuboCDController) CreateK8sRelease(c *gin.Context, clusterID string, namespace string, params _api.CreateK8sReleaseParams) {
	var release model.Release

	if err := c.ShouldBindJSON(&release); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.JSON(resp.Status, resp)
		return
	}

	response := r.k8sService.CreateRelease(clusterID, namespace, &release, utils.OrFalse(params.DryRun))

	c.JSON(response.Status, response)

}

func (r IKuboCDController) UpdateK8sRelease(c *gin.Context, clusterID string, namespace string, params _api.UpdateK8sReleaseParams) {
	var release model.Release

	if err := c.ShouldBindJSON(&release); err != nil {
		resp := model.NewServerResponse(model.OkdpServerResponse).BadRequest("%+v", err.Error())
		c.JSON(resp.Status, resp)
		return
	}

	response := r.k8sService.UpdateRelease(clusterID, namespace, &release, utils.OrFalse(params.DryRun))
	c.JSON(response.Status, response)
}

func (r IKuboCDController) DeleteK8sRelease(c *gin.Context, clusterID string, namespace string, releaseName string) {
	response := r.k8sService.DeleteRelease(clusterID, namespace, releaseName)
	c.JSON(response.Status, response)
}
