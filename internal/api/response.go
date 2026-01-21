// Package api provides HTTP handlers for the API.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
)

// ===== Standard Response Envelope Types =====
//
// All API responses follow a consistent envelope format:
//
// Success: {"data": ..., "meta": {...}, "warnings": [...]}
// Error:   {"error": {"code": "...", "message": "...", "details": {...}}}
//
// This ensures predictable response structures for API clients.

// Response is the standard success response envelope.
// All successful API responses are wrapped in this structure.
type Response struct {
	Data     interface{} `json:"data"`
	Meta     *Meta       `json:"meta,omitempty"`
	Warnings []string    `json:"warnings,omitempty"`
}

// Meta contains optional metadata for responses.
type Meta struct {
	// Pagination fields (for list responses)
	Total   *int64 `json:"total,omitempty"`
	Limit   *int   `json:"limit,omitempty"`
	Offset  *int   `json:"offset,omitempty"`
	HasMore *bool  `json:"hasMore,omitempty"`
}

// ErrorDetail represents the structured error information.
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse represents the standard API error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ===== Pagination Utilities =====
//
// All list endpoints use consistent offset-based pagination:
//   - Query params: limit (default: 20, max: 100) and offset (default: 0)
//   - Response: includes meta with total, limit, offset, hasMore

// PaginationDefaults defines the default pagination values.
const (
	DefaultLimit = 20
	MaxLimit     = 100
	DefaultOffset = 0
)

// Pagination holds parsed pagination parameters.
type Pagination struct {
	Limit  int
	Offset int
}

// ParsePagination extracts limit and offset from query parameters.
// Returns pagination with defaults applied and limits enforced.
func ParsePagination(query QueryGetter) Pagination {
	p := Pagination{
		Limit:  DefaultLimit,
		Offset: DefaultOffset,
	}

	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			p.Limit = parsed
			if p.Limit > MaxLimit {
				p.Limit = MaxLimit
			}
		}
	}

	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			p.Offset = parsed
		}
	}

	return p
}

// QueryGetter is an interface for getting query parameter values.
// This allows ParsePagination to work with url.Values or any similar type.
type QueryGetter interface {
	Get(key string) string
}

// ===== Legacy Response Types (deprecated, use standard envelope) =====

// PaginatedResponse wraps a list response with pagination metadata.
// Deprecated: Use Response with Meta for new code.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int64       `json:"totalPages"`
}

// ResponseWithWarnings wraps a response with optional warnings.
// Deprecated: Use Response with Warnings field for new code.
type ResponseWithWarnings struct {
	Data     interface{} `json:"data"`
	Warnings []string    `json:"warnings,omitempty"`
}

// ===== Response Helper Functions =====

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeData writes a success response with data wrapped in the standard envelope.
func writeData(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, Response{Data: data})
}

// writeDataWithWarnings writes a success response with data and warnings.
func writeDataWithWarnings(w http.ResponseWriter, status int, data interface{}, warnings []string) {
	writeJSON(w, status, Response{Data: data, Warnings: warnings})
}

// writePaginatedData writes a paginated list response with the standard envelope.
// It calculates hasMore based on offset, limit, and total.
func writePaginatedData(w http.ResponseWriter, status int, data interface{}, total int64, limit, offset int) {
	hasMore := int64(offset+limit) < total
	writeJSON(w, status, Response{
		Data: data,
		Meta: &Meta{
			Total:   &total,
			Limit:   &limit,
			Offset:  &offset,
			HasMore: &hasMore,
		},
	})
}

// writeError writes an error response using the standard error envelope.
func writeError(w http.ResponseWriter, status int, code, message string, details interface{}) {
	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	writeJSON(w, status, resp)
}

// WriteError is the exported version of writeError for use by middleware.
func WriteError(w http.ResponseWriter, status int, message string) {
	code := httpStatusToErrorCode(status)
	writeError(w, status, code, message, nil)
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
	code := domainErrorToCode(err)
	message := apperrors.GetMessage(err)

	// Log internal errors for debugging
	if apperrors.IsInternal(err) {
		log.Printf("Internal error: %v", err)
	}

	// Convert details slice to structured format if present
	var detailsObj interface{}
	if len(details) > 0 {
		detailsObj = map[string]interface{}{"validationErrors": details}
	}

	writeError(w, status, code, message, detailsObj)
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

// domainErrorToCode converts a domain error to a standardized error code string.
func domainErrorToCode(err error) string {
	switch {
	case apperrors.IsNotFound(err):
		return "NOT_FOUND"
	case apperrors.IsValidation(err):
		return "VALIDATION_ERROR"
	case apperrors.IsConflict(err):
		return "CONFLICT"
	case apperrors.IsForbidden(err):
		return "FORBIDDEN"
	case apperrors.IsUnauthorized(err):
		return "UNAUTHORIZED"
	case apperrors.IsBadRequest(err):
		return "BAD_REQUEST"
	default:
		return "INTERNAL_ERROR"
	}
}

// httpStatusToErrorCode converts an HTTP status code to a standardized error code.
func httpStatusToErrorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	default:
		return "INTERNAL_ERROR"
	}
}
