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
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	kubocdv1alpha1 "kubocd/api/v1alpha1"

	"github.com/okdp/okdp-server/internal/common/constants"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

func (c KubeClient) ListReleases(ctx context.Context, namespaces ...string) ([]*model.Release, *model.ServerResponse) {

	var releaseList kubocdv1alpha1.ReleaseList
	if err := c.List(ctx, &releaseList); err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to list KuboCD Releases '%s'", err.Error())
	}

	converted := model.ReleaseList(releaseList)

	filtered := utils.Filter2(converted.ToReleases(), func(k model.Release) bool {
		return len(namespaces) == 0 || utils.Contains(namespaces, k.Namespace)
	})

	return filtered, nil
}

func (c KubeClient) GetRelease(ctx context.Context, namespace string, releaseName string) (*model.Release, *model.ServerResponse) {
	releaseKey := ctrlclient.ObjectKey{
		Namespace: namespace,
		Name:      releaseName,
	}

	var release kubocdv1alpha1.Release
	if err := c.Get(ctx, releaseKey, &release); err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get KuboCD Release '%s'", err.Error())
	}

	converted := model.Release(release)
	converted.SanitizeMetadata()

	return &converted, nil
}

func (c KubeClient) GetReleaseStatus(ctx context.Context, namespace string, releaseName string) (*model.ReleaseStatus, *model.ServerResponse) {
	releaseKey := ctrlclient.ObjectKey{
		Namespace: namespace,
		Name:      releaseName,
	}

	var release kubocdv1alpha1.Release
	if err := c.Get(ctx, releaseKey, &release); err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get KuboCD Release Status '%s'", err.Error())
	}

	converted := model.ReleaseStatus(release.Status)

	return &converted, nil
}

func (c KubeClient) CreateRelease(ctx context.Context, namespace string, release *model.Release, dryRun bool) *model.ServerResponse {
	rel := kubocdv1alpha1.Release(*release)
	rel.Namespace = namespace
	var err error

	if dryRun {
		err = c.Create(ctx, &rel, &ctrlclient.CreateOptions{DryRun: []string{constants.All}})
	} else {
		err = c.Create(ctx, &rel)
	}
	if err != nil {
		return model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to create KuboCD Release '%s/%s (%s)', details: '%s'", release.Namespace, release.Name, dryRun, err.Error())
	}

	return model.NewServerResponse(model.K8sClusterResponse).Created("Successfuly created release %s/%s", release.Namespace, release.Name)
}

func (c KubeClient) UpdateRelease(ctx context.Context, namespace string, release *model.Release, dryRun bool) *model.ServerResponse {
	rel := kubocdv1alpha1.Release(*release)
	rel.Namespace = namespace
	var err error

	if dryRun {
		err = c.Update(ctx, &rel, &ctrlclient.UpdateOptions{DryRun: []string{constants.All}})
	} else {
		err = c.Update(ctx, &rel)
	}
	if err != nil {
		return model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to update KuboCD Release '%s/%s (%s)', details: '%s'", release.Namespace, release.Name, dryRun, err.Error())
	}

	return model.NewServerResponse(model.K8sClusterResponse).Updated("Successfuly updated release %s/%s", release.Namespace, release.Name)
}

func (c KubeClient) DeleteRelease(ctx context.Context, namespace string, releaseName string) *model.ServerResponse {
	rel := kubocdv1alpha1.Release{
		ObjectMeta: metav1.ObjectMeta{
			Name:      releaseName,
			Namespace: namespace,
		},
	}

	err := c.Delete(ctx, &rel)
	if err != nil {
		return model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to delete KuboCD Release '%s/%s', details: '%s'", namespace, releaseName, err.Error())
	}

	return model.NewServerResponse(model.K8sClusterResponse).Deleted("Successfuly deleted release %s", releaseName)

}

func (c KubeClient) ListEventsRelease(ctx context.Context, namespace, releaseName string) ([]map[string]interface{}, *model.ServerResponse) {
	fieldSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.name", releaseName),
		fields.OneTermEqualSelector("involvedObject.kind", "Release"),
	).String()

	eventList, err := c.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to list events for release '%s/%s': %s", namespace, releaseName, err.Error())
	}

	events := make([]map[string]interface{}, 0, len(eventList.Items))
	for _, event := range eventList.Items {
		b, err := json.Marshal(event)
		if err != nil {
			continue
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		events = append(events, m)
	}
	return events, nil
}
