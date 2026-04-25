package config

import "testing"

func TestPagerDutyV2Config_Defaults(t *testing.T) {
	cfg := pagerDutyV2FromConfig(map[string]interface{}{})
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.Endpoint != "https://events.pagerduty.com/v2/enqueue" {
		t.Errorf("endpoint: got %q", cfg.Endpoint)
	}
	if cfg.RoutingKey != "" {
		t.Errorf("routing_key: expected empty, got %q", cfg.RoutingKey)
	}
}

func TestPagerDutyV2Config_Fields(t *testing.T) {
	cfg := pagerDutyV2FromConfig(map[string]interface{}{
		"enabled":     true,
		"routing_key": "abc123",
		"endpoint":    "https://custom.example.com/enqueue",
	})
	if !cfg.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.RoutingKey != "abc123" {
		t.Errorf("routing_key: got %q, want abc123", cfg.RoutingKey)
	}
	if cfg.Endpoint != "https://custom.example.com/enqueue" {
		t.Errorf("endpoint: got %q", cfg.Endpoint)
	}
}

func TestPagerDutyV2Config_DefaultEndpointPreservedWhenEmpty(t *testing.T) {
	cfg := pagerDutyV2FromConfig(map[string]interface{}{
		"endpoint": "",
	})
	if cfg.Endpoint != pagerDutyV2Defaults.Endpoint {
		t.Errorf("expected default endpoint when empty string provided, got %q", cfg.Endpoint)
	}
}
