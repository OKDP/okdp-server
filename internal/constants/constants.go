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

package constants

const (
	OkdpServerBaseUrl = "/api/v1"
	// OAuth2
	OAuth2SessionName = "OKDP_OAUT2_SESSION"
	OAuth2State       = "state"
	OAuth2Nonce       = "nonce"
	OAuth2UserInfo    = "userInfo"
	OAuth2LoginUrl    = "/oauth_login"
	// Authorization
	CasbinAuthzModel  = "pkg/auth/authz/authz-model.conf"
	CasbinAuthzPolicy = "pkg/auth/authz/authz-policy.csv"
	CasbinRolePrefix  = "role:"
	// Swagger API Docs URI
	SwaggerApiDocsUri = OkdpServerBaseUrl + "/api-docs"
	// KAD
	ComponentURL = "/mycluster/component"
	ComponentReleaseURL = "/mycluster/component-release"
	TemplateReleaseURL = "/mycluster/template-release"
	CatalogURL = "/mycluster/catalog"
)
