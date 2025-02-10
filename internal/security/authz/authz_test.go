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

package authz

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/security/authc/model"
	"github.com/stretchr/testify/assert"
)

func Test_AuthZ_Succeed(t *testing.T) {
	// Given
	log.SetupGlobalLogger(config.Logging{})
	authzConfig := config.AuthZ{Provider: "file",
		File: config.FileAuthZ{
			ModelPath:  "testdata/authz-model.conf",
			PolicyPath: "testdata/authz-policy.csv",
		},
	}
	called := false
	// Create a ResponseRecorder to capture the response
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, router := gin.CreateTestContext(resp)
	router.Use(Set(constants.OAuth2UserInfo, &model.UserInfo{Roles: []string{"developers"}}))

	router.Use(Authorizer(authzConfig))
	router.GET("/api/v1/spaces/1/composition/1/deployment", func(_ *gin.Context) { called = true })

	// When
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/spaces/1/composition/1/deployment", nil)
	router.HandleContext(c)

	// Then - Ensure the user was authorized
	assert.Equal(t, http.StatusOK, resp.Code, "The response code should be 200")
	assert.True(t, called, "The endpoint should be called")
}

func Test_AuthZ_Failed(t *testing.T) {
	// Given
	log.SetupGlobalLogger(config.Logging{})
	authzConfig := config.AuthZ{Provider: "file",
		File: config.FileAuthZ{
			ModelPath:  "testdata/authz-model.conf",
			PolicyPath: "testdata/authz-policy.csv",
		},
	}
	called := false
	// Create a ResponseRecorder to capture the response
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, router := gin.CreateTestContext(resp)
	router.Use(Set(constants.OAuth2UserInfo, &model.UserInfo{Roles: []string{"developers"}}))

	router.Use(Authorizer(authzConfig))
	router.GET("/api/v1/spaces/2/composition/1/deployment", func(_ *gin.Context) { called = true })

	// When
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/spaces/2/composition/1/deployment", nil)
	router.HandleContext(c)

	// Then - Ensure the user was not authorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code, "The response code should be 401")
	assert.False(t, called, "The endpoint should not be called")
}

func Test_AuthZ_InLine(t *testing.T) {
	// Given
	log.SetupGlobalLogger(config.Logging{})
	authzConfig := config.AuthZ{Provider: "inline",
		InLine: config.InLineAuthZ{
			Policy: `
  p, role:viewers, /api/v1/spaces/1/composition/*/deployment, GET
  p, role:developers, /api/v1/spaces/1/composition/*/deployment, *
  p, role:admins, /api/v1/spaces/*/composition, *

  g, role:admins, role:developers
  g, role:developers, role:viewers
  `,

			Model: `
  # casbin AuthZ configuration file
  [request_definition]
  r = sub, obj, act

  [policy_definition]
  p = sub, obj, act

  [role_definition]
  g = _, _

  [policy_effect]
  e = some(where (p.eft == allow))

  [matchers]
  m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
  `},
	}

	called := false
	// Create a ResponseRecorder to capture the response
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, router := gin.CreateTestContext(resp)
	router.Use(Set(constants.OAuth2UserInfo, &model.UserInfo{Roles: []string{"developers"}}))

	router.Use(Authorizer(authzConfig))
	router.GET("/api/v1/spaces/2/composition/1/deployment", func(_ *gin.Context) { called = true })

	// When
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/spaces/2/composition/1/deployment", nil)
	router.HandleContext(c)

	// Then - Ensure the user was not authorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code, "The response code should be 401")
	assert.False(t, called, "The endpoint should not be called")
}

func Set(key string, value any) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, value)
	}
}
