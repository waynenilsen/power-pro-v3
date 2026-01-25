package event

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBus(t *testing.T) {
	bus := NewBus()
	if bus == nil {
		t.Fatal("expected non-nil bus")
	}
	if bus.handlers == nil {
		t.Error("handlers map should be initialized")
	}
}

func TestBus_Subscribe(t *testing.T) {
	bus := NewBus()

	handler := func(ctx context.Context, event StateEvent) error {
		return nil
	}

	bus.Subscribe(EventEnrolled, handler)

	if !bus.HasSubscribers(EventEnrolled) {
		t.Error("expected subscribers for EventEnrolled")
	}
	if bus.SubscriberCount(EventEnrolled) != 1 {
		t.Errorf("expected 1 subscriber, got %d", bus.SubscriberCount(EventEnrolled))
	}
}

func TestBus_SubscribeMultiple(t *testing.T) {
	bus := NewBus()

	callCount := 0
	handler := func(ctx context.Context, event StateEvent) error {
		callCount++
		return nil
	}

	eventTypes := []EventType{EventEnrolled, EventQuit, EventSetLogged}
	bus.SubscribeMultiple(eventTypes, handler)

	for _, et := range eventTypes {
		if !bus.HasSubscribers(et) {
			t.Errorf("expected subscribers for %s", et)
		}
	}
}

func TestBus_Publish(t *testing.T) {
	bus := NewBus()

	var receivedEvent StateEvent
	handler := func(ctx context.Context, event StateEvent) error {
		receivedEvent = event
		return nil
	}

	bus.Subscribe(EventEnrolled, handler)

	event := NewStateEvent(EventEnrolled, "user-123", "program-456")
	err := bus.Publish(context.Background(), event)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if receivedEvent.UserID != "user-123" {
		t.Errorf("expected userID user-123, got %s", receivedEvent.UserID)
	}
}

func TestBus_Publish_NoSubscribers(t *testing.T) {
	bus := NewBus()

	event := NewStateEvent(EventEnrolled, "user-123", "program-456")
	err := bus.Publish(context.Background(), event)

	if err != nil {
		t.Errorf("expected no error for publish with no subscribers, got %v", err)
	}
}

func TestBus_Publish_MultipleHandlers(t *testing.T) {
	bus := NewBus()

	var order []int
	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		idx := i
		bus.Subscribe(EventEnrolled, func(ctx context.Context, event StateEvent) error {
			mu.Lock()
			order = append(order, idx)
			mu.Unlock()
			return nil
		})
	}

	event := NewStateEvent(EventEnrolled, "user", "program")
	err := bus.Publish(context.Background(), event)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 3 {
		t.Errorf("expected 3 handlers to be called, got %d", len(order))
	}
	// Handlers should be called in order
	for i, v := range order {
		if v != i {
			t.Errorf("expected handler %d to be called at position %d, got %d", i, i, v)
		}
	}
}

func TestBus_Publish_HandlerError(t *testing.T) {
	bus := NewBus()

	expectedErr := errors.New("handler error")
	var callCount int

	bus.Subscribe(EventEnrolled, func(ctx context.Context, event StateEvent) error {
		callCount++
		return expectedErr
	})
	bus.Subscribe(EventEnrolled, func(ctx context.Context, event StateEvent) error {
		callCount++
		return nil
	})
	bus.Subscribe(EventEnrolled, func(ctx context.Context, event StateEvent) error {
		callCount++
		return errors.New("second error")
	})

	event := NewStateEvent(EventEnrolled, "user", "program")
	err := bus.Publish(context.Background(), event)

	// Should return first error but call all handlers
	if err != expectedErr {
		t.Errorf("expected first error to be returned, got %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected all 3 handlers to be called, got %d", callCount)
	}
}

