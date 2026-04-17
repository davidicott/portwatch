package config

// PushoverConfig holds settings for the Pushover notifier.
type PushoverConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	UserKey string `yaml:"user_key"`
}

func pushoverDefaults() PushoverConfig {
	return PushoverConfig{
		Enabled: false,
	}
}

func init() {
	registerNotifierDefaults("pushover", func(n *NotifiersConfig) {
		if n.Pushover == (PushoverConfig{}) {
			n.Pushover = pushoverDefaults()
		}
	})
}
