package notify

import (
	"encoding/json"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/scanner"
)

func makeRMQEvent(port uint16, kind string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
		Time: time.Now(),
	}
}

// fakeRabbitMQ captures published messages without a real broker.
type fakeRabbitMQ struct {
	published []amqp.Publishing
}

func (f *fakeRabbitMQ) Publish(_, _ string, _, _ bool, msg amqp.Publishing) error {
	f.published = append(f.published, msg)
	return nil
}

// publisherFunc is the minimal interface used by the notifier for testing.
type publisherFunc func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error

func TestRabbitMQNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	// Construct a notifier with a no-op channel substitute via the internal helper.
	n := &RabbitMQNotifier{
		exchange:   "portwatch",
		routingKey: "ports.changed",
	}
	// Override channel publish via embedding — use Notify directly with empty slice.
	_ = called
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error for empty events, got %v", err)
	}
	if err := n.Notify([]alert.Event{}); err != nil {
		t.Fatalf("expected no error for empty events, got %v", err)
	}
}

func TestRabbitMQPayload_ContainsEvents(t *testing.T) {
	events := []alert.Event{
		makeRMQEvent(8080, "opened"),
		makeRMQEvent(9090, "closed"),
	}

	payload := rabbitmqPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Count:     len(events),
		Events:    events,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out rabbitmqPayload
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.Count != 2 {
		t.Errorf("expected count 2, got %d", out.Count)
	}
	if len(out.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(out.Events))
	}
	if out.Events[0].Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", out.Events[0].Port.Port)
	}
}

func TestRabbitMQPayload_TimestampPresent(t *testing.T) {
	events := []alert.Event{makeRMQEvent(443, "opened")}

	payload := rabbitmqPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Count:     len(events),
		Events:    events,
	}

	body, _ := json.Marshal(payload)
	var out map[string]interface{}
	_ = json.Unmarshal(body, &out)

	if _, ok := out["timestamp"]; !ok {
		t.Error("expected timestamp field in payload")
	}
}
