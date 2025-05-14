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
	"path/filepath"
	"strings"
)

// IsYaml checks if the given filename has a .yaml or .yml extension (case-insensitive).
func IsYaml(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".yaml" || ext == ".yml"
}

// PathOrFallback returns the given `path` if it contains a slash ("/"),
// indicating it is a likely full or relative path.
// Otherwise, it returns the provided `fallback` value.
//
// This is useful in cases where you want to prefer a user-supplied
// file or directory path, but fall back to a default if not provided.
func PathOrFallback(path string, fallback string) string {
	if strings.Contains(path, "/") {
		return path
	}
	return fallback

}
