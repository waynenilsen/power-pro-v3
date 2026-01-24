package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// TestStage_Validate tests Stage validation.
func TestStage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		stage   Stage
		wantErr bool
	}{
		{
			name:    "valid stage",
			stage:   Stage{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			wantErr: false,
		},
		{
			name:    "valid non-AMRAP stage",
			stage:   Stage{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
			wantErr: false,
		},
		{
			name:    "missing name",
			stage:   Stage{Name: "", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			wantErr: true,
		},
		{
			name:    "zero sets",
			stage:   Stage{Name: "5x3+", Sets: 0, Reps: 3, IsAMRAP: true, MinVolume: 15},
			wantErr: true,
		},
		{
			name:    "negative sets",
			stage:   Stage{Name: "5x3+", Sets: -1, Reps: 3, IsAMRAP: true, MinVolume: 15},
			wantErr: true,
		},
		{
			name:    "zero reps",
			stage:   Stage{Name: "5x3+", Sets: 5, Reps: 0, IsAMRAP: true, MinVolume: 15},
			wantErr: true,
		},
		{
			name:    "zero minVolume",
			stage:   Stage{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stage.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestStage_ToSetScheme tests Stage to SetScheme conversion.
func TestStage_ToSetScheme(t *testing.T) {
	t.Run("AMRAP stage returns AMRAPSetScheme", func(t *testing.T) {
		stage := Stage{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15}
		scheme := stage.ToSetScheme()

		if scheme.Type() != setscheme.TypeAMRAP {
			t.Errorf("expected TypeAMRAP, got %s", scheme.Type())
		}
	})

	t.Run("non-AMRAP stage returns FixedSetScheme", func(t *testing.T) {
		stage := Stage{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30}
		scheme := stage.ToSetScheme()

		if scheme.Type() != setscheme.TypeFixed {
			t.Errorf("expected TypeFixed, got %s", scheme.Type())
		}
	})
}

// TestStageProgression_Type tests that StageProgression returns correct type.
func TestStageProgression_Type(t *testing.T) {
	s := &StageProgression{
		ID:   "prog-1",
		Name: "Test Stage",
		Stages: []Stage{
			{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
		},
		MaxTypeValue: TrainingMax,
	}
	if s.Type() != TypeStage {
		t.Errorf("expected %s, got %s", TypeStage, s.Type())
	}
}

// TestStageProgression_TriggerType tests that StageProgression returns ON_FAILURE trigger.
func TestStageProgression_TriggerType(t *testing.T) {
	s := &StageProgression{
		ID:   "prog-1",
		Name: "Test Stage",
		Stages: []Stage{
			{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
		},
		MaxTypeValue: TrainingMax,
	}
	if s.TriggerType() != TriggerOnFailure {
		t.Errorf("expected %s, got %s", TriggerOnFailure, s.TriggerType())
	}
}

// TestStageProgression_Validate tests StageProgression validation.
func TestStageProgression_Validate(t *testing.T) {
	validStages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
		{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
	}

	tests := []struct {
		name    string
		s       StageProgression
		wantErr bool
		errType error
	}{
		{
			name: "valid progression",
			s: StageProgression{
				ID:                "prog-1",
				Name:              "GZCLP T1",
				Stages:            validStages,
				CurrentStage:      0,
				ResetOnExhaustion: true,
				DeloadOnReset:     true,
				DeloadPercent:     0.15,
				MaxTypeValue:      TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "valid without deload",
			s: StageProgression{
				ID:                "prog-1",
				Name:              "GZCLP T2",
				Stages:            validStages,
				CurrentStage:      0,
				ResetOnExhaustion: true,
				DeloadOnReset:     false,
				MaxTypeValue:      TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			s: StageProgression{
				ID:           "",
				Name:         "Test",
				Stages:       validStages,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing name",
			s: StageProgression{
				ID:           "prog-1",
				Name:         "",
				Stages:       validStages,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "only one stage",
			s: StageProgression{
				ID:           "prog-1",
				Name:         "Test",
				Stages:       []Stage{{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15}},
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "no stages",
			s: StageProgression{
				ID:           "prog-1",
				Name:         "Test",
				Stages:       []Stage{},
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "currentStage out of bounds (negative)",
			s: StageProgression{
				ID:           "prog-1",
				Name:         "Test",
				Stages:       validStages,
				CurrentStage: -1,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "currentStage out of bounds (too high)",
			s: StageProgression{
				ID:           "prog-1",
				Name:         "Test",
				Stages:       validStages,
				CurrentStage: 5,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "invalid max type",
			s: StageProgression{
				ID:           "prog-1",
				Name:         "Test",
				Stages:       validStages,
				MaxTypeValue: "INVALID",
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "deloadOnReset without resetOnExhaustion",
			s: StageProgression{
				ID:                "prog-1",
				Name:              "Test",
				Stages:            validStages,
				ResetOnExhaustion: false,
				DeloadOnReset:     true,
				DeloadPercent:     0.15,
				MaxTypeValue:      TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "deloadOnReset with zero percent",
			s: StageProgression{
				ID:                "prog-1",
				Name:              "Test",
				Stages:            validStages,
				ResetOnExhaustion: true,
				DeloadOnReset:     true,
				DeloadPercent:     0,
				MaxTypeValue:      TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "deloadOnReset with percent > 1",
			s: StageProgression{
				ID:                "prog-1",
				Name:              "Test",
				Stages:            validStages,
				ResetOnExhaustion: true,
				DeloadOnReset:     true,
				DeloadPercent:     1.5,
				MaxTypeValue:      TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "invalid stage in list",
			s: StageProgression{
				ID:   "prog-1",
				Name: "Test",
				Stages: []Stage{
					{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
					{Name: "", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12}, // Invalid: empty name
				},
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected error type %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestNewStageProgression tests the factory function.
func TestNewStageProgression(t *testing.T) {
	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
	}

	t.Run("valid parameters", func(t *testing.T) {
		s, err := NewStageProgression("prog-1", "GZCLP T1", stages, true, true, 0.15, TrainingMax)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", s.ID)
		}
		if s.CurrentStage != 0 {
			t.Errorf("expected CurrentStage 0, got %d", s.CurrentStage)
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		_, err := NewStageProgression("", "Test", stages, true, true, 0.15, TrainingMax)
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
}

// makeStageFailureTriggerEvent creates a TriggerEvent for stage progression tests.
func makeStageFailureTriggerEvent() TriggerEvent {
	return TriggerEvent{
		Type:                TriggerOnFailure,
		Timestamp:           time.Now(),
		FailedSetID:         strPtr("set-1"),
		TargetReps:          intPtr(3),
		RepsPerformed:       intPtr(2),
		ConsecutiveFailures: intPtr(1),
		SetWeight:           floatPtr(200),
		ProgressionID:       strPtr("prog-1"),
	}
}

// TestStageProgression_Apply tests the Apply method.
func TestStageProgression_Apply(t *testing.T) {
	ctx := context.Background()

	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
		{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
	}

	t.Run("advance from stage 0 to stage 1", func(t *testing.T) {
		s := &StageProgression{
			ID:                "prog-1",
			Name:              "GZCLP T1",
			Stages:            stages,
			CurrentStage:      0,
			ResetOnExhaustion: true,
			DeloadOnReset:     true,
			DeloadPercent:     0.15,
			MaxTypeValue:      TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0 (weight stays same), got %f", result.Delta)
		}
		if result.NewValue != 200 {
			t.Errorf("expected NewValue 200, got %f", result.NewValue)
		}
		if s.CurrentStage != 1 {
			t.Errorf("expected CurrentStage 1, got %d", s.CurrentStage)
		}
	})

	t.Run("advance from stage 1 to stage 2", func(t *testing.T) {
		s := &StageProgression{
			ID:                "prog-1",
			Name:              "GZCLP T1",
			Stages:            stages,
			CurrentStage:      1,
			ResetOnExhaustion: true,
			DeloadOnReset:     true,
			DeloadPercent:     0.15,
			MaxTypeValue:      TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0, got %f", result.Delta)
		}
		if s.CurrentStage != 2 {
			t.Errorf("expected CurrentStage 2, got %d", s.CurrentStage)
		}
	})

	t.Run("reset on exhaustion with deload", func(t *testing.T) {
		s := &StageProgression{
			ID:                "prog-1",
			Name:              "GZCLP T1",
			Stages:            stages,
			CurrentStage:      2, // At last stage
			ResetOnExhaustion: true,
			DeloadOnReset:     true,
			DeloadPercent:     0.15,
			MaxTypeValue:      TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.NewValue != 170 { // 200 - 15% = 170
			t.Errorf("expected NewValue 170, got %f", result.NewValue)
		}
		if result.Delta != -30 {
			t.Errorf("expected Delta -30, got %f", result.Delta)
		}
		if s.CurrentStage != 0 {
			t.Errorf("expected CurrentStage reset to 0, got %d", s.CurrentStage)
		}
	})

	t.Run("reset on exhaustion without deload", func(t *testing.T) {
		s := &StageProgression{
			ID:                "prog-1",
			Name:              "GZCLP T2",
			Stages:            stages,
			CurrentStage:      2, // At last stage
			ResetOnExhaustion: true,
			DeloadOnReset:     false,
			MaxTypeValue:      TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.NewValue != 200 { // No deload
			t.Errorf("expected NewValue 200 (no deload), got %f", result.NewValue)
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0, got %f", result.Delta)
		}
		if s.CurrentStage != 0 {
			t.Errorf("expected CurrentStage reset to 0, got %d", s.CurrentStage)
		}
	})

	t.Run("no reset on exhaustion requires manual intervention", func(t *testing.T) {
		s := &StageProgression{
			ID:                "prog-1",
			Name:              "Manual Reset",
			Stages:            stages,
			CurrentStage:      2, // At last stage
			ResetOnExhaustion: false,
			MaxTypeValue:      TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when manual intervention required")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("trigger type mismatch", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			MaxTypeValue: TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession, // Wrong trigger type
				Timestamp: time.Now(),
			},
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for trigger type mismatch")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("max type mismatch", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			MaxTypeValue: TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      OneRM, // Wrong max type
			CurrentValue: 225,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		result, err := s.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for max type mismatch")
		}
	})

	t.Run("invalid context", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			MaxTypeValue: TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeStageFailureTriggerEvent(),
		}

		_, err := s.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})
}

// TestStageProgression_GetCurrentStage tests the GetCurrentStage method.
func TestStageProgression_GetCurrentStage(t *testing.T) {
	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
	}

	t.Run("returns correct stage", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			CurrentStage: 1,
			MaxTypeValue: TrainingMax,
		}

		stage := s.GetCurrentStage()
		if stage.Name != "6x2+" {
			t.Errorf("expected stage name '6x2+', got '%s'", stage.Name)
		}
	})

	t.Run("returns first stage as fallback for out of bounds", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			CurrentStage: 10, // Out of bounds
			MaxTypeValue: TrainingMax,
		}

		stage := s.GetCurrentStage()
		if stage.Name != "5x3+" {
			t.Errorf("expected fallback to first stage '5x3+', got '%s'", stage.Name)
		}
	})
}

// TestStageProgression_SetCurrentStage tests the SetCurrentStage method.
func TestStageProgression_SetCurrentStage(t *testing.T) {
	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
	}

	s := &StageProgression{
		ID:           "prog-1",
		Name:         "Test",
		Stages:       stages,
		CurrentStage: 0,
		MaxTypeValue: TrainingMax,
	}

	t.Run("valid stage index", func(t *testing.T) {
		err := s.SetCurrentStage(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.CurrentStage != 1 {
			t.Errorf("expected CurrentStage 1, got %d", s.CurrentStage)
		}
	})

	t.Run("negative stage index", func(t *testing.T) {
		err := s.SetCurrentStage(-1)
		if err == nil {
			t.Error("expected error for negative stage index")
		}
	})

	t.Run("stage index too high", func(t *testing.T) {
		err := s.SetCurrentStage(5)
		if err == nil {
			t.Error("expected error for stage index too high")
		}
	})
}

// TestStageProgression_StageCount tests the StageCount method.
func TestStageProgression_StageCount(t *testing.T) {
	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
		{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
	}

	s := &StageProgression{
		ID:           "prog-1",
		Name:         "Test",
		Stages:       stages,
		MaxTypeValue: TrainingMax,
	}

	if s.StageCount() != 3 {
		t.Errorf("expected StageCount 3, got %d", s.StageCount())
	}
}

// TestStageProgression_IsAtLastStage tests the IsAtLastStage method.
func TestStageProgression_IsAtLastStage(t *testing.T) {
	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
	}

	t.Run("not at last stage", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			CurrentStage: 0,
			MaxTypeValue: TrainingMax,
		}
		if s.IsAtLastStage() {
			t.Error("expected IsAtLastStage to return false at stage 0")
		}
	})

	t.Run("at last stage", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			CurrentStage: 1,
			MaxTypeValue: TrainingMax,
		}
		if !s.IsAtLastStage() {
			t.Error("expected IsAtLastStage to return true at stage 1")
		}
	})
}

// TestStageProgression_ShouldResetFailureCounter tests the ShouldResetFailureCounter method.
func TestStageProgression_ShouldResetFailureCounter(t *testing.T) {
	s := &StageProgression{}
	if !s.ShouldResetFailureCounter() {
		t.Error("expected ShouldResetFailureCounter to return true")
	}
}

// TestStageProgression_JSON tests JSON serialization roundtrip.
func TestStageProgression_JSON(t *testing.T) {
	s := &StageProgression{
		ID:   "prog-123",
		Name: "GZCLP T1 Default",
		Stages: []Stage{
			{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
			{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
		},
		CurrentStage:      1,
		ResetOnExhaustion: true,
		DeloadOnReset:     true,
		DeloadPercent:     0.15,
		MaxTypeValue:      TrainingMax,
	}

	// Marshal
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeStage) {
		t.Errorf("expected type %s, got %v", TypeStage, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["currentStage"] != float64(1) {
		t.Errorf("expected currentStage 1, got %v", parsed["currentStage"])
	}

	// Unmarshal back
	var restored StageProgression
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != s.ID {
		t.Errorf("ID mismatch: expected %s, got %s", s.ID, restored.ID)
	}
	if restored.Name != s.Name {
		t.Errorf("Name mismatch: expected %s, got %s", s.Name, restored.Name)
	}
	if restored.CurrentStage != s.CurrentStage {
		t.Errorf("CurrentStage mismatch: expected %d, got %d", s.CurrentStage, restored.CurrentStage)
	}
	if len(restored.Stages) != len(s.Stages) {
		t.Errorf("Stages length mismatch: expected %d, got %d", len(s.Stages), len(restored.Stages))
	}
	if restored.DeloadPercent != s.DeloadPercent {
		t.Errorf("DeloadPercent mismatch: expected %f, got %f", s.DeloadPercent, restored.DeloadPercent)
	}
}

// TestUnmarshalStageProgression tests deserialization function.
func TestUnmarshalStageProgression(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "GZCLP T1",
			"stages": [
				{"name": "5x3+", "sets": 5, "reps": 3, "isAmrap": true, "minVolume": 15},
				{"name": "6x2+", "sets": 6, "reps": 2, "isAmrap": true, "minVolume": 12}
			],
			"currentStage": 0,
			"resetOnExhaustion": true,
			"deloadOnReset": true,
			"deloadPercent": 0.15,
			"maxType": "TRAINING_MAX"
		}`)

		progression, err := UnmarshalStageProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		s, ok := progression.(*StageProgression)
		if !ok {
			t.Fatalf("expected *StageProgression, got %T", progression)
		}
		if s.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", s.ID)
		}
		if len(s.Stages) != 2 {
			t.Errorf("expected 2 stages, got %d", len(s.Stages))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalStageProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"stages": [{"name": "5x3+", "sets": 5, "reps": 3, "isAmrap": true, "minVolume": 15}],
			"maxType": "TRAINING_MAX"
		}`)

		_, err := UnmarshalStageProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})
}

// TestRegisterStageProgression tests factory registration.
func TestRegisterStageProgression(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(TypeStage) {
		t.Error("TypeStage should not be registered initially")
	}

	// Register
	RegisterStageProgression(factory)

	// Verify registered
	if !factory.IsRegistered(TypeStage) {
		t.Error("TypeStage should be registered after calling RegisterStageProgression")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "STAGE_PROGRESSION",
		"id": "prog-1",
		"name": "Factory Test",
		"stages": [
			{"name": "5x3+", "sets": 5, "reps": 3, "isAmrap": true, "minVolume": 15},
			{"name": "6x2+", "sets": 6, "reps": 2, "isAmrap": true, "minVolume": 12}
		],
		"currentStage": 0,
		"resetOnExhaustion": true,
		"deloadOnReset": false,
		"maxType": "TRAINING_MAX"
	}`)

	progression, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progression.Type() != TypeStage {
		t.Errorf("expected type %s, got %s", TypeStage, progression.Type())
	}
}

// TestStageProgression_Interface verifies that StageProgression implements Progression.
func TestStageProgression_Interface(t *testing.T) {
	var _ Progression = &StageProgression{}
}

// TestStageProgression_FullGZCLPT1Cycle tests a complete GZCLP T1 cycle.
func TestStageProgression_FullGZCLPT1Cycle(t *testing.T) {
	ctx := context.Background()

	// Create GZCLP T1 progression
	s, err := NewGZCLPT1DefaultProgression("gzclp-t1", "GZCLP T1 Squat")
	if err != nil {
		t.Fatalf("failed to create GZCLP T1 progression: %v", err)
	}

	// Starting weight
	currentWeight := 200.0

	// Stage 0: 5x3+ - fail
	if s.CurrentStage != 0 {
		t.Errorf("expected starting at stage 0, got %d", s.CurrentStage)
	}
	if s.GetCurrentStage().Name != "5x3+" {
		t.Errorf("expected stage name '5x3+', got '%s'", s.GetCurrentStage().Name)
	}

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "squat",
		MaxType:      TrainingMax,
		CurrentValue: currentWeight,
		TriggerEvent: makeStageFailureTriggerEvent(),
	}

	result1, err := s.Apply(ctx, params)
	if err != nil {
		t.Fatalf("Apply 1 failed: %v", err)
	}
	if !result1.Applied {
		t.Errorf("expected Apply 1 to be Applied, reason: %s", result1.Reason)
	}
	if s.CurrentStage != 1 {
		t.Errorf("expected stage 1 after first failure, got %d", s.CurrentStage)
	}
	if result1.Delta != 0 {
		t.Errorf("expected no weight change on stage advance, got delta %f", result1.Delta)
	}

	// Stage 1: 6x2+ - fail
	if s.GetCurrentStage().Name != "6x2+" {
		t.Errorf("expected stage name '6x2+', got '%s'", s.GetCurrentStage().Name)
	}

	result2, err := s.Apply(ctx, params)
	if err != nil {
		t.Fatalf("Apply 2 failed: %v", err)
	}
	if !result2.Applied {
		t.Errorf("expected Apply 2 to be Applied, reason: %s", result2.Reason)
	}
	if s.CurrentStage != 2 {
		t.Errorf("expected stage 2 after second failure, got %d", s.CurrentStage)
	}

	// Stage 2: 10x1+ - fail (reset with deload)
	if s.GetCurrentStage().Name != "10x1+" {
		t.Errorf("expected stage name '10x1+', got '%s'", s.GetCurrentStage().Name)
	}

	result3, err := s.Apply(ctx, params)
	if err != nil {
		t.Fatalf("Apply 3 failed: %v", err)
	}
	if !result3.Applied {
		t.Errorf("expected Apply 3 to be Applied, reason: %s", result3.Reason)
	}
	if s.CurrentStage != 0 {
		t.Errorf("expected reset to stage 0, got %d", s.CurrentStage)
	}
	expectedNewWeight := 170.0 // 200 - 15% = 170
	if result3.NewValue != expectedNewWeight {
		t.Errorf("expected NewValue %f after deload, got %f", expectedNewWeight, result3.NewValue)
	}
	if result3.Delta != -30 {
		t.Errorf("expected Delta -30, got %f", result3.Delta)
	}

	// Back to stage 0: 5x3+
	if s.GetCurrentStage().Name != "5x3+" {
		t.Errorf("expected stage name '5x3+' after reset, got '%s'", s.GetCurrentStage().Name)
	}
}

// TestStageProgression_GZCLPT2Cycle tests a GZCLP T2 cycle (no deload on reset).
func TestStageProgression_GZCLPT2Cycle(t *testing.T) {
	ctx := context.Background()

	// Create GZCLP T2 progression
	s, err := NewGZCLPT2DefaultProgression("gzclp-t2", "GZCLP T2 Lat Pulldown")
	if err != nil {
		t.Fatalf("failed to create GZCLP T2 progression: %v", err)
	}

	currentWeight := 100.0

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "lat-pulldown",
		MaxType:      TrainingMax,
		CurrentValue: currentWeight,
		TriggerEvent: makeStageFailureTriggerEvent(),
	}

	// Advance through all stages
	s.Apply(ctx, params) // 3x10 -> 3x8
	s.Apply(ctx, params) // 3x8 -> 3x6

	// Final failure - reset without deload
	result, err := s.Apply(ctx, params)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if !result.Applied {
		t.Errorf("expected Applied=true, reason: %s", result.Reason)
	}
	if result.NewValue != 100 {
		t.Errorf("expected no deload (NewValue 100), got %f", result.NewValue)
	}
	if result.Delta != 0 {
		t.Errorf("expected Delta 0 (no deload), got %f", result.Delta)
	}
	if s.CurrentStage != 0 {
		t.Errorf("expected reset to stage 0, got %d", s.CurrentStage)
	}
}

// TestGZCLPPresetProgressions tests the preset progression factory functions.
func TestGZCLPPresetProgressions(t *testing.T) {
	t.Run("T1 Default", func(t *testing.T) {
		s, err := NewGZCLPT1DefaultProgression("t1", "T1 Squat")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s.Stages) != 3 {
			t.Errorf("expected 3 stages, got %d", len(s.Stages))
		}
		if s.Stages[0].Name != "5x3+" {
			t.Errorf("expected first stage '5x3+', got '%s'", s.Stages[0].Name)
		}
		if !s.DeloadOnReset {
			t.Error("expected DeloadOnReset to be true")
		}
		if s.DeloadPercent != 0.15 {
			t.Errorf("expected DeloadPercent 0.15, got %f", s.DeloadPercent)
		}
	})

	t.Run("T1 Modified", func(t *testing.T) {
		s, err := NewGZCLPT1ModifiedProgression("t1m", "T1 Modified")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Stages[0].Name != "3x5+" {
			t.Errorf("expected first stage '3x5+', got '%s'", s.Stages[0].Name)
		}
	})

	t.Run("T2 Default", func(t *testing.T) {
		s, err := NewGZCLPT2DefaultProgression("t2", "T2 Rows")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s.Stages) != 3 {
			t.Errorf("expected 3 stages, got %d", len(s.Stages))
		}
		if s.Stages[0].Name != "3x10" {
			t.Errorf("expected first stage '3x10', got '%s'", s.Stages[0].Name)
		}
		if s.Stages[0].IsAMRAP {
			t.Error("expected T2 stages to not be AMRAP")
		}
		if s.DeloadOnReset {
			t.Error("expected DeloadOnReset to be false for T2")
		}
	})
}

// TestStageProgression_GetCurrentSetScheme tests the GetCurrentSetScheme method.
func TestStageProgression_GetCurrentSetScheme(t *testing.T) {
	stages := []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
	}

	t.Run("AMRAP stage returns AMRAP scheme", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			CurrentStage: 0,
			MaxTypeValue: TrainingMax,
		}

		scheme := s.GetCurrentSetScheme()
		if scheme.Type() != setscheme.TypeAMRAP {
			t.Errorf("expected TypeAMRAP, got %s", scheme.Type())
		}
	})

	t.Run("Fixed stage returns Fixed scheme", func(t *testing.T) {
		s := &StageProgression{
			ID:           "prog-1",
			Name:         "Test",
			Stages:       stages,
			CurrentStage: 1,
			MaxTypeValue: TrainingMax,
		}

		scheme := s.GetCurrentSetScheme()
		if scheme.Type() != setscheme.TypeFixed {
			t.Errorf("expected TypeFixed, got %s", scheme.Type())
		}
	})
}
