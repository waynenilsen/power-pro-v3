package progression

import (
	"testing"
	"time"
)

// TestNewFailureCounter_Valid tests creating a valid failure counter.
func TestNewFailureCounter_Valid(t *testing.T) {
	input := CreateFailureCounterInput{
		UserID:        "user-123",
		LiftID:        "lift-456",
		ProgressionID: "prog-789",
	}

	fc, result := NewFailureCounter(input, "fc-001")

	if !result.Valid {
		t.Fatalf("expected valid result, got errors: %v", result.Errors)
	}
	if fc == nil {
		t.Fatal("expected non-nil FailureCounter")
	}
	if fc.ID != "fc-001" {
		t.Errorf("expected ID 'fc-001', got %s", fc.ID)
	}
	if fc.UserID != input.UserID {
		t.Errorf("expected UserID %s, got %s", input.UserID, fc.UserID)
	}
	if fc.LiftID != input.LiftID {
		t.Errorf("expected LiftID %s, got %s", input.LiftID, fc.LiftID)
	}
	if fc.ProgressionID != input.ProgressionID {
		t.Errorf("expected ProgressionID %s, got %s", input.ProgressionID, fc.ProgressionID)
	}
	if fc.ConsecutiveFailures != 0 {
		t.Errorf("expected ConsecutiveFailures 0, got %d", fc.ConsecutiveFailures)
	}
	if fc.LastFailureAt != nil {
		t.Errorf("expected LastFailureAt nil, got %v", fc.LastFailureAt)
	}
	if fc.LastSuccessAt != nil {
		t.Errorf("expected LastSuccessAt nil, got %v", fc.LastSuccessAt)
	}
	if fc.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if fc.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

// TestNewFailureCounter_ValidationErrors tests validation during creation.
func TestNewFailureCounter_ValidationErrors(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateFailureCounterInput
		expectedErrs int
	}{
		{
			name: "missing user_id",
			input: CreateFailureCounterInput{
				UserID:        "",
				LiftID:        "lift-456",
				ProgressionID: "prog-789",
			},
			expectedErrs: 1,
		},
		{
			name: "missing lift_id",
			input: CreateFailureCounterInput{
				UserID:        "user-123",
				LiftID:        "",
				ProgressionID: "prog-789",
			},
			expectedErrs: 1,
		},
		{
			name: "missing progression_id",
			input: CreateFailureCounterInput{
				UserID:        "user-123",
				LiftID:        "lift-456",
				ProgressionID: "",
			},
			expectedErrs: 1,
		},
		{
			name: "all fields missing",
			input: CreateFailureCounterInput{
				UserID:        "",
				LiftID:        "",
				ProgressionID: "",
			},
			expectedErrs: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc, result := NewFailureCounter(tt.input, "fc-001")

			if result.Valid {
				t.Error("expected invalid result")
			}
			if fc != nil {
				t.Error("expected nil FailureCounter on validation error")
			}
			if len(result.Errors) != tt.expectedErrs {
				t.Errorf("expected %d errors, got %d: %v", tt.expectedErrs, len(result.Errors), result.Errors)
			}
		})
	}
}

// TestFailureCounter_Validate tests validation on existing counter.
func TestFailureCounter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		counter FailureCounter
		wantErr bool
	}{
		{
			name: "valid counter",
			counter: FailureCounter{
				ID:                  "fc-001",
				UserID:              "user-123",
				LiftID:              "lift-456",
				ProgressionID:       "prog-789",
				ConsecutiveFailures: 3,
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			counter: FailureCounter{
				ID:            "fc-001",
				UserID:        "",
				LiftID:        "lift-456",
				ProgressionID: "prog-789",
			},
			wantErr: true,
		},
		{
			name: "negative consecutive failures",
			counter: FailureCounter{
				ID:                  "fc-001",
				UserID:              "user-123",
				LiftID:              "lift-456",
				ProgressionID:       "prog-789",
				ConsecutiveFailures: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.counter.Validate()
			if tt.wantErr && result.Valid {
				t.Error("expected invalid result")
			}
			if !tt.wantErr && !result.Valid {
				t.Errorf("expected valid result, got errors: %v", result.Errors)
			}
		})
	}
}

