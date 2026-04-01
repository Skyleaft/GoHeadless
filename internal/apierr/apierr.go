package apierr

import "net/http"

// AppError is a structured API error that carries an HTTP status code,
// a client-facing message, and an optional internal cause for server-side logging.
type AppError struct {
	Code     int
	Message  string
	Internal error // logged server-side, never sent to client
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Internal }

// Constructors

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func Internal(msg string, cause error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Internal: cause}
}

func Wrap(err error, code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg, Internal: err}
}
