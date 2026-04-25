package config

import (
	"testing"
)

func TestGoogleChatConfig_Defaults(t *testing.T) {
	defaults := googlechatDefaults()

	if defaults.Enabled {
		t.Error("expected enabled to be false by default")
	}
	if defaults.WebhookURL != "" {
		t.Errorf("expected empty webhook URL, got %q", defaults.WebhookURL)
	}
}

func TestGoogleChatConfig_Fields(t *testing.T) {
	cfg := GoogleChatNotifierConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy",
	}

	if !cfg.Enabled {
		t.Error("expected enabled to be true")
	}
	if cfg.WebhookURL == "" {
		t.Error("expected non-empty webhook URL")
	}
}
