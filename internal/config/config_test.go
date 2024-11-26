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

func Test_LoadConfig_Server_Logging(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	logging := GetAppConfig().Logging
	// Then
	assert.Equal(t, "debug",   logging.Level, "Level")
	assert.Equal(t, "console", logging.Format, "Format")
}

func Test_LoadConfig_Server_Cors(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	security := GetAppConfig().Security
	// Then
	assert.Equal(t, []string{"*"}, security.Cors.AllowedOrigins, "AllowedOrigins")
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE", "PATCH"}, security.Cors.AllowedMethods, "AllowedMethods")
	assert.Equal(t, []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, security.Cors.AllowedHeaders, "AllowedHeaders")
	assert.Equal(t, []string{"Link"}, security.Cors.ExposedHeaders, "ExposedHeaders")
	assert.True(t, true, security.Cors.AllowCredentials, "AllowCredentials")
	assert.Equal(t, int64(3600), security.Cors.MaxAge, "MaxAge")
}

func Test_LoadConfig_Server_Headers(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	security := GetAppConfig().Security
	// Then
	assert.Equal(t, map[string]string{"x-frame-options": "DENY", "x-content-type-options": "nosniff",}, security.Headers, "Headers")
}

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

func Test_LoadConfig_AuthOpenId(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	config := GetAppConfig()
	// Then
	assert.Equal(t, "confidential-oidc-client", config.Security.AuthN.OpenID.ClientID, "ClientID")
	assert.Equal(t, "secret1", config.Security.AuthN.OpenID.ClientSecret, "ClientSecret")
	assert.Equal(t, "http://keycloak:7080/realms/master", config.Security.AuthN.OpenID.IssuerUri, "IssuerUri")
	assert.Equal(t, "http://localhost:8090/oauth2/callback", config.Security.AuthN.OpenID.RedirectUri, "RedirectUri")
	assert.Equal(t, "secret1!", config.Security.AuthN.OpenID.CookieSecret, "CookieSecret")
	assert.Equal(t, "openid+profile+email+roles", config.Security.AuthN.OpenID.Scope, "Scope")
	assert.Equal(t, "realm_access.roles", config.Security.AuthN.OpenID.RolesAttributePath, "RolesAttributePath")
	assert.Equal(t, "realm_access.groups", config.Security.AuthN.OpenID.GroupsAttributePath, "GroupsAttributePath")
}

func Test_LoadConfig_AuthBearer(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	config := GetAppConfig()
	// Then
	assert.Equal(t, "http://keycloak:7080/realms/master", config.Security.AuthN.Bearer.IssuerUri, "IssuerUri")
	assert.Equal(t, "http://keycloak:7080/realms/master/protocol/openid-connect/certs", config.Security.AuthN.Bearer.JwksURL, "JwksURL")
	assert.Equal(t, "realm_access.roles", config.Security.AuthN.Bearer.RolesAttributePath, "RolesAttributePath")
	assert.Equal(t, "realm_access.groups", config.Security.AuthN.Bearer.GroupsAttributePath, "GroupsAttributePath")
	assert.True(t, config.Security.AuthN.Bearer.SkipIssuerCheck, "SkipIssuerCheck")
	assert.False(t, config.Security.AuthN.Bearer.SkipSignatureCheck, "SkipSignatureCheck")
}

func Test_LoadConfig_AuthZProvider_File(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	config := GetAppConfig()
	// Then
	assert.Equal(t, "testdata/security/authz-model.conf", config.Security.AuthZ.File.ModelPath, "ModelPath")
	assert.Equal(t, "testdata/security/authz-policy.csv", config.Security.AuthZ.File.PolicyPath, "PolicyPath")
}

func Test_LoadConfig_AuthZProvider_Database(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	config := GetAppConfig()
	// Then
	assert.Equal(t, "okdp", config.Security.AuthZ.Database.Name, "Name")
	assert.Equal(t, "localhost", config.Security.AuthZ.Database.Host, "Host")
	assert.Equal(t, 5432, config.Security.AuthZ.Database.Port, "Port")
	assert.Equal(t, "adm", config.Security.AuthZ.Database.Username, "Username")
	assert.Equal(t, "passDB!", config.Security.AuthZ.Database.Password, "Password")
}

func Test_LoadConfig_ConfigFileNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recover panic, unable to parse the configuration file")
		}
	}()
	// Given
	viper.Set("config", "not-found/application.yaml")
	resetAppConfig()
	// When
	GetAppConfig()
	// Then
	t.Errorf("Panic was expected")
}
