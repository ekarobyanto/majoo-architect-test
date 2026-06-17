package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	CodeInternal      ErrorCode = "INTERNAL_ERROR"
	CodeNotFound      ErrorCode = "NOT_FOUND"
	CodeBadRequest    ErrorCode = "BAD_REQUEST"
	CodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	CodeForbidden     ErrorCode = "FORBIDDEN"
	CodeUnprocessable ErrorCode = "UNPROCESSABLE_ENTITY"
	CodeConflict      ErrorCode = "CONFLICT"
)

type AppError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAppError(status int, code ErrorCode, message string, details ...interface{}) *AppError {
	var d interface{}
	if len(details) > 0 {
		d = details[0]
	}
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Details:    d,
	}
}

func Internal(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusInternalServerError, CodeInternal, message, details...)
}

func NotFound(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusNotFound, CodeNotFound, message, details...)
}

func BadRequest(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusBadRequest, CodeBadRequest, message, details...)
}

func Unauthorized(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusUnauthorized, CodeUnauthorized, message, details...)
}

func Forbidden(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusForbidden, CodeForbidden, message, details...)
}

func Unprocessable(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, CodeUnprocessable, message, details...)
}

func Conflict(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusConflict, CodeConflict, message, details...)
}

func ValidationError(message string, details ...interface{}) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, CodeUnprocessable, message, details...)
}
