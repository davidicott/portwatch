package config

// WebexNotifierConfig holds configuration for the Cisco Webex notifier.
type WebexNotifierConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	RoomID  string `yaml:"room_id"`
}

func webexDefaults() WebexNotifierConfig {
	return WebexNotifierConfig{
		Enabled: false,
	}
}

func init() {
	notifierDefaults["webex"] = func(nc *NotifierConfig) {
		if nc.Webex == nil {
			d := webexDefaults()
			nc.Webex = &d
		}
	}
}
