package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
	"github.com/sanikasurose/surose-os/internal/analytics"
	"github.com/sanikasurose/surose-os/internal/guestbook"
	"github.com/sanikasurose/surose-os/internal/stats"
	"github.com/sanikasurose/surose-os/internal/ui"
	gossh "golang.org/x/crypto/ssh"
)

const (
	host        = "0.0.0.0"
	port        = "2222"
	hostKeyPath = ".ssh/id_ed25519"

	statsRefreshInterval = 5 * time.Minute
	defaultDataDir       = "data"
)

func main() {
	// Server runs headless under systemd with no TERM set, so lipgloss's
	// default profile detection (which reads the server's own env, not each
	// SSH client's) would fall back to no color. The palette is exact hex,
	// so force truecolor globally rather than per-session detection.
	lipgloss.SetColorProfile(termenv.TrueColor)

	dataDir := envOr("SUROSE_DATA_DIR", defaultDataDir)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatalf("could not create data dir %s: %v", dataDir, err)
	}

	visitorsPath := filepath.Join(dataDir, "visitors.count")
	visitors, err := stats.NewVisitors(visitorsPath)
	if err != nil {
		log.Fatalf("could not open visitors counter at %s: %v", visitorsPath, err)
	}
	log.Printf("visitor counter: %s (current=%d)", visitorsPath, visitors.Current())

	guestbookPath := filepath.Join(dataDir, "guestbook.db")
	gb, err := guestbook.Open(guestbookPath)
	if err != nil {
		log.Fatalf("could not open guestbook db at %s: %v", guestbookPath, err)
	}
	defer gb.Close()
	if n, err := gb.Count(); err == nil {
		log.Printf("guestbook: %s (%d messages)", guestbookPath, n)
	}

	analyticsPath := filepath.Join(dataDir, "sessions.db")
	an, err := analytics.Open(analyticsPath)
	if err != nil {
		log.Fatalf("could not open analytics db at %s: %v", analyticsPath, err)
	}
	defer an.Close()
	if n, err := an.Count(); err == nil {
		log.Printf("analytics: %s (%d sessions recorded)", analyticsPath, n)
	}

	githubUser := envOr("GITHUB_USER", "sanikasurose")
	githubToken := os.Getenv("GITHUB_TOKEN")
	cache := stats.NewCache(githubUser, githubToken)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	refresh := cache.Run(ctx, statsRefreshInterval)
	if githubToken == "" {
		log.Println("warning: GITHUB_TOKEN not set — contribution grid and streak require a token")
	}
	if refresh.CalendarErr != nil {
		log.Printf("warning: contribution calendar fetch failed: %v", refresh.CalendarErr)
	}
	if refresh.StarsErr != nil {
		log.Printf("warning: stars fetch failed: %v", refresh.StarsErr)
	}
	if refresh.OK() {
		snap := cache.Snapshot(0)
		log.Printf("stats cache ready — %d contributions, %dd streak, %d stars",
			snap.TotalContributions, snap.Streak, snap.Stars)
	}

	handler := ui.MakeHandler(ui.Deps{
		Stats:     cache,
		Visitors:  visitors,
		Guestbook: gb,
		Analytics: an,
	})

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
		// No auth handlers — charmbracelet/ssh sets NoClientAuth so visitors
		// connect straight into the app without a password prompt.
		// Advertise post-quantum and modern KEX algorithms so visitors don't
		// see the "not using post-quantum key exchange" warning.
		func(s *ssh.Server) error {
			s.ServerConfigCallback = func(_ ssh.Context) *gossh.ServerConfig {
				return &gossh.ServerConfig{
					Config: gossh.Config{
						KeyExchanges: []string{
							"mlkem768x25519-sha256", // post-quantum ML-KEM hybrid
							"curve25519-sha256",     // modern fallback
						},
					},
				}
			}
			return nil
		},
		wish.WithMiddleware(
			bubbletea.Middleware(handler),
		),
	)
	if err != nil {
		log.Fatalln("could not start server:", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("surose-os listening on %s:%s", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Fatalln("server error:", err)
		}
	}()

	<-done
	log.Println("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := s.Shutdown(shutdownCtx); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Fatalln("shutdown error:", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
