package config

func init() {
	registerNotifierDefaults("googlechat", googlechatDefaults())
}

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
