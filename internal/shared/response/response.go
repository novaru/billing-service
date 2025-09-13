package response

import (
	"encoding/json"
	"errors"
	"net/http"

	E "github.com/novaru/billing-service/internal/shared/errors"
)

type APIResponse struct {
	Data    any        `json:"data,omitempty"`
	Error   *ErrorData `json:"error,omitempty"`
	Success bool       `json:"success"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// WriteSuccess writes a successful JSON response
func WriteSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, &APIResponse{
		Data:    data,
		Success: true,
	})
}

// WriteCreated writes a created response
func WriteCreated(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, &APIResponse{
		Data:    data,
		Success: true,
	})
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, e error) {
	var appErr *E.AppError
	var statusCode int
	var errorData *ErrorData

	// Check if it's our custom error
	if errors.As(e, &appErr) {
		statusCode = appErr.HTTPStatus()
		errorData = &ErrorData{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}
	} else {
		// Generic error
		statusCode = http.StatusInternalServerError
		errorData = &ErrorData{
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occurred",
		}
	}

	writeJSON(w, statusCode, &APIResponse{
		Error:   errorData,
		Success: false,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
