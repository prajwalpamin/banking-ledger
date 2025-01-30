package kafka

import (
	"context"
	"encoding/json"

	"github.com/prajwalpamin/banking-ledger/internal/domain"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		RequiredAcks: int(kafka.RequireAll),
		Async:        false,
	})

	return &Producer{writer: writer}, nil
}

func (p *Producer) PublishMessage(ctx context.Context, msg *domain.Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: body,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
