// Package progression provides domain logic for progression strategies.
// This file defines trigger types and event structures that control when progressions fire.
package progression

import (
	"encoding/json"
	"fmt"
	"time"
)

// TriggerContext is the interface for type-specific trigger context.
// Each trigger type has its own context carrying relevant information about the event.
// Implementations must be JSON serializable for ProgressionLog storage.
type TriggerContext interface {
	// TriggerType returns the trigger type this context belongs to.
	TriggerType() TriggerType
	// Validate validates the context data.
	Validate() error
}

// SessionTriggerContext contains context for AFTER_SESSION triggers.
// This context is passed when a user completes a training day.
type SessionTriggerContext struct {
	// SessionID is the UUID of the completed session/workout.
	SessionID string `json:"sessionId"`
	// DaySlug identifies which day template was completed (e.g., "day-a", "heavy-day").
	DaySlug string `json:"daySlug"`
	// WeekNumber is the week number within the cycle when the session was completed.
	WeekNumber int `json:"weekNumber"`
	// LiftsPerformed contains the UUIDs of lifts performed in this session.
	// This enables per-lift progression decisions (Starting Strength progresses each lift independently).
	LiftsPerformed []string `json:"liftsPerformed"`
}

// TriggerType implements TriggerContext.
func (c SessionTriggerContext) TriggerType() TriggerType {
	return TriggerAfterSession
}

// Validate implements TriggerContext.
func (c SessionTriggerContext) Validate() error {
	if c.SessionID == "" {
		return fmt.Errorf("%w: sessionId is required for session trigger context", ErrInvalidParams)
	}
	if c.DaySlug == "" {
		return fmt.Errorf("%w: daySlug is required for session trigger context", ErrInvalidParams)
	}
	if c.WeekNumber < 1 {
		return fmt.Errorf("%w: weekNumber must be at least 1", ErrInvalidParams)
	}
	return nil
}

// WeekTriggerContext contains context for AFTER_WEEK triggers.
// This context is passed when a user advances from week N to week N+1.
type WeekTriggerContext struct {
	// PreviousWeek is the week number that was just completed.
	PreviousWeek int `json:"previousWeek"`
	// NewWeek is the week number being advanced to.
	NewWeek int `json:"newWeek"`
	// CycleIteration tracks which cycle the user is on (1-indexed).
	CycleIteration int `json:"cycleIteration"`
}

// TriggerType implements TriggerContext.
func (c WeekTriggerContext) TriggerType() TriggerType {
	return TriggerAfterWeek
}

// Validate implements TriggerContext.
func (c WeekTriggerContext) Validate() error {
	if c.PreviousWeek < 1 {
		return fmt.Errorf("%w: previousWeek must be at least 1", ErrInvalidParams)
	}
	if c.NewWeek < 1 {
		return fmt.Errorf("%w: newWeek must be at least 1", ErrInvalidParams)
	}
	if c.NewWeek <= c.PreviousWeek {
		return fmt.Errorf("%w: newWeek must be greater than previousWeek for week advancement", ErrInvalidParams)
	}
	if c.CycleIteration < 1 {
		return fmt.Errorf("%w: cycleIteration must be at least 1", ErrInvalidParams)
	}
	return nil
}

// CycleTriggerContext contains context for AFTER_CYCLE triggers.
// This context is passed when a user completes a cycle (wraps from week N to week 1).
type CycleTriggerContext struct {
	// CompletedCycle is the cycle iteration that was just completed (1-indexed).
	CompletedCycle int `json:"completedCycle"`
	// NewCycle is the cycle iteration being started (CompletedCycle + 1).
	NewCycle int `json:"newCycle"`
	// TotalWeeks is the number of weeks in the completed cycle.
	TotalWeeks int `json:"totalWeeks"`
}

// TriggerType implements TriggerContext.
func (c CycleTriggerContext) TriggerType() TriggerType {
	return TriggerAfterCycle
}

