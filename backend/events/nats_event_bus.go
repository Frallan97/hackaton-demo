package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// NATSEventBus implements EventBus interface using NATS
type NATSEventBus struct {
	nc          *nats.Conn
	js          nats.JetStreamContext
	subscribers map[string][]chan Event
	handlers    map[string][]EventHandler
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	serverURL   string
	streamName  string
}

// NewNATSEventBus creates a new NATS event bus instance
func NewNATSEventBus(serverURL string) (*NATSEventBus, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Connect to NATS server
	nc, err := nats.Connect(serverURL,
		nats.Name("hackaton-demo-event-bus"),
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(5),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to NATS server: %s", nc.ConnectedUrl())
		}),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Printf("Disconnected from NATS server")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("NATS error: %v", err)
		}),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context for persistence
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		cancel()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create stream for event persistence
	streamName := "EVENTS"
	stream, err := js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
		Storage:  nats.FileStorage,
		MaxAge:   24 * time.Hour, // Keep events for 24 hours
	})
	if err != nil {
		// Check if stream already exists (common case)
		if err.Error() == "stream name already in use" {
			log.Printf("NATS stream already exists: %s", streamName)
		} else {
			log.Printf("Warning: Failed to create stream: %v", err)
		}
	} else {
		log.Printf("Created NATS stream: %s", stream.Config.Name)
	}

	eb := &NATSEventBus{
		nc:          nc,
		js:          js,
		subscribers: make(map[string][]chan Event),
		handlers:    make(map[string][]EventHandler),
		ctx:         ctx,
		cancel:      cancel,
		serverURL:   serverURL,
		streamName:  streamName,
	}

	// Start the event processor
	go eb.processEvents()

	log.Printf("NATS Event Bus initialized successfully. Connected to: %s", serverURL)
	return eb, nil
}

// Publish sends an event to NATS
func (eb *NATSEventBus) Publish(topic string, eventType string, data map[string]interface{}, userID *int) error {
	event := Event{
		ID:        generateEventID(),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
		Source:    "backend",
		UserID:    userID,
	}

	// Publish to NATS subject
	subject := fmt.Sprintf("events.%s.%s", topic, eventType)
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to NATS
	if err := eb.nc.Publish(subject, eventJSON); err != nil {
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	// Also publish to JetStream for persistence
	if _, err := eb.js.Publish(subject, eventJSON); err != nil {
		log.Printf("Warning: Failed to persist event to JetStream: %v", err)
	}

	// Send to local subscribers (for immediate processing)
	eb.sendToLocalSubscribers(topic, event)

	// Process with handlers
	eb.processWithHandlers(event)

	log.Printf("Published event: %s (type: %s) to NATS subject: %s", event.ID, eventType, subject)
	return nil
}

// Subscribe creates a subscription to a topic
func (eb *NATSEventBus) Subscribe(topic string) (<-chan Event, error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, 100) // Buffer size of 100

	if eb.subscribers[topic] == nil {
		eb.subscribers[topic] = make([]chan Event, 0)
	}

	eb.subscribers[topic] = append(eb.subscribers[topic], ch)

	// Also subscribe to NATS for real-time updates
	subject := fmt.Sprintf("events.%s.>", topic)
	_, err := eb.nc.Subscribe(subject, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal NATS message: %v", err)
			return
		}

		// Send to local subscribers
		eb.sendToLocalSubscribers(topic, event)

		// Acknowledge the message
		msg.Ack()
	})
	if err != nil {
		log.Printf("Warning: Failed to subscribe to NATS subject %s: %v", subject, err)
	}

	log.Printf("New subscription to topic: %s (NATS subject: %s)", topic, subject)
	return ch, nil
}

// Unsubscribe removes a subscription
func (eb *NATSEventBus) Unsubscribe(topic string, ch <-chan Event) {
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

// RegisterHandler registers an event handler
func (eb *NATSEventBus) RegisterHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	log.Printf("Registered handler for event type: %s", eventType)
}

// UnregisterHandler removes an event handler
func (eb *NATSEventBus) UnregisterHandler(eventType string, handler EventHandler) {
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

// sendToLocalSubscribers sends events to local subscribers
func (eb *NATSEventBus) sendToLocalSubscribers(topic string, event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

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
}

// processWithHandlers processes events with registered handlers
func (eb *NATSEventBus) processWithHandlers(event Event) {
	eb.mu.RLock()
	handlers, exists := eb.handlers[event.Type]
	eb.mu.RUnlock()

	if !exists {
		return
	}

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

// processEvents handles event processing in the background
func (eb *NATSEventBus) processEvents() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-eb.ctx.Done():
			log.Println("NATS Event Bus shutting down")
			return
		case <-ticker.C:
			// Periodic cleanup and health check
			eb.cleanupClosedChannels()
			eb.checkNATSConnection()
		}
	}
}

// cleanupClosedChannels removes closed channels from subscribers
func (eb *NATSEventBus) cleanupClosedChannels() {
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

// checkNATSConnection checks if NATS connection is healthy
func (eb *NATSEventBus) checkNATSConnection() {
	if !eb.nc.IsConnected() {
		log.Printf("Warning: NATS connection lost. Attempting to reconnect...")
	}
}

// Shutdown gracefully shuts down the NATS event bus
func (eb *NATSEventBus) Shutdown() {
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

	// Close NATS connection
	if eb.nc != nil {
		eb.nc.Close()
	}

	log.Println("NATS Event Bus shutdown complete")
}

// GetEventStats returns event bus statistics
func (eb *NATSEventBus) GetEventStats() map[string]interface{} {
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

	// NATS connection info
	natsInfo := make(map[string]interface{})
	if eb.nc != nil {
		natsInfo["connected"] = eb.nc.IsConnected()
		natsInfo["server_url"] = eb.serverURL
		natsInfo["connection_id"] = eb.nc.ConnectedServerId()
	}

	stats["topics"] = len(eb.subscribers)
	stats["total_subscribers"] = len(eb.subscribers)
	stats["total_handlers"] = len(eb.handlers)
	stats["topic_subscribers"] = topicSubscribers
	stats["event_handlers"] = eventHandlers
	stats["nats"] = natsInfo

	return stats
}

// GetJetStreamInfo returns JetStream statistics
func (eb *NATSEventBus) GetJetStreamInfo() (map[string]interface{}, error) {
	if eb.js == nil {
		return nil, fmt.Errorf("JetStream not available")
	}

	info, err := eb.js.StreamInfo(eb.streamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream info: %w", err)
	}

	return map[string]interface{}{
		"stream_name":    info.Config.Name,
		"subjects":       info.Config.Subjects,
		"messages":       info.State.Msgs,
		"bytes":          info.State.Bytes,
		"first_sequence": info.State.FirstSeq,
		"last_sequence":  info.State.LastSeq,
		"consumer_count": info.State.Consumers,
	}, nil
}
