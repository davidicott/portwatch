package config

// SignalRConfig holds configuration for the SignalR notifier.
type SignalRConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Hub      string `yaml:"hub"`
	APIKey   string `yaml:"api_key"`
}

func signalRDefaults() SignalRConfig {
	return SignalRConfig{
		Enabled:  false,
		Endpoint: "",
		Hub:      "portwatch",
		APIKey:   "",
	}
}

func init() {
	defaultNotifierInits = append(defaultNotifierInits, func(n *NotifiersConfig) {
		if n.SignalR == nil {
			d := signalRDefaults()
			n.SignalR = &d
		}
	})
}