// Validate implements TriggerContext.
func (c CycleTriggerContext) Validate() error {
	if c.CompletedCycle < 1 {
		return fmt.Errorf("%w: completedCycle must be at least 1", ErrInvalidParams)
	}
	if c.NewCycle < 1 {
		return fmt.Errorf("%w: newCycle must be at least 1", ErrInvalidParams)
	}
	if c.NewCycle != c.CompletedCycle+1 {
		return fmt.Errorf("%w: newCycle must be completedCycle + 1", ErrInvalidParams)
	}
	if c.TotalWeeks < 1 {
		return fmt.Errorf("%w: totalWeeks must be at least 1", ErrInvalidParams)
	}
	return nil
}

// FailureTriggerContext contains context for ON_FAILURE triggers.
// This context is passed when a user fails to meet target reps on a set.
type FailureTriggerContext struct {
	// LoggedSetID is the UUID of the logged set that triggered the failure.
	LoggedSetID string `json:"loggedSetId"`
	// LiftID is the UUID of the lift for the failed set.
	LiftID string `json:"liftId"`
	// TargetReps is the number of reps that were prescribed.
	TargetReps int `json:"targetReps"`
	// RepsPerformed is the number of reps actually achieved.
	RepsPerformed int `json:"repsPerformed"`
	// RepsDifference is RepsPerformed - TargetReps (always negative for failures).
	RepsDifference int `json:"repsDifference"`
	// ConsecutiveFailures is the current count of consecutive failures for this lift/progression.
	ConsecutiveFailures int `json:"consecutiveFailures"`
	// Weight is the weight used for the failed set.
	Weight float64 `json:"weight"`
	// ProgressionID is the UUID of the progression this failure counter is tracking.
	ProgressionID string `json:"progressionId"`
}

// TriggerType implements TriggerContext.
func (c FailureTriggerContext) TriggerType() TriggerType {
	return TriggerOnFailure
}

// Validate implements TriggerContext.
func (c FailureTriggerContext) Validate() error {
	if c.LoggedSetID == "" {
		return fmt.Errorf("%w: loggedSetId is required for failure trigger context", ErrInvalidParams)
	}
	if c.LiftID == "" {
		return fmt.Errorf("%w: liftId is required for failure trigger context", ErrInvalidParams)
	}
	if c.TargetReps < 1 {
		return fmt.Errorf("%w: targetReps must be at least 1", ErrInvalidParams)
	}
	if c.RepsPerformed < 0 {
		return fmt.Errorf("%w: repsPerformed must be non-negative", ErrInvalidParams)
	}
	if c.RepsPerformed >= c.TargetReps {
		return fmt.Errorf("%w: repsPerformed must be less than targetReps for a failure", ErrInvalidParams)
	}
	if c.ConsecutiveFailures < 1 {
		return fmt.Errorf("%w: consecutiveFailures must be at least 1 for a failure trigger", ErrInvalidParams)
	}
	if c.ProgressionID == "" {
		return fmt.Errorf("%w: progressionId is required for failure trigger context", ErrInvalidParams)
	}
	return nil
}

// TriggerEventV2 contains all parameters for a trigger event.
// This is the new trigger event structure with strongly-typed context.
// The "V2" suffix distinguishes it from the existing flat TriggerEvent during migration.
type TriggerEventV2 struct {
	// Type is the type of trigger that fired.
	Type TriggerType `json:"type"`
	// UserID is the UUID of the user who triggered the event.
	UserID string `json:"userId"`
	// Timestamp is when the trigger event occurred.
	Timestamp time.Time `json:"timestamp"`
	// Context contains type-specific context for the trigger.
	// Use SessionTriggerContext, WeekTriggerContext, or CycleTriggerContext based on Type.
	Context TriggerContext `json:"-"`
	// RawContext holds the JSON representation of Context for serialization.
	RawContext json.RawMessage `json:"context,omitempty"`
}

