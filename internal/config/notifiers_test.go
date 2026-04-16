package config

import (
	"testing"
)

func TestNotifierConfig_Defaults(t *testing.T) {
	var nc NotifierConfig
	if nc.Stdout {
		t.Error("stdout should default to false")
	}
	if nc.Webhook != nil {
		t.Error("webhook should default to nil")
	}
	if nc.VictorOps != nil {
		t.Error("victorops should default to nil")
	}
}

func TestNotifierConfig_VictorOpsFields(t *testing.T) {
	nc := NotifierConfig{
		VictorOps: &VictorOpsConfig{
			RoutingKey:      "rk-abc",
			RestEndpointURL: "https://alert.victorops.com/integrations/generic",
		},
	}
	if nc.VictorOps.RoutingKey != "rk-abc" {
		t.Errorf("unexpected routing key: %s", nc.VictorOps.RoutingKey)
	}
	if nc.VictorOps.RestEndpointURL == "" {
		t.Error("rest endpoint url should not be empty")
	}
}

func TestNotifierConfig_EmailFields(t *testing.T) {
	nc := NotifierConfig{
		Email: &EmailConfig{
			SMTPHost: "smtp.example.com",
			SMTPPort: 587,
			From:     "alerts@example.com",
			To:       []string{"ops@example.com"},
		},
	}
	if nc.Email.SMTPPort != 587 {
		t.Errorf("unexpected smtp port: %d", nc.Email.SMTPPort)
	}
	if len(nc.Email.To) != 1 {
		t.Errorf("expected 1 recipient, got %d", len(nc.Email.To))
	}
}

func TestNotifierConfig_MultipleEnabled(t *testing.T) {
	nc := NotifierConfig{
		Stdout:  true,
		Slack:   &SlackConfig{WebhookURL: "https://hooks.slack.com/x"},
		Discord: &DiscordConfig{WebhookURL: "https://discord.com/api/webhooks/x"},
	}
	if !nc.Stdout {
		t.Error("stdout should be true")
	}
	if nc.Slack == nil {
		t.Error("slack should not be nil")
	}
	if nc.Discord == nil {
		t.Error("discord should not be nil")
	}
}
