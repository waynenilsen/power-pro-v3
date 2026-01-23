// Package api provides HTTP handlers for the API.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// ===== Filtering Utilities =====
//
// All list endpoints use consistent filter parameter naming:
//   - Simple filters: field names in snake_case (e.g., ?lift_id=123&user_id=456)
//   - Boolean filters: use "true"/"false" or "1"/"0" (e.g., ?is_competition_lift=true)
//   - Date ranges: use _after/_before suffixes (e.g., ?created_after=2024-01-01&created_before=2024-12-31)
//   - Numeric ranges: use _gte/_lte suffixes (e.g., ?weight_gte=100&weight_lte=200)
//
// Filter behavior:
//   - Multiple filters are combined with AND logic
//   - Unknown filter parameters are ignored
//   - Invalid filter values return validation errors

// FilterError represents a validation error for a filter parameter.
type FilterError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ParseFilterString extracts a string filter value from query parameters.
// Returns nil if the parameter is not present or empty.
func ParseFilterString(query QueryGetter, key string) *string {
	if v := query.Get(key); v != "" {
		return &v
	}
	return nil
}

// ParseFilterBool extracts a boolean filter value from query parameters.
// Accepts "true", "false", "1", "0" (case-insensitive).
// Returns nil if the parameter is not present or empty.
// Returns an error if the value is invalid.
func ParseFilterBool(query QueryGetter, key string) (*bool, error) {
	v := query.Get(key)
	if v == "" {
		return nil, nil
	}
	switch strings.ToLower(v) {
	case "true", "1":
		b := true
		return &b, nil
	case "false", "0":
		b := false
		return &b, nil
	default:
		return nil, apperrors.NewValidation(key, "must be true, false, 1, or 0")
	}
}

// ParseFilterDate extracts a date filter value from query parameters.
// Accepts ISO 8601 formats: RFC3339 (2024-01-15T10:00:00Z) or date-only (2024-01-15).
// Returns nil if the parameter is not present or empty.
// Returns an error if the value is invalid.
func ParseFilterDate(query QueryGetter, key string) (*time.Time, error) {
	v := query.Get(key)
	if v == "" {
		return nil, nil
	}
	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return &t, nil
	}
	// Try date-only format
	if t, err := time.Parse("2006-01-02", v); err == nil {
		return &t, nil
	}
	return nil, apperrors.NewValidation(key, "invalid format; use ISO 8601 format (e.g., 2024-01-15 or 2024-01-15T10:00:00Z)")
}

// ParseFilterDateEndOfDay extracts a date filter value and sets time to end of day for date-only format.
// This is useful for "before" or "until" date filters where you want to include the entire day.
// Returns nil if the parameter is not present or empty.
// Returns an error if the value is invalid.
func ParseFilterDateEndOfDay(query QueryGetter, key string) (*time.Time, error) {
	v := query.Get(key)
	if v == "" {
		return nil, nil
	}
	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return &t, nil
	}
	// Try date-only format and set to end of day
	if t, err := time.Parse("2006-01-02", v); err == nil {
		t = t.Add(24*time.Hour - time.Second)
		return &t, nil
	}
	return nil, apperrors.NewValidation(key, "invalid format; use ISO 8601 format (e.g., 2024-01-15 or 2024-01-15T10:00:00Z)")
}

// ParseFilterInt extracts an integer filter value from query parameters.
// Returns nil if the parameter is not present or empty.
// Returns an error if the value is invalid.
func ParseFilterInt(query QueryGetter, key string) (*int, error) {
	v := query.Get(key)
	if v == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return nil, apperrors.NewValidation(key, "must be a valid integer")
	}
	return &parsed, nil
}

// ParseFilterFloat extracts a float filter value from query parameters.
// Returns nil if the parameter is not present or empty.
// Returns an error if the value is invalid.
func ParseFilterFloat(query QueryGetter, key string) (*float64, error) {
	v := query.Get(key)
	if v == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, apperrors.NewValidation(key, "must be a valid number")
	}
	return &parsed, nil
}

// ParseFilterEnum extracts a string filter value and validates it against allowed values.
// The value is normalized to uppercase before validation.
// Returns nil if the parameter is not present or empty.
// Returns an error if the value is not in the allowed list.
func ParseFilterEnum(query QueryGetter, key string, allowedValues []string) (*string, error) {
	v := query.Get(key)
	if v == "" {
		return nil, nil
	}
	normalized := strings.ToUpper(v)
	for _, allowed := range allowedValues {
		if normalized == allowed {
			return &normalized, nil
		}
	}
	return nil, apperrors.NewValidation(key, "invalid value; valid values: "+strings.Join(allowedValues, ", "))
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
	_ = json.NewEncoder(w).Encode(data)
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
