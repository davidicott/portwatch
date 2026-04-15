package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a persisted set of ports captured at a point in time.
type Snapshot struct {
	CapturedAt time.Time        `json:"captured_at"`
	Ports      []scanner.Port   `json:"ports"`
}

// Store handles reading and writing snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a Store that persists snapshots at the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the given ports as a new snapshot, overwriting any existing file.
func (s *Store) Save(ports []scanner.Port) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Ports:      ports,
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), ".portwatch-snap-*")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	f.Close()
	return os.Rename(tmpName, s.path)
}

// Load reads the most recent snapshot from disk.
// Returns an empty Snapshot (no ports) if the file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
