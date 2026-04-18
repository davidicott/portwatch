package config

import (
	"testing"
)

func TestDatadogConfig_Defaults(t *testing.T) {
	defaults := datadogDefaults()
	if defaults["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", defaults["enabled"])
	}
	if defaults["api_key"] != "" {
		t.Errorf("expected empty api_key, got %v", defaults["api_key"])
	}
	if defaults["host"] != "" {
		t.Errorf("expected empty host, got %v", defaults["host"])
	}
}

func TestDatadogConfig_Fields(t *testing.T) {
	cfg := DatadogConfig{
		Enabled: true,
		APIKey:  "abc123",
		Host:    "prod-server",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.APIKey != "abc123" {
		t.Errorf("expected APIKey=abc123, got %s", cfg.APIKey)
	}
	if cfg.Host != "prod-server" {
		t.Errorf("expected Host=prod-server, got %s", cfg.Host)
	}
}
