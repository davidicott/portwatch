package config

import (
	"testing"
)

func TestMattermostConfig_Defaults(t *testing.T) {
	cfg := mattermostDefaults()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Channel != "" {
		t.Errorf("expected empty channel, got %q", cfg.Channel)
	}
	if cfg.WebhookURL != "" {
		t.Errorf("expected empty webhook_url, got %q", cfg.WebhookURL)
	}
}

func TestMattermostConfig_Fields(t *testing.T) {
	cfg := MattermostConfig{
		Enabled:    true,
		WebhookURL: "https://mattermost.example.com/hooks/abc",
		Channel:    "#ops",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled true")
	}
	if cfg.WebhookURL == "" {
		t.Error("expected non-empty WebhookURL")
	}
	if cfg.Channel != "#ops" {
		t.Errorf("unexpected channel: %q", cfg.Channel)
	}
}
