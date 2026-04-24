package config

// GoogleChatNotifierConfig holds configuration for the Google Chat notifier.
type GoogleChatNotifierConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

func googlechatDefaults() GoogleChatNotifierConfig {
	return GoogleChatNotifierConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

func init() {
	registerNotifierDefaults("googlechat", func(n *NotifiersConfig) {
		if n.GoogleChat == nil {
			defaults := googlechatDefaults()
			n.GoogleChat = &defaults
		}
	})
}
