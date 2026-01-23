package service

import (
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// TestFailureService_CheckForFailure tests failure detection logic.
func TestFailureService_CheckForFailure(t *testing.T) {
	svc := &FailureService{}

	tests := []struct {
		name          string
		loggedSet     *loggedset.LoggedSet
		expectedFail  bool
	}{
		{
			name: "failure - reps less than target",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-1",
				TargetReps:    5,
				RepsPerformed: 3,
			},
			expectedFail: true,
		},
		{
			name: "success - reps equal target",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-2",
				TargetReps:    5,
				RepsPerformed: 5,
			},
			expectedFail: false,
		},
		{
			name: "success - reps exceed target",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-3",
				TargetReps:    5,
				RepsPerformed: 8,
			},
			expectedFail: false,
		},
		{
			name: "failure - zero reps",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-4",
				TargetReps:    5,
				RepsPerformed: 0,
			},
			expectedFail: true,
		},
		{
			name: "failure - one less rep",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-5",
				TargetReps:    5,
				RepsPerformed: 4,
			},
			expectedFail: true,
		},
		{
			name:          "nil logged set",
			loggedSet:     nil,
			expectedFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.CheckForFailure(tt.loggedSet)
			if result != tt.expectedFail {
				t.Errorf("CheckForFailure() = %v, want %v", result, tt.expectedFail)
			}
		})
	}
}

// TestFailureService_IsSuccess tests success detection logic.
func TestFailureService_IsSuccess(t *testing.T) {
	svc := &FailureService{}

	tests := []struct {
		name           string
		loggedSet      *loggedset.LoggedSet
		expectedSuccess bool
	}{
		{
			name: "success - reps equal target",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-1",
				TargetReps:    5,
				RepsPerformed: 5,
			},
			expectedSuccess: true,
		},
		{
			name: "success - reps exceed target",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-2",
				TargetReps:    5,
				RepsPerformed: 8,
			},
			expectedSuccess: true,
		},
		{
			name: "failure - reps less than target",
			loggedSet: &loggedset.LoggedSet{
				ID:            "set-3",
				TargetReps:    5,
				RepsPerformed: 3,
			},
			expectedSuccess: false,
		},
		{
			name:            "nil logged set",
			loggedSet:       nil,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.IsSuccess(tt.loggedSet)
			if result != tt.expectedSuccess {
				t.Errorf("IsSuccess() = %v, want %v", result, tt.expectedSuccess)
			}
		})
	}
}

// TestFailureService_BuildFailureTriggerContext tests trigger context construction.
func TestFailureService_BuildFailureTriggerContext(t *testing.T) {
	svc := &FailureService{}

	ls := &loggedset.LoggedSet{
		ID:            "set-123",
		UserID:        "user-456",
		LiftID:        "lift-789",
		TargetReps:    5,
		RepsPerformed: 3,
		Weight:        225.0,
		CreatedAt:     time.Now(),
	}

	ctx := svc.BuildFailureTriggerContext(ls, 2, "prog-abc")

	if ctx.LoggedSetID != ls.ID {
		t.Errorf("LoggedSetID = %s, want %s", ctx.LoggedSetID, ls.ID)
	}
	if ctx.LiftID != ls.LiftID {
		t.Errorf("LiftID = %s, want %s", ctx.LiftID, ls.LiftID)
	}
	if ctx.TargetReps != ls.TargetReps {
		t.Errorf("TargetReps = %d, want %d", ctx.TargetReps, ls.TargetReps)
	}
	if ctx.RepsPerformed != ls.RepsPerformed {
		t.Errorf("RepsPerformed = %d, want %d", ctx.RepsPerformed, ls.RepsPerformed)
	}
	if ctx.RepsDifference != ls.RepsDifference() {
		t.Errorf("RepsDifference = %d, want %d", ctx.RepsDifference, ls.RepsDifference())
	}
	if ctx.ConsecutiveFailures != 2 {
		t.Errorf("ConsecutiveFailures = %d, want 2", ctx.ConsecutiveFailures)
	}
	if ctx.Weight != ls.Weight {
		t.Errorf("Weight = %f, want %f", ctx.Weight, ls.Weight)
	}
	if ctx.ProgressionID != "prog-abc" {
		t.Errorf("ProgressionID = %s, want prog-abc", ctx.ProgressionID)
	}
}

// TestFailureService_CreateFailureTriggerEvent tests trigger event construction.
func TestFailureService_CreateFailureTriggerEvent(t *testing.T) {
	svc := &FailureService{}

	ls := &loggedset.LoggedSet{
		ID:            "set-123",
		UserID:        "user-456",
		LiftID:        "lift-789",
		TargetReps:    5,
		RepsPerformed: 3,
		Weight:        225.0,
		CreatedAt:     time.Now(),
	}

	event := svc.CreateFailureTriggerEvent(ls, 3, "prog-abc")

	if event == nil {
		t.Fatal("expected non-nil event")
	}
	if event.Type != progression.TriggerOnFailure {
		t.Errorf("Type = %s, want %s", event.Type, progression.TriggerOnFailure)
	}
	if event.UserID != ls.UserID {
		t.Errorf("UserID = %s, want %s", event.UserID, ls.UserID)
	}
	if event.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set")
	}

	ctx, ok := event.Context.(progression.FailureTriggerContext)
	if !ok {
		t.Fatalf("expected FailureTriggerContext, got %T", event.Context)
	}
	if ctx.LoggedSetID != ls.ID {
		t.Errorf("Context.LoggedSetID = %s, want %s", ctx.LoggedSetID, ls.ID)
	}
	if ctx.ConsecutiveFailures != 3 {
		t.Errorf("Context.ConsecutiveFailures = %d, want 3", ctx.ConsecutiveFailures)
	}
}

// TestFailureCheckResult_Structure tests result struct usage.
func TestFailureCheckResult_Structure(t *testing.T) {
	result := FailureCheckResult{
		IsFailure:           true,
		ConsecutiveFailures: 2,
		ProgressionID:       "prog-123",
		TriggerFired:        true,
	}

	if !result.IsFailure {
		t.Error("expected IsFailure to be true")
	}
	if result.ConsecutiveFailures != 2 {
		t.Errorf("ConsecutiveFailures = %d, want 2", result.ConsecutiveFailures)
	}
	if result.ProgressionID != "prog-123" {
		t.Errorf("ProgressionID = %s, want prog-123", result.ProgressionID)
	}
	if !result.TriggerFired {
		t.Error("expected TriggerFired to be true")
	}
}

// TestProcessSetResult_Structure tests result struct usage.
func TestProcessSetResult_Structure(t *testing.T) {
	result := ProcessSetResult{
		LoggedSetID: "set-123",
		Results: []FailureCheckResult{
			{
				IsFailure:           true,
				ConsecutiveFailures: 1,
				ProgressionID:       "prog-1",
			},
			{
				IsFailure:           false,
				ConsecutiveFailures: 0,
				ProgressionID:       "prog-2",
			},
		},
	}

	if result.LoggedSetID != "set-123" {
		t.Errorf("LoggedSetID = %s, want set-123", result.LoggedSetID)
	}
	if len(result.Results) != 2 {
		t.Errorf("Results length = %d, want 2", len(result.Results))
	}
	if !result.Results[0].IsFailure {
		t.Error("first result should be a failure")
	}
	if result.Results[1].IsFailure {
		t.Error("second result should not be a failure")
	}
}
