package errors

import (
	"fmt"
	"net/http"

	"github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type ServerError _api.ServerError

const (
	// Error types
	Kad        = "kad"
	OkdpServer = "okdp_server"
)

func OfType(errorType string) *ServerError {
	return &ServerError{
		Type: errorType,
	}
}

func (s *ServerError) NotFoundError(messages ...interface{}) *ServerError {
	s.Message = toError(messages...)
	s.Status = http.StatusNotFound
	return s
}

func (s *ServerError) Forbidden(messages ...interface{}) *ServerError {
	s.Message = toError(messages...)
	s.Status = http.StatusForbidden
	return s
}

func (s *ServerError) GenericError(statusCode int, messages ...interface{}) *ServerError {
	s.Message = toError(messages...)
	s.Status = statusCode
	return s
}

func toError(messages ...interface{}) string {
	if len(messages) == 1 {
		return fmt.Errorf("%+v", messages...).Error()
	}
	return fmt.Errorf(messages[0].(string), messages[1:]...).Error()
}
