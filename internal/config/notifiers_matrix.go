package config

// MatrixConfig holds settings for the Matrix notifier.
type MatrixConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Homeserver string `yaml:"homeserver"`
	Token      string `yaml:"token"`
	RoomID     string `yaml:"room_id"`
}

func matrixDefaults() MatrixConfig {
	return MatrixConfig{
		Enabled:    false,
		Homeserver: "https://matrix.org",
	}
}

func init() {
	notifierDefaults["matrix"] = func(nc *NotifierConfig) {
		if nc.Matrix.Homeserver == "" {
			nc.Matrix = matrixDefaults()
		}
	}
}
