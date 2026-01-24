// Package rpechart provides domain logic for the RPE (Rate of Perceived Exertion) Chart.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
//
// The RPE Chart maps (reps, RPE) combinations to a percentage of 1RM.
// This is the core lookup table used by RTS (Reactive Training Systems) programs.
package rpechart

import (
	"errors"
	"fmt"
)

// Validation errors
var (
	ErrEntriesRequired    = errors.New("entries are required")
	ErrRepsInvalid        = errors.New("reps must be between 1 and 12")
	ErrRPEInvalid         = errors.New("RPE must be between 7.0 and 10.0")
	ErrPercentageInvalid  = errors.New("percentage must be between 0.0 and 1.0")
	ErrDuplicateEntry     = errors.New("duplicate (reps, RPE) combination in entries")
	ErrEntryNotFound      = errors.New("no entry found for the given reps and RPE")
)

// Valid RPE values (7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0)
var validRPEValues = []float64{7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0}

// RPEChartEntry represents an entry in the RPE chart.
// It maps a (reps, RPE) combination to a percentage of 1RM.
type RPEChartEntry struct {
	TargetReps int     `json:"targetReps"` // 1-12
	TargetRPE  float64 `json:"targetRpe"`  // 7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0
	Percentage float64 `json:"percentage"` // 0.0-1.0 (e.g., 0.82 for 82%)
}

// RPEChart represents an RPE chart domain entity.
// It provides lookup functionality for converting (reps, RPE) to percentage of 1RM.
type RPEChart struct {
	Entries []RPEChartEntry
}

// isValidRPE checks if the given RPE value is valid (7.0-10.0 in 0.5 increments).
func isValidRPE(rpe float64) bool {
	for _, valid := range validRPEValues {
		if rpe == valid {
			return true
		}
	}
	return false
}

// ValidateEntry validates a single RPE chart entry.
// Returns an error if validation fails, nil otherwise.
func ValidateEntry(entry RPEChartEntry) error {
	if entry.TargetReps < 1 || entry.TargetReps > 12 {
		return ErrRepsInvalid
	}
	if !isValidRPE(entry.TargetRPE) {
		return ErrRPEInvalid
	}
	if entry.Percentage < 0.0 || entry.Percentage > 1.0 {
		return ErrPercentageInvalid
	}
	return nil
}

// ValidateEntries validates a slice of RPE chart entries.
// Returns an error if validation fails, nil otherwise.
func ValidateEntries(entries []RPEChartEntry) error {
	if len(entries) == 0 {
		return ErrEntriesRequired
	}

	// Track seen (reps, RPE) combinations for duplicate detection
	seen := make(map[string]bool)
	for _, entry := range entries {
		if err := ValidateEntry(entry); err != nil {
			return err
		}
		key := fmt.Sprintf("%d:%.1f", entry.TargetReps, entry.TargetRPE)
		if seen[key] {
			return ErrDuplicateEntry
		}
		seen[key] = true
	}
	return nil
}

// NewRPEChart creates a new RPEChart with the given entries.
// Returns an error if validation fails.
func NewRPEChart(entries []RPEChartEntry) (*RPEChart, error) {
	if err := ValidateEntries(entries); err != nil {
		return nil, err
	}
	return &RPEChart{Entries: entries}, nil
}

// Validate performs full validation on an existing RPE chart.
func (c *RPEChart) Validate() error {
	return ValidateEntries(c.Entries)
}

// GetPercentage returns the percentage of 1RM for the given reps and RPE.
// Returns an error if no matching entry is found.
func (c *RPEChart) GetPercentage(reps int, rpe float64) (float64, error) {
	for _, entry := range c.Entries {
		if entry.TargetReps == reps && entry.TargetRPE == rpe {
			return entry.Percentage, nil
		}
	}
	return 0, ErrEntryNotFound
}

