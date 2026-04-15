package filter

import (
	"github.com/user/portwatch/internal/scanner"
)

// Filter decides which ports should be excluded from alerting.
type Filter struct {
	ignorePorts map[uint16]struct{}
	ignoreProto map[string]struct{}
}

// New creates a Filter from lists of ports and protocols to ignore.
func New(ports []uint16, protocols []string) *Filter {
	f := &Filter{
		ignorePorts: make(map[uint16]struct{}, len(ports)),
		ignoreProto: make(map[string]struct{}, len(protocols)),
	}
	for _, p := range ports {
		f.ignorePorts[p] = struct{}{}
	}
	for _, proto := range protocols {
		f.ignoreProto[proto] = struct{}{}
	}
	return f
}

// Apply returns only the ports that are NOT suppressed by the filter.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	out := ports[:0:0]
	for _, p := range ports {
		if _, skip := f.ignorePorts[p.Number]; skip {
			continue
		}
		if _, skip := f.ignoreProto[p.Proto]; skip {
			continue
		}
		out = append(out, p)
	}
	return out
}

// ShouldIgnore reports whether a single port is suppressed.
func (f *Filter) ShouldIgnore(p scanner.Port) bool {
	_, byPort := f.ignorePorts[p.Number]
	_, byProto := f.ignoreProto[p.Proto]
	return byPort || byProto
}
