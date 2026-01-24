package e1rm

import (
	"errors"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/domain/rpechart"
)

func TestNewCalculator(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	if calc == nil {
		t.Fatal("expected calculator to be created, got nil")
	}
	if calc.rpeChart != chart {
		t.Error("expected calculator to have the provided RPE chart")
	}
}

func TestCalculator_Calculate_BasicCalculation(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	// Test case from README: 315 × 5 @ RPE 8 → 315 / 0.77 = 408.4 → rounded to 407.5
	// RPE chart shows 5 reps @ RPE 8 = 0.77 (77%)
	e1rm, err := calc.Calculate(315, 5, 8.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 315 / 0.77 = 409.09... → rounds to 410.0 with 2.5 increment
	expected := 410.0
	if e1rm != expected {
		t.Errorf("expected E1RM of %.1f, got %.1f", expected, e1rm)
	}
}

func TestCalculator_Calculate_OneRepAtRPE10(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	// 1 rep @ RPE 10 means 100% of 1RM, so E1RM equals weight
	e1rm, err := calc.Calculate(400, 1, 10.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 400 / 1.0 = 400 → rounds to 400 (already on 2.5 increment)
	expected := 400.0
	if e1rm != expected {
		t.Errorf("expected E1RM of %.1f, got %.1f", expected, e1rm)
	}
}

func TestCalculator_Calculate_RoundingTo2_5(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	// Test that result is properly rounded to 2.5 increments
	// 200 × 3 @ RPE 9 → 200 / 0.89 = 224.72... → rounds to 225.0
	e1rm, err := calc.Calculate(200, 3, 9.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 225.0
	if e1rm != expected {
		t.Errorf("expected E1RM of %.1f, got %.1f", expected, e1rm)
	}
}

func TestCalculator_Calculate_InvalidWeight(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	tests := []struct {
		name   string
		weight float64
	}{
		{"zero weight", 0},
		{"negative weight", -100},
		{"very negative weight", -0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Calculate(tt.weight, 5, 8.0)
			if err == nil {
				t.Error("expected error for invalid weight, got nil")
			}
			if !errors.Is(err, ErrWeightMustBePositive) {
				t.Errorf("expected ErrWeightMustBePositive, got: %v", err)
			}
		})
	}
}

func TestCalculator_Calculate_InvalidReps(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	tests := []struct {
		name string
		reps int
	}{
		{"zero reps", 0},
		{"negative reps", -1},
		{"reps too high", 13},
		{"reps way too high", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Calculate(315, tt.reps, 8.0)
			if err == nil {
				t.Error("expected error for invalid reps, got nil")
			}
			if !errors.Is(err, ErrRepsOutOfRange) {
				t.Errorf("expected ErrRepsOutOfRange, got: %v", err)
			}
		})
	}
}

func TestCalculator_Calculate_InvalidRPE(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	tests := []struct {
		name string
		rpe  float64
	}{
		{"RPE too low", 6.5},
		{"RPE way too low", 5.0},
		{"RPE too high", 10.5},
		{"negative RPE", -8.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Calculate(315, 5, tt.rpe)
			if err == nil {
				t.Error("expected error for invalid RPE, got nil")
			}
			if !errors.Is(err, ErrRPEOutOfRange) {
				t.Errorf("expected ErrRPEOutOfRange, got: %v", err)
			}
		})
	}
}

func TestCalculator_Calculate_RPEChartLookupFailure(t *testing.T) {
	// Create a minimal chart that doesn't have all entries
	entries := []rpechart.RPEChartEntry{
		{TargetReps: 1, TargetRPE: 10.0, Percentage: 1.0},
	}
	chart, err := rpechart.NewRPEChart(entries)
	if err != nil {
		t.Fatalf("failed to create chart: %v", err)
	}
	calc := NewCalculator(chart)

	// Try to calculate with a (reps, RPE) combination that doesn't exist in the chart
	// 5 reps @ RPE 8 won't be found in this minimal chart
	_, err = calc.Calculate(315, 5, 8.0)
	if err == nil {
		t.Error("expected error for RPE chart lookup failure, got nil")
	}

	// Should contain the underlying rpechart error
	if !errors.Is(err, rpechart.ErrEntryNotFound) {
		t.Errorf("expected error to wrap ErrEntryNotFound, got: %v", err)
	}
}

func TestCalculator_Calculate_VariousRPEValues(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	// Test various valid RPE values including half values
	tests := []struct {
		name   string
		weight float64
		reps   int
		rpe    float64
	}{
		{"RPE 7.0", 200, 5, 7.0},
		{"RPE 7.5", 200, 5, 7.5},
		{"RPE 8.0", 200, 5, 8.0},
		{"RPE 8.5", 200, 5, 8.5},
		{"RPE 9.0", 200, 5, 9.0},
		{"RPE 9.5", 200, 5, 9.5},
		{"RPE 10.0", 200, 5, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e1rm, err := calc.Calculate(tt.weight, tt.reps, tt.rpe)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// E1RM should always be >= the weight lifted
			if e1rm < tt.weight {
				t.Errorf("E1RM (%.1f) should be >= weight lifted (%.1f)", e1rm, tt.weight)
			}
		})
	}
}

func TestCalculator_Calculate_EdgeReps(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calc := NewCalculator(chart)

	// Test edge cases for reps (1 and 12)
	tests := []struct {
		name string
		reps int
	}{
		{"1 rep (minimum)", 1},
		{"12 reps (maximum)", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e1rm, err := calc.Calculate(200, tt.reps, 8.0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if e1rm <= 0 {
				t.Errorf("E1RM should be positive, got %.1f", e1rm)
			}
		})
	}
}
