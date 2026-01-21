/*
Package validation provides shared validation utilities for domain entities.

This package consolidates common validation patterns that were previously duplicated
across multiple domain packages (lift, day, cycle, week, program, prescription, etc.),
ensuring consistency and reducing code duplication.

# Validation Result

The Result type tracks validation outcomes including errors and warnings:

	result := validation.NewResult()
	if err := validateSomething(value); err != nil {
		result.AddError(err)
	}
	if !result.Valid {
		return nil, result
	}

Result supports warnings for soft validation issues that don't prevent operation:

	result.AddWarning("Value is unusually high")

# Slug Validation and Generation

The package provides slug validation and generation utilities:

	// Validate a slug with max length
	err := validation.ValidateSlug(slug, 100)

	// Generate a slug from a name
	slug := validation.GenerateSlug("Bench Press") // returns "bench-press"

Slug rules:
  - Must contain only lowercase alphanumeric characters and hyphens
  - Cannot be empty
  - Cannot start or end with a hyphen
  - Cannot have consecutive hyphens
  - Must not exceed the specified max length

# Usage in Domain Packages

Domain packages should create type aliases for backward compatibility:

	type ValidationResult = validation.Result

	func NewValidationResult() *ValidationResult {
		return validation.NewResult()
	}

For slug errors, domain packages can export references to the shared errors:

	var (
		ErrSlugEmpty   = validation.ErrSlugEmpty
		ErrSlugInvalid = validation.ErrSlugInvalid
		ErrSlugTooLong = validation.SlugTooLongError(MaxSlugLength)
	)
*/
package validation
