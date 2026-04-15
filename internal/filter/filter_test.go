package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Proto: proto}
}

func TestApply_FiltersIgnoredPorts(t *testing.T) {
	f := filter.New([]uint16{22, 80}, nil)
	input := []scanner.Port{
		makePort(22, "tcp"),
		makePort(443, "tcp"),
		makePort(80, "tcp"),
		makePort(8080, "tcp"),
	}
	got := f.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
	if got[0].Number != 443 || got[1].Number != 8080 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestApply_FiltersIgnoredProtocols(t *testing.T) {
	f := filter.New(nil, []string{"udp"})
	input := []scanner.Port{
		makePort(53, "udp"),
		makePort(53, "tcp"),
		makePort(123, "udp"),
	}
	got := f.Apply(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
	if got[0].Proto != "tcp" {
		t.Errorf("expected tcp port, got %s", got[0].Proto)
	}
}

func TestApply_EmptyFilter(t *testing.T) {
	f := filter.New(nil, nil)
	input := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	got := f.Apply(input)
	if len(got) != len(input) {
		t.Errorf("expected %d ports, got %d", len(input), len(got))
	}
}

func TestShouldIgnore(t *testing.T) {
	f := filter.New([]uint16{22}, []string{"udp"})

	if !f.ShouldIgnore(makePort(22, "tcp")) {
		t.Error("port 22 should be ignored")
	}
	if !f.ShouldIgnore(makePort(9999, "udp")) {
		t.Error("udp port should be ignored")
	}
	if f.ShouldIgnore(makePort(443, "tcp")) {
		t.Error("port 443/tcp should not be ignored")
	}
}
