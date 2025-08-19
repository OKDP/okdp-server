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

func OrFalse(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// DefaultIfEmpty returns `value` if it is not an empty string,
// otherwise it returns `defaultValue`.
// Useful for setting defaults in configurations.
//
// Example:
//
//	DefaultIfEmpty("", "default") // returns "default"
//	DefaultIfEmpty("foo", "default") // returns "foo"
func DefaultIfEmpty(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// EmptyToNil returns a pointer to the given string, unless the string is empty.
// If the input string is "", it returns nil. Otherwise, it returns a pointer to the string.
// This is useful when you want to omit empty fields in JSON serialization with the 'omitempty' tag.
func EmptyToNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NilToEmptySlice returns the input slice if not nil, or an empty slice if input is nil.
// Useful to ensure you never return a nil slice (for API responses, etc).
func NilToEmptySlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
