package stats

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Visitors is a thread-safe, file-persisted monotonically-increasing counter.
// Each call to Increment() returns a unique sequential number — exactly the
// visitor number you want to show on the home screen.
//
// Persistence is best-effort: if the write fails (e.g. read-only fs in a test),
// the in-memory counter still works for the current process.
type Visitors struct {
	mu    sync.Mutex
	path  string
	count int
}

// NewVisitors opens (or creates) the counter file at `path`. The parent
// directory is created automatically. Reads the existing count if present.
func NewVisitors(path string) (*Visitors, error) {
	v := &Visitors{path: path}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if data, err := os.ReadFile(path); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			v.count = n
		}
	}
	return v, nil
}

// Increment bumps the counter by 1, persists it, and returns the new value.
// This new value IS the current visitor's number — the first SSH session
// after a fresh install gets back 1.
func (v *Visitors) Increment() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.count++
	n := v.count
	// Best-effort persist. Atomic write via temp file + rename so a crash
	// mid-write doesn't leave an empty file.
	tmp := v.path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strconv.Itoa(n)), 0o644); err == nil {
		_ = os.Rename(tmp, v.path)
	}
	return n
}

// Current returns the counter without incrementing — useful for admin tools.
func (v *Visitors) Current() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.count
}
