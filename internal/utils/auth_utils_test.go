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
	"encoding/base64"
	"testing"
)

func TestToLoginPassword(t *testing.T) {
	// Setup test values
	username := "myuser"
	password := "mypassword"
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	dockerJSON := `{
		"auths": {
			"https://index.docker.io/v1/": {
				"auth": "` + auth + `"
			}
		}
	}`

	encodedDockerJSON := base64.StdEncoding.EncodeToString([]byte(dockerJSON))

	user, pass, err := ToLoginPassword(encodedDockerJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user != username {
		t.Errorf("expected username '%s', got '%s'", username, user)
	}
	if pass != password {
		t.Errorf("expected password '%s', got '%s'", password, pass)
	}
}

func TestToLoginPassword_InvalidBase64(t *testing.T) {
	_, _, err := ToLoginPassword("not-base64!")
	if err == nil {
		t.Fatal("expected error for invalid base64, got none")
	}
}

func TestToLoginPassword_InvalidJSON(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("not json"))
	_, _, err := ToLoginPassword(encoded)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got none")
	}
}

func TestToLoginPassword_MissingAuth(t *testing.T) {
	dockerJSON := `{"auths": {}}`
	encoded := base64.StdEncoding.EncodeToString([]byte(dockerJSON))
	_, _, err := ToLoginPassword(encoded)
	if err == nil {
		t.Fatal("expected error for missing auths, got none")
	}
}
