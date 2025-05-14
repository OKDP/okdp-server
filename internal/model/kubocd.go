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
	"sigs.k8s.io/yaml"

	kubocdv1alpha1 "kubocd/api/v1alpha1"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type Release kubocdv1alpha1.Release
type ReleaseList kubocdv1alpha1.ReleaseList
type ReleaseStatus kubocdv1alpha1.ReleaseStatus

type ReleaseInfo _api.ReleaseInfo

func (r *ReleaseList) ToReleases() []*Release {
	converted := make([]*Release, len(r.Items))
	for i, release := range r.Items {
		c := Release(release)
		c.SanitizeMetadata()
		converted[i] = &c
	}
	return converted
}

func (r *Release) SanitizeMetadata() *Release {
	r.ObjectMeta.ManagedFields = nil
	r.ObjectMeta.Finalizers = nil
	r.ObjectMeta.UID = ""
	r.ObjectMeta.Generation = 0
	// r.ObjectMeta.CreationTimestamp = metav1.Time{}
	return r
}

func (r *Release) SanitizeStatus() *Release {
	r.ObjectMeta.ResourceVersion = ""
	// r.Status = kubocdv1alpha1.ReleaseStatus{}
	return r
}

func (r *Release) ToYAML() (string, error) {
	data, err := yaml.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func KuboCDReleaseNotFoundError(clusterID string, namespace string, releaseName string) *ServerResponse {
	return NewServerResponse(OkdpServerResponse).
		NotFoundError("Unable to find KuboCD release '%s' in the kubernetes cluster '%s' on the namespace '%s'.", releaseName, clusterID, namespace)
}

func KuboCDReleaseCreated(clusterID string, namespace string, releaseName string) *ServerResponse {
	return NewServerResponse(OkdpServerResponse).
		NotFoundError("Unable to find KuboCD release '%s' in the kubernetes cluster '%s' on the namespace '%s'.", releaseName, clusterID, namespace)
}
