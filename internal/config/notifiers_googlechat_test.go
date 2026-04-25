package config

import "testing"

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
	cfg := GoogleChatConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/XXX/messages?key=YYY",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.WebhookURL == "" {
		t.Error("expected non-empty WebhookURL")
	}
}
