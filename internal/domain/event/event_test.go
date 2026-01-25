package event

import (
	"testing"
	"time"
)

func TestNewStateEvent(t *testing.T) {
	before := time.Now()
	event := NewStateEvent(EventEnrolled, "user-123", "program-456")
	after := time.Now()

	if event.Type != EventEnrolled {
		t.Errorf("expected type %s, got %s", EventEnrolled, event.Type)
	}
	if event.UserID != "user-123" {
		t.Errorf("expected userID user-123, got %s", event.UserID)
	}
	if event.ProgramID != "program-456" {
		t.Errorf("expected programID program-456, got %s", event.ProgramID)
	}
	if event.Timestamp.Before(before) || event.Timestamp.After(after) {
		t.Error("timestamp should be between before and after test execution")
	}
	if event.Payload == nil {
		t.Error("payload should be initialized")
	}
}

func TestStateEvent_WithPayload(t *testing.T) {
	event := NewStateEvent(EventSetLogged, "user-123", "program-456").
		WithPayload("key1", "value1").
		WithPayload("key2", 42)

	if event.Payload["key1"] != "value1" {
		t.Errorf("expected key1=value1, got %v", event.Payload["key1"])
	}
	if event.Payload["key2"] != 42 {
		t.Errorf("expected key2=42, got %v", event.Payload["key2"])
	}
}

func TestStateEvent_GetString(t *testing.T) {
	event := NewStateEvent(EventSetLogged, "user", "program").
		WithPayload("strKey", "hello").
		WithPayload("intKey", 123)

	if got := event.GetString("strKey"); got != "hello" {
		t.Errorf("expected 'hello', got '%s'", got)
	}
	if got := event.GetString("intKey"); got != "" {
		t.Errorf("expected empty string for int value, got '%s'", got)
	}
	if got := event.GetString("missing"); got != "" {
		t.Errorf("expected empty string for missing key, got '%s'", got)
	}

	// Test with nil payload
	emptyEvent := StateEvent{}
	if got := emptyEvent.GetString("any"); got != "" {
		t.Errorf("expected empty string for nil payload, got '%s'", got)
	}
}

func TestStateEvent_GetInt(t *testing.T) {
	event := NewStateEvent(EventSetLogged, "user", "program").
		WithPayload("intKey", 42).
		WithPayload("strKey", "notanint")

	if got := event.GetInt("intKey"); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
	if got := event.GetInt("strKey"); got != 0 {
		t.Errorf("expected 0 for string value, got %d", got)
	}
	if got := event.GetInt("missing"); got != 0 {
		t.Errorf("expected 0 for missing key, got %d", got)
	}

	// Test with nil payload
	emptyEvent := StateEvent{}
	if got := emptyEvent.GetInt("any"); got != 0 {
		t.Errorf("expected 0 for nil payload, got %d", got)
	}
}

func TestStateEvent_GetFloat64(t *testing.T) {
	event := NewStateEvent(EventSetLogged, "user", "program").
		WithPayload("floatKey", 3.14).
		WithPayload("strKey", "notafloat")

	if got := event.GetFloat64("floatKey"); got != 3.14 {
		t.Errorf("expected 3.14, got %f", got)
	}
	if got := event.GetFloat64("strKey"); got != 0.0 {
		t.Errorf("expected 0.0 for string value, got %f", got)
	}
	if got := event.GetFloat64("missing"); got != 0.0 {
		t.Errorf("expected 0.0 for missing key, got %f", got)
	}

	// Test with nil payload
	emptyEvent := StateEvent{}
	if got := emptyEvent.GetFloat64("any"); got != 0.0 {
		t.Errorf("expected 0.0 for nil payload, got %f", got)
	}
}

func TestStateEvent_GetBool(t *testing.T) {
	event := NewStateEvent(EventSetLogged, "user", "program").
		WithPayload("boolTrue", true).
		WithPayload("boolFalse", false).
		WithPayload("strKey", "notabool")

	if got := event.GetBool("boolTrue"); !got {
		t.Error("expected true, got false")
	}
	if got := event.GetBool("boolFalse"); got {
		t.Error("expected false, got true")
	}
	if got := event.GetBool("strKey"); got {
		t.Error("expected false for string value, got true")
	}
	if got := event.GetBool("missing"); got {
		t.Error("expected false for missing key, got true")
	}

	// Test with nil payload
	emptyEvent := StateEvent{}
	if got := emptyEvent.GetBool("any"); got {
		t.Error("expected false for nil payload, got true")
	}
}

func TestValidEventTypes(t *testing.T) {
	expectedTypes := []EventType{
		EventEnrolled,
		EventCycleBoundaryReached,
		EventQuit,
		EventCycleStarted,
		EventCycleCompleted,
		EventWeekStarted,
		EventWeekCompleted,
		EventWorkoutStarted,
		EventWorkoutCompleted,
		EventWorkoutAbandoned,
		EventSetLogged,
	}

	for _, et := range expectedTypes {
		if !ValidEventTypes[et] {
			t.Errorf("expected %s to be a valid event type", et)
		}
	}

	if ValidEventTypes["INVALID_EVENT_TYPE"] {
		t.Error("INVALID_EVENT_TYPE should not be valid")
	}
}

func TestStateEvent_WithPayload_NilPayload(t *testing.T) {
	// Start with an event that has nil payload
	event := StateEvent{
		Type:   EventSetLogged,
		UserID: "user",
	}

	// WithPayload should initialize the map
	event = event.WithPayload("key", "value")

	if event.Payload == nil {
		t.Error("payload should be initialized after WithPayload")
	}
	if event.Payload["key"] != "value" {
		t.Errorf("expected payload[key]=value, got %v", event.Payload["key"])
	}
}
