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

package services

import (
	"github.com/okdp/okdp-server/internal/integrations/k8s"
	"github.com/okdp/okdp-server/internal/model"
)

type ClusterService struct {
	cluster *k8s.K8S
}

func NewClusterService() *ClusterService {
	return &ClusterService{
		cluster: k8s.NewK8S(),
	}
}

func (s ClusterService) ListClusters() []*model.Cluster {
	return s.cluster.ListClusters()
}

func (s ClusterService) GetCluster(clusterID string) (*model.Cluster, *model.ServerResponse) {
	return s.cluster.GetCluster(clusterID)
}

func (s ClusterService) ListNamespaces(clusterID string) ([]string, *model.ServerResponse) {
	return s.cluster.ListNamespaces(clusterID)
}

func (s ClusterService) GetNamespaceByName(clusterID string, namespace string) (string, *model.ServerResponse) {
	return "Not Implemented: " + clusterID + "/" + namespace, nil
}