// Validate validates the TriggerEventV2.
func (e *TriggerEventV2) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("%w: trigger type is required", ErrInvalidParams)
	}
	if !ValidTriggerTypes[e.Type] {
		return fmt.Errorf("%w: %s", ErrUnknownTriggerType, e.Type)
	}
	if e.UserID == "" {
		return fmt.Errorf("%w: userId is required for trigger event", ErrInvalidParams)
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("%w: trigger timestamp is required", ErrInvalidParams)
	}
	if e.Context == nil {
		return fmt.Errorf("%w: trigger context is required", ErrInvalidParams)
	}
	// Verify context type matches trigger type
	if e.Context.TriggerType() != e.Type {
		return fmt.Errorf("%w: context type %s does not match trigger type %s",
			ErrInvalidParams, e.Context.TriggerType(), e.Type)
	}
	return e.Context.Validate()
}

// MarshalJSON implements json.Marshaler for TriggerEventV2.
// This serializes the Context field to RawContext for storage.
func (e TriggerEventV2) MarshalJSON() ([]byte, error) {
	type Alias TriggerEventV2
	alias := Alias(e)

	// Marshal context to RawContext
	if e.Context != nil {
		contextData, err := json.Marshal(e.Context)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal trigger context: %w", err)
		}
		alias.RawContext = contextData
	}

	return json.Marshal(alias)
}

// UnmarshalJSON implements json.Unmarshaler for TriggerEventV2.
// This deserializes RawContext to the appropriate Context type based on Type.
func (e *TriggerEventV2) UnmarshalJSON(data []byte) error {
	type Alias TriggerEventV2
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal trigger event: %w", err)
	}

	e.Type = alias.Type
	e.UserID = alias.UserID
	e.Timestamp = alias.Timestamp
	e.RawContext = alias.RawContext

	// Deserialize context based on type
	if len(alias.RawContext) > 0 {
		ctx, err := UnmarshalTriggerContext(alias.Type, alias.RawContext)
		if err != nil {
			return err
		}
		e.Context = ctx
	}

	return nil
}

// UnmarshalTriggerContext deserializes a TriggerContext from JSON based on the trigger type.
func UnmarshalTriggerContext(triggerType TriggerType, data json.RawMessage) (TriggerContext, error) {
	switch triggerType {
	case TriggerAfterSession:
		var ctx SessionTriggerContext
		if err := json.Unmarshal(data, &ctx); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session trigger context: %w", err)
		}
		return ctx, nil
	case TriggerAfterWeek:
		var ctx WeekTriggerContext
		if err := json.Unmarshal(data, &ctx); err != nil {
			return nil, fmt.Errorf("failed to unmarshal week trigger context: %w", err)
		}
		return ctx, nil
	case TriggerAfterCycle:
		var ctx CycleTriggerContext
		if err := json.Unmarshal(data, &ctx); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cycle trigger context: %w", err)
		}
		return ctx, nil
	case TriggerOnFailure:
		var ctx FailureTriggerContext
		if err := json.Unmarshal(data, &ctx); err != nil {
			return nil, fmt.Errorf("failed to unmarshal failure trigger context: %w", err)
		}
		return ctx, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownTriggerType, triggerType)
	}
}

// NewSessionTriggerEvent creates a new AFTER_SESSION trigger event.
func NewSessionTriggerEvent(userID, sessionID, daySlug string, weekNumber int, liftsPerformed []string) *TriggerEventV2 {
	return &TriggerEventV2{
		Type:      TriggerAfterSession,
		UserID:    userID,
		Timestamp: time.Now(),
		Context: SessionTriggerContext{
			SessionID:      sessionID,
			DaySlug:        daySlug,
			WeekNumber:     weekNumber,
			LiftsPerformed: liftsPerformed,
		},
	}
}

// NewWeekTriggerEvent creates a new AFTER_WEEK trigger event.
func NewWeekTriggerEvent(userID string, previousWeek, newWeek, cycleIteration int) *TriggerEventV2 {
	return &TriggerEventV2{
		Type:      TriggerAfterWeek,
		UserID:    userID,
		Timestamp: time.Now(),
		Context: WeekTriggerContext{
			PreviousWeek:   previousWeek,
			NewWeek:        newWeek,
			CycleIteration: cycleIteration,
		},
	}
}

