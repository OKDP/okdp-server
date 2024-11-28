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

package model

import (
	response "github.com/okdp/okdp-server/internal/kad/response"
)

type Component struct {
	Spec   Spec            `json:"spec"`
	Status Status          `json:"status"`
}

type Spec struct {
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	Suspended bool     `json:"suspended"`
	Protected bool     `json:"protected"`
	Catalogs  []string `json:"catalogs"`
}

type Status struct {
	Error    string   `json:"error"`
	File     string   `json:"file"`
	Path     string   `json:"path"`
	Title    string   `json:"title"`
	Releases []string `json:"releases"`
}

func ToComponent(jsonString string) (response.KadResponseWrapper[Component], error) {
	return response.ParseJson[Component](jsonString)
}

func ToComponents(jsonString string) ([]response.KadResponseWrapper[Component], error) {
	return response.ParseJsonArray[Component](jsonString)
}

