package event

import (
	"context"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

func TestNewHandlerRegistry(t *testing.T) {
	bus := NewBus()
	registry := NewHandlerRegistry(bus)

	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
	if registry.Bus() != bus {
		t.Error("expected registry to wrap the same bus")
	}
}

func TestHandlerRegistry_RegisterHandler(t *testing.T) {
	bus := NewBus()
	registry := NewHandlerRegistry(bus)

	called := false
	registry.RegisterHandler(EventEnrolled, func(ctx context.Context, event StateEvent) error {
		called = true
		return nil
	})

	event := NewStateEvent(EventEnrolled, "user", "program")
	_ = bus.Publish(context.Background(), event)

	if !called {
		t.Error("expected handler to be called")
	}
}

func TestHandlerRegistry_RegisterMultiple(t *testing.T) {
	bus := NewBus()
	registry := NewHandlerRegistry(bus)

	callCount := 0
	registry.RegisterMultiple(
		[]EventType{EventEnrolled, EventQuit, EventSetLogged},
		func(ctx context.Context, event StateEvent) error {
			callCount++
			return nil
		},
	)

	_ = bus.Publish(context.Background(), NewStateEvent(EventEnrolled, "user", "program"))
	_ = bus.Publish(context.Background(), NewStateEvent(EventQuit, "user", "program"))
	_ = bus.Publish(context.Background(), NewStateEvent(EventSetLogged, "user", "program"))

	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestMapEventToTriggerType(t *testing.T) {
	testCases := []struct {
		eventType    EventType
		expectedType progression.TriggerType
	}{
		{EventSetLogged, progression.TriggerAfterSet},
		{EventWorkoutCompleted, progression.TriggerAfterSession},
		{EventWeekCompleted, progression.TriggerAfterWeek},
		{EventCycleCompleted, progression.TriggerAfterCycle},
		{EventEnrolled, ""},
		{EventQuit, ""},
		{EventWorkoutStarted, ""},
	}

	for _, tc := range testCases {
		t.Run(string(tc.eventType), func(t *testing.T) {
			got := MapEventToTriggerType(tc.eventType)
			if got != tc.expectedType {
				t.Errorf("expected %s, got %s", tc.expectedType, got)
			}
		})
	}
}

func TestGetTriggerTypeForSetEvent(t *testing.T) {
	// SET_LOGGED with isFailure=true should return ON_FAILURE
	failureEvent := NewStateEvent(EventSetLogged, "user", "program").
		WithPayload(PayloadIsFailure, true)
	if got := GetTriggerTypeForSetEvent(failureEvent); got != progression.TriggerOnFailure {
		t.Errorf("expected ON_FAILURE for failure set, got %s", got)
	}

	// SET_LOGGED with isFailure=false should return AFTER_SET
	successEvent := NewStateEvent(EventSetLogged, "user", "program").
		WithPayload(PayloadIsFailure, false)
	if got := GetTriggerTypeForSetEvent(successEvent); got != progression.TriggerAfterSet {
		t.Errorf("expected AFTER_SET for success set, got %s", got)
	}

	// SET_LOGGED without isFailure should return AFTER_SET (default)
	defaultEvent := NewStateEvent(EventSetLogged, "user", "program")
	if got := GetTriggerTypeForSetEvent(defaultEvent); got != progression.TriggerAfterSet {
		t.Errorf("expected AFTER_SET by default, got %s", got)
	}

	// Non-SET_LOGGED event should return empty
	otherEvent := NewStateEvent(EventEnrolled, "user", "program")
	if got := GetTriggerTypeForSetEvent(otherEvent); got != "" {
		t.Errorf("expected empty for non-SET_LOGGED event, got %s", got)
	}
}

func TestProgressionEventTypes(t *testing.T) {
	types := ProgressionEventTypes()

	expected := map[EventType]bool{
		EventSetLogged:        true,
		EventWorkoutCompleted: true,
		EventWeekCompleted:    true,
		EventCycleCompleted:   true,
	}

	if len(types) != len(expected) {
		t.Errorf("expected %d event types, got %d", len(expected), len(types))
	}

	for _, et := range types {
		if !expected[et] {
			t.Errorf("unexpected event type %s", et)
		}
	}
}

func TestProgressionTriggerAdapter_RegisterProgressionHandler(t *testing.T) {
	bus := NewBus()
	registry := NewHandlerRegistry(bus)
	adapter := NewProgressionTriggerAdapter(registry)

	receivedEvents := make([]struct {
		event       StateEvent
		triggerType progression.TriggerType
	}, 0)

	adapter.RegisterProgressionHandler(func(ctx context.Context, event StateEvent, triggerType progression.TriggerType) error {
		receivedEvents = append(receivedEvents, struct {
			event       StateEvent
			triggerType progression.TriggerType
		}{event, triggerType})
		return nil
	})

	// Test SET_LOGGED without failure
	_ = bus.Publish(context.Background(), NewStateEvent(EventSetLogged, "user", "program"))
	// Test SET_LOGGED with failure
	failureEvent := NewStateEvent(EventSetLogged, "user", "program").WithPayload(PayloadIsFailure, true)
	_ = bus.Publish(context.Background(), failureEvent)
	// Test WORKOUT_COMPLETED
	_ = bus.Publish(context.Background(), NewStateEvent(EventWorkoutCompleted, "user", "program"))
	// Test WEEK_COMPLETED
	_ = bus.Publish(context.Background(), NewStateEvent(EventWeekCompleted, "user", "program"))
	// Test CYCLE_COMPLETED
	_ = bus.Publish(context.Background(), NewStateEvent(EventCycleCompleted, "user", "program"))

	if len(receivedEvents) != 5 {
		t.Errorf("expected 5 events, got %d", len(receivedEvents))
	}

	// Verify trigger type mappings
	expectedTriggerTypes := []progression.TriggerType{
		progression.TriggerAfterSet,     // SET_LOGGED (no failure)
		progression.TriggerOnFailure,    // SET_LOGGED (failure)
		progression.TriggerAfterSession, // WORKOUT_COMPLETED
		progression.TriggerAfterWeek,    // WEEK_COMPLETED
		progression.TriggerAfterCycle,   // CYCLE_COMPLETED
	}

	for i, expected := range expectedTriggerTypes {
		if receivedEvents[i].triggerType != expected {
			t.Errorf("event %d: expected trigger type %s, got %s", i, expected, receivedEvents[i].triggerType)
		}
	}
}

func TestStateEventBuilder(t *testing.T) {
	event := NewEventBuilder(EventSetLogged, "user-123", "program-456").
		WithSession("session-789", "day-a", 2).
		WithLifts([]string{"lift-1", "lift-2"}).
		WithLoggedSet("set-001", "lift-1", 8, 5, 100.0, true, false).
		Build()

	if event.Type != EventSetLogged {
		t.Errorf("expected type %s, got %s", EventSetLogged, event.Type)
	}
	if event.UserID != "user-123" {
		t.Errorf("expected userID user-123, got %s", event.UserID)
	}
	if event.ProgramID != "program-456" {
		t.Errorf("expected programID program-456, got %s", event.ProgramID)
	}
	if event.GetString(PayloadSessionID) != "session-789" {
		t.Errorf("expected sessionId session-789, got %s", event.GetString(PayloadSessionID))
	}
	if event.GetString(PayloadDaySlug) != "day-a" {
		t.Errorf("expected daySlug day-a, got %s", event.GetString(PayloadDaySlug))
	}
	if event.GetInt(PayloadWeekNumber) != 2 {
		t.Errorf("expected weekNumber 2, got %d", event.GetInt(PayloadWeekNumber))
	}
	if event.GetString(PayloadLoggedSetID) != "set-001" {
		t.Errorf("expected loggedSetId set-001, got %s", event.GetString(PayloadLoggedSetID))
	}
	if event.GetInt(PayloadRepsPerformed) != 8 {
		t.Errorf("expected repsPerformed 8, got %d", event.GetInt(PayloadRepsPerformed))
	}
	if event.GetInt(PayloadTargetReps) != 5 {
		t.Errorf("expected targetReps 5, got %d", event.GetInt(PayloadTargetReps))
	}
	if event.GetFloat64(PayloadWeight) != 100.0 {
		t.Errorf("expected weight 100.0, got %f", event.GetFloat64(PayloadWeight))
	}
	if !event.GetBool(PayloadIsAMRAP) {
		t.Error("expected isAMRAP true")
	}
	if event.GetBool(PayloadIsFailure) {
		t.Error("expected isFailure false")
	}
}

func TestStateEventBuilder_WithWeekAdvancement(t *testing.T) {
	event := NewEventBuilder(EventWeekCompleted, "user", "program").
		WithWeekAdvancement(1, 2, 1).
		Build()

	if event.GetInt(PayloadPreviousWeek) != 1 {
		t.Errorf("expected previousWeek 1, got %d", event.GetInt(PayloadPreviousWeek))
	}
	if event.GetInt(PayloadNewWeek) != 2 {
		t.Errorf("expected newWeek 2, got %d", event.GetInt(PayloadNewWeek))
	}
	if event.GetInt(PayloadCycleIteration) != 1 {
		t.Errorf("expected cycleIteration 1, got %d", event.GetInt(PayloadCycleIteration))
	}
}

func TestStateEventBuilder_WithCycleAdvancement(t *testing.T) {
	event := NewEventBuilder(EventCycleCompleted, "user", "program").
		WithCycleAdvancement(1, 2, 4).
		Build()

	if event.GetInt(PayloadCompletedCycle) != 1 {
		t.Errorf("expected completedCycle 1, got %d", event.GetInt(PayloadCompletedCycle))
	}
	if event.GetInt(PayloadNewCycle) != 2 {
		t.Errorf("expected newCycle 2, got %d", event.GetInt(PayloadNewCycle))
	}
	if event.GetInt(PayloadTotalWeeks) != 4 {
		t.Errorf("expected totalWeeks 4, got %d", event.GetInt(PayloadTotalWeeks))
	}
}

func TestStateEventBuilder_WithPayload(t *testing.T) {
	event := NewEventBuilder(EventEnrolled, "user", "program").
		WithPayload("customKey", "customValue").
		Build()

	if event.GetString("customKey") != "customValue" {
		t.Errorf("expected customKey=customValue, got %s", event.GetString("customKey"))
	}
}
