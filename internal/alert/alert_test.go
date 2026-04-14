package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Number: port}
}

func TestBuildEvents_Opened(t *testing.T) {
	opened := []scanner.Port{makePort("tcp", 8080)}
	events := alert.BuildEvents(opened, nil)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelAlert {
		t.Errorf("expected ALERT level, got %s", events[0].Level)
	}
	if !strings.Contains(events[0].Message, "opened") {
		t.Errorf("expected message to contain 'opened', got %q", events[0].Message)
	}
}

func TestBuildEvents_Closed(t *testing.T) {
	closed := []scanner.Port{makePort("tcp", 22)}
	events := alert.BuildEvents(nil, closed)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelWarn {
		t.Errorf("expected WARN level, got %s", events[0].Level)
	}
	if !strings.Contains(events[0].Message, "closed") {
		t.Errorf("expected message to contain 'closed', got %q", events[0].Message)
	}
}

func TestBuildEvents_Empty(t *testing.T) {
	events := alert.BuildEvents(nil, nil)
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestLogNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.LogNotifier{Out: &buf}

	events := alert.BuildEvents([]scanner.Port{makePort("udp", 53)}, nil)
	for _, e := range events {
		if err := n.Notify(e); err != nil {
			t.Fatalf("Notify returned error: %v", err)
		}
	}

	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected output to contain ALERT, got %q", output)
	}
	if !strings.Contains(output, "opened") {
		t.Errorf("expected output to contain 'opened', got %q", output)
	}
}
