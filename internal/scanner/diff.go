package scanner

import "sort"

// Diff holds the changes detected between two port snapshots.
type Diff struct {
	Opened []Port
	Closed []Port
}

// HasChanges returns true if any ports were opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare computes the difference between a previous and current
// set of ports, returning newly opened and newly closed ports.
// Results are sorted by port number for deterministic output.
func Compare(previous, current []Port) Diff {
	prevMap := toMap(previous)
	currMap := toMap(current)

	var diff Diff

	for key, port := range currMap {
		if _, exists := prevMap[key]; !exists {
			diff.Opened = append(diff.Opened, port)
		}
	}

	for key, port := range prevMap {
		if _, exists := currMap[key]; !exists {
			diff.Closed = append(diff.Closed, port)
		}
	}

	sort.Slice(diff.Opened, func(i, j int) bool { return diff.Opened[i].Number < diff.Opened[j].Number })
	sort.Slice(diff.Closed, func(i, j int) bool { return diff.Closed[i].Number < diff.Closed[j].Number })

	return diff
}

// toMap converts a slice of Ports into a map keyed by Port.Key().
func toMap(ports []Port) map[string]Port {
	m := make(map[string]Port, len(ports))
	for _, p := range ports {
		m[p.Key()] = p
	}
	return m
}
