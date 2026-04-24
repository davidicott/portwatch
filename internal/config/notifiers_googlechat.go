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
	registerNotifierDefaults("googlechat", func(nc *NotifiersConfig) {
		if nc.GoogleChat == (GoogleChatNotifierConfig{}) {
			nc.GoogleChat = googlechatDefaults()
		}
	})
}
