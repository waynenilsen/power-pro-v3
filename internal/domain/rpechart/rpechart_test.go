package rpechart

import (
	"errors"
	"testing"
)

// ==================== Entry Validation Tests ====================

func TestValidateEntry_Valid(t *testing.T) {
	tests := []struct {
		name  string
		entry RPEChartEntry
	}{
		{
			name:  "minimum valid values",
			entry: RPEChartEntry{TargetReps: 1, TargetRPE: 7.0, Percentage: 0.0},
		},
		{
			name:  "maximum valid values",
			entry: RPEChartEntry{TargetReps: 12, TargetRPE: 10.0, Percentage: 1.0},
		},
		{
			name:  "mid-range values",
			entry: RPEChartEntry{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
		},
		{
			name:  "half RPE value 7.5",
			entry: RPEChartEntry{TargetReps: 3, TargetRPE: 7.5, Percentage: 0.81},
		},
		{
			name:  "half RPE value 8.5",
			entry: RPEChartEntry{TargetReps: 5, TargetRPE: 8.5, Percentage: 0.785},
		},
		{
			name:  "half RPE value 9.5",
			entry: RPEChartEntry{TargetReps: 1, TargetRPE: 9.5, Percentage: 0.975},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntry(tt.entry)
			if err != nil {
				t.Errorf("ValidateEntry() = %v, want nil", err)
			}
		})
	}
}

