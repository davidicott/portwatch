package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// Daemon orchestrates periodic port scanning and alerting.
type Daemon struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	dispatch *alert.Dispatcher
}

// New creates a new Daemon with the provided configuration.
func New(cfg *config.Config, s *scanner.Scanner, d *alert.Dispatcher) *Daemon {
	return &Daemon{
		cfg:      cfg,
		scanner:  s,
		dispatch: d,
	}
}

// Run starts the daemon loop, scanning at the configured interval until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon starting (interval: %s)", d.cfg.Interval)

	prev, err := d.scanner.Scan()
	if err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}
	log.Printf("initial scan: %d open ports detected", len(prev))

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopped")
			return nil
		case <-ticker.C:
			curr, err := d.scanner.Scan()
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}

			events := alert.BuildEvents(prev, curr)
			if len(events) > 0 {
				if err := d.dispatch.Notify(ctx, events); err != nil {
					log.Printf("notify error: %v", err)
				}
			}
			prev = curr
		}
	}
}
