package controllers

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

type Group struct {
	*gin.RouterGroup
}

func (g *Group) RegisterControllers() {

}

