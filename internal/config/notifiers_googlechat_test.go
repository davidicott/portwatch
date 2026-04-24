package config

import "testing"

func TestGoogleChatConfig_Defaults(t *testing.T) {
	defaults := googlechatDefaults()

	if defaults["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", defaults["enabled"])
	}
	if defaults["webhook_url"] != "" {
		t.Errorf("expected empty webhook_url, got %v", defaults["webhook_url"])
	}
}

func TestGoogleChatConfig_Fields(t *testing.T) {
	cfg := GoogleChatConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/abc/messages?key=xyz",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.WebhookURL == "" {
		t.Error("expected non-empty WebhookURL")
	}
}
