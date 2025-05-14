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

type KuboCDService struct {
	kubocd *k8s.K8S
}

func NewKuboCDService() *KuboCDService {
	return &KuboCDService{
		kubocd: k8s.NewK8S(),
	}
}

func (s KuboCDService) ListReleases(clusterID string, namespace string) ([]*model.Release, *model.ServerResponse) {
	return s.kubocd.ListReleases(clusterID, namespace)

}

func (s KuboCDService) GetRelease(clusterID string, namespace string, releaseName string) (*model.Release, *model.ServerResponse) {
	return s.kubocd.GetRelease(clusterID, namespace, releaseName)
}

func (s KuboCDService) GetReleaseStatus(clusterID string, namespace string, releaseName string) (*model.ReleaseStatus, *model.ServerResponse) {
	return s.kubocd.GetReleaseStatus(clusterID, namespace, releaseName)
}

func (s KuboCDService) CreateRelease(clusterID string, namespace string, release *model.Release, dryRun bool) *model.ServerResponse {
	return s.kubocd.CreateRelease(clusterID, namespace, release, dryRun)
}

func (s KuboCDService) UpdateRelease(clusterID string, namespace string, release *model.Release, dryRun bool) *model.ServerResponse {
	return s.kubocd.UpdateRelease(clusterID, namespace, release, dryRun)
}

func (s KuboCDService) DeleteRelease(clusterID string, namespace string, releaseName string) *model.ServerResponse {
	return s.kubocd.DeleteRelease(clusterID, namespace, releaseName)
}
