/*
 *    Copyright 2025 okdp.io
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

package k8s

import (
	"strings"

	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/model"
)

func (r K8S) ListClusters() []*model.Cluster {
	clusters := config.GetAppConfig().Clusters
	if clusters == nil {
		return []*model.Cluster{}
	}
	return clusters
}

func (r K8S) GetCluster(clusterID string) (*model.Cluster, *model.ServerResponse) {
	clusters := config.GetAppConfig().Clusters
	for _, cluster := range clusters {
		if strings.EqualFold(cluster.ID, clusterID) {
			return cluster, nil
		}
	}
	return nil, model.ClusterNotFoundError(clusterID)
}
