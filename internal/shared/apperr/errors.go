package apperr

import (
	"errors"
	"net/http"
)

// AppError es un error con código HTTP asociado.
// Los services lo usan para señalar si el fallo es del cliente (4xx) o del servidor (5xx),
// permitiendo que los handlers elijan el status code sin inspeccionar el mensaje.
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

func NotFound(msg string) *AppError   { return &AppError{Code: http.StatusNotFound, Message: msg} }
func Conflict(msg string) *AppError   { return &AppError{Code: http.StatusConflict, Message: msg} }
func BadRequest(msg string) *AppError { return &AppError{Code: http.StatusBadRequest, Message: msg} }
func Forbidden(msg string) *AppError  { return &AppError{Code: http.StatusForbidden, Message: msg} }
func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

// StatusCode devuelve el HTTP status code del error.
// Si el error no es un AppError, asume fallo interno (500).
func StatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}
