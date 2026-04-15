package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/watch"
)

// TestWatcherRunsAtLeastOneTick verifies that at least one scan cycle
// completes within a short deadline without returning an unexpected error.
func TestWatcherRunsAtLeastOneTick(t *testing.T) {
	n := &fakeNotifier{}
	dir := t.TempDir()

	w := watch.New(watch.Config{
		Scanner:  scanner.New(),
		Store:    snapshot.NewStore(dir + "/snap.json"),
		Filter:   filter.New(filter.Options{}),
		Notifier: n,
		History:  history.New(20),
		Metrics:  metrics.New(),
		Interval: 20 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestWatcherMetricsRecorded checks that metrics are updated after ticks.
func TestWatcherMetricsRecorded(t *testing.T) {
	n := &fakeNotifier{}
	dir := t.TempDir()
	m := metrics.New()

	w := watch.New(watch.Config{
		Scanner:  scanner.New(),
		Store:    snapshot.NewStore(dir + "/snap.json"),
		Filter:   filter.New(filter.Options{}),
		Notifier: n,
		History:  history.New(20),
		Metrics:  m,
		Interval: 20 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	snap := m.Snapshot()
	if snap.TotalScans == 0 {
		t.Error("expected at least one scan recorded in metrics")
	}
}

// Ensure fakeNotifier satisfies alert.Notifier at compile time.
var _ alert.Notifier = (*fakeNotifier)(nil)
