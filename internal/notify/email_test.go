package notify

import (
	"net"
	"net/smtp"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEmailEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{
			Proto:  proto,
			Addr:   "127.0.0.1",
			Number: port,
		},
	}
}

func TestBuildEmailBody_ContainsEvents(t *testing.T) {
	events := []alert.Event{
		makeEmailEvent("opened", "tcp", 8080),
		makeEmailEvent("closed", "udp", 53),
	}
	body := buildEmailBody(events)

	if !contains(body, "OPENED") {
		t.Error("expected body to contain OPENED")
	}
	if !contains(body, "CLOSED") {
		t.Error("expected body to contain CLOSED")
	}
	if !contains(body, "8080") {
		t.Error("expected body to contain port 8080")
	}
	if !contains(body, "udp") {
		t.Error("expected body to contain protocol udp")
	}
}

func TestBuildEmailBody_Empty(t *testing.T) {
	body := buildEmailBody(nil)
	if body == "" {
		t.Error("expected non-empty body even with no events")
	}
}

func TestEmailNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewEmailNotifier(EmailConfig{
		Host: "localhost",
		Port: 25,
		From: "portwatch@example.com",
		To:   []string{"admin@example.com"},
	})
	// Should return nil without attempting a connection.
	if err := n.Notify([]alert.Event{}); err != nil {
		t.Errorf("unexpected error for empty events: %v", err)
	}
}

func TestEmailNotifier_SendsToSMTP(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skip("cannot bind local port:", err)
	}
	defer ln.Close()

	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	// Spin up a minimal fake SMTP server that accepts and closes.
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write([]byte("220 fake smtp\r\n"))
		buf := make([]byte, 4096)
		conn.Read(buf)
	}()

	_ = smtp.SendMail // ensure import used
	n := NewEmailNotifier(EmailConfig{
		Host: "127.0.0.1",
		Port: port,
		From: "portwatch@example.com",
		To:   []string{"admin@example.com"},
	})
	// We only verify no panic; real SMTP handshake will fail on fake server.
	_ = n.Notify([]alert.Event{makeEmailEvent("opened", "tcp", 9090)})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
