package config

import (
	"testing"
)

func TestGoogleChatConfig_Defaults(t *testing.T) {
	d := googlechatDefaults()
	if d.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if d.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", d.WebhookURL)
	}
}

func TestGoogleChatConfig_Fields(t *testing.T) {
	cfg := GoogleChatNotifierConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/ABC/messages?key=xyz",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled true")
	}
	if cfg.WebhookURL == "" {
		t.Error("expected non-empty WebhookURL")
	}
}
