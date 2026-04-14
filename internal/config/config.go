package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"scan_interval"`
	LogFile      string        `yaml:"log_file"`
	AlertOnNew   bool          `yaml:"alert_on_new"`
	AlertOnClose bool          `yaml:"alert_on_close"`
	IgnorePorts  []int         `yaml:"ignore_ports"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval: 30 * time.Second,
		LogFile:      "",
		AlertOnNew:   true,
		AlertOnClose: true,
		IgnorePorts:  []int{},
	}
}

// Load reads a YAML config file from path and merges it over the defaults.
// If path is empty the default config is returned without error.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// IgnoreSet returns the ignore_ports list as a map for O(1) lookup.
func (c *Config) IgnoreSet() map[int]struct{} {
	s := make(map[int]struct{}, len(c.IgnorePorts))
	for _, p := range c.IgnorePorts {
		s[p] = struct{}{}
	}
	return s
}
