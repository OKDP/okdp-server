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

package client

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/errors"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/utils"
)

var (
	instance *KadClients
	once     sync.Once
)

type KadClients struct {
	clients map[string]*KadClient
}

type KadClient struct {
	*resty.Client
}

type Request struct {
	*resty.Request
}

// KAD errors response are plain text
// type KadError struct {
// 	message     string
//  statusCode  int
// }

func GetClients() *KadClients {
	once.Do(func() {
		clients := make(map[string]*KadClient)
		kadsConf := config.GetAppConfig().Kad
		for _, kadConf := range kadsConf {
			log.Info("KAD configuration: %s", kadConf)
			client := resty.New()
			client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: kadConf.InsecureSkipVerify}).
				SetAuthToken(kadConf.AuthBearer).
				SetHeader("Content-Type", "application/json").
				SetBaseURL(kadConf.ApiUrl)
			clients[kadConf.Id] = &KadClient{client}
		}
		instance = &KadClients{clients: clients}
	})
	return instance
}

func (c *KadClients) ID(id string) (*KadClient, *errors.ServerError) {
	client, found := c.clients[id]
	if !found {
		return nil, invalidInstanceError(id)
	}
	return client, nil
}

func ListInstances() []config.KadInstance {
	return config.GetAppConfig().Kad
}

func GetInstanceById(id string) (config.KadInstance, *errors.ServerError) {
	instances := ListInstances()
	for _, i := range instances {
		if i.Id == id {
			return i, nil
		}
	}
	return config.KadInstance{}, errors.OfType(errors.OkdpServer).
		NotFoundError("kad instance with id %s not found", id)
}

func DoGet[T any](request *resty.Request) (*T, *errors.ServerError) {
	request.Method = resty.MethodGet
	return doExecute[T](request)
}

func (c *KadClient) NewRequest(url string) *resty.Request {
	req := c.R()
	req.URL = url
	return req
}

func doExecute[T any](request *resty.Request) (*T, *errors.ServerError) {
	var object T
	//request.SetError(&KadError{})
	resp, err := request.Send()

	if err != nil {
		return nil, errors.OfType(errors.Kad).Forbidden(err)
	}

	if resp.IsError() {
		// KAD errors response are plain text
		return nil, errors.OfType(errors.Kad).
			GenericError(resp.StatusCode(), "Kad rejected the request, reason: %s", resp.String())
	}

	err = json.Unmarshal([]byte(resp.String()), &object)

	if err != nil {
		return &object, errors.OfType(errors.OkdpServer).
			GenericError(http.StatusUnprocessableEntity, "Unable to process kad response, reason: %+v", err)
	}

	return &object, nil
}

func invalidInstanceError(provided string) *errors.ServerError {
	instances := utils.Map(ListInstances(), func(k config.KadInstance) string { return k.Id })
	return errors.OfType(errors.OkdpServer).
		NotFoundError("kad instance with id %s not found, valid ones: %+v", provided, instances)
}
