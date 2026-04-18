package config

// NewRelicConfig holds configuration for the New Relic notifier.
type NewRelicConfig struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"api_key"`
	URL     string `yaml:"url"`
}

func newRelicDefaults() NewRelicConfig {
	return NewRelicConfig{
		Enabled: false,
		URL:     "https://log-api.newrelic.com/log/v1",
	}
}

func init() {
	registerNotifierDefault("newrelic", func(n *NotifierConfig) {
		if n.NewRelic.URL == "" {
			n.NewRelic.URL = newRelicDefaults().URL
		}
	})
}
