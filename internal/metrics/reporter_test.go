package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReporterContainsKeys(t *testing.T) {
	r := New()
	r.RecordScan()
	r.RecordAlerts(2, 1)

	var buf bytes.Buffer
	rp := NewReporter(&buf)
	if err := rp.Report(r.Snapshot()); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}

	out := buf.String()
	expectedKeys := []string{
		"uptime", "scans_total", "alerts_total",
		"ports_opened", "ports_closed", "last_scan",
	}
	for _, key := range expectedKeys {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain %q, got:\n%s", key, out)
		}
	}
}

func TestReporterNoScans(t *testing.T) {
	r := New()
	var buf bytes.Buffer
	rp := NewReporter(&buf)
	if err := rp.Report(r.Snapshot()); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "never") {
		t.Error("expected 'never' for last_scan when no scans recorded")
	}
}

func TestReporterLastScanFormatted(t *testing.T) {
	r := New()
	r.RecordScan()
	time.Sleep(5 * time.Millisecond)

	var buf bytes.Buffer
	rp := NewReporter(&buf)
	if err := rp.Report(r.Snapshot()); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}
	if strings.Contains(buf.String(), "never") {
		t.Error("expected a real timestamp, not 'never'")
	}
}

func TestReporterAlertCounts(t *testing.T) {
	r := New()
	r.RecordAlerts(3, 5)

	var buf bytes.Buffer
	rp := NewReporter(&buf)
	if err := rp.Report(r.Snapshot()); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "3") {
		t.Errorf("expected output to contain ports_opened count 3, got:\n%s", out)
	}
	if !strings.Contains(out, "5") {
		t.Errorf("expected output to contain ports_closed count 5, got:\n%s", out)
	}
}
