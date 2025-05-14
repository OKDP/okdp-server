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
	// OkdpServerBaseURL is the API URI
	OkdpServerBaseURL = "/api/v1"
	// OAuth2LoginURL is the OAuth2 login URI
	OAuth2LoginURL = "/oauth_login"
	// OAuth2SessionName is name of the OAuth2 user session
	OAuth2SessionName = "OKDP_OAUT2_SESSION"
	OAuth2State       = "state"
	OAuth2Nonce       = "nonce"
	OAuth2UserInfo    = "userInfo"
	// CasbinRolePrefix is used to prefix the roles/groups in the casbin policy (p, role:viewers, /api/v1/users/myprofile, *)
	CasbinRolePrefix = "role:"
	// SwaggerAPIDocsURI is the swagger API Docs public URI
	SwaggerAPIDocsURI = OkdpServerBaseURL + "/api-docs"
	HealthzURI        = "/healthz"
	ReadinessURI      = "/readiness"
	All               = "All"
	GitRepository     = "GitRepository"
)
