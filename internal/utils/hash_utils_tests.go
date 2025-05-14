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
	"testing"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic Git URL",
			input:    "ssh://git@github.com/kubocd/kubocd-infra-ii",
			expected: "ff7b5b726d45",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  ssh://git@github.com/kubocd/kubocd-infra-ii  ",
			expected: "ff7b5b726d45",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "da39a3ee5e6b",
		},
		{
			name:     "Different strings produce different IDs",
			input:    "ssh://git@github.com/kubocd/another-repo",
			expected: "c79dc88f5600",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromString(tt.input)
			if got != tt.expected {
				t.Errorf("FromString(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
