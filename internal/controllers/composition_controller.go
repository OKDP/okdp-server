package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
	"github.com/okdp/okdp-server/internal/services"
)

type ICompositionController struct {
	compositionService *services.CompositionService
}

func CompositionController() *ICompositionController {
	compositionService, err := services.NewCompositionService()
	if err != nil {
		return nil
	}
	return &ICompositionController{
		compositionService: compositionService,
	}
}

func (r ICompositionController) ListCompositions(c *gin.Context, spaceid string) {
	c.JSON(http.StatusOK, r.compositionService.List())
}

func (r ICompositionController) GetCompositions(c *gin.Context, spaceid string, compositionid string) {
	composition := r.compositionService.Get(compositionid)
	if composition == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Composition with id %s not found", compositionid))
		return
	}
	c.JSON(http.StatusOK, composition)

}

func (r ICompositionController) CreateComposition(c *gin.Context, spaceid string) {
	var composition api.Composition
	err := c.BindJSON(&composition)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.compositionService.Create(composition)
	c.JSON(http.StatusCreated, composition)
}

func (r ICompositionController) UpdateComposition(c *gin.Context, spaceid string, compositionid string) {
	var composition api.Composition
	err := c.BindJSON(&composition)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	updatedComposition := r.compositionService.Update(compositionid, composition)
	if updatedComposition == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Composition with id %s not found", compositionid))
		return
	}
	c.JSON(http.StatusOK, composition)
}

func (r ICompositionController) DeleteComposition(c *gin.Context, spaceid string, compositionid string) {
	composition := r.compositionService.Delete(compositionid)
	if composition == nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Composition with id %s not found", compositionid))
		return
	}
	c.JSON(http.StatusOK, composition)
}
