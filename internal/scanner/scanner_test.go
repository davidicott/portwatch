package scanner

import (
	"net"
	"testing"
)

func TestPortKey(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "localhost", Port: 8080, State: "open"}
	want := "tcp/localhost:8080"
	if got := p.Key(); got != want {
		t.Errorf("Key() = %q, want %q", got, want)
	}
}

func TestPortString(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "localhost", Port: 8080, State: "open"}
	want := "localhost:8080 (tcp) [open]"
	if got := p.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestScanDetectsOpenPort(t *testing.T) {
	// Start a real TCP listener on an ephemeral port.
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	port := addr.Port

	// Verify getListeners detects it (only if port <= 1024).
	if port > 1024 {
		t.Skipf("ephemeral port %d is outside scan range, skipping", port)
	}

	ports, err := getListeners("tcp")
	if err != nil {
		t.Fatalf("getListeners error: %v", err)
	}

	found := false
	for _, p := range ports {
		if p.Port == port {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected port %d to be detected, but it was not", port)
	}
}

func TestNewScanner(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if len(s.Protocols) == 0 {
		t.Error("expected at least one protocol, got none")
	}
}
