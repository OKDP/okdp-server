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
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func RandomString() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type dockerConfig struct {
	Auths map[string]dockerAuthEntry `json:"auths"`
}

type dockerAuthEntry struct {
	Auth string `json:"auth"` // base64("username:password")
}

// ToLoginPassword decodes a base64-encoded .dockerconfigjson value and extracts the first
// username and password found in the `auths` section.
// It returns the username, password, and an error if decoding or parsing fails.
func ToLoginPassword(encodedDockerJSON string) (string, string, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(encodedDockerJSON)
	if err != nil {
		return "", "", fmt.Errorf("failed to base64-decode dockerjson: %w", err)
	}

	var config dockerConfig
	if err := json.Unmarshal(jsonBytes, &config); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal docker config json: %w", err)
	}

	if len(config.Auths) == 0 {
		return "", "", errors.New("no auth entries found in docker config")
	}

	for registry, entry := range config.Auths {
		if entry.Auth == "" {
			continue
		}

		decodedAuth, err := base64.StdEncoding.DecodeString(entry.Auth)
		if err != nil {
			return "", "", fmt.Errorf("failed to decode auth for registry '%s': %w", registry, err)
		}

		parts := strings.SplitN(string(decodedAuth), ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid auth format for registry '%s'", registry)
		}

		return parts[0], parts[1], nil
	}

	return "", "", errors.New("no valid auth credentials found")
}
