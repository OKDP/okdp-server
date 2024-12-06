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
	"fmt"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/okdp/okdp-server/internal/config"
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

type KadError struct {
	Error error
	StatusCode int
}

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

func (c *KadClients) ID(id string) *KadClient {
	return c.clients[id]
}

func ListInstances() []config.KadInstance {
	return config.GetAppConfig().Kad
}

func GetInstanceById(id string) (config.KadInstance, *KadError) {
	instances := ListInstances()
	for _, i := range instances {
		if i.Id == id {
			return i, nil
		}
	}
	return config.KadInstance{}, KadInstanceNotFound(id)
}

func DoGet[T any](request *resty.Request) (*T, error) {
	request.Method = resty.MethodGet
	return doExecute[T](request)
}

func (c *KadClient) NewRequest(url string) *resty.Request {
	req := c.R()
	req.URL = url
	return req
}

func doExecute[T any](request *resty.Request) (*T, error) {
	var object T
	resp, err := request.Send()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp.String()), &object)

	return &object, err
}

func InvalidInstanceError(provided string) error {
	instances := utils.Map(ListInstances(), func(k config.KadInstance) string { return k.Id })
	return fmt.Errorf("invalid kad instance provided: %s, valid ones: %+v", provided, instances)
}


func KadInstanceNotFound(kadInstance string) *KadError {
	return &KadError {
		Error: fmt.Errorf("kad instance (%s) not found", kadInstance),
		StatusCode: http.StatusNotFound,
	 }
}

