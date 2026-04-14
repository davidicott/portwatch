package daemon_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/scanner"
)

func defaultCfg(interval time.Duration) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Interval = interval
	return cfg
}

func TestDaemonStopsOnContextCancel(t *testing.T) {
	cfg := defaultCfg(50 * time.Millisecond)
	s := scanner.New(cfg)
	notifier := alert.NewLogNotifier()
	dispatch := alert.NewDispatcher(notifier)

	d := daemon.New(cfg, s, dispatch)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != nil {
		t.Fatalf("expected nil error on context cancel, got: %v", err)
	}
}

func TestDaemonNew(t *testing.T) {
	cfg := defaultCfg(1 * time.Second)
	s := scanner.New(cfg)
	notifier := alert.NewLogNotifier()
	dispatch := alert.NewDispatcher(notifier)

	d := daemon.New(cfg, s, dispatch)
	if d == nil {
		t.Fatal("expected non-nil Daemon")
	}
}
