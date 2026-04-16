package ratelimit_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ratelimit"
)

func TestFromConfig_Disabled(t *testing.T) {
	cfg := config.RateLimitConfig{Enabled: false, MaxEvents: 5, WindowSeconds: 60}
	if ratelimit.FromConfig(cfg) != nil {
		t.Fatal("expected nil limiter when disabled")
	}
}

func TestFromConfig_ZeroMax(t *testing.T) {
	cfg := config.RateLimitConfig{Enabled: true, MaxEvents: 0, WindowSeconds: 60}
	if ratelimit.FromConfig(cfg) != nil {
		t.Fatal("expected nil limiter when max is zero")
	}
}

func TestFromConfig_Valid(t *testing.T) {
	cfg := config.RateLimitConfig{Enabled: true, MaxEvents: 3, WindowSeconds: 30}
	l := ratelimit.FromConfig(cfg)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	for i := 0; i < 3; i++ {
		if !l.Allow("key") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if l.Allow("key") {
		t.Fatal("expected deny after max events")
	}
}

func TestFromConfig_DefaultWindow(t *testing.T) {
	cfg := config.RateLimitConfig{Enabled: true, MaxEvents: 2, WindowSeconds: 0}
	l := ratelimit.FromConfig(cfg)
	if l == nil {
		t.Fatal("expected non-nil limiter with default window")
	}
}