func TestBus_PublishAsync(t *testing.T) {
	bus := NewBus()

	var called atomic.Bool
	bus.Subscribe(EventEnrolled, func(ctx context.Context, event StateEvent) error {
		time.Sleep(10 * time.Millisecond)
		called.Store(true)
		return nil
	})

	event := NewStateEvent(EventEnrolled, "user", "program")
	bus.PublishAsync(context.Background(), event)

	// Should return immediately
	if called.Load() {
		t.Error("handler should not have completed yet")
	}

	// Wait for handler to complete
	time.Sleep(50 * time.Millisecond)
	if !called.Load() {
		t.Error("handler should have completed")
	}
}

func TestBus_PublishAll(t *testing.T) {
	bus := NewBus()

	var receivedEvents []EventType
	var mu sync.Mutex

	handler := func(ctx context.Context, event StateEvent) error {
		mu.Lock()
		receivedEvents = append(receivedEvents, event.Type)
		mu.Unlock()
		return nil
	}

	bus.Subscribe(EventEnrolled, handler)
	bus.Subscribe(EventQuit, handler)
	bus.Subscribe(EventSetLogged, handler)

	events := []StateEvent{
		NewStateEvent(EventEnrolled, "user", "program"),
		NewStateEvent(EventQuit, "user", "program"),
		NewStateEvent(EventSetLogged, "user", "program"),
	}

	err := bus.PublishAll(context.Background(), events)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(receivedEvents) != 3 {
		t.Errorf("expected 3 events received, got %d", len(receivedEvents))
	}
}

func TestBus_PublishAll_Error(t *testing.T) {
	bus := NewBus()

	expectedErr := errors.New("first error")
	var callCount int

	bus.Subscribe(EventEnrolled, func(ctx context.Context, event StateEvent) error {
		callCount++
		return expectedErr
	})
	bus.Subscribe(EventQuit, func(ctx context.Context, event StateEvent) error {
		callCount++
		return nil
	})

	events := []StateEvent{
		NewStateEvent(EventEnrolled, "user", "program"),
		NewStateEvent(EventQuit, "user", "program"),
	}

	err := bus.PublishAll(context.Background(), events)

	// Should return first error but publish all events
	if err != expectedErr {
		t.Errorf("expected first error, got %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected all events to be published, got %d", callCount)
	}
}

func TestBus_Clear(t *testing.T) {
	bus := NewBus()

	handler := func(ctx context.Context, event StateEvent) error {
		return nil
	}

	bus.Subscribe(EventEnrolled, handler)
	bus.Subscribe(EventQuit, handler)

	bus.Clear()

	if bus.HasSubscribers(EventEnrolled) {
		t.Error("expected no subscribers after Clear")
	}
	if bus.HasSubscribers(EventQuit) {
		t.Error("expected no subscribers after Clear")
	}
}

func TestBus_ClearEventType(t *testing.T) {
	bus := NewBus()

	handler := func(ctx context.Context, event StateEvent) error {
		return nil
	}

	bus.Subscribe(EventEnrolled, handler)
	bus.Subscribe(EventQuit, handler)

	bus.ClearEventType(EventEnrolled)

	if bus.HasSubscribers(EventEnrolled) {
		t.Error("expected no subscribers for EventEnrolled")
	}
	if !bus.HasSubscribers(EventQuit) {
		t.Error("expected subscribers for EventQuit to remain")
	}
}

func TestBus_ThreadSafety(t *testing.T) {
	bus := NewBus()

	var wg sync.WaitGroup
	var callCount atomic.Int64

	handler := func(ctx context.Context, event StateEvent) error {
		callCount.Add(1)
		return nil
	}

	// Concurrent subscribe and publish
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			bus.Subscribe(EventEnrolled, handler)
		}()
		go func() {
			defer wg.Done()
			bus.Publish(context.Background(), NewStateEvent(EventEnrolled, "user", "program"))
		}()
	}

	wg.Wait()

	// Should have at least some successful publishes
	if callCount.Load() == 0 {
		t.Error("expected at least some handlers to be called")
	}
}