func TestValidateEntry_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		entry       RPEChartEntry
		expectedErr error
	}{
		{
			name:        "reps zero",
			entry:       RPEChartEntry{TargetReps: 0, TargetRPE: 8.0, Percentage: 0.5},
			expectedErr: ErrRepsInvalid,
		},
		{
			name:        "reps negative",
			entry:       RPEChartEntry{TargetReps: -1, TargetRPE: 8.0, Percentage: 0.5},
			expectedErr: ErrRepsInvalid,
		},
		{
			name:        "reps too high",
			entry:       RPEChartEntry{TargetReps: 13, TargetRPE: 8.0, Percentage: 0.5},
			expectedErr: ErrRepsInvalid,
		},
		{
			name:        "RPE too low",
			entry:       RPEChartEntry{TargetReps: 5, TargetRPE: 6.5, Percentage: 0.5},
			expectedErr: ErrRPEInvalid,
		},
		{
			name:        "RPE too high",
			entry:       RPEChartEntry{TargetReps: 5, TargetRPE: 10.5, Percentage: 0.5},
			expectedErr: ErrRPEInvalid,
		},
		{
			name:        "RPE not on 0.5 increment",
			entry:       RPEChartEntry{TargetReps: 5, TargetRPE: 8.3, Percentage: 0.5},
			expectedErr: ErrRPEInvalid,
		},
		{
			name:        "percentage negative",
			entry:       RPEChartEntry{TargetReps: 5, TargetRPE: 8.0, Percentage: -0.1},
			expectedErr: ErrPercentageInvalid,
		},
		{
			name:        "percentage too high",
			entry:       RPEChartEntry{TargetReps: 5, TargetRPE: 8.0, Percentage: 1.1},
			expectedErr: ErrPercentageInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntry(tt.entry)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateEntry() = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

// ==================== Entries Validation Tests ====================

func TestValidateEntries_Valid(t *testing.T) {
	tests := []struct {
		name    string
		entries []RPEChartEntry
	}{
		{
			name: "single entry",
			entries: []RPEChartEntry{
				{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
			},
		},
		{
			name: "multiple entries different reps same RPE",
			entries: []RPEChartEntry{
				{TargetReps: 1, TargetRPE: 8.0, Percentage: 0.91},
				{TargetReps: 2, TargetRPE: 8.0, Percentage: 0.88},
				{TargetReps: 3, TargetRPE: 8.0, Percentage: 0.82},
			},
		},
		{
			name: "multiple entries same reps different RPE",
			entries: []RPEChartEntry{
				{TargetReps: 5, TargetRPE: 7.0, Percentage: 0.74},
				{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
				{TargetReps: 5, TargetRPE: 9.0, Percentage: 0.80},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntries(tt.entries)
			if err != nil {
				t.Errorf("ValidateEntries() = %v, want nil", err)
			}
		})
	}
}

func TestValidateEntries_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		entries     []RPEChartEntry
		expectedErr error
	}{
		{
			name:        "empty entries",
			entries:     []RPEChartEntry{},
			expectedErr: ErrEntriesRequired,
		},
		{
			name:        "nil entries",
			entries:     nil,
			expectedErr: ErrEntriesRequired,
		},
		{
			name: "duplicate reps and RPE combination",
			entries: []RPEChartEntry{
				{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
				{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.80},
			},
			expectedErr: ErrDuplicateEntry,
		},
		{
			name: "invalid entry in list",
			entries: []RPEChartEntry{
				{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
				{TargetReps: 0, TargetRPE: 8.0, Percentage: 0.77}, // invalid reps
			},
			expectedErr: ErrRepsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntries(tt.entries)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateEntries() = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

// ==================== NewRPEChart Tests ====================

func TestNewRPEChart_Valid(t *testing.T) {
	entries := []RPEChartEntry{
		{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
		{TargetReps: 5, TargetRPE: 9.0, Percentage: 0.80},
	}

	chart, err := NewRPEChart(entries)
	if err != nil {
		t.Errorf("NewRPEChart() error = %v, want nil", err)
	}
	if chart == nil {
		t.Fatal("NewRPEChart() returned nil chart")
	}
	if len(chart.Entries) != 2 {
		t.Errorf("len(chart.Entries) = %d, want 2", len(chart.Entries))
	}
}

func TestNewRPEChart_Invalid(t *testing.T) {
	entries := []RPEChartEntry{
		{TargetReps: 0, TargetRPE: 8.0, Percentage: 0.77}, // invalid reps
	}

	chart, err := NewRPEChart(entries)
	if err == nil {
		t.Error("NewRPEChart() expected error, got nil")
	}
	if chart != nil {
		t.Error("NewRPEChart() expected nil chart on error")
	}
}

// ==================== GetPercentage Tests ====================

func TestRPEChart_GetPercentage_Found(t *testing.T) {
	chart := NewDefaultRPEChart()

	tests := []struct {
		name       string
		reps       int
		rpe        float64
		wantPct    float64
	}{
		{"1 rep at RPE 10", 1, 10.0, 1.00},
		{"5 reps at RPE 8", 5, 8.0, 0.77},
		{"3 reps at RPE 9", 3, 9.0, 0.89},
		{"12 reps at RPE 7", 12, 7.0, 0.56},
		{"1 rep at RPE 7", 1, 7.0, 0.88},
		{"5 reps at RPE 7.5", 5, 7.5, 0.755},
		{"3 reps at RPE 8.5", 3, 8.5, 0.855},
		{"1 rep at RPE 9.5", 1, 9.5, 0.975},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pct, err := chart.GetPercentage(tt.reps, tt.rpe)
			if err != nil {
				t.Errorf("GetPercentage(%d, %.1f) error = %v", tt.reps, tt.rpe, err)
				return
			}
			if pct != tt.wantPct {
				t.Errorf("GetPercentage(%d, %.1f) = %.3f, want %.3f", tt.reps, tt.rpe, pct, tt.wantPct)
			}
		})
	}
}

func TestRPEChart_GetPercentage_NotFound(t *testing.T) {
	entries := []RPEChartEntry{
		{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
	}
	chart, _ := NewRPEChart(entries)

	tests := []struct {
		name string
		reps int
		rpe  float64
	}{
		{"wrong reps", 3, 8.0},
		{"wrong RPE", 5, 9.0},
		{"both wrong", 3, 9.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := chart.GetPercentage(tt.reps, tt.rpe)
			if !errors.Is(err, ErrEntryNotFound) {
				t.Errorf("GetPercentage(%d, %.1f) error = %v, want %v", tt.reps, tt.rpe, err, ErrEntryNotFound)
			}
		})
	}
}

// ==================== NewDefaultRPEChart Tests ====================

func TestNewDefaultRPEChart_ContainsAllExpectedEntries(t *testing.T) {
	chart := NewDefaultRPEChart()

	// 7 RPE values * 12 rep values = 84 entries
	expectedCount := 7 * 12
	if len(chart.Entries) != expectedCount {
		t.Errorf("len(chart.Entries) = %d, want %d", len(chart.Entries), expectedCount)
	}
}

func TestNewDefaultRPEChart_ValidatesSuccessfully(t *testing.T) {
	chart := NewDefaultRPEChart()

	err := chart.Validate()
	if err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestNewDefaultRPEChart_SpecificValues(t *testing.T) {
	chart := NewDefaultRPEChart()

	// Test specific known values from the RTS chart
	tests := []struct {
		reps    int
		rpe     float64
		wantPct float64
	}{
		// RPE 10 row
		{1, 10.0, 1.00},
		{2, 10.0, 0.95},
		{3, 10.0, 0.92},
		{4, 10.0, 0.88},
		{5, 10.0, 0.82},
		{6, 10.0, 0.80},
		{7, 10.0, 0.74},
		{8, 10.0, 0.71},
		{9, 10.0, 0.68},
		{10, 10.0, 0.66},
		{11, 10.0, 0.64},
		{12, 10.0, 0.62},

		// RPE 9 row
		{1, 9.0, 0.95},
		{2, 9.0, 0.91},
		{3, 9.0, 0.89},
		{4, 9.0, 0.82},
		{5, 9.0, 0.80},
		{6, 9.0, 0.74},

		// RPE 8 row
		{1, 8.0, 0.91},
		{2, 8.0, 0.88},
		{3, 8.0, 0.82},
		{4, 8.0, 0.80},
		{5, 8.0, 0.77},
		{6, 8.0, 0.71},

		// RPE 7 row
		{1, 7.0, 0.88},
		{2, 7.0, 0.82},
		{3, 7.0, 0.80},
		{4, 7.0, 0.74},
		{5, 7.0, 0.74},
		{6, 7.0, 0.68},
	}

	for _, tt := range tests {
		pct, err := chart.GetPercentage(tt.reps, tt.rpe)
		if err != nil {
			t.Errorf("GetPercentage(%d, %.1f) error = %v", tt.reps, tt.rpe, err)
			continue
		}
		if pct != tt.wantPct {
			t.Errorf("GetPercentage(%d, %.1f) = %.3f, want %.3f", tt.reps, tt.rpe, pct, tt.wantPct)
		}
	}
}

func TestNewDefaultRPEChart_HalfRPEValues(t *testing.T) {
	chart := NewDefaultRPEChart()

	// Verify half-RPE values exist
	halfRPEValues := []float64{7.5, 8.5, 9.5}
	for _, rpe := range halfRPEValues {
		for reps := 1; reps <= 12; reps++ {
			_, err := chart.GetPercentage(reps, rpe)
			if err != nil {
				t.Errorf("GetPercentage(%d, %.1f) error = %v; half-RPE entry missing", reps, rpe, err)
			}
		}
	}
}

func TestNewDefaultRPEChart_PercentagesDecreaseWithReps(t *testing.T) {
	chart := NewDefaultRPEChart()

	// For a given RPE, percentages should generally decrease as reps increase
	// (except for some quirks in the RTS chart where some values are equal)
	rpeValues := []float64{7.0, 8.0, 9.0, 10.0}
	for _, rpe := range rpeValues {
		prevPct := 2.0 // Start higher than any valid percentage
		for reps := 1; reps <= 12; reps++ {
			pct, err := chart.GetPercentage(reps, rpe)
			if err != nil {
				t.Errorf("GetPercentage(%d, %.1f) error = %v", reps, rpe, err)
				continue
			}
			// Allow equal or decreasing (some RTS chart values are equal)
			if pct > prevPct {
				t.Errorf("At RPE %.1f, percentage for %d reps (%.3f) > percentage for %d reps (%.3f)",
					rpe, reps, pct, reps-1, prevPct)
			}
			prevPct = pct
		}
	}
}

func TestNewDefaultRPEChart_PercentagesIncreaseWithRPE(t *testing.T) {
	chart := NewDefaultRPEChart()

	// For a given rep count, percentages should increase as RPE increases
	rpeValues := []float64{7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0}
	for reps := 1; reps <= 12; reps++ {
		prevPct := 0.0 // Start lower than any valid percentage
		for _, rpe := range rpeValues {
			pct, err := chart.GetPercentage(reps, rpe)
			if err != nil {
				t.Errorf("GetPercentage(%d, %.1f) error = %v", reps, rpe, err)
				continue
			}
			if pct < prevPct {
				t.Errorf("At %d reps, percentage for RPE %.1f (%.3f) < percentage for previous RPE (%.3f)",
					reps, rpe, pct, prevPct)
			}
			prevPct = pct
		}
	}
}

// ==================== Validate Tests ====================

func TestRPEChart_Validate_Valid(t *testing.T) {
	entries := []RPEChartEntry{
		{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
	}
	chart := &RPEChart{Entries: entries}

	err := chart.Validate()
	if err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestRPEChart_Validate_Invalid(t *testing.T) {
	chart := &RPEChart{Entries: []RPEChartEntry{}}

	err := chart.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty entries")
	}
}
