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

import "testing"

func TestIsYaml(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"Valid .yaml file", "config.yaml", true},
		{"Valid .yml file", "config.yml", true},
		{"Uppercase .YAML", "CONFIG.YAML", true},
		{"Uppercase .YML", "CONFIG.YML", true},
		{"Not a YAML file (.json)", "data.json", false},
		{"No extension", "README", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsYaml(tt.filename); result != tt.expected {
				t.Errorf("IsYaml(%q) = %v; want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestPathOrFallback(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		fallback string
		want     string
	}{
		{
			name:     "Path contains slash",
			path:     "./config/file.yaml",
			fallback: "default.yaml",
			want:     "./config/file.yaml",
		},
		{
			name:     "Path does not contain slash",
			path:     "file.yaml",
			fallback: "default.yaml",
			want:     "default.yaml",
		},
		{
			name:     "Empty path returns fallback",
			path:     "",
			fallback: "default.yaml",
			want:     "default.yaml",
		},
		{
			name:     "Path is slash only",
			path:     "/",
			fallback: "default.yaml",
			want:     "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PathOrFallback(tt.path, tt.fallback)
			if got != tt.want {
				t.Errorf("PathOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}
