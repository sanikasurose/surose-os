package ui

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/sanikasurose/surose-os/internal/analytics"
	"github.com/sanikasurose/surose-os/internal/guestbook"
	"github.com/sanikasurose/surose-os/internal/stats"
	gossh "golang.org/x/crypto/ssh"
)

// Deps bundles cross-cutting dependencies that the Handler injects into each
// new session's root model.
type Deps struct {
	Stats     *stats.Cache     // background-refreshed GitHub stats
	Visitors  *stats.Visitors  // persistent per-session counter
	Guestbook *guestbook.Store // SQLite-backed message store (shared across sessions)
	Analytics *analytics.Store // anonymous session-metadata store (shared across sessions)
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

		var recorder *screenRecorder
		if d.Analytics != nil {
			recorder = newScreenRecorder()
			sessionID, err := d.Analytics.RecordConnect(sessionIdentifier(s))
			if err != nil {
				log.Printf("analytics: record connect failed: %v", err)
			} else {
				go func() {
					<-s.Context().Done()
					if err := d.Analytics.RecordDisconnect(sessionID, recorder.Snapshot()); err != nil {
						log.Printf("analytics: record disconnect failed: %v", err)
					}
				}()
			}
		}

		return NewRootModel(snap, d.Guestbook, recorder), []tea.ProgramOption{tea.WithAltScreen()}
	}
}

// sessionIdentifier returns a stable string to hash for analytics — the
// marshaled public key when the client presented one, otherwise the remote
// IP (port stripped, since the OS assigns a new ephemeral source port on
// every connection, which would make every visit look "unique"). Never
// stored raw; only ever passed to analytics.RecordConnect, which hashes it
// before touching disk.
func sessionIdentifier(s ssh.Session) string {
	if pub := s.PublicKey(); pub != nil {
		return string(gossh.MarshalAuthorizedKey(pub))
	}
	if host, _, err := net.SplitHostPort(s.RemoteAddr().String()); err == nil {
		return host
	}
	return s.RemoteAddr().String()
}

// Handler is kept for backwards compatibility with the existing Wish
// middleware wiring. It builds a no-op model with no stats or store.
//
// Deprecated: replace `bubbletea.Middleware(ui.Handler)` with
// `bubbletea.Middleware(ui.MakeHandler(ui.Deps{...}))`.
func Handler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return NewRootModel(stats.Snapshot{}, nil, nil), []tea.ProgramOption{tea.WithAltScreen()}
}
