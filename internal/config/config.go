package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	// Interval between port scans.
	Interval time.Duration `yaml:"interval"`

	// SnapshotPath is the file used to persist the last known port state.
	SnapshotPath string `yaml:"snapshot_path"`

	// Ignore lists ports and protocols to suppress from alerting.
	Ignore IgnoreConfig `yaml:"ignore"`

	// LogLevel controls verbosity (debug, info, warn, error).
	LogLevel string `yaml:"log_level"`
}

// IgnoreConfig holds exclusion rules.
type IgnoreConfig struct {
	Ports     []int    `yaml:"ports"`
	Protocols []string `yaml:"protocols"`
}

// IgnoreSet is a fast-lookup representation of IgnoreConfig.
type IgnoreSet struct {
	Ports     map[int]struct{}
	Protocols map[string]struct{}
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:     15 * time.Second,
		SnapshotPath: "/var/lib/portwatch/snapshot.json",
		LogLevel:     "info",
	}
}

// Load reads a YAML config file from path.
// If path is empty the default config is returned.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// IgnoreSetFrom builds an IgnoreSet from a Config for O(1) lookups.
func IgnoreSetFrom(cfg Config) IgnoreSet {
	is := IgnoreSet{
		Ports:     make(map[int]struct{}, len(cfg.Ignore.Ports)),
		Protocols: make(map[string]struct{}, len(cfg.Ignore.Protocols)),
	}
	for _, p := range cfg.Ignore.Ports {
		is.Ports[p] = struct{}{}
	}
	for _, proto := range cfg.Ignore.Protocols {
		is.Protocols[proto] = struct{}{}
	}
	return is
}
