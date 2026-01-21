// Package validation provides shared validation utilities for domain entities.
// This package contains common validation patterns used across multiple domain packages,
// ensuring consistency and reducing code duplication.
package validation

import (
	"fmt"
	"strings"
)

// Result contains the outcome of validating an entity.
// It tracks whether validation passed, any errors encountered, and optional warnings.
type Result struct {
	Valid    bool
	Errors   []error
	Warnings []string
}

// NewResult creates a new validation result, initially valid.
func NewResult() *Result {
	return &Result{Valid: true, Errors: []error{}, Warnings: []string{}}
}

// AddError adds an error to the validation result and marks it invalid.
func (r *Result) AddError(err error) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

// AddWarning adds a warning to the validation result without marking it invalid.
func (r *Result) AddWarning(warning string) {
	r.Warnings = append(r.Warnings, warning)
}

// HasWarnings returns true if there are any warnings.
func (r *Result) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// Error returns a combined error message if there are validation errors.
// Returns nil if validation passed.
func (r *Result) Error() error {
	if r.Valid {
		return nil
	}
	var msgs []string
	for _, err := range r.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Errorf("validation failed: %s", strings.Join(msgs, "; "))
}

// Merge combines another validation result into this one.
// All errors and warnings from the other result are added to this result.
func (r *Result) Merge(other *Result) {
	if other == nil {
		return
	}
	for _, err := range other.Errors {
		r.AddError(err)
	}
	for _, warning := range other.Warnings {
		r.AddWarning(warning)
	}
}
