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
	"io"

	"github.com/okdp/okdp-server/internal/integrations/k8s"
	"github.com/okdp/okdp-server/internal/model"
)

type PodService struct {
	pod *k8s.K8S
}

func NewPodService() *PodService {
	return &PodService{
		pod: k8s.NewK8S(),
	}
}

func (s PodService) GetPods(clusterID, namespace, releaseName string) ([]*model.PodInfo, *model.ServerResponse) {
	return s.pod.GetPods(clusterID, namespace, releaseName)
}

func (s PodService) StreamLogs(clusterID, namespace, pod, container string, tailLines *int64, isSSE bool) (io.ReadCloser, *model.ServerResponse) {
	return s.pod.StreamLogs(clusterID, namespace, pod, container, tailLines, isSSE)
}
