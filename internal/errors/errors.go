// Package errors provides standardized error types for the application.
// This package defines domain-specific error types that can be used across
// all packages to ensure consistent error handling and API responses.
package errors

import (
	"errors"
	"fmt"
)

// Standard error categories for HTTP status code mapping.
// These allow handlers to determine the appropriate HTTP response.
var (
	// ErrNotFound indicates a resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrValidation indicates input validation failed.
	ErrValidation = errors.New("validation failed")

	// ErrConflict indicates a conflict with existing data (e.g., duplicate slug).
	ErrConflict = errors.New("conflict")

	// ErrForbidden indicates the user lacks permission for the operation.
	ErrForbidden = errors.New("forbidden")

	// ErrUnauthorized indicates the user is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInternal indicates an internal server error.
	ErrInternal = errors.New("internal error")

	// ErrBadRequest indicates a malformed request.
	ErrBadRequest = errors.New("bad request")
)

// DomainError represents a domain-specific error with context.
// It implements the error interface and supports error wrapping.
type DomainError struct {
	// Category is the base error category for HTTP status code mapping.
	Category error
	// Message is the user-facing error message.
	Message string
	// Field is the optional field name for validation errors.
	Field string
	// Cause is the underlying error, if any.
	Cause error
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// Unwrap returns the underlying error for errors.Is and errors.As support.
func (e *DomainError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Category
}

// Is implements error comparison for errors.Is support.
func (e *DomainError) Is(target error) bool {
	if target == e.Category {
		return true
	}
	return false
}

// NewNotFound creates a not found error.
func NewNotFound(resource, identifier string) *DomainError {
	return &DomainError{
		Category: ErrNotFound,
		Message:  fmt.Sprintf("%s not found: %s", resource, identifier),
	}
}

// NewValidation creates a validation error for a specific field.
func NewValidation(field, message string) *DomainError {
	return &DomainError{
		Category: ErrValidation,
		Message:  message,
		Field:    field,
	}
}

// NewValidationMsg creates a validation error without a specific field.
func NewValidationMsg(message string) *DomainError {
	return &DomainError{
		Category: ErrValidation,
		Message:  message,
	}
}

// NewConflict creates a conflict error.
func NewConflict(message string) *DomainError {
	return &DomainError{
		Category: ErrConflict,
		Message:  message,
	}
}

// NewForbidden creates a forbidden error.
func NewForbidden(message string) *DomainError {
	return &DomainError{
		Category: ErrForbidden,
		Message:  message,
	}
}

// NewUnauthorized creates an unauthorized error.
func NewUnauthorized(message string) *DomainError {
	return &DomainError{
		Category: ErrUnauthorized,
		Message:  message,
	}
}

// NewInternal creates an internal error with an underlying cause.
func NewInternal(message string, cause error) *DomainError {
	return &DomainError{
		Category: ErrInternal,
		Message:  message,
		Cause:    cause,
	}
}

// NewBadRequest creates a bad request error.
func NewBadRequest(message string) *DomainError {
	return &DomainError{
		Category: ErrBadRequest,
		Message:  message,
	}
}

// Wrap wraps an error with additional context while preserving the category.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return &DomainError{
			Category: domainErr.Category,
			Message:  message,
			Cause:    err,
		}
	}

	return &DomainError{
		Category: ErrInternal,
		Message:  message,
		Cause:    err,
	}
}

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsValidation checks if an error is a validation error.
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsConflict checks if an error is a conflict error.
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsForbidden checks if an error is a forbidden error.
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsUnauthorized checks if an error is an unauthorized error.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsInternal checks if an error is an internal error.
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

// IsBadRequest checks if an error is a bad request error.
func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

// GetCategory extracts the error category from an error.
// Returns ErrInternal if the error is not a DomainError.
func GetCategory(err error) error {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Category
	}
	return ErrInternal
}

// GetMessage extracts the message from an error.
// Returns the error string if not a DomainError.
func GetMessage(err error) string {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Message
	}
	if err != nil {
		return err.Error()
	}
	return ""
}
