package config

// ZulipConfig holds configuration for the Zulip notifier.
type ZulipConfig struct {
	Enabled bool   `yaml:"enabled"`
	BaseURL string `yaml:"base_url"`
	Email   string `yaml:"email"`
	APIKey  string `yaml:"api_key"`
	Stream  string `yaml:"stream"`
	Topic   string `yaml:"topic"`
}

func zulipDefaults() ZulipConfig {
	return ZulipConfig{
		Enabled: false,
		Stream:  "general",
		Topic:   "portwatch alerts",
	}
}

func init() {
	registerNotifierDefault("zulip", func(nc *NotifierConfig) {
		if nc.Zulip.Stream == "" {
			nc.Zulip.Stream = zulipDefaults().Stream
		}
		if nc.Zulip.Topic == "" {
			nc.Zulip.Topic = zulipDefaults().Topic
		}
	})
}
