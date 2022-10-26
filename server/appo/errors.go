package appo

import (
	"errors"
	"net/http"
)

const (
	Internal           = "000001"
	ServiceUnavailable = "000002"
	NotFound           = "000003"
)

func (r *FailResponse) Error() string {
	return r.Message
}

func (r *FailResponse) Status() int {
	switch r.ErrCode {
	case Internal:
		return http.StatusInternalServerError
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *FailResponse
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}
