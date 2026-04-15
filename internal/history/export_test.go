package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeExportEvent(kind alert.EventKind, proto, addr string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Time: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Port: scanner.Port{Proto: proto, Addr: addr, Port: port, PID: 42},
	}
}

func TestExportJSON_Empty(t *testing.T) {
	h := New(10)
	var buf bytes.Buffer
	if err := h.ExportJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var events []alert.Event
	if err := json.Unmarshal(buf.Bytes(), &events); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty slice, got %d events", len(events))
	}
}

func TestExportJSON_ContainsEvents(t *testing.T) {
	h := New(10)
	h.Record(makeExportEvent(alert.Opened, "tcp", "0.0.0.0", 8080))
	h.Record(makeExportEvent(alert.Closed, "tcp", "0.0.0.0", 9090))

	var buf bytes.Buffer
	if err := h.ExportJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var events []alert.Event
	if err := json.Unmarshal(buf.Bytes(), &events); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestExportJSON_PreservesEventFields(t *testing.T) {
	h := New(10)
	h.Record(makeExportEvent(alert.Opened, "tcp", "127.0.0.1", 8080))

	var buf bytes.Buffer
	if err := h.ExportJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var events []alert.Event
	if err := json.Unmarshal(buf.Bytes(), &events); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ev := events[0]
	if ev.Kind != alert.Opened {
		t.Errorf("expected Kind=Opened, got %v", ev.Kind)
	}
	if ev.Port.Proto != "tcp" || ev.Port.Addr != "127.0.0.1" || ev.Port.Port != 8080 || ev.Port.PID != 42 {
		t.Errorf("unexpected port fields: %+v", ev.Port)
	}
}

func TestExportTable_ContainsHeaders(t *testing.T) {
	h := New(10)
	h.Record(makeExportEvent(alert.Opened, "tcp", "127.0.0.1", 3000))

	var buf bytes.Buffer
	if err := h.ExportTable(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"TIME", "EVENT", "PROTO", "ADDR", "PORT", "PID"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in table output", hdr)
		}
	}
}

func TestExportTable_ContainsEventData(t *testing.T) {
	h := New(10)
	h.Record(makeExportEvent(alert.Opened, "udp", "0.0.0.0", 5353))

	var buf bytes.Buffer
	if err := h.ExportTable(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"OPENED", "udp", "0.0.0.0", "5353", "42"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in table output\n%s", want, out)
		}
	}
}
