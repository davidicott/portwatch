// Package metrics provides lightweight runtime counters and reporting
// for the portwatch daemon.
//
// A Recorder accumulates statistics such as the number of scans performed
// and port-change events detected. Snapshots of the counters can be taken
// at any time for safe, lock-free reads.
//
// A Reporter formats a Counters snapshot into a human-readable tabular
// summary suitable for logging or writing to a status file.
//
// Typical usage:
//
//	rec := metrics.New()
//
//	// inside the scan loop:
//	rec.RecordScan()
//	rec.RecordAlerts(openedCount, closedCount)
//
//	// periodically or on SIGUSR1:
//	rp := metrics.NewReporter(os.Stdout)
//	rp.Report(rec.Snapshot())
package metrics
