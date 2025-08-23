package events

// EventBus defines the interface for event bus implementations
type EventBus interface {
	// Publish sends an event to a topic
	Publish(topic string, eventType string, data map[string]interface{}, userID *int) error

	// Subscribe creates a subscription to a topic
	Subscribe(topic string) (<-chan Event, error)

	// Unsubscribe removes a subscription
	Unsubscribe(topic string, ch <-chan Event)

	// RegisterHandler registers an event handler
	RegisterHandler(eventType string, handler EventHandler)

	// UnregisterHandler removes an event handler
	UnregisterHandler(eventType string, handler EventHandler)

	// Shutdown gracefully shuts down the event bus
	Shutdown()

	// GetEventStats returns event bus statistics
	GetEventStats() map[string]interface{}
}
