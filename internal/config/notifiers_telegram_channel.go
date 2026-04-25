package config

// TelegramChannelNotifierConfig holds settings for the Telegram channel notifier.
type TelegramChannelNotifierConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channel_id"`
}

func telegramChannelDefaults() TelegramChannelNotifierConfig {
	return TelegramChannelNotifierConfig{
		Enabled:   false,
		Token:     "",
		ChannelID: "",
	}
}

func init() {
	registerNotifierDefaults("telegram_channel", func(n *NotifiersConfig) {
		if n.TelegramChannel == (TelegramChannelNotifierConfig{}) {
			n.TelegramChannel = telegramChannelDefaults()
		}
	})
}
