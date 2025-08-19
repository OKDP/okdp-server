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
	"context"
	"io"

	"github.com/okdp/okdp-server/internal/model"
)

func (r K8S) GetPods(clusterID, namespace, releaseName string) ([]*model.PodInfo, *model.ServerResponse) {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return kubeClient.GetPods(context.Background(), namespace, releaseName)
}

func (r K8S) StreamLogs(clusterID, namespace, pod, container string, tailLines *int64, isSSE bool) (io.ReadCloser, *model.ServerResponse) {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return kubeClient.StreamLogs(context.Background(), namespace, pod, container, tailLines, isSSE)
}
