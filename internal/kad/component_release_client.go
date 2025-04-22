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
	"bytes"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-yaml"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/kad/client"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/servererrors"
)

type ComponentReleaseClient struct {
	KAD *client.KadClients
}

func NewComponentReleaseClient() *ComponentReleaseClient {
	return &ComponentReleaseClient{
		KAD: client.GetClients(),
	}
}

func (c ComponentReleaseClient) Get(kadInstanceID string, name string, catalog *string) (*model.ComponentReleaseResponse, *servererrors.ServerError) {
	kadClient, err := c.KAD.ID(kadInstanceID)
	if err != nil {
		return nil, err
	}
	req := kadClient.NewRequest(constants.ComponentReleaseURL + "/" + name)
	if catalog != nil {
		req = req.SetQueryParam("catalog", *catalog)
	}
	return client.DoGet[model.ComponentReleaseResponse](req)
}

func (c ComponentReleaseClient) List(kadInstanceID string, catalog *string) (*model.ComponentReleasesResponse, *servererrors.ServerError) {
	kadClient, err := c.KAD.ID(kadInstanceID)
	if err != nil {
		return nil, err
	}
	req := kadClient.NewRequest(constants.ComponentReleaseURL)
	if catalog != nil {
		req = req.SetQueryParam("catalog", *catalog)
	}
	return client.DoGet[model.ComponentReleasesResponse](req)
}

func (c ComponentReleaseClient) UploadAsYaml(kadInstanceID string, name string,
	componentReleaseRequest model.ComponentReleaseRequest,
	commitData map[string]string) (*model.GitCommit, *servererrors.ServerError) {
	kadClient, err := c.KAD.ID(kadInstanceID)
	if err != nil {
		return nil, err
	}

	// Convert the ComponentReleases into YAML
	result := map[string]interface{}{
		"componentReleases": componentReleaseRequest.ComponentReleases,
	}
	componentReleasesYAML, err2 := yaml.Marshal(result)
	if err2 != nil {
		log.Error("Error marshaling ComponentReleases to YAML: %v", err)
		return nil, servererrors.OfType(servererrors.OkdpServer).
			GenericError(http.StatusBadRequest, "Unable to parse json to yaml: %+v", err2)
	}

	log.Info("Uploading Component release (Name: %s, Git Path: %s): \n%s, \n Author Info: %s", name, componentReleaseRequest.GitRepoFolder,
		string(componentReleasesYAML), commitData)

	req := kadClient.NewRequest(constants.GitURL + "/" + componentReleaseRequest.GitRepoFolder + "/" + name + ".yaml").
		SetMultipartFields(
			&resty.MultipartField{
				Param:       "kadfile",
				FileName:    name + ".yaml",
				ContentType: "application/x-yaml",
				Reader:      bytes.NewReader(componentReleasesYAML),
			},
		).
		SetFormData(commitData)

	return client.DoPut[model.GitCommit](req)
}
