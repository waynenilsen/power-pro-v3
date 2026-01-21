package api_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/api"
)

// TestParseFilterString tests the string filter parsing utility.
func TestParseFilterString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		query    url.Values
		expected *string
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "lift_id",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "returns nil for empty parameter",
			key:      "lift_id",
			query:    url.Values{"lift_id": {""}},
			expected: nil,
		},
		{
			name:     "returns value for present parameter",
			key:      "lift_id",
			query:    url.Values{"lift_id": {"abc123"}},
			expected: strPtr("abc123"),
		},
		{
			name:     "returns first value when multiple provided",
			key:      "lift_id",
			query:    url.Values{"lift_id": {"abc123", "def456"}},
			expected: strPtr("abc123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := api.ParseFilterString(tt.query, tt.key)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %q", *result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %q, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("expected %q, got %q", *tt.expected, *result)
				}
			}
		})
	}
}

// TestParseFilterBool tests the boolean filter parsing utility.
func TestParseFilterBool(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		query       url.Values
		expected    *bool
		expectError bool
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "is_active",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "returns nil for empty parameter",
			key:      "is_active",
			query:    url.Values{"is_active": {""}},
			expected: nil,
		},
		{
			name:     "parses true",
			key:      "is_active",
			query:    url.Values{"is_active": {"true"}},
			expected: boolPtr(true),
		},
		{
			name:     "parses TRUE (case insensitive)",
			key:      "is_active",
			query:    url.Values{"is_active": {"TRUE"}},
			expected: boolPtr(true),
		},
		{
			name:     "parses True (mixed case)",
			key:      "is_active",
			query:    url.Values{"is_active": {"True"}},
			expected: boolPtr(true),
		},
		{
			name:     "parses 1 as true",
			key:      "is_active",
			query:    url.Values{"is_active": {"1"}},
			expected: boolPtr(true),
		},
		{
			name:     "parses false",
			key:      "is_active",
			query:    url.Values{"is_active": {"false"}},
			expected: boolPtr(false),
		},
		{
			name:     "parses FALSE (case insensitive)",
			key:      "is_active",
			query:    url.Values{"is_active": {"FALSE"}},
			expected: boolPtr(false),
		},
		{
			name:     "parses 0 as false",
			key:      "is_active",
			query:    url.Values{"is_active": {"0"}},
			expected: boolPtr(false),
		},
		{
			name:        "returns error for invalid value",
			key:         "is_active",
			query:       url.Values{"is_active": {"maybe"}},
			expectError: true,
		},
		{
			name:        "returns error for numeric other than 0/1",
			key:         "is_active",
			query:       url.Values{"is_active": {"2"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseFilterBool(tt.query, tt.key)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.expected == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", *result)
					}
				} else {
					if result == nil {
						t.Errorf("expected %v, got nil", *tt.expected)
					} else if *result != *tt.expected {
						t.Errorf("expected %v, got %v", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// TestParseFilterDate tests the date filter parsing utility.
func TestParseFilterDate(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		query       url.Values
		expected    *time.Time
		expectError bool
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "start_date",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "returns nil for empty parameter",
			key:      "start_date",
			query:    url.Values{"start_date": {""}},
			expected: nil,
		},
		{
			name:     "parses RFC3339 format",
			key:      "start_date",
			query:    url.Values{"start_date": {"2024-01-15T10:30:00Z"}},
			expected: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
		},
		{
			name:     "parses date-only format",
			key:      "start_date",
			query:    url.Values{"start_date": {"2024-01-15"}},
			expected: timePtr(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:        "returns error for invalid format",
			key:         "start_date",
			query:       url.Values{"start_date": {"01/15/2024"}},
			expectError: true,
		},
		{
			name:        "returns error for garbage input",
			key:         "start_date",
			query:       url.Values{"start_date": {"not-a-date"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseFilterDate(tt.query, tt.key)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.expected == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", *result)
					}
				} else {
					if result == nil {
						t.Errorf("expected %v, got nil", *tt.expected)
					} else if !result.Equal(*tt.expected) {
						t.Errorf("expected %v, got %v", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// TestParseFilterDateEndOfDay tests the end-of-day date filter parsing utility.
func TestParseFilterDateEndOfDay(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		query       url.Values
		expected    *time.Time
		expectError bool
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "end_date",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "RFC3339 is returned as-is",
			key:      "end_date",
			query:    url.Values{"end_date": {"2024-01-15T10:30:00Z"}},
			expected: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
		},
		{
			name:     "date-only is adjusted to end of day",
			key:      "end_date",
			query:    url.Values{"end_date": {"2024-01-15"}},
			expected: timePtr(time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC)),
		},
		{
			name:        "returns error for invalid format",
			key:         "end_date",
			query:       url.Values{"end_date": {"invalid"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseFilterDateEndOfDay(tt.query, tt.key)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.expected == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", *result)
					}
				} else {
					if result == nil {
						t.Errorf("expected %v, got nil", *tt.expected)
					} else if !result.Equal(*tt.expected) {
						t.Errorf("expected %v, got %v", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// TestParseFilterInt tests the integer filter parsing utility.
func TestParseFilterInt(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		query       url.Values
		expected    *int
		expectError bool
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "count",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "returns nil for empty parameter",
			key:      "count",
			query:    url.Values{"count": {""}},
			expected: nil,
		},
		{
			name:     "parses positive integer",
			key:      "count",
			query:    url.Values{"count": {"42"}},
			expected: intPtr(42),
		},
		{
			name:     "parses zero",
			key:      "count",
			query:    url.Values{"count": {"0"}},
			expected: intPtr(0),
		},
		{
			name:     "parses negative integer",
			key:      "count",
			query:    url.Values{"count": {"-5"}},
			expected: intPtr(-5),
		},
		{
			name:        "returns error for non-integer",
			key:         "count",
			query:       url.Values{"count": {"abc"}},
			expectError: true,
		},
		{
			name:        "returns error for float",
			key:         "count",
			query:       url.Values{"count": {"3.14"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseFilterInt(tt.query, tt.key)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.expected == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", *result)
					}
				} else {
					if result == nil {
						t.Errorf("expected %v, got nil", *tt.expected)
					} else if *result != *tt.expected {
						t.Errorf("expected %v, got %v", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// TestParseFilterFloat tests the float filter parsing utility.
func TestParseFilterFloat(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		query       url.Values
		expected    *float64
		expectError bool
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "weight",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "returns nil for empty parameter",
			key:      "weight",
			query:    url.Values{"weight": {""}},
			expected: nil,
		},
		{
			name:     "parses float",
			key:      "weight",
			query:    url.Values{"weight": {"3.14"}},
			expected: floatPtr(3.14),
		},
		{
			name:     "parses integer as float",
			key:      "weight",
			query:    url.Values{"weight": {"42"}},
			expected: floatPtr(42.0),
		},
		{
			name:     "parses negative float",
			key:      "weight",
			query:    url.Values{"weight": {"-2.5"}},
			expected: floatPtr(-2.5),
		},
		{
			name:        "returns error for non-number",
			key:         "weight",
			query:       url.Values{"weight": {"abc"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseFilterFloat(tt.query, tt.key)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.expected == nil {
					if result != nil {
						t.Errorf("expected nil, got %v", *result)
					}
				} else {
					if result == nil {
						t.Errorf("expected %v, got nil", *tt.expected)
					} else if *result != *tt.expected {
						t.Errorf("expected %v, got %v", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// TestParseFilterEnum tests the enum filter parsing utility.
func TestParseFilterEnum(t *testing.T) {
	allowedValues := []string{"ACTIVE", "INACTIVE", "PENDING"}

	tests := []struct {
		name        string
		key         string
		query       url.Values
		expected    *string
		expectError bool
	}{
		{
			name:     "returns nil for missing parameter",
			key:      "status",
			query:    url.Values{},
			expected: nil,
		},
		{
			name:     "returns nil for empty parameter",
			key:      "status",
			query:    url.Values{"status": {""}},
			expected: nil,
		},
		{
			name:     "parses valid enum value (uppercase)",
			key:      "status",
			query:    url.Values{"status": {"ACTIVE"}},
			expected: strPtr("ACTIVE"),
		},
		{
			name:     "normalizes to uppercase",
			key:      "status",
			query:    url.Values{"status": {"active"}},
			expected: strPtr("ACTIVE"),
		},
		{
			name:     "normalizes mixed case to uppercase",
			key:      "status",
			query:    url.Values{"status": {"PeNdInG"}},
			expected: strPtr("PENDING"),
		},
		{
			name:        "returns error for invalid value",
			key:         "status",
			query:       url.Values{"status": {"UNKNOWN"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseFilterEnum(tt.query, tt.key, allowedValues)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.expected == nil {
					if result != nil {
						t.Errorf("expected nil, got %q", *result)
					}
				} else {
					if result == nil {
						t.Errorf("expected %q, got nil", *tt.expected)
					} else if *result != *tt.expected {
						t.Errorf("expected %q, got %q", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}
