package apperrors

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Code string

const (
	Internal           Code = "000000"
	NotFound           Code = "000001"
	BadRequest         Code = "000002"
	Authorization      Code = "000003"
	ServiceUnavailable Code = "000004"
)

type Error struct {
	T       int64  `json:"t"`
	Success bool   `json:"success"`
	ErrCode Code   `json:"err_code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.ErrCode {
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case BadRequest:
		return http.StatusBadRequest
	case Authorization:
		return http.StatusUnauthorized
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

func New(errCode Code, message string) *Error {
	return &Error{
		T:       time.Now().UnixMilli(),
		Success: false,
		ErrCode: errCode,
		Message: message,
	}
}

func NewInternal() (e *Error) {
	e = New(Internal, ServerError)
	return
}

func NewNotFound(name string, value string) (e *Error) {
	e = New(NotFound, fmt.Sprintf("resource: (%v) with value (%v) not found", name, value))
	return
}

func NewBadRequest(reason string) (e *Error) {
	e = New(BadRequest, fmt.Sprintf("Bad request. Reason: (%v)", reason))
	return
}

func NewAuthorization(reason string) (e *Error) {
	e = New(Authorization, reason)
	return
}

func NewServiceUnavailable() (e *Error) {
	e = New(ServiceUnavailable, "Service unavailable or timed out")
	return
}
