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
	"testing"
)

func TestResolveEnv(t *testing.T) {
	// Test the behavior when the environment variable is set
	os.Setenv("OCI_USERNAME", "myusername")
	defer os.Unsetenv("OCI_USERNAME")

	if got := ResolveEnv("$(OCI_USERNAME)"); got != "myusername" {
		t.Errorf("expected 'myusername', got %s", got)
	}

	// Test when the environment variable does not exist
	if got := ResolveEnv("$(NON_EXISTENT_VAR)"); got != "" {
		t.Errorf("expected '', got %s", got)
	}

	// Test non-placeholder value
	if got := ResolveEnv("JustSomeOtherValue"); got != "JustSomeOtherValue" {
		t.Errorf("expected 'JustSomeOtherValue', got %s", got)
	}
}
