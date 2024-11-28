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

package kad

import (
	"crypto/tls"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/okdp/okdp-server/internal/config"
	log "github.com/okdp/okdp-server/internal/logging"
)

var (
	instance *KadClient
	once     sync.Once
)

type KadClient struct {
	*resty.Client
}

func GetClient() *KadClient {
	once.Do(func() {
		kadConf := config.GetAppConfig().Kad
		log.Info("KAD configuration: %s", kadConf)
		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: kadConf.InsecureSkipVerify}).
			SetAuthToken(kadConf.AuthBearer).
			SetHeader("Content-Type", "application/json").
			SetBaseURL(kadConf.ApiUrl)
		instance = &KadClient{client}
	})
	return instance
}

func (client *KadClient) Get(id string) (string, error) {
	resp, err := client.R().Get("/mycluster/component")

	return resp.String(), err
}

func (client *KadClient) List() (string, error) {
	resp, err := client.R().Get("/mycluster/component")

	return resp.String(), err
}

