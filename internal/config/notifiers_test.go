package config

import (
	"testing"
)

func TestNotifierConfig_Defaults(t *testing.T) {
	var nc NotifierConfig
	if nc.Stdout.Enabled {
		t.Error("stdout should be disabled by default")
	}
	if nc.SNS.Enabled {
		t.Error("sns should be disabled by default")
	}
}

func TestNotifierConfig_VictorOpsFields(t *testing.T) {
	nc := NotifierConfig{
		VictorOps: VictorOpsConfig{
			Enabled:    true,
			URL:        "https://alert.victorops.com/integrations/generic",
			RoutingKey: "default",
		},
	}
	if !nc.VictorOps.Enabled {
		t.Error("expected VictorOps to be enabled")
	}
	if nc.VictorOps.RoutingKey != "default" {
		t.Errorf("unexpected routing key: %s", nc.VictorOps.RoutingKey)
	}
}

func TestNotifierConfig_EmailFields(t *testing.T) {
	nc := NotifierConfig{
		Email: EmailConfig{
			Enabled:  true,
			SMTPHost: "smtp.example.com",
			SMTPPort: 587,
			From:     "alerts@example.com",
			To:       "ops@example.com",
		},
	}
	if nc.Email.SMTPPort != 587 {
		t.Errorf("expected port 587, got %d", nc.Email.SMTPPort)
	}
}

func TestNotifierConfig_SNSFields(t *testing.T) {
	nc := NotifierConfig{
		SNS: SNSConfig{
			Enabled:  true,
			TopicARN: "arn:aws:sns:us-east-1:123456789012:portwatch",
			Region:   "us-east-1",
		},
	}
	if !nc.SNS.Enabled {
		t.Error("expected SNS to be enabled")
	}
	if nc.SNS.Region != "us-east-1" {
		t.Errorf("unexpected region: %s", nc.SNS.Region)
	}
	if nc.SNS.TopicARN == "" {
		t.Error("expected non-empty topic ARN")
	}
}

func TestNotifierConfig_MultipleEnabled(t *testing.T) {
	nc := NotifierConfig{
		Stdout:  StdoutConfig{Enabled: true},
		Slack:   SlackConfig{Enabled: true, WebhookURL: "https://hooks.slack.com/x"},
		SNS:     SNSConfig{Enabled: true, TopicARN: "arn:aws:sns:us-east-1:123:test", Region: "us-east-1"},
		Discord: DiscordConfig{Enabled: false},
	}
	count := 0
	if nc.Stdout.Enabled {
		count++
	}
	if nc.Slack.Enabled {
		count++
	}
	if nc.SNS.Enabled {
		count++
	}
	if nc.Discord.Enabled {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 enabled notifiers, got %d", count)
	}
}
