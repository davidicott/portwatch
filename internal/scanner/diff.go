package scanner

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
