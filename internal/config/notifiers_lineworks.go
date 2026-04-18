package config

// LineWorksConfig holds configuration for the LINE WORKS notifier.
type LineWorksConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

func lineWorksDefaults() LineWorksConfig {
	return LineWorksConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

func init() {
	notifierDefaults["lineworks"] = func(nc *NotifierConfig) {
		if nc.LineWorks == (LineWorksConfig{}) {
			nc.LineWorks = lineWorksDefaults()
		}
	}
}
