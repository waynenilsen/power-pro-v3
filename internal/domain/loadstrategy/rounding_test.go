package loadstrategy

import (
	"errors"
	"math"
	"testing"
)

func TestRoundWeight(t *testing.T) {
	tests := []struct {
		name      string
		weight    float64
		increment float64
		direction RoundingDirection
		expected  float64
		wantErr   error
	}{
		// NEAREST rounding tests
		{
			name:      "nearest: exact multiple",
			weight:    100.0,
			increment: 5.0,
			direction: RoundNearest,
			expected:  100.0,
		},
		{
			name:      "nearest: round up",
			weight:    102.6,
			increment: 5.0,
			direction: RoundNearest,
			expected:  105.0,
		},
		{
			name:      "nearest: round down",
			weight:    102.4,
			increment: 5.0,
			direction: RoundNearest,
			expected:  100.0,
		},
		{
			name:      "nearest: exactly at midpoint rounds up",
			weight:    102.5,
			increment: 5.0,
			direction: RoundNearest,
			expected:  105.0, // math.Round rounds .5 to nearest even, but 102.5/5=20.5 -> 21*5=105
		},
		{
			name:      "nearest: 2.5 increment",
			weight:    143.75,
			increment: 2.5,
			direction: RoundNearest,
			expected:  145.0,
		},
		{
			name:      "nearest: 2.5 increment round down",
			weight:    143.24,
			increment: 2.5,
			direction: RoundNearest,
			expected:  142.5,
		},
		{
			name:      "nearest: small increment 1.0",
			weight:    143.4,
			increment: 1.0,
			direction: RoundNearest,
			expected:  143.0,
		},
		{
			name:      "nearest: 10.0 increment",
			weight:    147.0,
			increment: 10.0,
			direction: RoundNearest,
			expected:  150.0,
		},

		// DOWN rounding tests
		{
			name:      "down: exact multiple",
			weight:    100.0,
			increment: 5.0,
			direction: RoundDown,
			expected:  100.0,
		},
		{
			name:      "down: basic floor",
			weight:    107.0,
			increment: 5.0,
			direction: RoundDown,
			expected:  105.0,
		},
		{
			name:      "down: just above multiple",
			weight:    100.01,
			increment: 5.0,
			direction: RoundDown,
			expected:  100.0,
		},
		{
			name:      "down: just below next multiple",
			weight:    104.99,
			increment: 5.0,
			direction: RoundDown,
			expected:  100.0,
		},
		{
			name:      "down: 2.5 increment",
			weight:    147.4,
			increment: 2.5,
			direction: RoundDown,
			expected:  145.0,
		},

		// UP rounding tests
		{
			name:      "up: exact multiple",
			weight:    100.0,
			increment: 5.0,
			direction: RoundUp,
			expected:  100.0,
		},
		{
			name:      "up: basic ceil",
			weight:    101.0,
			increment: 5.0,
			direction: RoundUp,
			expected:  105.0,
		},
		{
			name:      "up: just above multiple",
			weight:    100.01,
			increment: 5.0,
			direction: RoundUp,
			expected:  105.0,
		},
		{
			name:      "up: just below next multiple",
			weight:    104.99,
			increment: 5.0,
			direction: RoundUp,
			expected:  105.0,
		},
		{
			name:      "up: 2.5 increment",
			weight:    142.6,
			increment: 2.5,
			direction: RoundUp,
			expected:  145.0,
		},

		// Edge cases
		{
			name:      "zero weight",
			weight:    0.0,
			increment: 5.0,
			direction: RoundNearest,
			expected:  0.0,
		},
		{
			name:      "very small increment 0.5",
			weight:    142.8,
			increment: 0.5,
			direction: RoundNearest,
			expected:  143.0,
		},
		{
			name:      "large weight",
			weight:    567.3,
			increment: 5.0,
			direction: RoundNearest,
			expected:  565.0,
		},

		// Error cases
		{
			name:      "negative weight",
			weight:    -100.0,
			increment: 5.0,
			direction: RoundNearest,
			wantErr:   ErrNegativeWeight,
		},
		{
			name:      "zero increment",
			weight:    100.0,
			increment: 0.0,
			direction: RoundNearest,
			wantErr:   ErrInvalidIncrement,
		},
		{
			name:      "negative increment",
			weight:    100.0,
			increment: -5.0,
			direction: RoundNearest,
			wantErr:   ErrInvalidIncrement,
		},
		{
			name:      "invalid direction",
			weight:    100.0,
			increment: 5.0,
			direction: "INVALID",
			wantErr:   ErrInvalidRoundingDirection,
		},
		{
			name:      "empty direction",
			weight:    100.0,
			increment: 5.0,
			direction: "",
			wantErr:   ErrInvalidRoundingDirection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RoundWeight(tt.weight, tt.increment, tt.direction)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

func TestRoundWeightNearest(t *testing.T) {
	result, err := RoundWeightNearest(142.5, 5.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 145.0 {
		t.Errorf("expected 145.0, got %f", result)
	}
}

func TestRoundWeightDown(t *testing.T) {
	result, err := RoundWeightDown(147.9, 5.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 145.0 {
		t.Errorf("expected 145.0, got %f", result)
	}
}

func TestRoundWeightUp(t *testing.T) {
	result, err := RoundWeightUp(142.1, 5.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 145.0 {
		t.Errorf("expected 145.0, got %f", result)
	}
}

func TestValidateRoundingDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction RoundingDirection
		wantErr   bool
	}{
		{"valid NEAREST", RoundNearest, false},
		{"valid DOWN", RoundDown, false},
		{"valid UP", RoundUp, false},
		{"empty is allowed", "", false},
		{"invalid direction", "INVALID", true},
		{"lowercase", "nearest", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRoundingDirection(tt.direction)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrInvalidRoundingDirection) {
					t.Errorf("expected ErrInvalidRoundingDirection, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNormalizeRoundingDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction RoundingDirection
		expected  RoundingDirection
	}{
		{"empty returns default", "", RoundNearest},
		{"NEAREST unchanged", RoundNearest, RoundNearest},
		{"DOWN unchanged", RoundDown, RoundDown},
		{"UP unchanged", RoundUp, RoundUp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRoundingDirection(tt.direction)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNormalizeRoundingIncrement(t *testing.T) {
	tests := []struct {
		name      string
		increment float64
		expected  float64
	}{
		{"zero returns default", 0, DefaultRoundingIncrement},
		{"negative returns default", -5.0, DefaultRoundingIncrement},
		{"positive unchanged", 2.5, 2.5},
		{"large value unchanged", 10.0, 10.0},
		{"small value unchanged", 0.5, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRoundingIncrement(tt.increment)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestValidRoundingDirections(t *testing.T) {
	expectedDirections := []RoundingDirection{
		RoundNearest,
		RoundDown,
		RoundUp,
	}

	for _, dir := range expectedDirections {
		if !ValidRoundingDirections[dir] {
			t.Errorf("expected %s to be in ValidRoundingDirections", dir)
		}
	}

	if len(ValidRoundingDirections) != len(expectedDirections) {
		t.Errorf("expected %d directions in ValidRoundingDirections, got %d",
			len(expectedDirections), len(ValidRoundingDirections))
	}
}

func TestRoundingConstants(t *testing.T) {
	if DefaultRoundingIncrement != 5.0 {
		t.Errorf("expected DefaultRoundingIncrement to be 5.0, got %f", DefaultRoundingIncrement)
	}

	if DefaultRoundingDirection != RoundNearest {
		t.Errorf("expected DefaultRoundingDirection to be NEAREST, got %s", DefaultRoundingDirection)
	}
}

// TestRoundingPrecision tests that rounding handles floating-point precision correctly.
func TestRoundingPrecision(t *testing.T) {
	tests := []struct {
		name      string
		weight    float64
		increment float64
		direction RoundingDirection
		expected  float64
	}{
		{
			name:      "handles 0.1 + 0.2 type precision issues",
			weight:    0.1 + 0.2, // 0.30000000000000004
			increment: 0.1,
			direction: RoundNearest,
			expected:  0.3,
		},
		{
			name:      "handles small fractions",
			weight:    267.749999999,
			increment: 2.5,
			direction: RoundNearest,
			expected:  267.5,
		},
		{
			name:      "handles large numbers",
			weight:    1000.123,
			increment: 5.0,
			direction: RoundNearest,
			expected:  1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RoundWeight(tt.weight, tt.increment, tt.direction)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Use a small epsilon for floating-point comparison
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %.10f, got %.10f", tt.expected, result)
			}
		})
	}
}

// Benchmark tests
func BenchmarkRoundWeight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RoundWeight(267.75, 5.0, RoundNearest)
	}
}

func BenchmarkRoundWeightDown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RoundWeightDown(267.75, 5.0)
	}
}

func BenchmarkRoundWeightUp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RoundWeightUp(267.75, 5.0)
	}
}
