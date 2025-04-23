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

package model

import (
	"strings"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type Catalog struct {
	_api.Catalog `mapstructure:",squash"`
}

type Package = _api.Package

func (c Catalog) RepoHost() string {
	parts := strings.SplitN(c.RepoURL, "/", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func (c Catalog) IsAuthenticated() bool {
	return c.Credentials != nil
}
