package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.ScanInterval)
	}
	if !cfg.AlertOnNew || !cfg.AlertOnClose {
		t.Error("expected alert flags to be true by default")
	}
	if len(cfg.IgnorePorts) != 0 {
		t.Error("expected empty ignore_ports by default")
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadValidFile(t *testing.T) {
	const yaml = `
scan_interval: 10s
log_file: /tmp/portwatch.log
alert_on_new: true
alert_on_close: false
ignore_ports:
  - 22
  - 80
`
	f, err := os.CreateTemp("", "portwatch-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(yaml)
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.ScanInterval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.ScanInterval)
	}
	if cfg.LogFile != "/tmp/portwatch.log" {
		t.Errorf("unexpected log_file: %s", cfg.LogFile)
	}
	if cfg.AlertOnClose {
		t.Error("expected alert_on_close to be false")
	}
	if len(cfg.IgnorePorts) != 2 {
		t.Errorf("expected 2 ignore_ports, got %d", len(cfg.IgnorePorts))
	}
}

func TestIgnoreSet(t *testing.T) {
	cfg := &Config{IgnorePorts: []int{22, 443, 8080}}
	set := cfg.IgnoreSet()
	for _, p := range []int{22, 443, 8080} {
		if _, ok := set[p]; !ok {
			t.Errorf("port %d missing from ignore set", p)
		}
	}
	if _, ok := set[80]; ok {
		t.Error("port 80 should not be in ignore set")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
