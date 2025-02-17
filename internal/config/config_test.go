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
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig_Server(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	server := GetAppConfig().Server
	// Then
	assert.Equal(t, "0.0.0.0", server.ListenAddress, "ListenAddress")
	assert.Equal(t, 8090, server.Port, "Port")
	assert.Equal(t, "debug", server.Mode, "Mode")
}

func Test_LoadConfig_Server_Logging(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	logging := GetAppConfig().Logging
	// Then
	assert.Equal(t, "debug", logging.Level, "Level")
	assert.Equal(t, "console", logging.Format, "Format")
}

func Test_LoadConfig_Server_Cors(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	cors := GetAppConfig().Security.Cors
	// Then
	assert.Equal(t, []string{"*"}, cors.AllowedOrigins, "AllowedOrigins")
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE", "PATCH"}, cors.AllowedMethods, "AllowedMethods")
	assert.Equal(t, []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, cors.AllowedHeaders, "AllowedHeaders")
	assert.Equal(t, []string{"Link"}, cors.ExposedHeaders, "ExposedHeaders")
	assert.True(t, true, cors.AllowCredentials, "AllowCredentials")
	assert.Equal(t, int64(3600), cors.MaxAge, "MaxAge")
}

func Test_LoadConfig_Server_Headers(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	security := GetAppConfig().Security
	// Then
	assert.Equal(t, map[string]string{"x-frame-options": "DENY", "x-content-type-options": "nosniff"}, security.Headers, "Headers")
}

func Test_LoadConfig_AuthBasic(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	authn := GetAppConfig().Security.AuthN
	// Then
	assert.Equal(t, []BasicAuth{{Login: "dev1",
		Password:  "passW!",
		FirstName: "dev1",
		LastName:  "dev",
		Email:     "dev1.dev@example.org",
		Roles:     []string{"developers", "team1"},
	},
	}, authn.Basic, "basic users")
}

func Test_LoadConfig_AuthOpenId(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	openID := GetAppConfig().Security.AuthN.OpenID
	// Then
	assert.Equal(t, "confidential-oidc-client", openID.ClientID, "ClientID")
	assert.Equal(t, "secret1", openID.ClientSecret, "ClientSecret")
	assert.Equal(t, "http://keycloak:7080/realms/master", openID.IssuerURI, "IssuerUri")
	assert.Equal(t, "http://localhost:8090/oauth2/callback", openID.RedirectURI, "RedirectUri")
	assert.Equal(t, "secret1!", openID.CookieSecret, "CookieSecret")
	assert.Equal(t, "openid+profile+email+roles", openID.Scope, "Scope")
	assert.Equal(t, "realm_access.roles", openID.RolesAttributePath, "RolesAttributePath")
	assert.Equal(t, "realm_access.groups", openID.GroupsAttributePath, "GroupsAttributePath")
}

func Test_LoadConfig_AuthBearer(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	bearer := GetAppConfig().Security.AuthN.Bearer
	// Then
	assert.Equal(t, "http://keycloak:7080/realms/master", bearer.IssuerURI, "IssuerUri")
	assert.Equal(t, "http://keycloak:7080/realms/master/protocol/openid-connect/certs", bearer.JwksURL, "JwksURL")
	assert.Equal(t, "realm_access.roles", bearer.RolesAttributePath, "RolesAttributePath")
	assert.Equal(t, "realm_access.groups", bearer.GroupsAttributePath, "GroupsAttributePath")
	assert.True(t, bearer.SkipIssuerCheck, "SkipIssuerCheck")
	assert.False(t, bearer.SkipSignatureCheck, "SkipSignatureCheck")
}

func Test_LoadConfig_AuthZProvider_File(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	autz := GetAppConfig().Security.AuthZ
	// Then
	assert.Equal(t, "testdata/security/authz-model.conf", autz.File.ModelPath, "ModelPath")
	assert.Equal(t, "testdata/security/authz-policy.csv", autz.File.PolicyPath, "PolicyPath")
}

func Test_LoadConfig_AuthZProvider_Database(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	autz := GetAppConfig().Security.AuthZ
	// Then
	assert.Equal(t, "okdp", autz.Database.Name, "Name")
	assert.Equal(t, "localhost", autz.Database.Host, "Host")
	assert.Equal(t, 5432, autz.Database.Port, "Port")
	assert.Equal(t, "adm", autz.Database.Username, "Username")
	assert.Equal(t, "passDB!", autz.Database.Password, "Password")
}

func Test_LoadConfig_Swagger(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")

	// When
	swagger := GetAppConfig().Swagger

	// Then
	oauth2Scheme, exists := swagger.SecuritySchemes["oauth2"]
	assert.True(t, exists, "Expected oauth2 scheme to be defined")
	fmt.Println(oauth2Scheme)
	assert.Equal(t, "oauth2", oauth2Scheme.Type, "Type should be oauth2")
	assert.NotNil(t, oauth2Scheme.Flows.AuthorizationCode, "AuthorizationCode flow should not be nil")
	assert.Equal(t, "http://keycloak:7080/realms/master/protocol/openid-connect/auth", oauth2Scheme.Flows.AuthorizationCode.AuthorizationURL)
	assert.Equal(t, "http://keycloak:7080/realms/master/protocol/openid-connect/token", oauth2Scheme.Flows.AuthorizationCode.TokenURL)

	// Check Scopes
	scopes := oauth2Scheme.Flows.AuthorizationCode.Scopes
	assert.Equal(t, "OpenId Authentication", scopes["openid"], "OpenId scope should be defined")
	assert.Equal(t, "User Email", scopes["email"], "Email scope should be defined")
	assert.Equal(t, "User Profile", scopes["profile"], "Profile scope should be defined")
	assert.Equal(t, "User Roles", scopes["roles"], "Roles scope should be defined")

	// Check Security
	assert.Len(t, swagger.Security, 1, "Expected one security requirement")
	assert.Contains(t, swagger.Security[0], "oauth2", "Security should include oauth2")
	assert.Equal(t, []string{"openid", "email", "profile", "roles"}, swagger.Security[0]["oauth2"], "OAuth2 scopes should match")

}

func Test_LoadConfig_Kad(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	KadInstances := GetAppConfig().Kad
	// Then
	assert.Equal(t, "sandbox", KadInstances[0].ID, "Id")
	assert.Equal(t, "Sandbox de idir", KadInstances[0].Name, "Name")
	assert.Equal(t, "https://host.docker.internal:6553/api/kad/v1", KadInstances[0].APIURL, "ApiUrl")
	assert.Equal(t, "JUDtoP55C2dLfeaXqSbehhKKRdmAWTfj", KadInstances[0].AuthBearer, "AuthBearer")
	assert.True(t, KadInstances[0].InsecureSkipVerify, "InsecureSkipVerify should be true")
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
