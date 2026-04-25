package config

// GoogleChatConfig holds configuration for the Google Chat notifier.
type GoogleChatConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

func googlechatDefaults() GoogleChatConfig {
	return GoogleChatConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

func init() {
	registerNotifierDefaults("googlechat", func(n *NotifierConfig) {
		if n.GoogleChat == nil {
			d := googlechatDefaults()
			n.GoogleChat = &d
		}
	})
}
