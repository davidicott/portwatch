package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makePorts() []scanner.Port {
	return []scanner.Port{
		{Protocol: "tcp", Port: 22, PID: 100, Process: "sshd"},
		{Protocol: "tcp", Port: 80, PID: 200, Process: "nginx"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state", "ports.json")
	store := snapshot.NewStore(path)

	ports := makePorts()
	if err := store.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
	if time.Since(snap.CapturedAt) > 5*time.Second {
		t.Error("CapturedAt is unexpectedly old")
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(filepath.Join(dir, "nonexistent.json"))

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports slice, got %d", len(snap.Ports))
	}
}

func TestSaveOverwritesPrevious(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(filepath.Join(dir, "ports.json"))

	_ = store.Save(makePorts())
	newPorts := []scanner.Port{{Protocol: "udp", Port: 53, PID: 300, Process: "systemd-resolved"}}
	if err := store.Save(newPorts); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != 1 {
		t.Fatalf("expected 1 port after overwrite, got %d", len(snap.Ports))
	}
	if snap.Ports[0].Port != 53 {
		t.Errorf("expected port 53, got %d", snap.Ports[0].Port)
	}
}

func TestSaveCreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	deepPath := filepath.Join(dir, "a", "b", "c", "ports.json")
	store := snapshot.NewStore(deepPath)

	if err := store.Save(makePorts()); err != nil {
		t.Fatalf("Save to deep path: %v", err)
	}
	if _, err := os.Stat(deepPath); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
