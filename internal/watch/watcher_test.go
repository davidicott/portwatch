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
	"github.com/user/portwatch/internal/watch"
)

type fakeNotifier struct {
	called int
	last   []alert.Event
}

func (f *fakeNotifier) Notify(_ context.Context, events []alert.Event) {
	f.called++
	f.last = events
}

func defaultWatcherConfig(t *testing.T, notifier alert.Notifier) watch.Config {
	t.Helper()
	dir := t.TempDir()
	return watch.Config{
		Scanner:  scanner.New(),
		Store:    snapshot.NewStore(dir + "/snap.json"),
		Filter:   filter.New(filter.Options{}),
		Notifier: notifier,
		History:  history.New(10),
		Metrics:  metrics.New(),
		Interval: 50 * time.Millisecond,
	}
}

func TestWatcherStopsOnContextCancel(t *testing.T) {
	n := &fakeNotifier{}
	cfg := defaultWatcherConfig(t, n)
	w := watch.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("expected context error, got %v", err)
	}
}

func TestWatcherNew(t *testing.T) {
	n := &fakeNotifier{}
	cfg := defaultWatcherConfig(t, n)
	w := watch.New(cfg)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}
