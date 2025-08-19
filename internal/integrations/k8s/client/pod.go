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
package client

import (
	"context"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

func (c KubeClient) GetPods(ctx context.Context, namespace string, releaseName string) ([]*model.PodInfo, *model.ServerResponse) {
	release, err := c.GetRelease(ctx, namespace, releaseName)
	if err != nil {
		return nil, err
	}

	podList, er := c.CoreV1().Pods(utils.DefaultIfEmpty(release.Spec.TargetNamespace, namespace)).List(ctx, metav1.ListOptions{})
	if er != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("failed to list pods from cluster '%s' matching release name '%s/%s', details: '%s'", c.clusterID, namespace, releaseName, er.Error())
	}

	var result []*model.PodInfo

	for _, pod := range podList.Items {
		if strings.HasPrefix(pod.Name, releaseName) {
			var containers []struct {
				Image   string  `json:"image"`
				Message *string `json:"message,omitempty"`
				Name    string  `json:"name"`
				Reason  *string `json:"reason,omitempty"`
				State   string  `json:"state"`
			}

			for _, c := range pod.Spec.Containers {
				cs := utils.GetContainerState(&pod, c.Name)
				containers = append(containers, struct {
					Image   string  `json:"image"`
					Message *string `json:"message,omitempty"`
					Name    string  `json:"name"`
					Reason  *string `json:"reason,omitempty"`
					State   string  `json:"state"`
				}{
					Image:   c.Image,
					Name:    c.Name,
					State:   cs.State,
					Reason:  utils.EmptyToNil(cs.Reason),
					Message: utils.EmptyToNil(cs.Message),
				})
			}

			result = append(result, &model.PodInfo{
				Name:       pod.Name,
				Namespace:  pod.Namespace,
				CreatedAt:  pod.CreationTimestamp.Time,
				State:      string(pod.Status.Phase),
				Health:     utils.GetPodHealth(&pod),
				Containers: containers,
			})
		}
	}

	return utils.NilToEmptySlice(result), nil
}

func (c KubeClient) StreamLogs(ctx context.Context, namespace, pod, container string, tailLines *int64, isSSE bool) (io.ReadCloser, *model.ServerResponse) {
	podLogOpts := &corev1.PodLogOptions{
		Container:  container,
		Timestamps: false,
		TailLines:  tailLines,
		Follow:     isSSE,
	}

	req := c.CoreV1().Pods(namespace).GetLogs(pod, podLogOpts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to fetch logs from cluster '%s' for '%s/%s (%s)', details: '%s'",
				c.clusterID, namespace, pod, container, err.Error())
	}

	return stream, nil
}
