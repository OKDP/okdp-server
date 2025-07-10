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
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMergeLabels_AddAndOverwrite(t *testing.T) {
	const newValue = "newValue"

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"existing": "oldValue",
			},
		},
	}
	newLabels := map[string]string{
		"existing": newValue,     // should overwrite
		"added":    "addedValue", // should be added
	}

	MergeLabels(ns, &newLabels)

	if val, ok := ns.ObjectMeta.Labels["existing"]; !ok || val != newValue {
		t.Errorf("expected 'existing' label to be %q, got %q", newValue, val)
	}
	if val, ok := ns.ObjectMeta.Labels["added"]; !ok || val != "addedValue" {
		t.Errorf("expected 'added' label to be 'addedValue', got %q", val)
	}
}

func TestMergeLabels_NilNewLabels(t *testing.T) {
	const oldValue = "oldValue"

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"existing": oldValue,
			},
		},
	}

	MergeLabels(ns, nil)

	if val, ok := ns.ObjectMeta.Labels["existing"]; !ok || val != oldValue {
		t.Errorf("expected 'existing' label to remain %q, got %q", oldValue, val)
	}
}

func TestMergeLabels_NilLabelsInNamespace(t *testing.T) {
	const newValue = "newValue"

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Labels: nil,
		},
	}
	newLabels := map[string]string{
		"newKey": newValue,
	}

	MergeLabels(ns, &newLabels)

	if val, ok := ns.ObjectMeta.Labels["newKey"]; !ok || val != newValue {
		t.Errorf("expected 'newKey' label to be %q, got %q", newValue, val)
	}
}

func TestMergeAnnotations_AddAndOverwrite(t *testing.T) {
	const newValue = "newValue"

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"existing": "oldValue",
			},
		},
	}
	newAnnotations := map[string]string{
		"existing": newValue,     // should overwrite
		"added":    "addedValue", // should be added
	}

	MergeAnnotations(ns, &newAnnotations)

	if val, ok := ns.ObjectMeta.Annotations["existing"]; !ok || val != newValue {
		t.Errorf("expected 'existing' annotation to be %q, got %q", newValue, val)
	}
	if val, ok := ns.ObjectMeta.Annotations["added"]; !ok || val != "addedValue" {
		t.Errorf("expected 'added' annotation to be 'addedValue', got %q", val)
	}
}

func TestMergeAnnotations_NilNewAnnotations(t *testing.T) {
	const oldValue = "oldValue"

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"existing": oldValue,
			},
		},
	}

	MergeAnnotations(ns, nil)

	if val, ok := ns.ObjectMeta.Annotations["existing"]; !ok || val != oldValue {
		t.Errorf("expected 'existing' annotation to remain %q, got %q", oldValue, val)
	}
}

func TestMergeAnnotations_NilAnnotationsInNamespace(t *testing.T) {
	const newValue = "newValue"

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: nil,
		},
	}
	newAnnotations := map[string]string{
		"newKey": newValue,
	}

	MergeAnnotations(ns, &newAnnotations)

	if val, ok := ns.ObjectMeta.Annotations["newKey"]; !ok || val != newValue {
		t.Errorf("expected 'newKey' annotation to be %q, got %q", newValue, val)
	}
}
