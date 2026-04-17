package config

// MattermostConfig holds configuration for the Mattermost notifier.
type MattermostConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
}

// mattermostDefaults returns a MattermostConfig with sensible defaults.
func mattermostDefaults() MattermostConfig {
	return MattermostConfig{
		Enabled: false,
		Channel: "",
	}
}