// NewCycleTriggerEvent creates a new AFTER_CYCLE trigger event.
func NewCycleTriggerEvent(userID string, completedCycle, totalWeeks int) *TriggerEventV2 {
	return &TriggerEventV2{
		Type:      TriggerAfterCycle,
		UserID:    userID,
		Timestamp: time.Now(),
		Context: CycleTriggerContext{
			CompletedCycle: completedCycle,
			NewCycle:       completedCycle + 1,
			TotalWeeks:     totalWeeks,
		},
	}
}

// NewFailureTriggerEvent creates a new ON_FAILURE trigger event.
func NewFailureTriggerEvent(userID string, ctx FailureTriggerContext) *TriggerEventV2 {
	return &TriggerEventV2{
		Type:      TriggerOnFailure,
		UserID:    userID,
		Timestamp: time.Now(),
		Context:   ctx,
	}
}

// TriggerContextEnvelope is a wrapper for polymorphic TriggerContext serialization.
// It includes the trigger type to enable proper deserialization.
type TriggerContextEnvelope struct {
	// Type identifies which trigger context type is stored.
	Type TriggerType `json:"type"`
	// Data is the JSON-encoded context.
	Data json.RawMessage `json:"data"`
}

// MarshalTriggerContext serializes a TriggerContext with its type discriminator.
// This format is suitable for ProgressionLog storage.
func MarshalTriggerContext(ctx TriggerContext) ([]byte, error) {
	data, err := json.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trigger context: %w", err)
	}
	envelope := TriggerContextEnvelope{
		Type: ctx.TriggerType(),
		Data: data,
	}
	return json.Marshal(envelope)
}

// UnmarshalTriggerContextEnvelope deserializes a TriggerContext from an envelope.
func UnmarshalTriggerContextEnvelope(data []byte) (TriggerContext, error) {
	var envelope TriggerContextEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trigger context envelope: %w", err)
	}
	return UnmarshalTriggerContext(envelope.Type, envelope.Data)
}

// ManualTriggerContext wraps another TriggerContext and adds manual trigger metadata.
// This is used when progressions are triggered manually via the API.
type ManualTriggerContext struct {
	// Manual indicates this was a manual trigger (always true for this type).
	Manual bool `json:"manual"`
	// Force indicates the idempotency check was bypassed.
	Force bool `json:"force"`
	// LiftID is the specific lift targeted, if any.
	LiftID string `json:"liftId,omitempty"`
	// UnderlyingContext contains the synthetic trigger context.
	UnderlyingContext TriggerContext `json:"-"`
	// InnerContext holds the serialized underlying context.
	InnerContext json.RawMessage `json:"context,omitempty"`
	// InnerTriggerType indicates what type of trigger was synthesized.
	InnerTriggerType TriggerType `json:"triggerType"`
}

// TriggerType implements TriggerContext.
// Returns the underlying trigger type for compatibility with progression logic.
func (c ManualTriggerContext) TriggerType() TriggerType {
	return c.InnerTriggerType
}

// Validate implements TriggerContext.
func (c ManualTriggerContext) Validate() error {
	// ManualTriggerContext is always valid since it's synthesized by the system
	return nil
}

// MarshalJSON implements json.Marshaler for ManualTriggerContext.
func (c ManualTriggerContext) MarshalJSON() ([]byte, error) {
	type Alias ManualTriggerContext
	alias := Alias(c)

	// Marshal the underlying context
	if c.UnderlyingContext != nil {
		contextData, err := json.Marshal(c.UnderlyingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal underlying context: %w", err)
		}
		alias.InnerContext = contextData
		alias.InnerTriggerType = c.UnderlyingContext.TriggerType()
	}

	return json.Marshal(alias)
}

// NewManualTriggerContext creates a new ManualTriggerContext wrapping the given context.
func NewManualTriggerContext(underlyingContext TriggerContext, liftID string, force bool) *ManualTriggerContext {
	return &ManualTriggerContext{
		Manual:            true,
		Force:             force,
		LiftID:            liftID,
		UnderlyingContext: underlyingContext,
		InnerTriggerType:  underlyingContext.TriggerType(),
	}
}
