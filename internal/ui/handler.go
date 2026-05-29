package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/sanikasurose/surose-os/internal/guestbook"
	"github.com/sanikasurose/surose-os/internal/stats"
)

// Deps bundles cross-cutting dependencies that the Handler injects into each
// new session's root model.
type Deps struct {
	Stats     *stats.Cache     // background-refreshed GitHub stats
	Visitors  *stats.Visitors  // persistent per-session counter
	Guestbook *guestbook.Store // SQLite-backed message store (shared across sessions)
}

// MakeHandler returns a Wish bubbletea handler that captures a per-session
// Snapshot at connect time. The visitor counter is incremented exactly once
// per session (here), not on every render. The guestbook store is shared —
// all sessions see the same messages.
//
// Stats are read from the cache atomically after the server's initial refresh.
// If GitHub data is unavailable, Snapshot.Available will be false and the home
// screen shows a ghost-text message instead of an empty grid animation.
func MakeHandler(d Deps) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		visitorN := d.Visitors.Increment()
		snap := d.Stats.Snapshot(visitorN)
		return NewRootModel(snap, d.Guestbook), []tea.ProgramOption{tea.WithAltScreen()}
	}
}

// Handler is kept for backwards compatibility with the existing Wish
// middleware wiring. It builds a no-op model with no stats or store.
//
// Deprecated: replace `bubbletea.Middleware(ui.Handler)` with
// `bubbletea.Middleware(ui.MakeHandler(ui.Deps{...}))`.
func Handler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return NewRootModel(stats.Snapshot{}, nil), []tea.ProgramOption{tea.WithAltScreen()}
}
