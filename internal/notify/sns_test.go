package notify

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
)

type mockSNSClient struct {
	called bool
	lastInput *sns.PublishInput
	err       error
}

func (m *mockSNSClient) Publish(_ context.Context, input *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.called = true
	m.lastInput = input
	return &sns.PublishOutput{}, m.err
}

func makeSNSEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Port: port},
	}
}

func TestSNSNotifier_SkipsEmptyEvents(t *testing.T) {
	client := &mockSNSClient{}
	n := NewSNSNotifierWithClient(client, "arn:aws:sns:us-east-1:123456789012:portwatch")
	if err := n.Notify(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.called {
		t.Error("expected Publish not to be called for empty events")
	}
}

func TestSNSNotifier_PostsPayload(t *testing.T) {
	client := &mockSNSClient{}
	n := NewSNSNotifierWithClient(client, "arn:aws:sns:us-east-1:123456789012:portwatch")
	events := []alert.Event{
		makeSNSEvent("opened", "tcp", 8080),
		makeSNSEvent("closed", "tcp", 9090),
	}
	if err := n.Notify(context.Background(), events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !client.called {
		t.Fatal("expected Publish to be called")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(*client.lastInput.Message), &body); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if body["source"] != "portwatch" {
		t.Errorf("expected source=portwatch, got %v", body["source"])
	}
	if *client.lastInput.TopicArn != "arn:aws:sns:us-east-1:123456789012:portwatch" {
		t.Errorf("unexpected topic ARN: %s", *client.lastInput.TopicArn)
	}
}

func TestSNSNotifier_SubjectContainsCount(t *testing.T) {
	client := &mockSNSClient{}
	n := NewSNSNotifierWithClient(client, "arn:aws:sns:us-east-1:123456789012:portwatch")
	events := []alert.Event{makeSNSEvent("opened", "udp", 53)}
	_ = n.Notify(context.Background(), events)
	if client.lastInput == nil {
		t.Fatal("Publish not called")
	}
	subject := *client.lastInput.Subject
	if subject == "" {
		t.Error("expected non-empty subject")
	}
}
