package metrics

import (
	"sync"
	"time"
)

// Counters holds runtime statistics for the portwatch daemon.
type Counters struct {
	mu           sync.RWMutex
	ScansTotal   int64
	AlertsTotal  int64
	OpenedPorts  int64
	ClosedPorts  int64
	LastScanTime time.Time
	StartTime    time.Time
}

// Recorder records daemon metrics.
type Recorder struct {
	counters Counters
}

// New returns a new Recorder initialised with the current time.
func New() *Recorder {
	return &Recorder{
		counters: Counters{
			StartTime: time.Now(),
		},
	}
}

// RecordScan increments the scan counter and updates the last scan timestamp.
func (r *Recorder) RecordScan() {
	r.counters.mu.Lock()
	defer r.counters.mu.Unlock()
	r.counters.ScansTotal++
	r.counters.LastScanTime = time.Now()
}

// RecordAlerts increments alert and port-change counters.
func (r *Recorder) RecordAlerts(opened, closed int) {
	r.counters.mu.Lock()
	defer r.counters.mu.Unlock()
	r.counters.AlertsTotal += int64(opened + closed)
	r.counters.OpenedPorts += int64(opened)
	r.counters.ClosedPorts += int64(closed)
}

// Snapshot returns a point-in-time copy of the current counters.
func (r *Recorder) Snapshot() Counters {
	r.counters.mu.RLock()
	defer r.counters.mu.RUnlock()
	return Counters{
		ScansTotal:   r.counters.ScansTotal,
		AlertsTotal:  r.counters.AlertsTotal,
		OpenedPorts:  r.counters.OpenedPorts,
		ClosedPorts:  r.counters.ClosedPorts,
		LastScanTime: r.counters.LastScanTime,
		StartTime:    r.counters.StartTime,
	}
}

// Uptime returns the duration since the recorder was created.
func (r *Recorder) Uptime() time.Duration {
	return time.Since(r.counters.StartTime)
}

// Reset zeroes all counters while preserving the original start time.
func (r *Recorder) Reset() {
	r.counters.mu.Lock()
	defer r.counters.mu.Unlock()
	start := r.counters.StartTime
	r.counters = Counters{
		StartTime: start,
	}
}

// ScanRate returns the average number of scans performed per minute since
// the recorder was created. It returns 0 if no scans have been recorded yet.
func (r *Recorder) ScanRate() float64 {
	r.counters.mu.RLock()
	defer r.counters.mu.RUnlock()
	uptime := time.Since(r.counters.StartTime).Minutes()
	if uptime <= 0 || r.counters.ScansTotal == 0 {
		return 0
	}
	return float64(r.counters.ScansTotal) / uptime
}
