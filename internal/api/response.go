// Package api provides HTTP handlers for the API.
package api

import (
	"encoding/json"
	"log"
	"net/http"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
)

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int64       `json:"totalPages"`
}

// ResponseWithWarnings wraps a response with optional warnings.
// Used when an operation succeeds but has informational warnings.
type ResponseWithWarnings struct {
	Data     interface{} `json:"data"`
	Warnings []string    `json:"warnings,omitempty"`
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, message string, details ...string) {
	resp := ErrorResponse{
		Error:   message,
		Details: details,
	}
	writeJSON(w, status, resp)
}

// WriteError is the exported version of writeError for use by middleware.
func WriteError(w http.ResponseWriter, status int, message string) {
	writeError(w, status, message)
}

// readJSON reads JSON from the request body.
func readJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// writeDomainError writes an error response based on a domain error.
// It automatically maps the error category to the appropriate HTTP status code
// and logs internal errors for debugging.
func writeDomainError(w http.ResponseWriter, err error, details ...string) {
	status := mapErrorToStatus(err)
	message := apperrors.GetMessage(err)

	// Log internal errors for debugging
	if apperrors.IsInternal(err) {
		log.Printf("Internal error: %v", err)
	}

	writeError(w, status, message, details...)
}

// mapErrorToStatus maps a domain error to an HTTP status code.
func mapErrorToStatus(err error) int {
	switch {
	case apperrors.IsNotFound(err):
		return http.StatusNotFound
	case apperrors.IsValidation(err):
		return http.StatusBadRequest
	case apperrors.IsConflict(err):
		return http.StatusConflict
	case apperrors.IsForbidden(err):
		return http.StatusForbidden
	case apperrors.IsUnauthorized(err):
		return http.StatusUnauthorized
	case apperrors.IsBadRequest(err):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
