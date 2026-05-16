package httperr

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// AppError is the unified error response format (CLAUDE.md ルール4).
// JSON shape: { "code": "...", "message": "...", "details": {...} }
type AppError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	status  int
}

// New creates an AppError with the given code, message, and HTTP status.
func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, status: status}
}

// WithDetails attaches field-level validation errors.
func (e *AppError) WithDetails(d map[string]string) *AppError {
	e.Details = d
	return e
}

// Write writes the error as JSON to the response writer.
func (e *AppError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		slog.Error("httperr: failed to encode error response", slog.Any("error", err))
	}
}

func NotFound(code, message string) *AppError {
	return New(code, message, http.StatusNotFound)
}

func Unauthorized(message string) *AppError {
	return New("UNAUTHORIZED", message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New("FORBIDDEN", message, http.StatusForbidden)
}

func BadRequest(code, message string) *AppError {
	return New(code, message, http.StatusBadRequest)
}

func UnprocessableEntity(code, message string) *AppError {
	return New(code, message, http.StatusUnprocessableEntity)
}

func InternalError(message string) *AppError {
	return New("INTERNAL_ERROR", message, http.StatusInternalServerError)
}
