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
	"github.com/okdp/okdp-server/internal/common/constants"
	corev1 "k8s.io/api/core/v1"
)

type ContainerStateInfo struct {
	State   string `json:"state"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

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

// GetPodHealth returns a high-level health status string for a given Kubernetes Pod.
// - If the Pod phase is "Running" and all containers are ready, returns "Healthy".
// - If the Pod is "Running" but not all containers are ready, returns "NotReady".
// - For other phases, returns "Pending", "Completed", or "Failed" accordingly.
// - If the phase is unknown, returns "Unknown".
func GetPodHealth(pod *corev1.Pod) string {
	switch pod.Status.Phase {
	case corev1.PodRunning:
		if AreAllContainersReady(pod) {
			return constants.StateHealthy
		}
		return constants.StateNotReady
	case corev1.PodPending:
		return constants.StatePending
	case corev1.PodSucceeded:
		return constants.StateCompleted
	case corev1.PodFailed:
		return constants.StateFailed
	default:
		return constants.StateUnknown
	}
}

// AreAllContainersReady returns true if all containers in the pod have their Ready status set to true.
// Otherwise, it returns false.
func AreAllContainersReady(pod *corev1.Pod) bool {
	for _, cs := range pod.Status.ContainerStatuses {
		if !cs.Ready {
			return false
		}
	}
	return true
}

// GetContainerState returns the high-level state, reason, and message of a container in a given pod.
// - State is one of "Running", "Waiting", "Terminated", or "Unknown" (exported constant recommended).
// - Reason and Message are only set for Waiting or Terminated states, otherwise empty strings.
// If the container is not found, returns "Unknown" state and empty reason/message.
func GetContainerState(pod *corev1.Pod, containerName string) ContainerStateInfo {
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == containerName {
			if status.State.Running != nil {
				return ContainerStateInfo{
					State:   constants.StateRunning,
					Reason:  "",
					Message: "",
				}
			}
			if status.State.Waiting != nil {
				return ContainerStateInfo{
					State:   constants.StateWaiting,
					Reason:  status.State.Waiting.Reason,
					Message: status.State.Waiting.Message,
				}
			}
			if status.State.Terminated != nil {
				return ContainerStateInfo{
					State:   constants.StateTerminated,
					Reason:  status.State.Terminated.Reason,
					Message: status.State.Terminated.Message,
				}
			}
			return ContainerStateInfo{State: constants.StateUnknown}
		}
	}
	return ContainerStateInfo{State: constants.StateUnknown}
}
