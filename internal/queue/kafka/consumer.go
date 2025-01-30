package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/prajwalpamin/banking-ledger/internal/domain"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader    *kafka.Reader
	processor domain.MessageProcessor
}

func NewConsumer(brokers []string, topic string, groupID string, processor domain.MessageProcessor) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader:    reader,
		processor: processor,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, err := c.reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}

				var msg domain.Message
				if err := json.Unmarshal(message.Value, &msg); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}

				if err := c.processor.ProcessMessage(ctx, &msg); err != nil {
					log.Printf("Error processing message: %v", err)
					continue
				}

				if err := c.reader.CommitMessages(ctx, message); err != nil {
					log.Printf("Error committing message: %v", err)
				}
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
