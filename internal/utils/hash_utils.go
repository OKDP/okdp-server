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

package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// FromString generates a short deterministic ID from a string
func FromString(input string) string {
	normalized := strings.TrimSpace(input)
	hash := sha1.Sum([]byte(normalized))
	return hex.EncodeToString(hash[:6]) // 6 bytes -> 12 hex chars
}
