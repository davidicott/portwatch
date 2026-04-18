package config

import "testing"

func TestZulipConfig_Defaults(t *testing.T) {
	d := zulipDefaults()
	if d.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if d.Stream != "general" {
		t.Errorf("expected stream=general, got %q", d.Stream)
	}
	if d.Topic != "portwatch alerts" {
		t.Errorf("expected topic='portwatch alerts', got %q", d.Topic)
	}
}

func TestZulipConfig_Fields(t *testing.T) {
	cfg := ZulipConfig{
		Enabled: true,
		BaseURL: "https://zulip.example.com",
		Email:   "bot@example.com",
		APIKey:  "abc123",
		Stream:  "ops",
		Topic:   "ports",
	}

	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.BaseURL != "https://zulip.example.com" {
		t.Errorf("unexpected BaseURL: %q", cfg.BaseURL)
	}
	if cfg.Email != "bot@example.com" {
		t.Errorf("unexpected Email: %q", cfg.Email)
	}
	if cfg.APIKey != "abc123" {
		t.Errorf("unexpected APIKey: %q", cfg.APIKey)
	}
	if cfg.Stream != "ops" {
		t.Errorf("unexpected Stream: %q", cfg.Stream)
	}
	if cfg.Topic != "ports" {
		t.Errorf("unexpected Topic: %q", cfg.Topic)
	}
}
