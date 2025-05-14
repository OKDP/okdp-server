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

	"github.com/stretchr/testify/assert"
)

func TestMapfunction(t *testing.T) {

	// Given
	numbers := []int{1, 2, 3}
	f := func(i int) int { return i * 2 }
	// When
	result := Map(numbers, f)
	// Then
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestFilter(t *testing.T) {
	type MyObject struct {
		Name  string
		Value int
	}
	// Test case 1: Filtering even numbers
	objects := []MyObject{
		{Name: "Object 1", Value: 10},
		{Name: "Object 2", Value: 15},
		{Name: "Object 3", Value: 20},
	}
	isEven := func(obj MyObject) bool {
		return obj.Value%2 == 0
	}

	filtered := Filter(objects, isEven)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 objects, but got %d", len(filtered))
	}
	if filtered[0].Value != 10 || filtered[1].Value != 20 {
		t.Errorf("Expected filtered objects to be [10, 20], but got %v", filtered)
	}

	// Test case 2: Filtering objects with value greater than 15
	isGreaterThan15 := func(obj MyObject) bool {
		return obj.Value > 15
	}

	filtered = Filter(objects, isGreaterThan15)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 object, but got %d", len(filtered))
	}
	if filtered[0].Value != 20 {
		t.Errorf("Expected filtered object to be 20, but got %v", filtered[0])
	}

	// Test case 3: Empty list
	filtered = Filter([]MyObject{}, isEven)
	if len(filtered) != 0 {
		t.Errorf("Expected 0 objects, but got %d", len(filtered))
	}
}

func TestContains(t *testing.T) {
	// Test case 1: Value exists in the slice
	namespaces := []string{"flux-system", "kube-system", "default"}
	value := "flux-system"
	if !Contains(namespaces, "flux-system") {
		t.Errorf("Expected '%s' to be in the namespaces slice", value)
	}

	// Test case 2: Value does not exist in the slice
	value = "dev-namespace"
	if Contains(namespaces, value) {
		t.Errorf("Expected '%s' to NOT be in the namespaces slice", value)
	}

	// Test case 3: Empty slice
	namespaces = []string{}
	value = "release-system"
	if Contains(namespaces, value) {
		t.Errorf("Expected '%s' to NOT be in the namespaces slice", value)
	}

	// Test case 4: Single element slice, value exists
	namespaces = []string{"default"}
	value = "default"
	if !Contains(namespaces, value) {
		t.Errorf("Expected '%s' to be in the namespaces slice", value)
	}

	// Test case 5: Single element slice, value does not exist
	namespaces = []string{"default"}
	value = "dev1"
	if Contains(namespaces, value) {
		t.Errorf("Expected '%s' to NOT be in the namespaces slice", value)
	}
}
