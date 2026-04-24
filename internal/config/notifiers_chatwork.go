package config

// ChatworkConfig holds configuration for the Chatwork notifier.
type ChatworkConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	RoomID  string `yaml:"room_id"`
}

func chatworkDefaults() ChatworkConfig {
	return ChatworkConfig{
		Enabled: false,
	}
}

func init() {
	registerNotifierDefaults("chatwork", func(n *NotifiersConfig) {
		if n.Chatwork == (ChatworkConfig{}) {
			n.Chatwork = chatworkDefaults()
		}
	})
}
