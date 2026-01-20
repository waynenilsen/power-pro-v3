// Package api provides HTTP handlers for the API.
package api

import (
	"encoding/json"
	"net/http"
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

// readJSON reads JSON from the request body.
func readJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
