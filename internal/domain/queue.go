package domain

import (
	"context"

	"github.com/google/uuid"
)

// Message represents a generic message in the queue
type Message struct {
	AccountID uuid.UUID       `json:"account_id"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
}

// MessageProcessor defines the interface for processing messages
type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg *Message) error
}

// MessageProducer defines the interface for producing messages
type MessageProducer interface {
	PublishMessage(ctx context.Context, msg *Message) error
	Close() error
}

// MessageConsumer defines the interface for consuming messages
type MessageConsumer interface {
	Start(ctx context.Context) error
	Close() error
}
