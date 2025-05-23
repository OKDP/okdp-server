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
	"github.com/stretchr/testify/require"
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

func Test_LoadConfig_Catalog(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	catalogs := GetAppConfig().Catalogs
	// Then
	require.NotEmpty(t, catalogs, "Catalogs should not be empty")
	catalog := catalogs[0]
	assert.True(t, catalog.IsAuthenticated(), "The catalog should be authenticated")
	assert.Equal(t, "infra01", catalog.ID, "ID")
	assert.Equal(t, "infra01 catalog", catalog.Name, "Name")
	assert.Equal(t, "My infrastructure components", catalog.Description, "Description")
	assert.Equal(t, "quay.io/okdp/applications", catalog.RepoURL, "RepoURL")
	assert.Equal(t, "$(OCI_USERNAME)", *catalog.Credentials.RobotAccountName, "Credentials.RobotAccountName")
	assert.Equal(t, "$(OCI_PASSWORD)", *catalog.Credentials.RobotAccountToken, "Credentials.RobotAccountToken")
	assert.Equal(t, "quay.io", catalog.RepoHost(), "RepoHost")

	require.NotEmpty(t, catalog.Packages, "Packages catalogs should not be empty")
	require.Len(t, catalog.Packages, 3, "The catalog should contain exactly 3 Packages")
	assert.Equal(t, "redis", catalog.Packages[0].Name, "Packages")
	assert.Equal(t, "podinfo", catalog.Packages[1].Name, "Packages")
	assert.Equal(t, "cert-manager", catalog.Packages[2].Name, "Packages")

	catalog = catalogs[1]
	assert.Equal(t, "infra02", catalog.ID, "ID")
	assert.False(t, catalog.IsAuthenticated(), fmt.Sprintf("The catalog '%s' should not be authenticated", catalog.ID))
}

func Test_LoadConfig_Clusters(t *testing.T) {
	// Given
	viper.Set("config", "testdata/application.yaml")
	// When
	clusters := GetAppConfig().Clusters
	// Then
	require.NotEmpty(t, clusters, "K8S clusters should not be empty")
	cluster := clusters[0]
	assert.Equal(t, "kubo03dev", cluster.ID, "ID")
	assert.Equal(t, "k8s infra dev", cluster.Name, "Name")
	assert.Equal(t, "dev", cluster.Env, "Env")
	assert.Equal(t, "/path/to/kubeconfig", cluster.Auth.Kubeconfig.Path, "cluster.Auth.Kubeconfig")
	assert.Equal(t, "dev-context", cluster.Auth.Kubeconfig.Context, "cluster.Auth.Context")
	assert.Equal(t, "https://host.docker.internal:56660", cluster.Auth.Kubeconfig.APIServer, "cluster.Auth.Kubeconfig.APIServer")
	assert.True(t, cluster.Auth.Kubeconfig.InsecureSkipTlsVerify, "cluster.Auth.Kubeconfig.InsecureSkipTlsVerify")

	assert.Equal(t, "https://k8s-api-server-url:6443", cluster.Auth.Certificate.APIServer, "cluster.Auth.Certificate.ApiServer")
	assert.Equal(t, "/path/to/client-key.pem", cluster.Auth.Certificate.ClientKey, "cluster.Auth.Certificate.ClientKey")
	assert.Equal(t, "/path/to/client-cert.pem", cluster.Auth.Certificate.ClientCert, "cluster.Auth.Certificate.ClientCert")
	assert.Equal(t, "/path/to/ca-cert.pem", cluster.Auth.Certificate.CACert, "cluster.Auth.Certificate.CaCert")

	assert.Equal(t, "https://k8s-api-server-url:6443", cluster.Auth.Bearer.APIServer, "cluster.Auth.Bearer.ApiServer")
	assert.Equal(t, "$(BEARER_TOKEN)", cluster.Auth.Bearer.BearerToken, "cluster.Auth.Bearer.BearerToken")

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
