package models

import "time"

// Message represents one row in messages.
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageInput represents the input for creating a message.
type MessageInput struct {
	Content string `json:"content"`
}
