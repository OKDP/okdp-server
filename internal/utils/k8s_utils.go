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

package utils

import (
	corev1 "k8s.io/api/core/v1"
)

// MergeLabels merges new labels into an existing Namespace's labels.
// It adds new keys and updates the values of existing keys.
func MergeLabels(ns *corev1.Namespace, newLabels *map[string]string) {
	if newLabels == nil {
		return
	}
	if ns.Labels == nil {
		ns.Labels = map[string]string{}
	}
	for k, v := range *newLabels {
		ns.Labels[k] = v // Overwrite or add
	}
}

// MergeAnnotations merges new annotations into an existing Namespace's annotations.
// It adds new keys and updates the values of existing keys.
func MergeAnnotations(ns *corev1.Namespace, newAnnotations *map[string]string) {
	if newAnnotations == nil {
		return
	}
	if ns.Annotations == nil {
		ns.Annotations = map[string]string{}
	}
	for k, v := range *newAnnotations {
		ns.Annotations[k] = v // Overwrite or add
	}
}
