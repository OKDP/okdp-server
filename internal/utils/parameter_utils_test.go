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

func TestOrFalse(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name     string
		input    *bool
		expected bool
	}{
		{
			name:     "nil input returns false",
			input:    nil,
			expected: false,
		},
		{
			name:     "true input returns true",
			input:    &trueVal,
			expected: true,
		},
		{
			name:     "false input returns false",
			input:    &falseVal,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OrFalse(tt.input)
			if result != tt.expected {
				t.Errorf("OrFalse(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultIfEmpty(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue string
		expected     string
	}{
		{"Empty value", "", "default", "default"},
		{"Non-empty value", "foo", "default", "foo"},
		{"Both empty", "", "", ""},
		{"Default value ignored", "bar", "ignored", "bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultIfEmpty(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
