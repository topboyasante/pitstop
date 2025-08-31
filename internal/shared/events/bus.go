package events

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/topboyasante/pitstop/internal/core/logger"
)

// EventHandler represents a function that handles events
type EventHandler func(event Event) error

// EventBus provides event publishing and subscription functionality
type EventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

// NewEventBus creates a new event bus instance
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers an event handler for a specific event type
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if _, exists := eb.handlers[eventType]; !exists {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	logger.Info("Event handler registered", "event_type", eventType)
}

// SubscribeToEvent registers an event handler for a specific event struct type
func (eb *EventBus) SubscribeToEvent(event Event, handler EventHandler) {
	eventType := reflect.TypeOf(event).Elem().Name()
	eb.Subscribe(eventType, handler)
}

// Publish publishes an event to all registered handlers
func (eb *EventBus) Publish(event Event) error {
	eb.mutex.RLock()
	eventType := reflect.TypeOf(event).Elem().Name()
	handlers, exists := eb.handlers[eventType]
	eb.mutex.RUnlock()

	if !exists {
		logger.Debug("No handlers registered for event", "event_type", eventType)
		return nil
	}

	logger.Info("Publishing event", "event_type", eventType, "handlers_count", len(handlers))

	var errs []error
	for i, handler := range handlers {
		if err := handler(event); err != nil {
			logger.Error("Event handler failed", 
				"event_type", eventType, 
				"handler_index", i, 
				"error", err)
			errs = append(errs, fmt.Errorf("handler %d failed: %w", i, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("event publishing failed with %d errors: %v", len(errs), errs)
	}

	logger.Debug("Event published successfully", "event_type", eventType)
	return nil
}

// PublishAsync publishes an event asynchronously
func (eb *EventBus) PublishAsync(event Event) {
	go func() {
		if err := eb.Publish(event); err != nil {
			logger.Error("Async event publishing failed", "error", err)
		}
	}()
}

// Unsubscribe removes all handlers for a specific event type
func (eb *EventBus) Unsubscribe(eventType string) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	delete(eb.handlers, eventType)
	logger.Info("All handlers unsubscribed", "event_type", eventType)
}

// GetHandlerCount returns the number of handlers for an event type
func (eb *EventBus) GetHandlerCount(eventType string) int {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	if handlers, exists := eb.handlers[eventType]; exists {
		return len(handlers)
	}
	return 0
}

// GetAllEventTypes returns all registered event types
func (eb *EventBus) GetAllEventTypes() []string {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	var eventTypes []string
	for eventType := range eb.handlers {
		eventTypes = append(eventTypes, eventType)
	}
	return eventTypes
}

// Global event bus instance
var GlobalEventBus = NewEventBus()

// Convenience functions for global event bus
func Subscribe(eventType string, handler EventHandler) {
	GlobalEventBus.Subscribe(eventType, handler)
}

func SubscribeToEvent(event Event, handler EventHandler) {
	GlobalEventBus.SubscribeToEvent(event, handler)
}

func Publish(event Event) error {
	return GlobalEventBus.Publish(event)
}

func PublishAsync(event Event) {
	GlobalEventBus.PublishAsync(event)
}