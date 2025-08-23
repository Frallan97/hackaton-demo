package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Get NATS URL from environment, fallback to localhost
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	// Connect to NATS
	nc, err := nats.Connect(natsURL,
		nats.Name("test-client"),
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(5),
	)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Printf("Connected to NATS server: %s", nc.ConnectedUrl())

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v", err)
	}

	// Create stream
	stream, err := js.AddStream(&nats.StreamConfig{
		Name:     "TEST",
		Subjects: []string{"test.>"},
		Storage:  nats.FileStorage,
		MaxAge:   1 * time.Hour,
	})
	if err != nil {
		log.Printf("Stream creation result: %v", err)
	} else {
		log.Printf("Created stream: %s", stream.Config.Name)
	}

	// Subscribe to test topic
	sub, err := nc.Subscribe("test.message", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	// Publish a test message
	testData := map[string]interface{}{
		"message": "Hello from NATS!",
		"time":    time.Now().Format(time.RFC3339),
		"number":  42,
	}

	jsonData, _ := json.Marshal(testData)
	err = nc.Publish("test.message", jsonData)
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	log.Println("Published test message")

	// Wait a bit for message processing
	time.Sleep(2 * time.Second)

	// Get stream info
	info, err := js.StreamInfo("TEST")
	if err != nil {
		log.Printf("Failed to get stream info: %v", err)
	} else {
		log.Printf("Stream info: %d messages, %d bytes", info.State.Msgs, info.State.Bytes)
	}

	log.Println("NATS test completed successfully!")
}
