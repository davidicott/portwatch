package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Port represents an open port with its protocol and process info.
type Port struct {
	Protocol string
	Address  string
	Port     int
	State    string
}

// String returns a human-readable representation of a Port.
func (p Port) String() string {
	return fmt.Sprintf("%s:%d (%s) [%s]", p.Address, p.Port, p.Protocol, p.State)
}

// Key returns a unique identifier for the port.
func (p Port) Key() string {
	return fmt.Sprintf("%s/%s:%d", p.Protocol, p.Address, p.Port)
}

// Scanner is responsible for discovering open ports on the system.
type Scanner struct {
	Protocols []string
	// DialTimeout is the timeout used when probing each port.
	DialTimeout time.Duration
}

// New creates a new Scanner with default protocols.
func New() *Scanner {
	return &Scanner{
		Protocols:   []string{"tcp", "tcp6"},
		DialTimeout: 500 * time.Millisecond,
	}
}

// Scan returns a list of currently open ports by attempting to
// parse active listeners via net.Listen probing on a range.
func (s *Scanner) Scan() ([]Port, error) {
	var ports []Port
	for _, proto := range s.Protocols {
		listeners, err := s.getListeners(proto)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", proto, err)
		}
		ports = append(ports, listeners...)
	}
	return ports, nil
}

// getListeners probes ports 1–1024 for the given protocol.
func (s *Scanner) getListeners(proto string) ([]Port, error) {
	var open []Port
	baseProto := strings.TrimSuffix(proto, "6")
	dialer := net.Dialer{Timeout: s.DialTimeout}
	for p := 1; p <= 1024; p++ {
		addr := net.JoinHostPort("localhost", strconv.Itoa(p))
		conn, err := dialer.Dial(baseProto, addr)
		if err == nil {
			conn.Close()
			open = append(open, Port{
				Protocol: proto,
				Address:  "localhost",
				Port:     p,
				State:    "open",
			})
		}
	}
	return open, nil
}
