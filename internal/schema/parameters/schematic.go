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

package schema

type KuboSchema struct {
	ParametersSchema *KuboSchemaItem `json:"parametersSchema"`
	ContextSchema    *KuboSchemaItem `json:"contextSchema"`
}

type KuboSchemaItem struct {
	Description string                     `json:"description,omitempty"`
	Type        string                     `json:"type,omitempty"`
	Properties  map[string]*KuboSchemaItem `json:"properties,omitempty"`
	Items       *KuboSchemaItem            `json:"items,omitempty"`
	Required    bool                       `json:"required,omitempty"`
	Default     interface{}                `json:"default,omitempty"`
}
