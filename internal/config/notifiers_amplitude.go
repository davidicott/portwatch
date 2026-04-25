package config

// AmplitudeConfig holds configuration for the Amplitude notifier.
type AmplitudeConfig struct {
	Enabled  bool   `yaml:"enabled"`
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}

func amplitudeDefaults() AmplitudeConfig {
	return AmplitudeConfig{
		Enabled:  false,
		Endpoint: "https://api2.amplitude.com/2/httpapi",
	}
}

func init() {
	registerNotifierDefaults("amplitude", func(n *NotifiersConfig) {
		if n.Amplitude.Endpoint == "" {
			n.Amplitude.Endpoint = amplitudeDefaults().Endpoint
		}
	})
}
