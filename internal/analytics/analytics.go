// Package analytics records anonymous session metadata (connect/disconnect
// times, screens visited) to a local SQLite file so the site owner can see
// aggregate usage — total visits, unique visitors, average session length,
// most-read screens. Visitor identifiers are never stored raw.
package analytics

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite" // pure-Go driver, no cgo
)

// Store wraps the SQLite connection.
type Store struct {
	db   *sql.DB
	salt string
}

// Open opens (or creates) the SQLite file at `path` and loads/generates the
// hashing salt in the same directory. The parent directory must already exist.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("analytics: open %s: %w", path, err)
	}
	db.SetMaxOpenConns(4)

	pragmas := []string{
		`PRAGMA journal_mode = WAL`,
		`PRAGMA synchronous = NORMAL`,
		`PRAGMA busy_timeout = 5000`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, fmt.Errorf("analytics: %s: %w", p, err)
		}
	}

	schema := `
		CREATE TABLE IF NOT EXISTS visits (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			visitor_hash    TEXT    NOT NULL,
			connected_at    INTEGER NOT NULL,
			disconnected_at INTEGER,
			screens_visited TEXT
		);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("analytics: schema: %w", err)
	}

	saltPath := saltPathFor(path)
	salt, err := loadOrCreateSalt(saltPath)
	if err != nil {
		return nil, fmt.Errorf("analytics: salt: %w", err)
	}

	return &Store{db: db, salt: salt}, nil
}

// Close closes the underlying database.
func (s *Store) Close() error { return s.db.Close() }

// Count returns the total number of recorded sessions — used for the startup log line.
func (s *Store) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM visits`).Scan(&n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// RecordConnect hashes `identifier` with the server-side salt and inserts a
// new open session row. Returns the row id to pass to RecordDisconnect later.
func (s *Store) RecordConnect(identifier string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hash := s.hash(identifier)
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO visits(visitor_hash, connected_at, disconnected_at, screens_visited)
		 VALUES (?, ?, NULL, '')`,
		hash, time.Now().Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("analytics: record connect: %w", err)
	}
	return res.LastInsertId()
}

// RecordDisconnect closes out a session row with the disconnect time and the
// comma-joined list of screens visited during the session.
func (s *Store) RecordDisconnect(sessionID int64, screens []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx,
		`UPDATE visits SET disconnected_at = ?, screens_visited = ? WHERE id = ?`,
		time.Now().Unix(), strings.Join(screens, ","), sessionID,
	)
	if err != nil {
		return fmt.Errorf("analytics: record disconnect: %w", err)
	}
	return nil
}

func (s *Store) hash(identifier string) string {
	sum := sha256.Sum256([]byte(s.salt + identifier))
	return hex.EncodeToString(sum[:])
}

func saltPathFor(dbPath string) string {
	dir := dbPath
	if idx := strings.LastIndexByte(dbPath, '/'); idx != -1 {
		dir = dbPath[:idx]
	} else {
		dir = "."
	}
	return dir + "/analytics.salt"
}

// loadOrCreateSalt prefers SUROSE_ANALYTICS_SALT if set; otherwise it loads
// (or generates and persists) a random salt file so hashes stay stable across
// restarts even without the env var.
func loadOrCreateSalt(saltPath string) (string, error) {
	if v := os.Getenv("SUROSE_ANALYTICS_SALT"); v != "" {
		return v, nil
	}

	if data, err := os.ReadFile(saltPath); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	salt := hex.EncodeToString(buf)

	if err := os.WriteFile(saltPath, []byte(salt), 0o600); err != nil {
		return "", fmt.Errorf("persist salt: %w", err)
	}
	return salt, nil
}
