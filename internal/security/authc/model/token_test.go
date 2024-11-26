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
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_UserInfoFromAccessToken(t *testing.T) {
	// Given
	accessToken := "eyJh.eyJleHAiOjE3MzIxMjQxNDIsImlhdCI6MTczMjEyMzU0MiwiYXV0aF90aW1lIjoxNzMyMTIzNTQyLCJqdGkiOiJlYjgwNjhlOS1mODZjLTQ3ZTktOTFmYS0xODAzOGE5YjIxMGIiLCJpc3MiOiJodHRwOi8va2V5Y2xvYWs6NzA4MC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6Ijk2ZjM0ZjNjLTEzODYtNDdhYi05YTg4LTA3MDdhODcwYTQ3ZSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNvbmZpZGVudGlhbC1vaWRjLWNsaWVudCIsInNpZCI6ImEyMDcxYWE5LTU0MTYtNGE5Yy04MzM5LWQzYWM4YjA4OThiNSIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xvY2FsaG9zdDo4MDkwIiwiaHR0cDovL2xvY2FsaG9zdDo4MDkyIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLW1hc3RlciIsImRldmVsb3BlcnMiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiZGV2MSBkZXYxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiZGV2MSIsImdpdmVuX25hbWUiOiJkZXYxIiwiZmFtaWx5X25hbWUiOiJkZXYxIiwiZW1haWwiOiJkZXYxLmRldmVsb3BlcnNAZXhhbXBsZS5vcmcifQ.gc"
	token := &Token {AccessToken: accessToken}
	// When
	userInfo, err := token.GetUserInfo("realm_access.roles", "realm_access.groups")
	// Then
	assert.NoError(t, err)
	assert.Equal(t, "", userInfo.Login, "The user login does not match the expected result")
	assert.Equal(t, "dev1 dev1", userInfo.Name, "The user name does not match the expected result")
	assert.Equal(t, "96f34f3c-1386-47ab-9a88-0707a870a47e", userInfo.Subject, "The user Sub does not match the expected result")
	assert.Equal(t, "dev1.developers@example.org", userInfo.Email, "The user email does not match the expected result")
	assert.Equal(t, []string{"default-roles-master", "developers", "offline_access", "uma_authorization"}, userInfo.Roles, "The user roles does not match the expected result")
	assert.Equal(t, []string{}, userInfo.Groups, "The user groups does not match the expected result")
}

