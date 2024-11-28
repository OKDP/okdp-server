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

package basic

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	"github.com/okdp/okdp-server/internal/model/auth"
	"github.com/stretchr/testify/assert"
)

func Test_BasicAuth_Succeed(t *testing.T) {
    // Given
	basicProvider, err := NewProvider([]config.BasicAuth{
		{Login: "user1", Password: "secret1", Roles: []string{"admin"}},
	})
    // Create a ResponseRecorder to capture the response
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, router := gin.CreateTestContext(resp)
	router.Use(basicProvider.Auth()...)
	router.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, c.MustGet(gin.AuthUserKey).(string))
	})
	// When
	c.Request, _ = http.NewRequest(http.MethodGet, "/login", nil)
	c.Request.Header.Set("Authorization", basicAuthHeader("user1", "secret1"))
	router.HandleContext(c)

	// Then 
	// 1- Ensure the user was authenticated
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "user1", resp.Body.String())
	// 2- Ensure the userInfo was set in the context
	u, found := c.Get(constants.OAuth2UserInfo)
	assert.True(t, found, "The user was not found in the context")
	userInfo := u.(*model.UserInfo)
	assert.Equal(t, "user1", userInfo.Login)
	assert.Equal(t, []string{"admin"}, userInfo.Roles)

}

func Test_BasicAuth_Failed(t *testing.T) {
    // Given
	basicProvider, err := NewProvider([]config.BasicAuth{
		{Login: "user1", Password: "secret1", Roles: []string{"admin"}},
	})
    // Create a ResponseRecorder to capture the response
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, router := gin.CreateTestContext(resp)
	router.Use(basicProvider.Auth()...)
	router.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, c.MustGet(gin.AuthUserKey).(string))
	})
	// When
	c.Request, _ = http.NewRequest(http.MethodGet, "/login", nil)
	c.Request.Header.Set("Authorization", basicAuthHeader("user2", "secret1"))
	router.HandleContext(c)

	// Then 
	// 1- Ensure the user was not authorized
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	// 2- Ensure the userInfo was not set in the context
	_, found := c.Get(constants.OAuth2UserInfo)
	assert.False(t, found, "The user was found in the context")
}

func basicAuthHeader(login string, password string) string {
	return "Basic " +  base64.StdEncoding.EncodeToString([]byte(login + ":" + password))
}

