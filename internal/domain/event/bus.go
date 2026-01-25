package event

import (
	"context"
	"sync"
)

// EventHandler is a function that handles a StateEvent.
// Handlers should be quick and non-blocking.
// For long-running operations, handlers should spawn goroutines or queue work.
type EventHandler func(ctx context.Context, event StateEvent) error

// Bus is a thread-safe in-memory pub/sub event bus.
// It supports subscribing handlers to specific event types and publishing events.
type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]EventHandler
}

// NewBus creates a new event bus.
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type.
// Multiple handlers can be registered for the same event type.
// Handlers are called in the order they were registered.
func (b *Bus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// SubscribeMultiple registers a handler for multiple event types.
// This is a convenience method for subscribing to several events at once.
func (b *Bus) SubscribeMultiple(eventTypes []EventType, handler EventHandler) {
	for _, et := range eventTypes {
		b.Subscribe(et, handler)
	}
}

// Publish sends an event to all registered handlers for that event type.
// Handlers are called synchronously in order. If any handler returns an error,
// subsequent handlers are still called but the first error is returned.
// This allows for partial processing while still reporting failures.
func (b *Bus) Publish(ctx context.Context, event StateEvent) error {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	// Make a copy to release the lock before calling handlers
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)
	b.mu.RUnlock()

	var firstErr error
	for _, handler := range handlersCopy {
		if err := handler(ctx, event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// PublishAsync sends an event to all registered handlers asynchronously.
// Each handler is called in its own goroutine.
// Returns immediately without waiting for handlers to complete.
// Use this when you don't need to know if handlers succeeded.
func (b *Bus) PublishAsync(ctx context.Context, event StateEvent) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)
	b.mu.RUnlock()

	for _, handler := range handlersCopy {
		go func(h EventHandler) {
			_ = h(ctx, event)
		}(handler)
	}
}

// PublishAll publishes multiple events synchronously.
// Events are published in order. If any event's handlers return an error,
// subsequent events are still published but the first error is returned.
func (b *Bus) PublishAll(ctx context.Context, events []StateEvent) error {
	var firstErr error
	for _, event := range events {
		if err := b.Publish(ctx, event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// HasSubscribers returns true if there are any handlers registered for the event type.
func (b *Bus) HasSubscribers(eventType EventType) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType]) > 0
}

// SubscriberCount returns the number of handlers registered for an event type.
func (b *Bus) SubscriberCount(eventType EventType) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType])
}

// Clear removes all handlers for all event types.
// This is useful for testing or resetting the bus.
func (b *Bus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = make(map[EventType][]EventHandler)
}

// ClearEventType removes all handlers for a specific event type.
func (b *Bus) ClearEventType(eventType EventType) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, eventType)
}
