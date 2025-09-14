package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Domain error types
var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal server error")
)

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the appropriate HTTP status code
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case "NOT_FOUND":
		return http.StatusNotFound
	case "ALREADY_EXISTS":
		return http.StatusConflict
	case "INVALID_INPUT":
		return http.StatusBadRequest
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Constructor functions
func NewNotFoundError(resource string, details ...string) *AppError {
	msg := fmt.Sprintf("%s not found", resource)
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    "NOT_FOUND",
		Message: msg,
		Details: detail,
		Err:     ErrNotFound,
	}
}

func NewInvalidInputError(msg string, err error) *AppError {
	return &AppError{
		Code:    "INVALID_INPUT",
		Message: msg,
		Err:     err,
	}
}

func NewUnauthorizedError(msg string, err error) *AppError {
	return &AppError{
		Code:    "UNAUTHORIZED",
		Message: msg,
		Err:     err,
	}
}

func NewInternalError(msg string, err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: msg,
		Err:     err,
	}
}

func NewAlreadyExistsError(resource string, details string) *AppError {
	return &AppError{
		Code:    "ALREADY_EXISTS",
		Message: fmt.Sprintf("%s already exists", resource),
		Details: details,
		Err:     ErrAlreadyExists,
	}
}
