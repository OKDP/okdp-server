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
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// SortVersions sorts a list of version strings in descending order.
//
// - Versions are expected to follow semantic versioning (SemVer).
// - Versions that begin with a 'v' (e.g. "v1.2.3") are supported and normalized for sorting.
// - Invalid SemVer strings are placed at the end of the result list in the original order.
//
// Example:
//
//	input:  []string{"v2.0.0", "v1.0.0", "invalid"}
//	output: []string{"v2.0.0", "v1.0.0", "invalid"}
func SortVersions(versions []string) []string {
	type versionWithOriginal struct {
		Original string
		Version  *semver.Version
	}

	validVersions := make([]versionWithOriginal, 0, len(versions))
	var invalidVersions []string

	for _, v := range versions {
		parsed := strings.TrimPrefix(v, "v")
		sv, err := semver.NewVersion(parsed)
		if err != nil {
			invalidVersions = append(invalidVersions, v)
			continue
		}
		validVersions = append(validVersions, versionWithOriginal{
			Original: v,
			Version:  sv,
		})
	}

	sort.Slice(validVersions, func(i, j int) bool {
		return validVersions[i].Version.GreaterThan(validVersions[j].Version)
	})

	sorted := make([]string, 0, len(versions))
	for _, v := range validVersions {
		sorted = append(sorted, v.Original)
	}

	return append(sorted, invalidVersions...)
}
