package config

func init() {
	registerNotifierDefaults("googlechat", googlechatDefaults())
}

func googlechatDefaults() map[string]interface{} {
	return map[string]interface{}{
		"enabled":     false,
		"webhook_url": "",
	}
}

// GoogleChatConfig holds configuration for the Google Chat notifier.
type GoogleChatConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}
