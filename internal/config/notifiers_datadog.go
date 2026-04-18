package config

func init() {
	registerNotifierDefaults("datadog", datadogDefaults())
}

func datadogDefaults() map[string]interface{} {
	return map[string]interface{}{
		"enabled": false,
		"api_key": "",
		"host":    "",
	}
}

// DatadogConfig holds configuration for the Datadog notifier.
type DatadogConfig struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"api_key"`
	Host    string `yaml:"host"`
}
