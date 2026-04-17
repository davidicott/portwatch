package config

import (
	"testing"
)

func TestGoogleChatConfig_Defaults(t *testing.T) {
	cfg := googlechatDefaults()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", cfg.WebhookURL)
	}
}

func TestGoogleChatConfig_Fields(t *testing.T) {
	cfg := GoogleChatConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/xyz/messages?key=abc",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.WebhookURL == "" {
		t.Error("expected non-empty WebhookURL")
	}
}
