package ratelimit

import (
	"time"

	"github.com/user/portwatch/internal/config"
)

// FromConfig builds a Limiter from the application config.
// If rate limiting is disabled or max is zero, FromConfig returns nil.
func FromConfig(cfg config.RateLimitConfig) *Limiter {
	if !cfg.Enabled || cfg.MaxEvents <= 0 {
		return nil
	}
	window := time.Duration(cfg.WindowSeconds) * time.Second
	if window <= 0 {
		window = time.Minute
	}
	return New(cfg.MaxEvents, window)
}
