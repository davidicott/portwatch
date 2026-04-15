package notify_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func makeEvent(kind, port string) alert.Event {
	return alert.Event{Kind: kind, Port: port}
}

func TestStdoutNotifier_WritesLines(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewStdoutNotifier(&buf)
	events := []alert.Event{
		makeEvent("opened", "tcp:8080"),
		makeEvent("closed", "tcp:9090"),
	}
	if err := n.Notify(context.Background(), events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %q", out)
	}
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %q", out)
	}
	if !strings.Contains(out, "tcp:8080") {
		t.Errorf("expected port in output, got: %q", out)
	}
}

func TestStdoutNotifier_EmptyEvents(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewStdoutNotifier(&buf)
	if err := n.Notify(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty events, got %q", buf.String())
	}
}

func TestMultiNotifier_CallsAll(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	n1 := notify.NewStdoutNotifier(&buf1)
	n2 := notify.NewStdoutNotifier(&buf2)
	multi := notify.NewMultiNotifier(n1, n2)

	events := []alert.Event{makeEvent("opened", "udp:53")}
	if err := multi.Notify(context.Background(), events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf1.String(), "udp:53") {
		t.Error("n1 did not receive event")
	}
	if !strings.Contains(buf2.String(), "udp:53") {
		t.Error("n2 did not receive event")
	}
}
