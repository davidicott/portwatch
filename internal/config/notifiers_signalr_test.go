package config

import "testing"

func TestSignalRConfig_Defaults(t *testing.T) {
	d := signalRDefaults()
	if d.Enabled {
		t.Error("expected disabled by default")
	}
	if d.Hub != "portwatch" {
		t.Errorf("expected hub 'portwatch', got %s", d.Hub)
	}
	if d.Endpoint != "" {
		t.Errorf("expected empty endpoint, got %s", d.Endpoint)
	}
}

func TestSignalRConfig_Fields(t *testing.T) {
	cfg := SignalRConfig{
		Enabled:  true,
		Endpoint: "https://my.signalr.example.com",
		Hub:      "alerts",
		APIKey:   "secret",
	}
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.Hub != "alerts" {
		t.Errorf("unexpected hub: %s", cfg.Hub)
	}
	if cfg.APIKey != "secret" {
		t.Errorf("unexpected api key: %s", cfg.APIKey)
	}
}
