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

package kad

import (
	"encoding/json"
)

// Struct that holds parsed fields and raw JSON
type KadResponseWrapper[T any] struct {
	Parsed T                 `json:"parsed"`
	Raw    json.RawMessage   `json:"raw,omitempty"`
}

// Generic function that takes JSON input and parses it into a struct of type T
func ParseJson[T any](jsonString string) (KadResponseWrapper[T], error) {
	wrappers, err := ParseJsonArray[T]("[" + jsonString + "]")
	if err != nil {
		return KadResponseWrapper[T]{}, err
	}

	return wrappers[0], nil
}

// Generic function that takes JSON input and parses it into a struct array of type T
func ParseJsonArray[T any](jsonString string) ([]KadResponseWrapper[T], error) {
	var rawArray []json.RawMessage
	var wrappers []KadResponseWrapper[T]

	err := json.Unmarshal([]byte(jsonString), &rawArray)
	if err != nil {
		return nil, err
	}
	for _, raw := range rawArray {
		var parsed T
		err := json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, err
		}
		wrappers = append(wrappers, KadResponseWrapper[T]{Parsed: parsed, Raw: raw})
	}

	return wrappers, nil
}


