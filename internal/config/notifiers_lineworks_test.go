package config

import "testing"

func TestLineWorksConfig_Defaults(t *testing.T) {
	defaults := lineWorksDefaults()
	if defaults.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if defaults.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", defaults.WebhookURL)
	}
}

func TestLineWorksConfig_Fields(t *testing.T) {
	cfg := LineWorksConfig{
		Enabled:    true,
		WebhookURL: "https://hooks.worksmobile.com/r/abc123",
	}
	if !cfg.Enabled {
		t.Error("expected Enabled true")
	}
	if cfg.WebhookURL != "https://hooks.worksmobile.com/r/abc123" {
		t.Errorf("unexpected WebhookURL: %s", cfg.WebhookURL)
	}
}
