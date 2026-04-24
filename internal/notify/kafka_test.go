package notify

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/patrickdappollonio/portwatch/internal/alert"
	"github.com/patrickdappollonio/portwatch/internal/scanner"
)

type mockKafkaWriter struct {
	written []kafkaMessage
	err     error
}

func (m *mockKafkaWriter) WriteMessages(_ context.Context, msgs ...kafkaMessage) error {
	if m.err != nil {
		return m.err
	}
	m.written = append(m.written, msgs...)
	return nil
}

func (m *mockKafkaWriter) Close() error { return nil }

func makeKafkaEvent(proto string, port uint16, kind alert.EventKind) alert.Event {
	return alert.Event{
		Port: scanner.Port{Proto: proto, Port: port},
		Kind: kind,
	}
}

func TestKafkaNotifier_SkipsEmptyEvents(t *testing.T) {
	w := &mockKafkaWriter{}
	n := &KafkaNotifier{writer: w, topic: "ports"}
	if err := n.Notify(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.written) != 0 {
		t.Fatalf("expected no messages, got %d", len(w.written))
	}
}

func TestKafkaNotifier_PostsPayload(t *testing.T) {
	w := &mockKafkaWriter{}
	n := &KafkaNotifier{writer: w, topic: "ports"}
	events := []alert.Event{
		makeKafkaEvent("tcp", 8080, alert.EventOpened),
		makeKafkaEvent("tcp", 9090, alert.EventClosed),
	}
	if err := n.Notify(context.Background(), events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.written) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(w.written))
	}
	var got alert.Event
	if err := json.Unmarshal(w.written[0].Value, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", got.Port.Port)
	}
}

func TestKafkaNotifier_NonSuccessStatus(t *testing.T) {
	w := &mockKafkaWriter{err: errors.New("broker unavailable")}
	n := &KafkaNotifier{writer: w, topic: "ports"}
	events := []alert.Event{makeKafkaEvent("tcp", 443, alert.EventOpened)}
	if err := n.Notify(context.Background(), events); err == nil {
		t.Fatal("expected error, got nil")
	}
}
