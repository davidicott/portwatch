package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

// kafkaWriter is a minimal interface for writing messages to Kafka.
type kafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafkaMessage) error
	Close() error
}

type kafkaMessage struct {
	Topic string
	Value []byte
}

// kafkaTCPWriter is a simple TCP-based Kafka producer for a single partition.
type kafkaTCPWriter struct {
	addr  string
	topic string
}

func (w *kafkaTCPWriter) WriteMessages(ctx context.Context, msgs ...kafkaMessage) error {
	conn, err := net.DialTimeout("tcp", w.addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("kafka: dial %s: %w", w.addr, err)
	}
	defer conn.Close()
	for _, m := range msgs {
		if _, err := fmt.Fprintf(conn, "%s\n", m.Value); err != nil {
			return fmt.Errorf("kafka: write message: %w", err)
		}
	}
	return nil
}

func (w *kafkaTCPWriter) Close() error { return nil }

// KafkaNotifier publishes port-change events to a Kafka topic.
type KafkaNotifier struct {
	writer kafkaWriter
	topic  string
}

// NewKafkaNotifier creates a KafkaNotifier that publishes to the given broker and topic.
func NewKafkaNotifier(broker, topic string) *KafkaNotifier {
	return &KafkaNotifier{
		writer: &kafkaTCPWriter{addr: broker, topic: topic},
		topic:  topic,
	}
}

// Notify publishes each event as a JSON message to the configured Kafka topic.
func (n *KafkaNotifier) Notify(ctx context.Context, events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	msgs := make([]kafkaMessage, 0, len(events))
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("kafka: marshal event: %w", err)
		}
		msgs = append(msgs, kafkaMessage{Topic: n.topic, Value: b})
	}
	if err := n.writer.WriteMessages(ctx, msgs...); err != nil {
		return fmt.Errorf("kafka: write messages: %w", err)
	}
	return nil
}
