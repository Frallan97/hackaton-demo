package events

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Event represents a generic event in the system
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	UserID    *int                   `json:"user_id,omitempty"`
}

// EventHandler is a function that processes events
type EventHandler func(ctx context.Context, event Event) error

// CustomEventBus manages event publishing and subscription
type CustomEventBus struct {
	subscribers map[string][]chan Event
	handlers    map[string][]EventHandler
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewEventBus creates a new custom event bus instance
func NewEventBus() *CustomEventBus {
	ctx, cancel := context.WithCancel(context.Background())

	eb := &CustomEventBus{
		subscribers: make(map[string][]chan Event),
		handlers:    make(map[string][]EventHandler),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Start the event processor
	go eb.processEvents()

	return eb
}

// Publish sends an event to all subscribers
func (eb *CustomEventBus) Publish(topic string, eventType string, data map[string]interface{}, userID *int) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	event := Event{
		ID:        generateEventID(),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
		Source:    "backend",
		UserID:    userID,
	}

	// Send to subscribers
	if chans, exists := eb.subscribers[topic]; exists {
		for _, ch := range chans {
			select {
			case ch <- event:
				// Event sent successfully
			default:
				// Channel is full, skip this subscriber
				log.Printf("Warning: subscriber channel is full for topic: %s", topic)
			}
		}
	}

	// Process with handlers
	if handlers, exists := eb.handlers[eventType]; exists {
		for _, handler := range handlers {
			go func(h EventHandler, e Event) {
				ctx, cancel := context.WithTimeout(eb.ctx, 30*time.Second)
				defer cancel()

				if err := h(ctx, e); err != nil {
					log.Printf("Error in event handler for %s: %v", e.Type, err)
				}
			}(handler, event)
		}
	}

	log.Printf("Published event: %s (type: %s) to topic: %s", event.ID, eventType, topic)
	return nil
}

// Subscribe creates a subscription to a topic
func (eb *CustomEventBus) Subscribe(topic string) (<-chan Event, error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, 100) // Buffer size of 100

	if eb.subscribers[topic] == nil {
		eb.subscribers[topic] = make([]chan Event, 0)
	}

	eb.subscribers[topic] = append(eb.subscribers[topic], ch)

	log.Printf("New subscription to topic: %s", topic)
	return ch, nil
}

// Unsubscribe removes a subscription
func (eb *CustomEventBus) Unsubscribe(topic string, ch <-chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if chans, exists := eb.subscribers[topic]; exists {
		for i, subscriber := range chans {
			if subscriber == ch {
				// Remove the channel
				eb.subscribers[topic] = append(chans[:i], chans[i+1:]...)
				close(subscriber)
				log.Printf("Unsubscribed from topic: %s", topic)
				return
			}
		}
	}
}

// RegisterHandler registers an event handler for a specific event type
func (eb *CustomEventBus) RegisterHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	log.Printf("Registered handler for event type: %s", eventType)
}

// UnregisterHandler removes an event handler
func (eb *CustomEventBus) UnregisterHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.handlers[eventType]; exists {
		for i, h := range handlers {
			if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
				eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				log.Printf("Unregistered handler for event type: %s", eventType)
				return
			}
		}
	}
}

// processEvents handles event processing in the background
func (eb *CustomEventBus) processEvents() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-eb.ctx.Done():
			log.Println("Event bus shutting down")
			return
		case <-ticker.C:
			// Periodic cleanup of closed channels
			eb.cleanupClosedChannels()
		}
	}
}

// cleanupClosedChannels removes closed channels from subscribers
func (eb *CustomEventBus) cleanupClosedChannels() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for topic, chans := range eb.subscribers {
		var activeChans []chan Event
		for _, ch := range chans {
			select {
			case _, ok := <-ch:
				if ok {
					// Channel is still open, keep it
					activeChans = append(activeChans, ch)
				}
			default:
				// Channel is open and not full, keep it
				activeChans = append(activeChans, ch)
			}
		}
		eb.subscribers[topic] = activeChans
	}
}

// Shutdown gracefully shuts down the event bus
func (eb *CustomEventBus) Shutdown() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Close all subscriber channels
	for topic, chans := range eb.subscribers {
		for _, ch := range chans {
			close(ch)
		}
		delete(eb.subscribers, topic)
	}

	// Cancel context
	eb.cancel()

	log.Println("Event bus shutdown complete")
}

// generateEventID creates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// GetEventStats returns statistics about the event bus
func (eb *CustomEventBus) GetEventStats() map[string]interface{} {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	stats := make(map[string]interface{})

	// Count subscribers per topic
	topicSubscribers := make(map[string]int)
	for topic, chans := range eb.subscribers {
		topicSubscribers[topic] = len(chans)
	}

	// Count handlers per event type
	eventHandlers := make(map[string]int)
	for eventType, handlers := range eb.handlers {
		eventHandlers[eventType] = len(handlers)
	}

	stats["topics"] = len(eb.subscribers)
	stats["total_subscribers"] = len(eb.subscribers)
	stats["total_handlers"] = len(eb.handlers)
	stats["topic_subscribers"] = topicSubscribers
	stats["event_handlers"] = eventHandlers

	return stats
}
