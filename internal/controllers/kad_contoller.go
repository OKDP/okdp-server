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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/kad/client"
)

type IKadController struct {
	clients *client.KadClients
}

func KadController() *IKadController {
	clients := client.GetClients()
	return &IKadController{
		clients: clients,
	}
}

func (r IKadController) GetKadInstance(c *gin.Context, kadInstanceID string) {
	instance, err := client.GetInstanceByID(kadInstanceID)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, instance)
	}
}

func (r IKadController) ListKadInstances(c *gin.Context) {
	instances := client.ListInstances()
	c.JSON(http.StatusOK, instances)
}
