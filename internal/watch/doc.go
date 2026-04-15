// Package watch provides the Watcher type, which orchestrates the
// port-monitoring lifecycle: scanning the host for open ports, diffing
// against the previous snapshot, dispatching alerts for any changes, and
// persisting the latest snapshot for the next cycle.
//
// # Architecture
//
// The Watcher composes several sub-components:
//
//   - Scanner: discovers currently open ports on the host.
//   - Store: persists and retrieves port snapshots between cycles.
//   - Filter: narrows the set of ports that are considered relevant.
//   - Notifier: dispatches alerts when ports are opened or closed.
//   - History: retains a rolling window of recent change events.
//   - Metrics: records counters and gauges for observability.
//
// # Typical usage
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
//
// Run blocks until ctx is cancelled, executing a scan-diff-notify cycle on
// each tick of the configured Interval. Any error encountered during a cycle
// is logged and the watcher continues to the next interval rather than
// terminating.
package watch
