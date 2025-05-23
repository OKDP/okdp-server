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

	"github.com/okdp/okdp-server/internal/model"
)

func (r K8S) ListReleases(clusterID string, namespaces ...string) ([]*model.Release, *model.ServerResponse) {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return kubeClient.ListReleases(context.Background(), namespaces...)
}

func (r K8S) GetRelease(clusterID string, namespace string, releaseName string) (*model.Release, *model.ServerResponse) {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return kubeClient.GetRelease(context.Background(), namespace, releaseName)
}

func (r K8S) GetReleaseStatus(clusterID string, namespace string, releaseName string) (*model.ReleaseStatus, *model.ServerResponse) {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	return kubeClient.GetReleaseStatus(context.Background(), namespace, releaseName)
}

func (r K8S) CreateRelease(clusterID string, namespace string, release *model.Release, dryRun bool) *model.ServerResponse {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return err
	}
	return kubeClient.CreateRelease(context.Background(), namespace, release, dryRun)
}

func (r K8S) UpdateRelease(clusterID string, namespace string, release *model.Release, dryRun bool) *model.ServerResponse {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return err
	}
	return kubeClient.UpdateRelease(context.Background(), namespace, release, dryRun)
}

func (r K8S) DeleteRelease(clusterID string, namespace string, releaseName string) *model.ServerResponse {
	kubeClient, err := r.GetClient(clusterID)
	if err != nil {
		return err
	}
	return kubeClient.DeleteRelease(context.Background(), namespace, releaseName)
}
