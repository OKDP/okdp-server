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

package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortVersions(t *testing.T) {
	tests := []struct {
		versions []string
		expected []string
	}{
		{
			versions: []string{"1.0.0", "1.2.0", "2.0.0", "1.1.1", "invalid", "v2.1.0"},
			expected: []string{"v2.1.0", "2.0.0", "1.2.0", "1.1.1", "1.0.0", "invalid"},
		},
		{
			versions: []string{"v2.1.0", "1.0.0", "v1.1.0", "v3.0.0", "1.1.1", "v2.0.0"},
			expected: []string{"v3.0.0", "v2.1.0", "v2.0.0", "1.1.1", "v1.1.0", "1.0.0"},
		},
		{
			versions: []string{"v1.0.0", "invalid", "v2.2.0"},
			expected: []string{"v2.2.0", "v1.0.0", "invalid"},
		},
		{
			versions: []string{"invalid", "v1.0.0", "v2.0.0"},
			expected: []string{"v2.0.0", "v1.0.0", "invalid"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Sorting %v", test.versions), func(t *testing.T) {
			result := SortVersions(test.versions)
			assert.Equal(t, test.expected, result)
		})
	}
}
