// Package event provides an in-memory event bus for decoupling state transitions
// from progression triggers and other side effects.
package event

import "time"

// EventType identifies the type of state event that occurred.
type EventType string

const (
	// EventEnrolled fires when a user enrolls in a program.
	EventEnrolled EventType = "ENROLLED"
	// EventCycleBoundaryReached fires when a user reaches the end of a cycle.
	EventCycleBoundaryReached EventType = "CYCLE_BOUNDARY_REACHED"
	// EventQuit fires when a user quits a program.
	EventQuit EventType = "QUIT"
	// EventCycleStarted fires when a user starts a new cycle.
	EventCycleStarted EventType = "CYCLE_STARTED"
	// EventCycleCompleted fires when a user completes a cycle.
	EventCycleCompleted EventType = "CYCLE_COMPLETED"
	// EventWeekStarted fires when a user starts a new week.
	EventWeekStarted EventType = "WEEK_STARTED"
	// EventWeekCompleted fires when a user completes a week.
	EventWeekCompleted EventType = "WEEK_COMPLETED"
	// EventWorkoutStarted fires when a user starts a workout session.
	EventWorkoutStarted EventType = "WORKOUT_STARTED"
	// EventWorkoutCompleted fires when a user completes a workout session.
	EventWorkoutCompleted EventType = "WORKOUT_COMPLETED"
	// EventWorkoutAbandoned fires when a user abandons a workout session.
	EventWorkoutAbandoned EventType = "WORKOUT_ABANDONED"
	// EventSetLogged fires when a user logs a set.
	EventSetLogged EventType = "SET_LOGGED"
)

// ValidEventTypes contains all valid event types for validation.
var ValidEventTypes = map[EventType]bool{
	EventEnrolled:             true,
	EventCycleBoundaryReached: true,
	EventQuit:                 true,
	EventCycleStarted:         true,
	EventCycleCompleted:       true,
	EventWeekStarted:          true,
	EventWeekCompleted:        true,
	EventWorkoutStarted:       true,
	EventWorkoutCompleted:     true,
	EventWorkoutAbandoned:     true,
	EventSetLogged:            true,
}

// StateEvent represents an event that occurred during a state transition.
// Events carry contextual information about what changed.
type StateEvent struct {
	// Type identifies the kind of event.
	Type EventType
	// UserID is the UUID of the user who triggered the event.
	UserID string
	// ProgramID is the UUID of the program associated with the event.
	ProgramID string
	// Timestamp is when the event occurred.
	Timestamp time.Time
	// Payload contains event-specific data.
	// Keys and values depend on the event type.
	Payload map[string]interface{}
}

// NewStateEvent creates a new StateEvent with the given type, user ID, and program ID.
// The timestamp is set to the current time.
func NewStateEvent(eventType EventType, userID, programID string) StateEvent {
	return StateEvent{
		Type:      eventType,
		UserID:    userID,
		ProgramID: programID,
		Timestamp: time.Now(),
		Payload:   make(map[string]interface{}),
	}
}

// WithPayload adds payload data to the event and returns the event for chaining.
func (e StateEvent) WithPayload(key string, value interface{}) StateEvent {
	if e.Payload == nil {
		e.Payload = make(map[string]interface{})
	}
	e.Payload[key] = value
	return e
}

// GetString retrieves a string value from the payload.
// Returns empty string if the key doesn't exist or isn't a string.
func (e StateEvent) GetString(key string) string {
	if e.Payload == nil {
		return ""
	}
	if v, ok := e.Payload[key].(string); ok {
		return v
	}
	return ""
}

// GetInt retrieves an int value from the payload.
// Returns 0 if the key doesn't exist or isn't an int.
func (e StateEvent) GetInt(key string) int {
	if e.Payload == nil {
		return 0
	}
	if v, ok := e.Payload[key].(int); ok {
		return v
	}
	return 0
}

// GetFloat64 retrieves a float64 value from the payload.
// Returns 0.0 if the key doesn't exist or isn't a float64.
func (e StateEvent) GetFloat64(key string) float64 {
	if e.Payload == nil {
		return 0.0
	}
	if v, ok := e.Payload[key].(float64); ok {
		return v
	}
	return 0.0
}

// GetBool retrieves a bool value from the payload.
// Returns false if the key doesn't exist or isn't a bool.
func (e StateEvent) GetBool(key string) bool {
	if e.Payload == nil {
		return false
	}
	if v, ok := e.Payload[key].(bool); ok {
		return v
	}
	return false
}

// Payload keys for common event data.
const (
	// PayloadSessionID is the key for workout session ID.
	PayloadSessionID = "sessionId"
	// PayloadDaySlug is the key for the day template slug.
	PayloadDaySlug = "daySlug"
	// PayloadWeekNumber is the key for the week number.
	PayloadWeekNumber = "weekNumber"
	// PayloadCycleIteration is the key for the cycle iteration.
	PayloadCycleIteration = "cycleIteration"
	// PayloadLiftsPerformed is the key for lift IDs performed in a session.
	PayloadLiftsPerformed = "liftsPerformed"
	// PayloadPreviousWeek is the key for the previous week number.
	PayloadPreviousWeek = "previousWeek"
	// PayloadNewWeek is the key for the new week number.
	PayloadNewWeek = "newWeek"
	// PayloadCompletedCycle is the key for the completed cycle number.
	PayloadCompletedCycle = "completedCycle"
	// PayloadNewCycle is the key for the new cycle number.
	PayloadNewCycle = "newCycle"
	// PayloadTotalWeeks is the key for total weeks in a cycle.
	PayloadTotalWeeks = "totalWeeks"
	// PayloadLoggedSetID is the key for logged set ID.
	PayloadLoggedSetID = "loggedSetId"
	// PayloadLiftID is the key for lift ID.
	PayloadLiftID = "liftId"
	// PayloadRepsPerformed is the key for reps performed.
	PayloadRepsPerformed = "repsPerformed"
	// PayloadTargetReps is the key for target reps.
	PayloadTargetReps = "targetReps"
	// PayloadWeight is the key for weight.
	PayloadWeight = "weight"
	// PayloadIsAMRAP is the key for whether a set is AMRAP.
	PayloadIsAMRAP = "isAMRAP"
	// PayloadIsFailure is the key for whether a set was a failure.
	PayloadIsFailure = "isFailure"
	// PayloadEnrolledAt is the key for when the user enrolled.
	PayloadEnrolledAt = "enrolledAt"
	// PayloadCyclesCompleted is the key for total cycles completed.
	PayloadCyclesCompleted = "cyclesCompleted"
	// PayloadWeeksCompleted is the key for total weeks completed.
	PayloadWeeksCompleted = "weeksCompleted"
)
