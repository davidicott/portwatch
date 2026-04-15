package filter

import (
	"github.com/user/portwatch/internal/config"
)

// FromConfig constructs a Filter from the application configuration.
// It reads IgnorePorts and IgnoreProtocols from cfg.
func FromConfig(cfg *config.Config) *Filter {
	ports := make([]uint16, 0, len(cfg.IgnorePorts))
	for _, p := range cfg.IgnorePorts {
		if p >= 0 && p <= 65535 {
			ports = append(ports, uint16(p))
		}
	}
	return New(ports, cfg.IgnoreProtocols)
}
