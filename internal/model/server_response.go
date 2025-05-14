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
	"fmt"
	"net/http"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type ServerResponse _api.ServerResponse
type ServerResponseType _api.ServerResponseType

const (
	OkdpServerResponse ServerResponseType = "okdp_server"
	RegistryResponse   ServerResponseType = "registry"
	GitRepoResponse    ServerResponseType = "git_repo"
	K8sClusterResponse ServerResponseType = "k8s_cluster"
)

func NewServerResponse(errorType ServerResponseType) *ServerResponse {
	return &ServerResponse{
		Type: _api.ServerResponseType(errorType),
	}
}

func (s *ServerResponse) NotFoundError(messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = http.StatusNotFound
	return s
}

func (s *ServerResponse) BadRequest(messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = http.StatusBadRequest
	return s
}

func (s *ServerResponse) Unauthorized(messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = http.StatusUnauthorized
	return s
}

func (s *ServerResponse) Forbidden(messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = http.StatusForbidden
	return s
}

func (s *ServerResponse) ConflictError(messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = http.StatusConflict
	return s
}

func (s *ServerResponse) UnprocessableEntity(messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = http.StatusUnprocessableEntity
	return s
}

func (s *ServerResponse) GenericError(statusCode int, messages ...interface{}) *ServerResponse {
	s.Message = toError(messages...)
	s.Status = statusCode
	return s
}

func (s *ServerResponse) Created(messages ...interface{}) *ServerResponse {
	s.Message = toSuccess(messages...)
	s.Status = http.StatusCreated
	return s
}

func (s *ServerResponse) Updated(messages ...interface{}) *ServerResponse {
	s.Message = toSuccess(messages...)
	s.Status = http.StatusOK
	return s
}

func (s *ServerResponse) Deleted(messages ...interface{}) *ServerResponse {
	s.Message = toSuccess(messages...)
	s.Status = http.StatusOK
	return s
}

func (s *ServerResponse) IsNotfound() bool {
	return s.Status == http.StatusNotFound
}

func toError(messages ...interface{}) string {
	if len(messages) == 1 {
		return fmt.Errorf("%+v", messages...).Error()
	}
	return fmt.Errorf(messages[0].(string), messages[1:]...).Error()
}

func toSuccess(messages ...interface{}) string {
	if len(messages) == 1 {
		return messages[0].(string)
	}
	return fmt.Sprintf(messages[0].(string), messages[1:]...)
}
