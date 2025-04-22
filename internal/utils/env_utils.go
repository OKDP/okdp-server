/*
 *    Copyright 2025 okdp.io
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

package utils

import (
	"os"
	"strings"
)

// GetEnv reads an environment variable. If the value is a placeholder like $(VAR_NAME),
// it will replace it with the actual value from the environment variable.
// If the value is not a placeholder, it returns the value as is.
func GetEnv(key string) string {
	// Check if the key has the $(...) format
	if strings.HasPrefix(key, "$(") && strings.HasSuffix(key, ")") {
		varName := strings.TrimPrefix(key, "$(")
		varName = strings.TrimSuffix(varName, ")")
		if value, exists := os.LookupEnv(varName); exists {
			return value
		}
		return ""
	}
	return key
}
