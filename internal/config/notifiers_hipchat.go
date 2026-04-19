package config

func init() {
	registerNotifierDefaults("hipchat", hipChatDefaults())
}

// HipChatConfig holds configuration for the HipChat notifier.
type HipChatConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
	RoomID    string `yaml:"room_id"`
	Token     string `yaml:"token"`
}

func hipChatDefaults() HipChatConfig {
	return HipChatConfig{
		Enabled:   false,
		ServerURL: "https://api.hipchat.com",
	}
}
