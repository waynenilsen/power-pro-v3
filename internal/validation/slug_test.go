package validation

import (
	"testing"
)

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name      string
		slug      string
		maxLength int
		wantErr   error
	}{
		{
			name:      "valid simple slug",
			slug:      "bench-press",
			maxLength: 50,
			wantErr:   nil,
		},
		{
			name:      "valid single word",
			slug:      "squat",
			maxLength: 50,
			wantErr:   nil,
		},
		{
			name:      "valid with numbers",
			slug:      "week-1",
			maxLength: 50,
			wantErr:   nil,
		},
		{
			name:      "valid just numbers",
			slug:      "123",
			maxLength: 50,
			wantErr:   nil,
		},
		{
			name:      "empty slug",
			slug:      "",
			maxLength: 50,
			wantErr:   ErrSlugEmpty,
		},
		{
			name:      "too long",
			slug:      "this-slug-is-way-too-long-for-the-max-length-allowed-here",
			maxLength: 50,
			wantErr:   SlugTooLongError(50),
		},
		{
			name:      "invalid uppercase",
			slug:      "Bench-Press",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
		{
			name:      "invalid underscore",
			slug:      "bench_press",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
		{
			name:      "invalid special chars",
			slug:      "bench!press",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
		{
			name:      "invalid consecutive hyphens",
			slug:      "bench--press",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
		{
			name:      "invalid leading hyphen",
			slug:      "-bench-press",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
		{
			name:      "invalid trailing hyphen",
			slug:      "bench-press-",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
		{
			name:      "just hyphen",
			slug:      "-",
			maxLength: 50,
			wantErr:   ErrSlugInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug, tt.maxLength)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidateSlug(%q) = %v, want nil", tt.slug, err)
				}
			} else {
				if err == nil {
					t.Errorf("ValidateSlug(%q) = nil, want error", tt.slug)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("ValidateSlug(%q) = %v, want %v", tt.slug, err, tt.wantErr)
				}
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "Bench Press",
			expected: "bench-press",
		},
		{
			name:     "already lowercase",
			input:    "squat",
			expected: "squat",
		},
		{
			name:     "with numbers",
			input:    "Week 1",
			expected: "week-1",
		},
		{
			name:     "underscores",
			input:    "bench_press_variation",
			expected: "bench-press-variation",
		},
		{
			name:     "multiple spaces",
			input:    "Bench  Press",
			expected: "bench-press",
		},
		{
			name:     "special characters",
			input:    "Bench (Press) & Stuff",
			expected: "bench-press-stuff",
		},
		{
			name:     "apostrophe",
			input:    "Texas Method's Variant",
			expected: "texas-methods-variant",
		},
		{
			name:     "dots",
			input:    "5.3.1",
			expected: "5-3-1",
		},
		{
			name:     "leading trailing spaces",
			input:    "  Bench Press  ",
			expected: "bench-press",
		},
		{
			name:     "slashes",
			input:    "Push/Pull/Legs",
			expected: "push-pull-legs",
		},
		{
			name:     "brackets",
			input:    "[Heavy] Day",
			expected: "heavy-day",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special chars",
			input:    "!@#$%",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSlugTooLongError(t *testing.T) {
	err := SlugTooLongError(50)
	expected := "slug must be 50 characters or less"
	if err.Error() != expected {
		t.Errorf("SlugTooLongError(50) = %q, want %q", err.Error(), expected)
	}

	err = SlugTooLongError(100)
	expected = "slug must be 100 characters or less"
	if err.Error() != expected {
		t.Errorf("SlugTooLongError(100) = %q, want %q", err.Error(), expected)
	}
}
