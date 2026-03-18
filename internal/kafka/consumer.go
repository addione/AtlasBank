package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(broker, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "atlasbank-events",
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader: reader,
	}
}

// Consume starts consuming messages from Kafka
func (c *Consumer) Consume(ctx context.Context, logger interface{}) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		// Process the message
		fmt.Printf("Received message: key=%s, value=%s, topic=%s, partition=%d, offset=%d\n",
			string(msg.Key), string(msg.Value), msg.Topic, msg.Partition, msg.Offset)

		// Here you can add custom message processing logic
		// For now, we just log the message
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
