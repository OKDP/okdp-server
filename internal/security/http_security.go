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

package security

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/okdp/okdp-server/internal/config"
)

func HttpSecurity(securityConfig config.Security) []gin.HandlerFunc {
	var handlers = []gin.HandlerFunc{}
	return append(handlers, Cors(securityConfig.Cors), SecurityHeaders(securityConfig.Headers))
}

func Cors(corsConfig config.Cors) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     corsConfig.AllowedMethods,
		AllowHeaders:     corsConfig.AllowedHeaders,
		ExposeHeaders:    corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           time.Duration(corsConfig.MaxAge) * time.Second,
	})
}

func SecurityHeaders(headersConf map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for header, value := range headersConf {
			c.Header(header, value)
		}
		c.Next()
	}
}
