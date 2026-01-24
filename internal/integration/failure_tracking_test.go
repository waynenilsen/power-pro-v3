// Package integration provides integration tests for cross-component behavior.
// This file tests the failure tracking system that powers ON_FAILURE progressions.
package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// TestFailureCounterIncrement tests that failure counters increment correctly.
func TestFailureCounterIncrement(t *testing.T) {
	t.Run("first failure sets counter to 1", func(t *testing.T) {
		counter, result := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        uuid.New().String(),
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		if !result.Valid {
			t.Fatalf("Failed to create failure counter: %v", result.Errors)
		}

		if counter.ConsecutiveFailures != 0 {
			t.Errorf("Initial counter should be 0, got %d", counter.ConsecutiveFailures)
		}

		count := counter.IncrementFailure()
		if count != 1 {
			t.Errorf("After first increment, expected 1, got %d", count)
		}
		if counter.LastFailureAt == nil {
			t.Error("LastFailureAt should be set after increment")
		}
	})

	t.Run("consecutive failures accumulate", func(t *testing.T) {
		counter, _ := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        uuid.New().String(),
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		for i := 1; i <= 5; i++ {
			count := counter.IncrementFailure()
			if count != i {
				t.Errorf("After increment %d, expected %d, got %d", i, i, count)
			}
		}

		if counter.ConsecutiveFailures != 5 {
			t.Errorf("Final counter should be 5, got %d", counter.ConsecutiveFailures)
		}
	})
}

// TestFailureCounterReset tests that success resets the counter.
func TestFailureCounterReset(t *testing.T) {
	t.Run("success resets counter to 0", func(t *testing.T) {
		counter, _ := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        uuid.New().String(),
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		// Accumulate failures
		counter.IncrementFailure()
		counter.IncrementFailure()
		counter.IncrementFailure()

		if counter.ConsecutiveFailures != 3 {
			t.Fatalf("Expected 3 failures, got %d", counter.ConsecutiveFailures)
		}

		// Reset on success
		counter.ResetOnSuccess()

		if counter.ConsecutiveFailures != 0 {
			t.Errorf("After success, counter should be 0, got %d", counter.ConsecutiveFailures)
		}
		if counter.LastSuccessAt == nil {
			t.Error("LastSuccessAt should be set after success")
		}
	})

	t.Run("reset updates timestamps correctly", func(t *testing.T) {
		counter, _ := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        uuid.New().String(),
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		createdAt := counter.CreatedAt
		time.Sleep(10 * time.Millisecond)

		counter.IncrementFailure()
		failureTime := *counter.LastFailureAt
		updatedAfterFailure := counter.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		counter.ResetOnSuccess()
		successTime := *counter.LastSuccessAt
		updatedAfterSuccess := counter.UpdatedAt

		if !failureTime.After(createdAt) {
			t.Error("Failure time should be after created time")
		}
		if !successTime.After(failureTime) {
			t.Error("Success time should be after failure time")
		}
		if !updatedAfterSuccess.After(updatedAfterFailure) {
			t.Error("Updated time should increase after each operation")
		}
	})
}

// TestFailureCounterThreshold tests threshold checking.
func TestFailureCounterThreshold(t *testing.T) {
	tests := []struct {
		name       string
		failures   int
		threshold  int
		meetsMeets bool
	}{
		{"0 failures, threshold 1", 0, 1, false},
		{"1 failure, threshold 1", 1, 1, true},
		{"1 failure, threshold 2", 1, 2, false},
		{"2 failures, threshold 2", 2, 2, true},
		{"3 failures, threshold 2", 3, 2, true},
		{"5 failures, threshold 3", 5, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter, _ := progression.NewFailureCounter(progression.CreateFailureCounterInput{
				UserID:        uuid.New().String(),
				LiftID:        uuid.New().String(),
				ProgressionID: uuid.New().String(),
			}, uuid.New().String())

			for i := 0; i < tt.failures; i++ {
				counter.IncrementFailure()
			}

			if counter.MeetsThreshold(tt.threshold) != tt.meetsMeets {
				t.Errorf("MeetsThreshold(%d) = %v, want %v",
					tt.threshold, counter.MeetsThreshold(tt.threshold), tt.meetsMeets)
			}
		})
	}
}

// TestFailureCounterHasFailures tests the HasFailures helper.
func TestFailureCounterHasFailures(t *testing.T) {
	counter, _ := progression.NewFailureCounter(progression.CreateFailureCounterInput{
		UserID:        uuid.New().String(),
		LiftID:        uuid.New().String(),
		ProgressionID: uuid.New().String(),
	}, uuid.New().String())

	if counter.HasFailures() {
		t.Error("New counter should not have failures")
	}

	counter.IncrementFailure()
	if !counter.HasFailures() {
		t.Error("Counter with 1 failure should have failures")
	}

	counter.ResetOnSuccess()
	if counter.HasFailures() {
		t.Error("Reset counter should not have failures")
	}
}

// TestFailureCounterValidation tests input validation.
func TestFailureCounterValidation(t *testing.T) {
	t.Run("requires user ID", func(t *testing.T) {
		_, result := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        "",
			LiftID:        uuid.New().String(),
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		if result.Valid {
			t.Error("Expected validation to fail for empty user ID")
		}
	})

	t.Run("requires lift ID", func(t *testing.T) {
		_, result := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        "",
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		if result.Valid {
			t.Error("Expected validation to fail for empty lift ID")
		}
	})

	t.Run("requires progression ID", func(t *testing.T) {
		_, result := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        uuid.New().String(),
			ProgressionID: "",
		}, uuid.New().String())

		if result.Valid {
			t.Error("Expected validation to fail for empty progression ID")
		}
	})

	t.Run("valid input creates counter", func(t *testing.T) {
		counter, result := progression.NewFailureCounter(progression.CreateFailureCounterInput{
			UserID:        uuid.New().String(),
			LiftID:        uuid.New().String(),
			ProgressionID: uuid.New().String(),
		}, uuid.New().String())

		if !result.Valid {
			t.Errorf("Expected valid counter, got errors: %v", result.Errors)
		}
		if counter == nil {
			t.Error("Expected counter to be created")
		}
	})
}
