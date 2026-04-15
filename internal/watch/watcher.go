package watch

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Watcher orchestrates a single scan-diff-alert cycle.
type Watcher struct {
	scanner  *scanner.Scanner
	store    *snapshot.Store
	filter   *filter.Filter
	notifier alert.Notifier
	history  *history.Ring
	metrics  *metrics.Metrics
	interval time.Duration
}

// Config holds dependencies for constructing a Watcher.
type Config struct {
	Scanner  *scanner.Scanner
	Store    *snapshot.Store
	Filter   *filter.Filter
	Notifier alert.Notifier
	History  *history.Ring
	Metrics  *metrics.Metrics
	Interval time.Duration
}

// New creates a Watcher from the provided Config.
func New(cfg Config) *Watcher {
	return &Watcher{
		scanner:  cfg.Scanner,
		store:    cfg.Store,
		filter:   cfg.Filter,
		notifier: cfg.Notifier,
		history:  cfg.History,
		metrics:  cfg.Metrics,
		interval: cfg.Interval,
	}
}

// Run starts the watch loop, ticking at w.interval until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				return err
			}
		}
	}
}

// tick performs one scan cycle.
func (w *Watcher) tick(ctx context.Context) error {
	current, err := w.scanner.Scan()
	if err != nil {
		return err
	}
	current = w.filter.Apply(current)

	previous, _ := w.store.Load()
	events := alert.BuildEvents(previous, current)

	w.metrics.RecordScan(len(events))

	for _, e := range events {
		w.history.Record(e)
	}

	if len(events) > 0 {
		w.notifier.Notify(ctx, events)
	}

	return w.store.Save(current)
}
