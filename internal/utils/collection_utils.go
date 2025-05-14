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

func Map[T, L any](array []T, f func(T) L) []L {
	result := make([]L, 0, len(array))
	for _, e := range array {
		result = append(result, f(e))
	}
	return result
}

func ArrayNullToEmpty[T any](a []T) []T {
	if len(a) > 0 {
		return a
	}
	return []T{}
}

func MapKey(keys ...string) string {
	var result string
	for _, key := range keys {
		result += key
	}
	return result
}

// Filter filters a slice of objects based on a predicate function.
// It returns a new slice containing only the elements that satisfy the predicate.
func Filter[T any](objects []T, predicate func(T) bool) []*T {
	var filtered []*T
	for _, obj := range objects {
		if predicate(obj) {
			filtered = append(filtered, &obj)
		}
	}
	return filtered
}

// Filter filters a slice of objects based on a predicate function.
// It returns a new slice containing only the elements that satisfy the predicate.
func Filter2[T any](objects []*T, predicate func(T) bool) []*T {
	var filtered []*T
	for _, obj := range objects {
		if predicate(*obj) {
			filtered = append(filtered, obj)
		}
	}
	return filtered
}

// Contains checks if a given string is in the namespaces slice.
func Contains(values []string, value string) bool {
	for _, ns := range values {
		if ns == value {
			return true
		}
	}
	return false
}
