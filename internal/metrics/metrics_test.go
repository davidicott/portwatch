package metrics

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("expected non-nil Recorder")
	}
	s := r.Snapshot()
	if s.ScansTotal != 0 {
		t.Errorf("expected 0 scans, got %d", s.ScansTotal)
	}
	if s.StartTime.IsZero() {
		t.Error("expected non-zero start time")
	}
}

func TestRecordScan(t *testing.T) {
	r := New()
	r.RecordScan()
	r.RecordScan()
	s := r.Snapshot()
	if s.ScansTotal != 2 {
		t.Errorf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.LastScanTime.IsZero() {
		t.Error("expected non-zero LastScanTime after scan")
	}
}

func TestRecordAlerts(t *testing.T) {
	r := New()
	r.RecordAlerts(3, 1)
	s := r.Snapshot()
	if s.OpenedPorts != 3 {
		t.Errorf("expected 3 opened ports, got %d", s.OpenedPorts)
	}
	if s.ClosedPorts != 1 {
		t.Errorf("expected 1 closed port, got %d", s.ClosedPorts)
	}
	if s.AlertsTotal != 4 {
		t.Errorf("expected 4 total alerts, got %d", s.AlertsTotal)
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	r := New()
	r.RecordScan()
	s1 := r.Snapshot()
	r.RecordScan()
	s2 := r.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Error("snapshots should be independent copies")
	}
}

func TestUptime(t *testing.T) {
	r := New()
	time.Sleep(10 * time.Millisecond)
	if r.Uptime() < 10*time.Millisecond {
		t.Error("uptime should be at least 10ms")
	}
}
