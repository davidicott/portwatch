package notify

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/yourorg/portwatch/internal/alert"
)

// RabbitMQNotifier publishes port-change events to a RabbitMQ exchange.
type RabbitMQNotifier struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// NewRabbitMQNotifier dials the given AMQP URL and declares the target exchange.
func NewRabbitMQNotifier(url, exchange, routingKey string) (*RabbitMQNotifier, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: dial %s: %w", url, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: declare exchange %q: %w", exchange, err)
	}

	return &RabbitMQNotifier{
		conn:       conn,
		channel:    ch,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

type rabbitmqPayload struct {
	Timestamp string        `json:"timestamp"`
	Count     int           `json:"count"`
	Events    []alert.Event `json:"events"`
}

// Notify publishes all events as a single JSON message to the configured exchange.
func (n *RabbitMQNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	payload := rabbitmqPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Count:     len(events),
		Events:    events,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rabbitmq: marshal payload: %w", err)
	}

	err = n.channel.Publish(
		n.exchange,
		n.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("rabbitmq: publish: %w", err)
	}
	return nil
}

// Close releases the channel and connection.
func (n *RabbitMQNotifier) Close() error {
	if err := n.channel.Close(); err != nil {
		return err
	}
	return n.conn.Close()
}
