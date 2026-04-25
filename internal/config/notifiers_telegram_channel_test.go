package config

import (
	"testing"
)

func TestTelegramChannelConfig_Defaults(t *testing.T) {
	defaults := telegramChannelDefaults()

	if defaults.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if defaults.Token != "" {
		t.Errorf("expected empty Token, got %q", defaults.Token)
	}
	if defaults.ChannelID != "" {
		t.Errorf("expected empty ChannelID, got %q", defaults.ChannelID)
	}
}

func TestTelegramChannelConfig_Fields(t *testing.T) {
	cfg := TelegramChannelNotifierConfig{
		Enabled:   true,
		Token:     "bot123:ABC",
		ChannelID: "@myalerts",
	}

	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.Token != "bot123:ABC" {
		t.Errorf("unexpected Token: %q", cfg.Token)
	}
	if cfg.ChannelID != "@myalerts" {
		t.Errorf("unexpected ChannelID: %q", cfg.ChannelID)
	}
}
