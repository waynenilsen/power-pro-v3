package progression

import (
	"encoding/json"
	"testing"
	"time"
)

// TestSessionTriggerContext_TriggerType tests that SessionTriggerContext returns correct type.
func TestSessionTriggerContext_TriggerType(t *testing.T) {
	ctx := SessionTriggerContext{}
	if ctx.TriggerType() != TriggerAfterSession {
		t.Errorf("expected %s, got %s", TriggerAfterSession, ctx.TriggerType())
	}
}

// TestWeekTriggerContext_TriggerType tests that WeekTriggerContext returns correct type.
func TestWeekTriggerContext_TriggerType(t *testing.T) {
	ctx := WeekTriggerContext{}
	if ctx.TriggerType() != TriggerAfterWeek {
		t.Errorf("expected %s, got %s", TriggerAfterWeek, ctx.TriggerType())
	}
}

// TestCycleTriggerContext_TriggerType tests that CycleTriggerContext returns correct type.
func TestCycleTriggerContext_TriggerType(t *testing.T) {
	ctx := CycleTriggerContext{}
	if ctx.TriggerType() != TriggerAfterCycle {
		t.Errorf("expected %s, got %s", TriggerAfterCycle, ctx.TriggerType())
	}
}

// TestSessionTriggerContext_Validate tests SessionTriggerContext validation.
func TestSessionTriggerContext_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ctx     SessionTriggerContext
		wantErr bool
	}{
		{
			name: "valid context",
			ctx: SessionTriggerContext{
				SessionID:      "session-123",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{"lift-1", "lift-2"},
			},
			wantErr: false,
		},
		{
			name: "valid context with empty lifts",
			ctx: SessionTriggerContext{
				SessionID:      "session-123",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{},
			},
			wantErr: false,
		},
		{
			name: "valid context with nil lifts",
			ctx: SessionTriggerContext{
				SessionID:  "session-123",
				DaySlug:    "day-a",
				WeekNumber: 1,
			},
			wantErr: false,
		},
		{
			name: "missing sessionId",
			ctx: SessionTriggerContext{
				SessionID:  "",
				DaySlug:    "day-a",
				WeekNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "missing daySlug",
			ctx: SessionTriggerContext{
				SessionID:  "session-123",
				DaySlug:    "",
				WeekNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "zero weekNumber",
			ctx: SessionTriggerContext{
				SessionID:  "session-123",
				DaySlug:    "day-a",
				WeekNumber: 0,
			},
			wantErr: true,
		},
		{
			name: "negative weekNumber",
			ctx: SessionTriggerContext{
				SessionID:  "session-123",
				DaySlug:    "day-a",
				WeekNumber: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestWeekTriggerContext_Validate tests WeekTriggerContext validation.
func TestWeekTriggerContext_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ctx     WeekTriggerContext
		wantErr bool
	}{
		{
			name: "valid context",
			ctx: WeekTriggerContext{
				PreviousWeek:   1,
				NewWeek:        2,
				CycleIteration: 1,
			},
			wantErr: false,
		},
		{
			name: "valid context week 3 to 4",
			ctx: WeekTriggerContext{
				PreviousWeek:   3,
				NewWeek:        4,
				CycleIteration: 2,
			},
			wantErr: false,
		},
		{
			name: "zero previousWeek",
			ctx: WeekTriggerContext{
				PreviousWeek:   0,
				NewWeek:        1,
				CycleIteration: 1,
			},
			wantErr: true,
		},
		{
			name: "zero newWeek",
			ctx: WeekTriggerContext{
				PreviousWeek:   1,
				NewWeek:        0,
				CycleIteration: 1,
			},
			wantErr: true,
		},
		{
			name: "newWeek equals previousWeek",
			ctx: WeekTriggerContext{
				PreviousWeek:   2,
				NewWeek:        2,
				CycleIteration: 1,
			},
			wantErr: true,
		},
		{
			name: "newWeek less than previousWeek",
			ctx: WeekTriggerContext{
				PreviousWeek:   3,
				NewWeek:        1,
				CycleIteration: 1,
			},
			wantErr: true,
		},
		{
			name: "zero cycleIteration",
			ctx: WeekTriggerContext{
				PreviousWeek:   1,
				NewWeek:        2,
				CycleIteration: 0,
			},
			wantErr: true,
		},
		{
			name: "negative cycleIteration",
			ctx: WeekTriggerContext{
				PreviousWeek:   1,
				NewWeek:        2,
				CycleIteration: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestCycleTriggerContext_Validate tests CycleTriggerContext validation.
func TestCycleTriggerContext_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ctx     CycleTriggerContext
		wantErr bool
	}{
		{
			name: "valid context first cycle",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       2,
				TotalWeeks:     4,
			},
			wantErr: false,
		},
		{
			name: "valid context fifth cycle",
			ctx: CycleTriggerContext{
				CompletedCycle: 5,
				NewCycle:       6,
				TotalWeeks:     3,
			},
			wantErr: false,
		},
		{
			name: "zero completedCycle",
			ctx: CycleTriggerContext{
				CompletedCycle: 0,
				NewCycle:       1,
				TotalWeeks:     4,
			},
			wantErr: true,
		},
		{
			name: "zero newCycle",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       0,
				TotalWeeks:     4,
			},
			wantErr: true,
		},
		{
			name: "newCycle not completedCycle plus 1",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       3,
				TotalWeeks:     4,
			},
			wantErr: true,
		},
		{
			name: "newCycle less than completedCycle",
			ctx: CycleTriggerContext{
				CompletedCycle: 3,
				NewCycle:       2,
				TotalWeeks:     4,
			},
			wantErr: true,
		},
		{
			name: "zero totalWeeks",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       2,
				TotalWeeks:     0,
			},
			wantErr: true,
		},
		{
			name: "negative totalWeeks",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       2,
				TotalWeeks:     -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestTriggerEventV2_Validate tests TriggerEventV2 validation.
func TestTriggerEventV2_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		event   TriggerEventV2
		wantErr bool
	}{
		{
			name: "valid session trigger event",
			event: TriggerEventV2{
				Type:      TriggerAfterSession,
				UserID:    "user-123",
				Timestamp: now,
				Context: SessionTriggerContext{
					SessionID:      "session-456",
					DaySlug:        "day-a",
					WeekNumber:     1,
					LiftsPerformed: []string{"lift-1"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid week trigger event",
			event: TriggerEventV2{
				Type:      TriggerAfterWeek,
				UserID:    "user-123",
				Timestamp: now,
				Context: WeekTriggerContext{
					PreviousWeek:   1,
					NewWeek:        2,
					CycleIteration: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "valid cycle trigger event",
			event: TriggerEventV2{
				Type:      TriggerAfterCycle,
				UserID:    "user-123",
				Timestamp: now,
				Context: CycleTriggerContext{
					CompletedCycle: 1,
					NewCycle:       2,
					TotalWeeks:     4,
				},
			},
			wantErr: false,
		},
		{
			name: "missing trigger type",
			event: TriggerEventV2{
				Type:      "",
				UserID:    "user-123",
				Timestamp: now,
				Context: SessionTriggerContext{
					SessionID:  "session-456",
					DaySlug:    "day-a",
					WeekNumber: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid trigger type",
			event: TriggerEventV2{
				Type:      "INVALID_TYPE",
				UserID:    "user-123",
				Timestamp: now,
				Context: SessionTriggerContext{
					SessionID:  "session-456",
					DaySlug:    "day-a",
					WeekNumber: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "missing userID",
			event: TriggerEventV2{
				Type:      TriggerAfterSession,
				UserID:    "",
				Timestamp: now,
				Context: SessionTriggerContext{
					SessionID:  "session-456",
					DaySlug:    "day-a",
					WeekNumber: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "missing timestamp",
			event: TriggerEventV2{
				Type:   TriggerAfterSession,
				UserID: "user-123",
				Context: SessionTriggerContext{
					SessionID:  "session-456",
					DaySlug:    "day-a",
					WeekNumber: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "nil context",
			event: TriggerEventV2{
				Type:      TriggerAfterSession,
				UserID:    "user-123",
				Timestamp: now,
				Context:   nil,
			},
			wantErr: true,
		},
		{
			name: "context type mismatch - session event with week context",
			event: TriggerEventV2{
				Type:      TriggerAfterSession,
				UserID:    "user-123",
				Timestamp: now,
				Context: WeekTriggerContext{
					PreviousWeek:   1,
					NewWeek:        2,
					CycleIteration: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "context type mismatch - week event with cycle context",
			event: TriggerEventV2{
				Type:      TriggerAfterWeek,
				UserID:    "user-123",
				Timestamp: now,
				Context: CycleTriggerContext{
					CompletedCycle: 1,
					NewCycle:       2,
					TotalWeeks:     4,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid context data",
			event: TriggerEventV2{
				Type:      TriggerAfterSession,
				UserID:    "user-123",
				Timestamp: now,
				Context: SessionTriggerContext{
					SessionID:  "",
					DaySlug:    "day-a",
					WeekNumber: 1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestSessionTriggerContext_JSON tests JSON serialization of SessionTriggerContext.
func TestSessionTriggerContext_JSON(t *testing.T) {
	ctx := SessionTriggerContext{
		SessionID:      "session-123",
		DaySlug:        "heavy-day",
		WeekNumber:     2,
		LiftsPerformed: []string{"squat-uuid", "bench-uuid", "row-uuid"},
	}

	// Marshal
	data, err := json.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["sessionId"] != "session-123" {
		t.Errorf("expected sessionId 'session-123', got %v", parsed["sessionId"])
	}
	if parsed["daySlug"] != "heavy-day" {
		t.Errorf("expected daySlug 'heavy-day', got %v", parsed["daySlug"])
	}
	if parsed["weekNumber"] != float64(2) {
		t.Errorf("expected weekNumber 2, got %v", parsed["weekNumber"])
	}

	// Unmarshal back
	var restored SessionTriggerContext
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.SessionID != ctx.SessionID {
		t.Errorf("SessionID mismatch: expected %s, got %s", ctx.SessionID, restored.SessionID)
	}
	if restored.DaySlug != ctx.DaySlug {
		t.Errorf("DaySlug mismatch: expected %s, got %s", ctx.DaySlug, restored.DaySlug)
	}
	if restored.WeekNumber != ctx.WeekNumber {
		t.Errorf("WeekNumber mismatch: expected %d, got %d", ctx.WeekNumber, restored.WeekNumber)
	}
	if len(restored.LiftsPerformed) != len(ctx.LiftsPerformed) {
		t.Errorf("LiftsPerformed length mismatch: expected %d, got %d", len(ctx.LiftsPerformed), len(restored.LiftsPerformed))
	}
}

// TestWeekTriggerContext_JSON tests JSON serialization of WeekTriggerContext.
func TestWeekTriggerContext_JSON(t *testing.T) {
	ctx := WeekTriggerContext{
		PreviousWeek:   3,
		NewWeek:        4,
		CycleIteration: 2,
	}

	// Marshal
	data, err := json.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var restored WeekTriggerContext
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.PreviousWeek != ctx.PreviousWeek {
		t.Errorf("PreviousWeek mismatch: expected %d, got %d", ctx.PreviousWeek, restored.PreviousWeek)
	}
	if restored.NewWeek != ctx.NewWeek {
		t.Errorf("NewWeek mismatch: expected %d, got %d", ctx.NewWeek, restored.NewWeek)
	}
	if restored.CycleIteration != ctx.CycleIteration {
		t.Errorf("CycleIteration mismatch: expected %d, got %d", ctx.CycleIteration, restored.CycleIteration)
	}
}

// TestCycleTriggerContext_JSON tests JSON serialization of CycleTriggerContext.
func TestCycleTriggerContext_JSON(t *testing.T) {
	ctx := CycleTriggerContext{
		CompletedCycle: 5,
		NewCycle:       6,
		TotalWeeks:     4,
	}

	// Marshal
	data, err := json.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var restored CycleTriggerContext
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.CompletedCycle != ctx.CompletedCycle {
		t.Errorf("CompletedCycle mismatch: expected %d, got %d", ctx.CompletedCycle, restored.CompletedCycle)
	}
	if restored.NewCycle != ctx.NewCycle {
		t.Errorf("NewCycle mismatch: expected %d, got %d", ctx.NewCycle, restored.NewCycle)
	}
	if restored.TotalWeeks != ctx.TotalWeeks {
		t.Errorf("TotalWeeks mismatch: expected %d, got %d", ctx.TotalWeeks, restored.TotalWeeks)
	}
}

// TestTriggerEventV2_JSON tests TriggerEventV2 JSON serialization roundtrip.
func TestTriggerEventV2_JSON(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name  string
		event TriggerEventV2
	}{
		{
			name: "session trigger event",
			event: TriggerEventV2{
				Type:      TriggerAfterSession,
				UserID:    "user-123",
				Timestamp: timestamp,
				Context: SessionTriggerContext{
					SessionID:      "session-456",
					DaySlug:        "day-a",
					WeekNumber:     1,
					LiftsPerformed: []string{"lift-1", "lift-2"},
				},
			},
		},
		{
			name: "week trigger event",
			event: TriggerEventV2{
				Type:      TriggerAfterWeek,
				UserID:    "user-456",
				Timestamp: timestamp,
				Context: WeekTriggerContext{
					PreviousWeek:   2,
					NewWeek:        3,
					CycleIteration: 1,
				},
			},
		},
		{
			name: "cycle trigger event",
			event: TriggerEventV2{
				Type:      TriggerAfterCycle,
				UserID:    "user-789",
				Timestamp: timestamp,
				Context: CycleTriggerContext{
					CompletedCycle: 1,
					NewCycle:       2,
					TotalWeeks:     4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tt.event)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			// Verify JSON contains expected fields
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("failed to parse JSON: %v", err)
			}
			if parsed["type"] != string(tt.event.Type) {
				t.Errorf("type mismatch: expected %s, got %v", tt.event.Type, parsed["type"])
			}
			if parsed["userId"] != tt.event.UserID {
				t.Errorf("userId mismatch: expected %s, got %v", tt.event.UserID, parsed["userId"])
			}
			if _, ok := parsed["context"]; !ok {
				t.Error("expected 'context' field in JSON")
			}

			// Unmarshal back
			var restored TriggerEventV2
			if err := json.Unmarshal(data, &restored); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			// Verify restored event
			if restored.Type != tt.event.Type {
				t.Errorf("Type mismatch: expected %s, got %s", tt.event.Type, restored.Type)
			}
			if restored.UserID != tt.event.UserID {
				t.Errorf("UserID mismatch: expected %s, got %s", tt.event.UserID, restored.UserID)
			}
			if !restored.Timestamp.Equal(tt.event.Timestamp) {
				t.Errorf("Timestamp mismatch: expected %v, got %v", tt.event.Timestamp, restored.Timestamp)
			}
			if restored.Context == nil {
				t.Fatal("expected Context to be restored")
			}
			if restored.Context.TriggerType() != tt.event.Context.TriggerType() {
				t.Errorf("Context type mismatch: expected %s, got %s",
					tt.event.Context.TriggerType(), restored.Context.TriggerType())
			}
		})
	}
}

// TestUnmarshalTriggerContext tests context deserialization by type.
func TestUnmarshalTriggerContext(t *testing.T) {
	tests := []struct {
		name        string
		triggerType TriggerType
		json        string
		wantErr     bool
	}{
		{
			name:        "session context",
			triggerType: TriggerAfterSession,
			json:        `{"sessionId":"s-1","daySlug":"day-a","weekNumber":1,"liftsPerformed":["l-1"]}`,
			wantErr:     false,
		},
		{
			name:        "week context",
			triggerType: TriggerAfterWeek,
			json:        `{"previousWeek":1,"newWeek":2,"cycleIteration":1}`,
			wantErr:     false,
		},
		{
			name:        "cycle context",
			triggerType: TriggerAfterCycle,
			json:        `{"completedCycle":1,"newCycle":2,"totalWeeks":4}`,
			wantErr:     false,
		},
		{
			name:        "unknown trigger type",
			triggerType: "UNKNOWN",
			json:        `{}`,
			wantErr:     true,
		},
		{
			name:        "invalid json for session",
			triggerType: TriggerAfterSession,
			json:        `{invalid}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := UnmarshalTriggerContext(tt.triggerType, []byte(tt.json))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ctx.TriggerType() != tt.triggerType {
				t.Errorf("TriggerType mismatch: expected %s, got %s", tt.triggerType, ctx.TriggerType())
			}
		})
	}
}

// TestNewSessionTriggerEvent tests session trigger event factory.
func TestNewSessionTriggerEvent(t *testing.T) {
	userID := "user-123"
	sessionID := "session-456"
	daySlug := "day-a"
	weekNumber := 2
	lifts := []string{"lift-1", "lift-2"}

	event := NewSessionTriggerEvent(userID, sessionID, daySlug, weekNumber, lifts)

	if event.Type != TriggerAfterSession {
		t.Errorf("expected type %s, got %s", TriggerAfterSession, event.Type)
	}
	if event.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, event.UserID)
	}
	if event.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set")
	}

	ctx, ok := event.Context.(SessionTriggerContext)
	if !ok {
		t.Fatalf("expected SessionTriggerContext, got %T", event.Context)
	}
	if ctx.SessionID != sessionID {
		t.Errorf("expected sessionID %s, got %s", sessionID, ctx.SessionID)
	}
	if ctx.DaySlug != daySlug {
		t.Errorf("expected daySlug %s, got %s", daySlug, ctx.DaySlug)
	}
	if ctx.WeekNumber != weekNumber {
		t.Errorf("expected weekNumber %d, got %d", weekNumber, ctx.WeekNumber)
	}
	if len(ctx.LiftsPerformed) != len(lifts) {
		t.Errorf("expected %d lifts, got %d", len(lifts), len(ctx.LiftsPerformed))
	}
}

// TestNewWeekTriggerEvent tests week trigger event factory.
func TestNewWeekTriggerEvent(t *testing.T) {
	userID := "user-123"
	previousWeek := 2
	newWeek := 3
	cycleIteration := 1

	event := NewWeekTriggerEvent(userID, previousWeek, newWeek, cycleIteration)

	if event.Type != TriggerAfterWeek {
		t.Errorf("expected type %s, got %s", TriggerAfterWeek, event.Type)
	}
	if event.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, event.UserID)
	}
	if event.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set")
	}

	ctx, ok := event.Context.(WeekTriggerContext)
	if !ok {
		t.Fatalf("expected WeekTriggerContext, got %T", event.Context)
	}
	if ctx.PreviousWeek != previousWeek {
		t.Errorf("expected previousWeek %d, got %d", previousWeek, ctx.PreviousWeek)
	}
	if ctx.NewWeek != newWeek {
		t.Errorf("expected newWeek %d, got %d", newWeek, ctx.NewWeek)
	}
	if ctx.CycleIteration != cycleIteration {
		t.Errorf("expected cycleIteration %d, got %d", cycleIteration, ctx.CycleIteration)
	}
}

// TestNewCycleTriggerEvent tests cycle trigger event factory.
func TestNewCycleTriggerEvent(t *testing.T) {
	userID := "user-123"
	completedCycle := 3
	totalWeeks := 4

	event := NewCycleTriggerEvent(userID, completedCycle, totalWeeks)

	if event.Type != TriggerAfterCycle {
		t.Errorf("expected type %s, got %s", TriggerAfterCycle, event.Type)
	}
	if event.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, event.UserID)
	}
	if event.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set")
	}

	ctx, ok := event.Context.(CycleTriggerContext)
	if !ok {
		t.Fatalf("expected CycleTriggerContext, got %T", event.Context)
	}
	if ctx.CompletedCycle != completedCycle {
		t.Errorf("expected completedCycle %d, got %d", completedCycle, ctx.CompletedCycle)
	}
	if ctx.NewCycle != completedCycle+1 {
		t.Errorf("expected newCycle %d, got %d", completedCycle+1, ctx.NewCycle)
	}
	if ctx.TotalWeeks != totalWeeks {
		t.Errorf("expected totalWeeks %d, got %d", totalWeeks, ctx.TotalWeeks)
	}
}

// TestMarshalTriggerContext tests MarshalTriggerContext with envelope.
func TestMarshalTriggerContext(t *testing.T) {
	tests := []struct {
		name string
		ctx  TriggerContext
	}{
		{
			name: "session context",
			ctx: SessionTriggerContext{
				SessionID:      "s-1",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{"l-1"},
			},
		},
		{
			name: "week context",
			ctx: WeekTriggerContext{
				PreviousWeek:   1,
				NewWeek:        2,
				CycleIteration: 1,
			},
		},
		{
			name: "cycle context",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       2,
				TotalWeeks:     4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalTriggerContext(tt.ctx)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			// Verify envelope structure
			var envelope TriggerContextEnvelope
			if err := json.Unmarshal(data, &envelope); err != nil {
				t.Fatalf("failed to parse envelope: %v", err)
			}
			if envelope.Type != tt.ctx.TriggerType() {
				t.Errorf("expected type %s, got %s", tt.ctx.TriggerType(), envelope.Type)
			}
			if len(envelope.Data) == 0 {
				t.Error("expected data to be present")
			}
		})
	}
}

// TestUnmarshalTriggerContextEnvelope tests UnmarshalTriggerContextEnvelope.
func TestUnmarshalTriggerContextEnvelope(t *testing.T) {
	tests := []struct {
		name         string
		ctx          TriggerContext
		expectedType TriggerType
	}{
		{
			name: "session context roundtrip",
			ctx: SessionTriggerContext{
				SessionID:      "s-1",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{"l-1"},
			},
			expectedType: TriggerAfterSession,
		},
		{
			name: "week context roundtrip",
			ctx: WeekTriggerContext{
				PreviousWeek:   1,
				NewWeek:        2,
				CycleIteration: 1,
			},
			expectedType: TriggerAfterWeek,
		},
		{
			name: "cycle context roundtrip",
			ctx: CycleTriggerContext{
				CompletedCycle: 1,
				NewCycle:       2,
				TotalWeeks:     4,
			},
			expectedType: TriggerAfterCycle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := MarshalTriggerContext(tt.ctx)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			// Unmarshal
			restored, err := UnmarshalTriggerContextEnvelope(data)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			if restored.TriggerType() != tt.expectedType {
				t.Errorf("TriggerType mismatch: expected %s, got %s",
					tt.expectedType, restored.TriggerType())
			}
		})
	}
}

// TestUnmarshalTriggerContextEnvelope_Invalid tests invalid envelope handling.
func TestUnmarshalTriggerContextEnvelope_Invalid(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "invalid json",
			data: `{invalid}`,
		},
		{
			name: "unknown type",
			data: `{"type":"UNKNOWN","data":{}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalTriggerContextEnvelope([]byte(tt.data))
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

// TestTriggerContext_Interface ensures all context types implement TriggerContext.
func TestTriggerContext_Interface(t *testing.T) {
	var _ TriggerContext = SessionTriggerContext{}
	var _ TriggerContext = WeekTriggerContext{}
	var _ TriggerContext = CycleTriggerContext{}
}