// NewDefaultRPEChart creates the standard RTS RPE chart.
// This chart is based on the Reactive Training Systems methodology.
func NewDefaultRPEChart() *RPEChart {
	// Standard RTS RPE Chart (from RTS Intermediate documentation)
	// Rows: RPE 7, 7.5, 8, 8.5, 9, 9.5, 10
	// Columns: Reps 1-12
	entries := []RPEChartEntry{
		// RPE 7
		{TargetReps: 1, TargetRPE: 7.0, Percentage: 0.88},
		{TargetReps: 2, TargetRPE: 7.0, Percentage: 0.82},
		{TargetReps: 3, TargetRPE: 7.0, Percentage: 0.80},
		{TargetReps: 4, TargetRPE: 7.0, Percentage: 0.74},
		{TargetReps: 5, TargetRPE: 7.0, Percentage: 0.74},
		{TargetReps: 6, TargetRPE: 7.0, Percentage: 0.68},
		{TargetReps: 7, TargetRPE: 7.0, Percentage: 0.66},
		{TargetReps: 8, TargetRPE: 7.0, Percentage: 0.64},
		{TargetReps: 9, TargetRPE: 7.0, Percentage: 0.62},
		{TargetReps: 10, TargetRPE: 7.0, Percentage: 0.60},
		{TargetReps: 11, TargetRPE: 7.0, Percentage: 0.58},
		{TargetReps: 12, TargetRPE: 7.0, Percentage: 0.56},

		// RPE 7.5 (interpolated: average of RPE 7 and RPE 8)
		{TargetReps: 1, TargetRPE: 7.5, Percentage: 0.895},
		{TargetReps: 2, TargetRPE: 7.5, Percentage: 0.85},
		{TargetReps: 3, TargetRPE: 7.5, Percentage: 0.81},
		{TargetReps: 4, TargetRPE: 7.5, Percentage: 0.77},
		{TargetReps: 5, TargetRPE: 7.5, Percentage: 0.755},
		{TargetReps: 6, TargetRPE: 7.5, Percentage: 0.695},
		{TargetReps: 7, TargetRPE: 7.5, Percentage: 0.67},
		{TargetReps: 8, TargetRPE: 7.5, Percentage: 0.65},
		{TargetReps: 9, TargetRPE: 7.5, Percentage: 0.63},
		{TargetReps: 10, TargetRPE: 7.5, Percentage: 0.61},
		{TargetReps: 11, TargetRPE: 7.5, Percentage: 0.59},
		{TargetReps: 12, TargetRPE: 7.5, Percentage: 0.57},

		// RPE 8
		{TargetReps: 1, TargetRPE: 8.0, Percentage: 0.91},
		{TargetReps: 2, TargetRPE: 8.0, Percentage: 0.88},
		{TargetReps: 3, TargetRPE: 8.0, Percentage: 0.82},
		{TargetReps: 4, TargetRPE: 8.0, Percentage: 0.80},
		{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.77},
		{TargetReps: 6, TargetRPE: 8.0, Percentage: 0.71},
		{TargetReps: 7, TargetRPE: 8.0, Percentage: 0.68},
		{TargetReps: 8, TargetRPE: 8.0, Percentage: 0.66},
		{TargetReps: 9, TargetRPE: 8.0, Percentage: 0.64},
		{TargetReps: 10, TargetRPE: 8.0, Percentage: 0.62},
		{TargetReps: 11, TargetRPE: 8.0, Percentage: 0.60},
		{TargetReps: 12, TargetRPE: 8.0, Percentage: 0.58},

		// RPE 8.5 (interpolated: average of RPE 8 and RPE 9)
		{TargetReps: 1, TargetRPE: 8.5, Percentage: 0.93},
		{TargetReps: 2, TargetRPE: 8.5, Percentage: 0.895},
		{TargetReps: 3, TargetRPE: 8.5, Percentage: 0.855},
		{TargetReps: 4, TargetRPE: 8.5, Percentage: 0.81},
		{TargetReps: 5, TargetRPE: 8.5, Percentage: 0.785},
		{TargetReps: 6, TargetRPE: 8.5, Percentage: 0.725},
		{TargetReps: 7, TargetRPE: 8.5, Percentage: 0.695},
		{TargetReps: 8, TargetRPE: 8.5, Percentage: 0.67},
		{TargetReps: 9, TargetRPE: 8.5, Percentage: 0.65},
		{TargetReps: 10, TargetRPE: 8.5, Percentage: 0.63},
		{TargetReps: 11, TargetRPE: 8.5, Percentage: 0.61},
		{TargetReps: 12, TargetRPE: 8.5, Percentage: 0.59},

		// RPE 9
		{TargetReps: 1, TargetRPE: 9.0, Percentage: 0.95},
		{TargetReps: 2, TargetRPE: 9.0, Percentage: 0.91},
		{TargetReps: 3, TargetRPE: 9.0, Percentage: 0.89},
		{TargetReps: 4, TargetRPE: 9.0, Percentage: 0.82},
		{TargetReps: 5, TargetRPE: 9.0, Percentage: 0.80},
		{TargetReps: 6, TargetRPE: 9.0, Percentage: 0.74},
		{TargetReps: 7, TargetRPE: 9.0, Percentage: 0.71},
		{TargetReps: 8, TargetRPE: 9.0, Percentage: 0.68},
		{TargetReps: 9, TargetRPE: 9.0, Percentage: 0.66},
		{TargetReps: 10, TargetRPE: 9.0, Percentage: 0.64},
		{TargetReps: 11, TargetRPE: 9.0, Percentage: 0.62},
		{TargetReps: 12, TargetRPE: 9.0, Percentage: 0.60},

		// RPE 9.5 (interpolated: average of RPE 9 and RPE 10)
		{TargetReps: 1, TargetRPE: 9.5, Percentage: 0.975},
		{TargetReps: 2, TargetRPE: 9.5, Percentage: 0.93},
		{TargetReps: 3, TargetRPE: 9.5, Percentage: 0.905},
		{TargetReps: 4, TargetRPE: 9.5, Percentage: 0.85},
		{TargetReps: 5, TargetRPE: 9.5, Percentage: 0.81},
		{TargetReps: 6, TargetRPE: 9.5, Percentage: 0.77},
		{TargetReps: 7, TargetRPE: 9.5, Percentage: 0.725},
		{TargetReps: 8, TargetRPE: 9.5, Percentage: 0.695},
		{TargetReps: 9, TargetRPE: 9.5, Percentage: 0.67},
		{TargetReps: 10, TargetRPE: 9.5, Percentage: 0.65},
		{TargetReps: 11, TargetRPE: 9.5, Percentage: 0.63},
		{TargetReps: 12, TargetRPE: 9.5, Percentage: 0.61},

		// RPE 10
		{TargetReps: 1, TargetRPE: 10.0, Percentage: 1.00},
		{TargetReps: 2, TargetRPE: 10.0, Percentage: 0.95},
		{TargetReps: 3, TargetRPE: 10.0, Percentage: 0.92},
		{TargetReps: 4, TargetRPE: 10.0, Percentage: 0.88},
		{TargetReps: 5, TargetRPE: 10.0, Percentage: 0.82},
		{TargetReps: 6, TargetRPE: 10.0, Percentage: 0.80},
		{TargetReps: 7, TargetRPE: 10.0, Percentage: 0.74},
		{TargetReps: 8, TargetRPE: 10.0, Percentage: 0.71},
		{TargetReps: 9, TargetRPE: 10.0, Percentage: 0.68},
		{TargetReps: 10, TargetRPE: 10.0, Percentage: 0.66},
		{TargetReps: 11, TargetRPE: 10.0, Percentage: 0.64},
		{TargetReps: 12, TargetRPE: 10.0, Percentage: 0.62},
	}

	// This is the default chart - no validation needed as we control the data
	return &RPEChart{Entries: entries}
}
