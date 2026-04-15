// Package watch provides the Watcher type, which orchestrates the
// port-monitoring lifecycle: scanning the host for open ports, diffing
// against the previous snapshot, dispatching alerts for any changes, and
// persisting the latest snapshot for the next cycle.
//
// Typical usage:
//
//	w := watch.New(watch.Config{
//		Scanner:  scanner.New(),
//		Store:    snapshot.NewStore(path),
//		Filter:   filter.New(opts),
//		Notifier: notifier,
//		History:  history.New(capacity),
//		Metrics:  metrics.New(),
//		Interval: 30 * time.Second,
//	})
//	w.Run(ctx)
package watch
