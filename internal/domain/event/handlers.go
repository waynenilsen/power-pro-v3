package event

import (
	"context"

	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// HandlerRegistry manages the registration of event handlers and provides
// a central place to wire up event-driven integrations.
type HandlerRegistry struct {
	bus *Bus
}

// NewHandlerRegistry creates a new handler registry backed by the given event bus.
func NewHandlerRegistry(bus *Bus) *HandlerRegistry {
	return &HandlerRegistry{bus: bus}
}

// Bus returns the underlying event bus.
func (r *HandlerRegistry) Bus() *Bus {
	return r.bus
}

// RegisterHandler registers a handler function for a specific event type.
func (r *HandlerRegistry) RegisterHandler(eventType EventType, handler EventHandler) {
	r.bus.Subscribe(eventType, handler)
}

// RegisterMultiple registers a handler for multiple event types.
func (r *HandlerRegistry) RegisterMultiple(eventTypes []EventType, handler EventHandler) {
	r.bus.SubscribeMultiple(eventTypes, handler)
}

// ProgressionTriggerAdapter adapts event types to progression trigger types.
// This enables the progression system to subscribe to relevant events.
type ProgressionTriggerAdapter struct {
	registry *HandlerRegistry
}

// NewProgressionTriggerAdapter creates a new adapter for connecting events to progressions.
func NewProgressionTriggerAdapter(registry *HandlerRegistry) *ProgressionTriggerAdapter {
	return &ProgressionTriggerAdapter{registry: registry}
}

// MapEventToTriggerType maps an event type to its corresponding progression trigger type.
// Returns empty string if there's no mapping.
func MapEventToTriggerType(eventType EventType) progression.TriggerType {
	switch eventType {
	case EventSetLogged:
		// SET_LOGGED can map to both AFTER_SET and ON_FAILURE depending on payload
		// Default to AFTER_SET; caller should check IsFailure payload
		return progression.TriggerAfterSet
	case EventWorkoutCompleted:
		return progression.TriggerAfterSession
	case EventWeekCompleted:
		return progression.TriggerAfterWeek
	case EventCycleCompleted:
		return progression.TriggerAfterCycle
	default:
		return ""
	}
}

// GetTriggerTypeForSetEvent returns the appropriate trigger type for a SET_LOGGED event.
// If the set was a failure, returns ON_FAILURE; otherwise returns AFTER_SET.
func GetTriggerTypeForSetEvent(event StateEvent) progression.TriggerType {
	if event.Type != EventSetLogged {
		return ""
	}
	if event.GetBool(PayloadIsFailure) {
		return progression.TriggerOnFailure
	}
	return progression.TriggerAfterSet
}

// ProgressionEventTypes returns the event types that are relevant for progression triggers.
func ProgressionEventTypes() []EventType {
	return []EventType{
		EventSetLogged,
		EventWorkoutCompleted,
		EventWeekCompleted,
		EventCycleCompleted,
	}
}

// ProgressionHandler is a function that handles progression-related events.
// It receives the event and the mapped trigger type.
type ProgressionHandler func(ctx context.Context, event StateEvent, triggerType progression.TriggerType) error

// RegisterProgressionHandler registers a handler that will be called for all
// progression-relevant events. The handler receives both the event and the
// appropriate trigger type based on the event.
func (a *ProgressionTriggerAdapter) RegisterProgressionHandler(handler ProgressionHandler) {
	// Create a wrapper that maps events to trigger types
	wrapper := func(ctx context.Context, event StateEvent) error {
		var triggerType progression.TriggerType

		// Special handling for SET_LOGGED which can map to two trigger types
		if event.Type == EventSetLogged {
			triggerType = GetTriggerTypeForSetEvent(event)
		} else {
			triggerType = MapEventToTriggerType(event.Type)
		}

		if triggerType == "" {
			return nil // No mapping for this event type
		}

		return handler(ctx, event, triggerType)
	}

	// Subscribe to all progression-relevant events
	a.registry.RegisterMultiple(ProgressionEventTypes(), wrapper)
}

// StateEventBuilder provides a fluent API for building StateEvents with
// common payload patterns.
type StateEventBuilder struct {
	event StateEvent
}

// NewEventBuilder creates a new event builder for the given event type.
func NewEventBuilder(eventType EventType, userID, programID string) *StateEventBuilder {
	return &StateEventBuilder{
		event: NewStateEvent(eventType, userID, programID),
	}
}

// WithSession adds session-related payload fields.
func (b *StateEventBuilder) WithSession(sessionID, daySlug string, weekNumber int) *StateEventBuilder {
	b.event = b.event.
		WithPayload(PayloadSessionID, sessionID).
		WithPayload(PayloadDaySlug, daySlug).
		WithPayload(PayloadWeekNumber, weekNumber)
	return b
}

// WithLifts adds the lifts performed to the payload.
func (b *StateEventBuilder) WithLifts(liftIDs []string) *StateEventBuilder {
	b.event = b.event.WithPayload(PayloadLiftsPerformed, liftIDs)
	return b
}

// WithWeekAdvancement adds week advancement payload fields.
func (b *StateEventBuilder) WithWeekAdvancement(previousWeek, newWeek, cycleIteration int) *StateEventBuilder {
	b.event = b.event.
		WithPayload(PayloadPreviousWeek, previousWeek).
		WithPayload(PayloadNewWeek, newWeek).
		WithPayload(PayloadCycleIteration, cycleIteration)
	return b
}

// WithCycleAdvancement adds cycle advancement payload fields.
func (b *StateEventBuilder) WithCycleAdvancement(completedCycle, newCycle, totalWeeks int) *StateEventBuilder {
	b.event = b.event.
		WithPayload(PayloadCompletedCycle, completedCycle).
		WithPayload(PayloadNewCycle, newCycle).
		WithPayload(PayloadTotalWeeks, totalWeeks)
	return b
}

// WithLoggedSet adds logged set payload fields.
func (b *StateEventBuilder) WithLoggedSet(loggedSetID, liftID string, repsPerformed, targetReps int, weight float64, isAMRAP, isFailure bool) *StateEventBuilder {
	b.event = b.event.
		WithPayload(PayloadLoggedSetID, loggedSetID).
		WithPayload(PayloadLiftID, liftID).
		WithPayload(PayloadRepsPerformed, repsPerformed).
		WithPayload(PayloadTargetReps, targetReps).
		WithPayload(PayloadWeight, weight).
		WithPayload(PayloadIsAMRAP, isAMRAP).
		WithPayload(PayloadIsFailure, isFailure)
	return b
}

// WithPayload adds a custom payload field.
func (b *StateEventBuilder) WithPayload(key string, value interface{}) *StateEventBuilder {
	b.event = b.event.WithPayload(key, value)
	return b
}

// Build returns the constructed event.
func (b *StateEventBuilder) Build() StateEvent {
	return b.event
}
