package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sanikasurose/surose-os/internal/content"
	"github.com/sanikasurose/surose-os/internal/stats"
)

// HomeModel — the entry screen. Renders:
//   1. Screen header (▌ home               hisanika ╱ home · visitor #NNN)
//   2. Name + tagline
//   3. Featured rail (2 work entries + 2 project highlights)
//   4. Nav menu (5 items)
//   5. Contribution grid — full GitHub-style heatmap, animates in on load
//   6. Legend (less … more  · total contributions  · streak)
type HomeModel struct {
	cursor       int
	snapshot     stats.Snapshot
	gridRevealed int // animation progress: 0..52 columns drawn
}

var homeItems = []struct {
	key   string
	label string
}{
	{"1", "about"},
	{"2", "experience"},
	{"3", "projects"},
	{"4", "guestbook"},
	{"5", "contact"},
}

type featuredEntry struct {
	num    string // "01"…"04"
	title  string
	meta   string // right-aligned subtitle
	detail string
}

var homeFeatured = []featuredEntry{
	{"01", "TranQuility Inc.", "may 2026 — present", "software engineer intern · toronto, on · remote"},
	{"02", "McMaster, Faculty of Engineering — Dr. Charles Welch", "may 2026 — present", "undergraduate research assistant · NLP, applied ML"},
	{"03", "Piratech (Hacktech 2026 @ Caltech)", "apr 2026", "AI-powered security analysis — automated vulnerability detection across codebases"},
	{"04", "neurosnake-rl", "2026", "CNN-based DQN agent · full experiment-tracking pipeline"},
}

// NewHomeModel constructs the home screen with a session-local stats snapshot.
// Pass an empty stats.Snapshot{} to render an empty grid (e.g. in tests).
func NewHomeModel(snap stats.Snapshot) HomeModel {
	return HomeModel{snapshot: snap}
}

// ─────────────────────────────────────────────────────────────────────────────
// Animation tick — fires every gridTickInterval until the grid is fully drawn.
// ─────────────────────────────────────────────────────────────────────────────

// GridTickMsg is sent on each animation tick. Exported so RootModel can route
// it back to HomeModel from its own Update.
type GridTickMsg struct{}

// gridTickInterval — 40ms per column × 52 columns ≈ 2.1s total.
const gridTickInterval = 40 * time.Millisecond

func gridTickCmd() tea.Cmd {
	return tea.Tick(gridTickInterval, func(time.Time) tea.Msg {
		return GridTickMsg{}
	})
}

// StatsAvailable reports whether the session snapshot has real GitHub data.
func (m HomeModel) StatsAvailable() bool { return m.snapshot.Available }

// StartGridAnimation returns the first tick command. Called by RootModel
// when transitioning from boot → home so the grid begins drawing in.
func (m HomeModel) StartGridAnimation() tea.Cmd {
	return gridTickCmd()
}

// Update handles the grid-draw tick. Returns (model, nextTickCmd) until done.
func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	if _, ok := msg.(GridTickMsg); ok {
		if m.gridRevealed < stats.GridCols {
			m.gridRevealed++
			if m.gridRevealed < stats.GridCols {
				return m, gridTickCmd()
			}
		}
	}
	return m, nil
}

func (m HomeModel) View(width int) string {
	var sb strings.Builder

	// Header — embed visitor number in the breadcrumb (replaces the old stats bar).
	bc := "hisanika ╱ home"
	if m.snapshot.Visitor > 0 {
		bc += fmt.Sprintf("  ·  visitor #%d", m.snapshot.Visitor)
	}
	sb.WriteString(strings.Repeat("\n", ScreenTopPad))
	sb.WriteString(ScreenHeader("home", bc, width))
	sb.WriteString("\n\n")

	pad := strings.Repeat(" ", ScreenLeftPad)

	// Name + tagline.
	sb.WriteString(pad + ScreenTitle.Render("Sanika Surose") + "\n")
	sb.WriteString(pad + MetadataLabel.Render("software engineer · mcmaster university") + "\n")
	sb.WriteString("\n\n")

	// Featured section.
	sb.WriteString(pad + SectionLabel("featured") + "\n\n")
	for _, f := range homeFeatured {
		titleLeft := MetadataLabel.Render(f.num) + "  " + ProjectTitle.Render(f.title)
		titleRow := fmtRightAligned(titleLeft, MetadataLabel.Render(f.meta), width-ScreenLeftPad-1)
		sb.WriteString(pad + titleRow + "\n")
		sb.WriteString(pad + strings.Repeat(" ", 4) + MetadataLabel.Render(f.detail) + "\n")
		sb.WriteString("\n")
	}

	// Navigate label + menu — selected key chip only ([1]…[5]), labels stay plain.
	sb.WriteString(pad + SectionLabel("navigate") + "\n\n")
	sb.WriteString(pad)
	for i, item := range homeItems {
		key := fmt.Sprintf("[%s]", item.key)
		if i == m.cursor {
			sb.WriteString(SelectedRow.Render(" " + key + " "))
		} else {
			sb.WriteString(KeyHint.Render(key))
		}
		sb.WriteString("  ")
		sb.WriteString(BodyText.Render(item.label))
		if i < len(homeItems)-1 {
			sb.WriteString("   ")
		}
	}
	sb.WriteString("\n\n\n")

	// Contribution grid — animates in when GitHub data is available.
	sb.WriteString(pad + SectionLabel("contributions") + "\n\n")
	if !m.snapshot.Available {
		sb.WriteString(pad + GhostText.Render("no contribution data — set GITHUB_TOKEN on the server"))
		sb.WriteString("\n")
	} else {
		sb.WriteString(RenderContributionGrid(m.snapshot.Grid, m.snapshot.GridStart, m.gridRevealed))
		sb.WriteString("\n")
		// Legend only renders once the grid finishes drawing — keeps the eye on
		// the animation while it runs.
		if m.gridRevealed >= stats.GridCols {
			sb.WriteString(RenderContributionLegend(
				m.snapshot.TotalContributions,
				m.snapshot.Streak,
			))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m HomeModel) CursorUp() HomeModel {
	if m.cursor > 0 {
		m.cursor--
	}
	return m
}

func (m HomeModel) CursorDown() HomeModel {
	if m.cursor < len(homeItems)-1 {
		m.cursor++
	}
	return m
}

// Selected returns the menu index the cursor is on
// (0=about, 1=experience, 2=projects, 3=guestbook, 4=contact).
func (m HomeModel) Selected() int { return m.cursor }

// _ keeps content imported so future featured-rail logic can pull from data.go.
var _ = content.Projects