// TestFailureCounter_IncrementFailure tests incrementing failure count.
func TestFailureCounter_IncrementFailure(t *testing.T) {
	fc := &FailureCounter{
		ID:                  "fc-001",
		UserID:              "user-123",
		LiftID:              "lift-456",
		ProgressionID:       "prog-789",
		ConsecutiveFailures: 0,
		CreatedAt:           time.Now().Add(-time.Hour),
		UpdatedAt:           time.Now().Add(-time.Hour),
	}

	beforeUpdate := fc.UpdatedAt

	// First increment
	count := fc.IncrementFailure()
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
	if fc.ConsecutiveFailures != 1 {
		t.Errorf("expected ConsecutiveFailures 1, got %d", fc.ConsecutiveFailures)
	}
	if fc.LastFailureAt == nil {
		t.Error("expected LastFailureAt to be set")
	}
	if !fc.UpdatedAt.After(beforeUpdate) {
		t.Error("expected UpdatedAt to be updated")
	}

	// Second increment
	count = fc.IncrementFailure()
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
	if fc.ConsecutiveFailures != 2 {
		t.Errorf("expected ConsecutiveFailures 2, got %d", fc.ConsecutiveFailures)
	}

	// Third increment
	count = fc.IncrementFailure()
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

// TestFailureCounter_ResetOnSuccess tests resetting failure count.
func TestFailureCounter_ResetOnSuccess(t *testing.T) {
	failureTime := time.Now().Add(-time.Hour)
	fc := &FailureCounter{
		ID:                  "fc-001",
		UserID:              "user-123",
		LiftID:              "lift-456",
		ProgressionID:       "prog-789",
		ConsecutiveFailures: 5,
		LastFailureAt:       &failureTime,
		CreatedAt:           time.Now().Add(-2 * time.Hour),
		UpdatedAt:           time.Now().Add(-time.Hour),
	}

	beforeUpdate := fc.UpdatedAt
	fc.ResetOnSuccess()

	if fc.ConsecutiveFailures != 0 {
		t.Errorf("expected ConsecutiveFailures 0, got %d", fc.ConsecutiveFailures)
	}
	if fc.LastSuccessAt == nil {
		t.Error("expected LastSuccessAt to be set")
	}
	if !fc.UpdatedAt.After(beforeUpdate) {
		t.Error("expected UpdatedAt to be updated")
	}
	// LastFailureAt should remain unchanged
	if fc.LastFailureAt == nil {
		t.Error("expected LastFailureAt to remain set")
	}
}

// TestFailureCounter_HasFailures tests failure check.
func TestFailureCounter_HasFailures(t *testing.T) {
	tests := []struct {
		name     string
		failures int
		expected bool
	}{
		{"zero failures", 0, false},
		{"one failure", 1, true},
		{"multiple failures", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FailureCounter{ConsecutiveFailures: tt.failures}
			if got := fc.HasFailures(); got != tt.expected {
				t.Errorf("HasFailures() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestFailureCounter_MeetsThreshold tests threshold checking.
func TestFailureCounter_MeetsThreshold(t *testing.T) {
	tests := []struct {
		name      string
		failures  int
		threshold int
		expected  bool
	}{
		{"zero failures, threshold 3", 0, 3, false},
		{"one failure, threshold 3", 1, 3, false},
		{"two failures, threshold 3", 2, 3, false},
		{"three failures, threshold 3", 3, 3, true},
		{"four failures, threshold 3", 4, 3, true},
		{"two failures, threshold 2", 2, 2, true},
		{"one failure, threshold 1", 1, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FailureCounter{ConsecutiveFailures: tt.failures}
			if got := fc.MeetsThreshold(tt.threshold); got != tt.expected {
				t.Errorf("MeetsThreshold(%d) = %v, want %v", tt.threshold, got, tt.expected)
			}
		})
	}
}

// TestValidation_Functions tests individual validation functions.
func TestValidation_Functions(t *testing.T) {
	t.Run("ValidateFailureCounterUserID", func(t *testing.T) {
		if err := ValidateFailureCounterUserID("user-123"); err != nil {
			t.Errorf("unexpected error for valid user ID: %v", err)
		}
		if err := ValidateFailureCounterUserID(""); err == nil {
			t.Error("expected error for empty user ID")
		}
	})

	t.Run("ValidateFailureCounterLiftID", func(t *testing.T) {
		if err := ValidateFailureCounterLiftID("lift-123"); err != nil {
			t.Errorf("unexpected error for valid lift ID: %v", err)
		}
		if err := ValidateFailureCounterLiftID(""); err == nil {
			t.Error("expected error for empty lift ID")
		}
	})

	t.Run("ValidateFailureCounterProgressionID", func(t *testing.T) {
		if err := ValidateFailureCounterProgressionID("prog-123"); err != nil {
			t.Errorf("unexpected error for valid progression ID: %v", err)
		}
		if err := ValidateFailureCounterProgressionID(""); err == nil {
			t.Error("expected error for empty progression ID")
		}
	})

	t.Run("ValidateConsecutiveFailures", func(t *testing.T) {
		if err := ValidateConsecutiveFailures(0); err != nil {
			t.Errorf("unexpected error for 0: %v", err)
		}
		if err := ValidateConsecutiveFailures(5); err != nil {
			t.Errorf("unexpected error for 5: %v", err)
		}
		if err := ValidateConsecutiveFailures(-1); err == nil {
			t.Error("expected error for negative value")
		}
	})
}
