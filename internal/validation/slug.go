// Package validation provides shared validation utilities for domain entities.
package validation

import (
	"errors"
	"regexp"
	"strings"
)

// SlugPattern matches valid slugs: lowercase alphanumeric with hyphens.
// Valid examples: "bench-press", "squat", "overhead-press-1"
// Invalid examples: "Bench_Press", "squat!", "overhead--press"
var SlugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Slug validation errors
var (
	ErrSlugEmpty   = errors.New("slug cannot be empty")
	ErrSlugInvalid = errors.New("slug must contain only lowercase alphanumeric characters and hyphens")
)

// SlugTooLongError returns an error for slugs exceeding the given max length.
func SlugTooLongError(maxLength int) error {
	return &slugTooLongError{maxLength: maxLength}
}

type slugTooLongError struct {
	maxLength int
}

func (e *slugTooLongError) Error() string {
	return "slug must be " + itoa(e.maxLength) + " characters or less"
}

// Is implements the errors.Is interface for comparing slugTooLongError instances.
// Two slugTooLongError instances are equal if they have the same maxLength.
func (e *slugTooLongError) Is(target error) bool {
	t, ok := target.(*slugTooLongError)
	if !ok {
		return false
	}
	return e.maxLength == t.maxLength
}

// itoa converts int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// ValidateSlug validates a slug against the standard slug rules.
// maxLength specifies the maximum allowed length (e.g., 50 or 100).
func ValidateSlug(slug string, maxLength int) error {
	if slug == "" {
		return ErrSlugEmpty
	}
	if len(slug) > maxLength {
		return SlugTooLongError(maxLength)
	}
	if !SlugPattern.MatchString(slug) {
		return ErrSlugInvalid
	}
	return nil
}

// GenerateSlug creates a URL-safe slug from a name.
// It converts to lowercase, replaces spaces and special characters with hyphens,
// removes consecutive hyphens, and trims leading/trailing hyphens.
func GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and common special characters with hyphens
	replacer := strings.NewReplacer(
		" ", "-",
		"_", "-",
		".", "-",
		",", "",
		"'", "",
		"\"", "",
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"/", "-",
		"\\", "-",
		"&", "-",
	)
	slug = replacer.Replace(slug)

	// Remove any characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
