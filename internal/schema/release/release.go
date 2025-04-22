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

// Release section
type Release struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

// Metadata section
type Metadata struct {
	Name string `json:"name"`
}

// Spec section
type Spec struct {
	Protected       bool                   `json:"protected"`
	TargetNamespace string                 `json:"targetNamespace"`
	CreateNamespace bool                   `json:"createNamespace"`
	Debug           Debug                  `json:"debug"`
	Parameters      map[string]interface{} `yaml:"parameters"`
	Application     Application            `json:"application"`
}

type Debug struct {
	DumpContext bool `json:"dumpContext"`
}

// Application section
type Application struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Interval   string `json:"interval"`
}
