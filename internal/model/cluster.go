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

package model

import (
	"github.com/okdp/okdp-server/api/openapi/v3/_api"
	"github.com/okdp/okdp-server/internal/common/constants"
)

type Cluster _api.Cluster

func ClusterNotFoundError(clusterID string) *ServerResponse {
	return NewServerResponse(OkdpServerResponse).
		NotFoundError("The cluster with id %s not found.", clusterID)
}

func (m Cluster) AuthType() string {
	if m.Auth.Kubeconfig != nil {
		return constants.K8SAuthKubeConfig
	}

	if m.Auth.Certificate != nil {
		return constants.K8SAuthCertificate
	}

	if m.Auth.Bearer != nil {
		return constants.K8SAuthBeaer
	}

	if *m.Auth.InCluster {
		return constants.K8SInCluster
	}

	return ""
}
