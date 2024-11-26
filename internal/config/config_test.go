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

package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func Test_LoadConfig_AuthBasic(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	config := GetAppConfig()
	// Then
	assert.Equal(t, []BasicAuth{{Login: "dev1",
		Password:  "passW!",
		FirstName: "dev1",
		LastName:  "dev",
		Email:     "dev1.dev@example.org",
		Roles:     []string{"developers", "team1"},
	},
	}, config.Security.AuthN.Basic, "basic users")
}
