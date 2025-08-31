package events

import "sync"

type EventHandler func(event Event)

type EventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *EventBus) Publish(eventType string, event Event) {
	eb.mutex.RLock()
	handlers := eb.handlers[eventType]
	eb.mutex.RUnlock()

	for _, handler := range handlers {
		go handler(event)
	}
}