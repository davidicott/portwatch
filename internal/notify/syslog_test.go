package notify

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeSyslogEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Port: port},
	}
}

func TestSyslogNotifier_SkipsEmptyEvents(t *testing.T) {
	n, err := NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()

	if err := n.Notify(nil); err != nil {
		t.Errorf("expected no error for empty events, got %v", err)
	}
}

func TestSyslogNotifier_DefaultTag(t *testing.T) {
	n, err := NewSyslogNotifier("")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()

	if n.tag != "portwatch" {
		t.Errorf("expected default tag 'portwatch', got %q", n.tag)
	}
}

func TestSyslogNotifier_WritesEvents(t *testing.T) {
	n, err := NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()

	events := []alert.Event{
		makeSyslogEvent("opened", "tcp", 8080),
		makeSyslogEvent("closed", "udp", 53),
	}

	if err := n.Notify(events); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
