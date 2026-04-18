package config

// TwilioConfig holds configuration for the Twilio SMS notifier.
type TwilioConfig struct {
	Enabled    bool   `yaml:"enabled"`
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
	From       string `yaml:"from"`
	To         string `yaml:"to"`
}

var twilioDefaults = TwilioConfig{
	Enabled: false,
}

func init() {
	registerNotifierDefaults("twilio", twilioDefaults)
}

// twilioFromConfig extracts TwilioConfig from the raw notifier map.
func twilioFromConfig(raw map[string]interface{}) TwilioConfig {
	cfg := twilioDefaults
	if v, ok := raw["enabled"].(bool); ok {
		cfg.Enabled = v
	}
	if v, ok := raw["account_sid"].(string); ok {
		cfg.AccountSID = v
	}
	if v, ok := raw["auth_token"].(string); ok {
		cfg.AuthToken = v
	}
	if v, ok := raw["from"].(string); ok {
		cfg.From = v
	}
	if v, ok := raw["to"].(string); ok {
		cfg.To = v
	}
	return cfg
}
