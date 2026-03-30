package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type errorCode string

const (
	errorCodeBadRequest errorCode = "bad_request"
	errorCodeNotFound   errorCode = "not_found"
	errorCodeConflict   errorCode = "conflict"
	errorCodeForbidden  errorCode = "forbidden"
)

type serviceError struct {
	code    errorCode
	message string
	err     error
}

func (e *serviceError) Error() string {
	if e.message != "" {
		return e.message
	}
	if e.err != nil {
		return e.err.Error()
	}
	return string(e.code)
}

func (e *serviceError) Unwrap() error {
	return e.err
}

func badRequestErrorf(format string, args ...interface{}) error {
	return &serviceError{code: errorCodeBadRequest, message: fmt.Sprintf(format, args...)}
}

func badRequestWrap(err error, format string, args ...interface{}) error {
	return &serviceError{code: errorCodeBadRequest, message: fmt.Sprintf(format, args...), err: err}
}

func notFoundErrorf(format string, args ...interface{}) error {
	return &serviceError{code: errorCodeNotFound, message: fmt.Sprintf(format, args...)}
}

func conflictErrorf(format string, args ...interface{}) error {
	return &serviceError{code: errorCodeConflict, message: fmt.Sprintf(format, args...)}
}

func forbiddenErrorf(format string, args ...interface{}) error {
	return &serviceError{code: errorCodeForbidden, message: fmt.Sprintf(format, args...)}
}

func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var svcErr *serviceError
	if errors.As(err, &svcErr) {
		switch svcErr.code {
		case errorCodeBadRequest:
			return http.StatusBadRequest
		case errorCodeNotFound:
			return http.StatusNotFound
		case errorCodeConflict:
			return http.StatusConflict
		case errorCodeForbidden:
			return http.StatusForbidden
		default:
			return http.StatusInternalServerError
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func ErrorTitle(err error) string {
	switch StatusCode(err) {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusForbidden:
		return "Forbidden"
	default:
		return "Internal Server Error"
	}
}
