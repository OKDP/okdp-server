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

package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParametersSchemaSimple(t *testing.T) {
	// Given
	jsonData := []byte(`{
		"parametersSchema": {
			"properties": {
				"hostname": {
					"type": "string",
					"required": false
				},
				"tls": {
					"type": "boolean",
					"default": false
				}
			}
		}
	}`)

	// When
	var s KuboSchema
	err := json.Unmarshal(jsonData, &s)

	// Then
	assert.NoError(t, err)

	assert.NotNil(t, s.ParametersSchema)
	props := s.ParametersSchema.Properties

	assert.Contains(t, props, "hostname")
	assert.Equal(t, "string", props["hostname"].Type)
	assert.Equal(t, false, props["hostname"].Required)

	assert.Contains(t, props, "tls")
	assert.Equal(t, "boolean", props["tls"].Type)
	assert.Equal(t, false, props["tls"].Default)
}

func TestParametersSchemaNested(t *testing.T) {
	// Given
	jsonData := []byte(`{
		"parametersSchema": {
			"properties": {
				"trust": {
					"properties": {
						"enabled": {
							"type": "boolean",
							"default": false
						}
					}
				},
				"issuers": {
					"properties": {
						"enabled": {
							"type": "boolean",
							"default": true
						},
						"caClusterIssuers": {
							"items": {
								"properties": {
									"name": {
										"type": "string",
										"required": true
									},
									"ca_crt": {
										"type": "string",
										"required": true
									},
									"ca_key": {
										"type": "string",
										"required": true
									}
								}
							}
						},
						"selfSignedClusterIssuers": {
							"items": {
								"type": "object",
								"properties": {
									"name": {
										"type": "string",
										"required": true
									}
								}
							}
						}
					}
				}
			}
		}
	}`)

	// When
	var s KuboSchema
	err := json.Unmarshal(jsonData, &s)

	// Then
	assert.NoError(t, err)

	// Validate top-level keys
	assert.NotNil(t, s.ParametersSchema)
	assert.Contains(t, s.ParametersSchema.Properties, "trust")
	assert.Contains(t, s.ParametersSchema.Properties, "issuers")

	// Validate nested fields
	trust := s.ParametersSchema.Properties["trust"]
	assert.Contains(t, trust.Properties, "enabled")
	assert.Equal(t, "boolean", trust.Properties["enabled"].Type)
	assert.Equal(t, false, trust.Properties["enabled"].Default)

	issuers := s.ParametersSchema.Properties["issuers"]
	assert.Contains(t, issuers.Properties, "enabled")
	assert.Equal(t, true, issuers.Properties["enabled"].Default)

	caClusterIssuers := issuers.Properties["caClusterIssuers"]
	assert.NotNil(t, caClusterIssuers.Items)
	assert.Contains(t, caClusterIssuers.Items.Properties, "name")
	assert.Equal(t, "string", caClusterIssuers.Items.Properties["name"].Type)
}

func TestParametersSchemaExtended(t *testing.T) {
	// Given
	jsonData := []byte(`{
		"parametersSchema": {
			"description": "Redis stack",
			"properties": {
				"redis": {
					"required": false,
					"properties": {
						"password": {
							"type": "string",
							"default": "redis123"
						},
						"replicaCount": {
							"type": "integer",
							"default": 1,
							"description": "The number of replicas"
						}
					}
				},
				"commander": {
					"required": true,
					"properties": {
						"enabled": {
							"type": "boolean",
							"default": true
						},
						"tls": {
							"type": "boolean",
							"default": false
						},
						"hostname": {
							"type": "string",
							"required": false
						}
					}
				}
			}
		}
	}`)

	// When
	var s KuboSchema
	err := json.Unmarshal(jsonData, &s)

	// Then
	assert.NoError(t, err)

	ps := s.ParametersSchema
	assert.Equal(t, "Redis stack", ps.Description)

	redis := ps.Properties["redis"]
	assert.NotNil(t, redis)
	assert.False(t, redis.Required)
	assert.Equal(t, "redis123", redis.Properties["password"].Default)
	assert.Equal(t, float64(1), redis.Properties["replicaCount"].Default)

	commander := ps.Properties["commander"]
	assert.NotNil(t, commander)
	assert.True(t, commander.Required)
	assert.Equal(t, true, commander.Properties["enabled"].Default)
	assert.Equal(t, false, commander.Properties["tls"].Default)
}

func TestContextSchema(t *testing.T) {
	// Given
	jsonData := []byte(`{
		"contextSchema": {
			"properties": {
				"ingress": {
					"required": true,
					"properties": {
						"className": {
							"type": "string",
							"default": "nginx"
						},
						"hostPostfix": {
							"type": "string",
							"required": true
						}
					}
				}
			}
		}
	}`)

	// When
	var s KuboSchema
	err := json.Unmarshal(jsonData, &s)

	// Then
	assert.NoError(t, err)

	ingress := s.ContextSchema.Properties["ingress"]
	assert.NotNil(t, ingress)
	assert.True(t, ingress.Required)

	className := ingress.Properties["className"]
	assert.NotNil(t, className)
	assert.Equal(t, "string", className.Type)
	assert.Equal(t, "nginx", className.Default)

	hostPostfix := ingress.Properties["hostPostfix"]
	assert.NotNil(t, hostPostfix)
	assert.True(t, hostPostfix.Required)
}
