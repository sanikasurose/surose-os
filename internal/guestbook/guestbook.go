// Package guestbook persists visitor messages to a local SQLite file.
// One table, append-only writes, indexed by created_at for fast "recent" lookups.
//
// Schema is created on Open() if absent — no migration system needed for v1.
package guestbook

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	_ "modernc.org/sqlite" // pure-Go driver, no cgo
)

// MaxBodyChars is the per-message character limit (runes, not bytes).
// Matches the spec: 120 chars max to prevent abuse.
const MaxBodyChars = 120

// Message is one guestbook entry.
type Message struct {
	ID        int64
	Handle    string    // e.g. "visitor_047"
	Body      string    // raw user text, <= MaxBodyChars runes
	CreatedAt time.Time // UTC
}

// Store wraps the SQLite connection.
type Store struct {
	db *sql.DB
}

// Open opens (or creates) the SQLite file at `path`. The parent directory
// must already exist. Sets pragmas for crash-safety + concurrency:
//   - WAL journaling — readers don't block writers (helpful for many
//     concurrent SSH sessions on the Insert path).
//   - synchronous=NORMAL — durable enough for a guestbook, faster than FULL.
//   - busy_timeout=5000 — retry briefly on lock contention rather than
//     erroring out immediately.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("guestbook: open %s: %w", path, err)
	}
	// SQLite is single-connection-friendly; cap at a few to avoid
	// SQLITE_BUSY under load.
	db.SetMaxOpenConns(4)

	pragmas := []string{
		`PRAGMA journal_mode = WAL`,
		`PRAGMA synchronous = NORMAL`,
		`PRAGMA busy_timeout = 5000`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, fmt.Errorf("guestbook: %s: %w", p, err)
		}
	}

	schema := `
		CREATE TABLE IF NOT EXISTS messages (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			handle     TEXT    NOT NULL,
			body       TEXT    NOT NULL,
			created_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_messages_created_at
			ON messages(created_at DESC);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("guestbook: schema: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the underlying database.
func (s *Store) Close() error { return s.db.Close() }

// Insert validates the body, then writes a new message. Returns an
// ErrEmpty/ErrTooLong sentinel for the two user-visible failure modes —
// the UI surfaces these as status text under the input.
func (s *Store) Insert(handle, body string) error {
	body = strings.TrimSpace(body)
	if body == "" {
		return ErrEmpty
	}
	if utf8.RuneCountInString(body) > MaxBodyChars {
		return ErrTooLong
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO messages(handle, body, created_at) VALUES (?, ?, ?)`,
		handle, body, time.Now().UTC().Unix(),
	)
	if err != nil {
		return fmt.Errorf("guestbook: insert: %w", err)
	}
	return nil
}

// Recent returns up to `limit` messages, most recent first.
func (s *Store) Recent(limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 50
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, handle, body, created_at
		   FROM messages
		  ORDER BY created_at DESC, id DESC
		  LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("guestbook: recent: %w", err)
	}
	defer rows.Close()

	var out []Message
	for rows.Next() {
		var m Message
		var ts int64
		if err := rows.Scan(&m.ID, &m.Handle, &m.Body, &ts); err != nil {
			return nil, fmt.Errorf("guestbook: scan: %w", err)
		}
		m.CreatedAt = time.Unix(ts, 0).UTC()
		out = append(out, m)
	}
	return out, rows.Err()
}

// Count returns the total number of messages — used in the screen header.
func (s *Store) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM messages`).Scan(&n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Sentinel errors. UI matches these to render friendly status text.
var (
	ErrEmpty   = fmt.Errorf("message is empty")
	ErrTooLong = fmt.Errorf("message is too long (max %d chars)", MaxBodyChars)
)
