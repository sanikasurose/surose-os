package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ─────────────────────────────────────────────────────────────────────────────
// Color palette — single source of truth.
// ─────────────────────────────────────────────────────────────────────────────

const (
	colorBackground    = "#1A1A1C"
	colorSurface       = "#222226"
	colorAccentPrimary = "#C8847A"
	colorTextPrimary   = "#E8E4DC"
	colorTextSecondary = "#7A7A82"
	colorTextGhost     = "#3D3D44"
	colorAccentSuccess = "#7A9E82"

	// Timeline interpolation stops — accent → secondary, evenly spaced.
	// Used by TimelineTitle() on the Experience screen.
	colorTimeline0 = "#C8847A" // most recent (= accent)
	colorTimeline1 = "#B47770"
	colorTimeline2 = "#A06A64"
	colorTimeline3 = "#8E5C56"
	colorTimeline4 = "#7A7A82" // oldest (= secondary)
)

// ─────────────────────────────────────────────────────────────────────────────
// Layout constants — every screen indents by these.
// ─────────────────────────────────────────────────────────────────────────────

const (
	ScreenLeftPad  = 3 // cols from terminal edge to content
	ScreenTopPad   = 2 // blank rows before the screen header
	NestedIndent   = 2 // extra cols for bullets / quote bars / metadata
	TagSpacing     = 2 // spaces between bracketed tags
	PromptRuleCols = 60
)

// Bar / glyph constants — the five identity marks used across the app.
const (
	GlyphAccentBar   = "▌" // section header bar, project title bar
	GlyphQuoteBar    = "▎" // blockquote / pull-quote indent
	GlyphCaret       = "▸" // selected list item marker
	GlyphDot         = "●" // timeline node
	GlyphArrow       = "↗" // outbound link
	GlyphSectionStub = "──"
)

// ─────────────────────────────────────────────────────────────────────────────
// Text styles — the atomic colored-text tokens.
// No lipgloss.NewStyle() calls anywhere else in the codebase.
// ─────────────────────────────────────────────────────────────────────────────

var (
	ScreenTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	ProjectTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	BodyText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	MetadataLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	GhostText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextGhost))

	AccentText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary))

	AccentTextBold = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorAccentPrimary))

	SuccessText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentSuccess))

	// The five identity glyphs as pre-styled spans.
	AccentBar   = AccentText
	QuoteBar    = AccentText
	Caret       = AccentText
	TimelineDot = AccentText
)

// ─────────────────────────────────────────────────────────────────────────────
// Row / container styles.
// ─────────────────────────────────────────────────────────────────────────────

var (
	NormalRow = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	// SelectedRow — full-width surface fill + accent foreground.
	// Render with .Width(w) when you know the column width, so the fill spans.
	SelectedRow = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary)).
			Background(lipgloss.Color(colorSurface)).
			Bold(true)

	// Menu items (Home nav, Projects menu).
	MenuItemNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary)).
			PaddingLeft(1).
			PaddingRight(1)

	MenuItemSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorAccentPrimary)).
				Background(lipgloss.Color(colorSurface)).
				PaddingLeft(1).
				PaddingRight(1)

	// Optional bordered containers (no longer used by default — the terminal
	// window IS the frame in the locked direction — kept for any future overlay).
	ContentBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorTextGhost)).
			Padding(0, 1)

	ActiveBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorAccentPrimary)).
			Padding(0, 1)

	// Root screen padding — applied to the entire screen render.
	ScreenPad = lipgloss.NewStyle().
			PaddingLeft(ScreenLeftPad).
			PaddingTop(ScreenTopPad)
)

// ─────────────────────────────────────────────────────────────────────────────
// Boot, prompt, hints.
// ─────────────────────────────────────────────────────────────────────────────

var (
	BootName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	BootSubtitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	PromptLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorAccentPrimary))

	PromptInput = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	PromptOutput = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	PromptCursor = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary))

	KeyHint = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextGhost))

	// TagPill — direction A "bracketed annotation".
	// No background, no padding — brackets are the visible chrome.
	// Use Tag() helper to wrap a tag string in brackets.
	TagPill = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextSecondary))
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers — composed rendering for the six recurring patterns.
// Each takes raw strings + a width, returns a fully-styled string.
// ─────────────────────────────────────────────────────────────────────────────

// ScreenHeader renders the screen title row + breadcrumb on a single line.
// e.g.   "   ▌ projects                              hisanika ╱ projects"
//
// width is the full terminal column width. The breadcrumb right-aligns to
// (width - ScreenLeftPad).
func ScreenHeader(label, breadcrumb string, width int) string {
	left := AccentBar.Render(GlyphAccentBar) + " " + ScreenTitle.Render(label)
	if breadcrumb == "" {
		return strings.Repeat(" ", ScreenLeftPad) + left
	}
	// Visible-length of left content: bar + space + label.
	leftLen := len(GlyphAccentBar) + 1 + len(label)
	gap := width - ScreenLeftPad - leftLen - len(breadcrumb)
	if gap < 2 {
		gap = 2
	}
	return strings.Repeat(" ", ScreenLeftPad) +
		left +
		strings.Repeat(" ", gap) +
		MetadataLabel.Render(breadcrumb)
}

// SectionLabel renders a subsection header — ghost stub + secondary label.
// e.g.   "── now"
func SectionLabel(text string) string {
	return GhostText.Render(GlyphSectionStub) + " " + MetadataLabel.Render(text)
}

// Tag wraps a tag string in brackets in the locked "bracketed annotation" style.
func Tag(text string) string {
	return TagPill.Render("[" + text + "]")
}

// Tags renders a slice of tags joined by TagSpacing spaces.
func Tags(tags []string) string {
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = Tag(t)
	}
	return strings.Join(parts, strings.Repeat(" ", TagSpacing))
}

// Bullet renders a leading "· " followed by the text in body.
func Bullet(text string) string {
	return MetadataLabel.Render("·") + " " + BodyText.Render(text)
}

// Quote renders a single quoted line — accent bar + 2 spaces + body text.
// Use one Quote() per visual line; the bar repeats top-to-bottom.
func Quote(text string) string {
	return QuoteBar.Render(GlyphQuoteBar) + "  " + BodyText.Render(text)
}

// SelectedItemLine — the canonical selected-row render.
// Caret + surface fill spanning to `width` cols, accent text inside.
// `inner` is the full inner content string (already laid out with spaces).
func SelectedItemLine(inner string, width int) string {
	// Leading caret is OUTSIDE the surface fill so the fill starts at col +2.
	caret := Caret.Render(GlyphCaret) + " "
	// Total visible columns available for the surface block:
	surfaceCols := width - ScreenLeftPad - 2 /* caret + space */
	bar := SelectedRow.Width(surfaceCols).Render(" " + inner + " ")
	return strings.Repeat(" ", ScreenLeftPad) + caret + bar
}

// TimelineTitle returns a lipgloss style with color interpolated from accent
// (idx=0, most recent) to secondary (idx=total-1, oldest).
// Hard-coded 5-stop ramp — clamps for total > 5.
func TimelineTitle(idx, total int) lipgloss.Style {
	stops := []string{colorTimeline0, colorTimeline1, colorTimeline2, colorTimeline3, colorTimeline4}
	if total <= 1 {
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(stops[0]))
	}
	pos := float64(idx) / float64(total-1) // 0.0 .. 1.0
	bucket := int(pos*float64(len(stops)-1) + 0.5)
	if bucket < 0 {
		bucket = 0
	}
	if bucket > len(stops)-1 {
		bucket = len(stops) - 1
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(stops[bucket]))
}

// TimelineNode renders the "●─── label ─────────────..." rule for the
// Experience screen group headers.
func TimelineNode(label string, width int) string {
	dot := TimelineDot.Render(GlyphDot)
	stub := GhostText.Render(strings.Repeat("─", 3))
	lab := " " + MetadataLabel.Render(label) + " "
	// Fill the rest of the row with a ghost rule.
	used := ScreenLeftPad + 1 /*dot*/ + 3 /*stub*/ + len(label) + 2 /*spaces*/
	rest := width - used
	if rest < 0 {
		rest = 0
	}
	tail := GhostText.Render(strings.Repeat("─", rest))
	return strings.Repeat(" ", ScreenLeftPad) + dot + stub + lab + tail
}

// PromptBar renders the persistent hisanika> prompt at screen bottom.
// state: "idle" | "focused" | "typing"
// input is shown when focused/typing; output appears on the line below if non-empty.
func PromptBar(state, input, output string) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
	sb.WriteString(GhostText.Render(strings.Repeat("─", PromptRuleCols)))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
	sb.WriteString(PromptLabel.Render("hisanika>"))
	sb.WriteString(" ")
	if state == "idle" {
		sb.WriteString(GhostText.Render("press : to type a command"))
	} else {
		sb.WriteString(PromptInput.Render(input))
		sb.WriteString(PromptCursor.Render("█"))
	}
	if output != "" {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
		sb.WriteString(PromptOutput.Render(output))
	}
	return sb.String()
}

// HelpHint formats a row of key hints, ghost-colored, with " · " separators.
// e.g. HelpHint("j/k scroll", "↵ open", "esc back", ": prompt")
func HelpHint(hints ...string) string {
	return strings.Repeat(" ", ScreenLeftPad) +
		KeyHint.Render(strings.Join(hints, "  ·  "))
}

// Breadcrumb formats "a ╱ b ╱ c" in MetadataLabel color.
func Breadcrumb(parts ...string) string {
	return MetadataLabel.Render(strings.Join(parts, " ╱ "))
}

// ─────────────────────────────────────────────────────────────────────────────
// Contribution grid — GitHub-style heatmap.
// Chart colors, intentionally OUTSIDE the warm rose palette: this is a heatmap,
// instantly recognizable. Level 0 is lifted from #1A1A1C so empty cells read
// as "cells" rather than vanishing into the background.
// ─────────────────────────────────────────────────────────────────────────────

// ContributionLevelColors — 5 stops from empty → bright green.
var ContributionLevelColors = [5]string{
	"#1F2329", // 0 - empty
	"#0E4429", // 1 - dim
	"#006D32", // 2
	"#26A641", // 3
	"#39D353", // 4 - bright
}

// Pre-built per-level cell styles so the grid render doesn't allocate per cell.
var contributionCellStyles = func() [5]lipgloss.Style {
	var s [5]lipgloss.Style
	for i, hex := range ContributionLevelColors {
		s[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
	}
	return s
}()

// Month label layout — removed; labels are derived from gridStart at render time.

// RenderContributionGrid produces the 8-line styled grid block (1 month label
// row + 7 weekday rows). `gridStart` is the UTC Sunday of column 0.
// `revealedCols` controls the draw-in animation — cells at column >= revealedCols
// render as 2 blank spaces.
//
// `grid` shape: [7 weekdays][52 weeks], values 0-4. Row 0 = Sunday.
func RenderContributionGrid(grid [7][52]int, gridStart time.Time, revealedCols int) string {
	if revealedCols < 0 {
		revealedCols = 0
	}
	if revealedCols > 52 {
		revealedCols = 52
	}

	var sb strings.Builder

	// Month label row — indented to sit above the cells (past the 3-char row label).
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+3))
	sb.WriteString(MetadataLabel.Render(monthLabelRow(gridStart)))
	sb.WriteString("\n")

	// Grid rows.
	for r := 0; r < 7; r++ {
		sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
		sb.WriteString(MetadataLabel.Render(weekdayLabel(r)))
		for c := 0; c < 52; c++ {
			if c >= revealedCols {
				sb.WriteString("  ")
				continue
			}
			level := grid[r][c]
			if level < 0 || level > 4 {
				level = 0
			}
			sb.WriteString(contributionCellStyles[level].Render("██"))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderContributionLegend renders the
//   "less ██ ██ ██ ██ ██ more     N contributions in the last year  ·  streak Nd"
// line that sits below the grid.
func RenderContributionLegend(total, streak int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+3))
	sb.WriteString(GhostText.Render("less  "))
	for i := 0; i < 5; i++ {
		sb.WriteString(contributionCellStyles[i].Render("██"))
		if i < 4 {
			sb.WriteString(" ")
		}
	}
	sb.WriteString(GhostText.Render("  more"))
	sb.WriteString(strings.Repeat(" ", 8))
	sb.WriteString(MetadataLabel.Render(fmt.Sprintf(
		"%s contributions in the last year  ·  current streak ",
		commaInt(total),
	)))
	sb.WriteString(BodyText.Render(fmt.Sprintf("%dd", streak)))
	return sb.String()
}

// weekdayLabel returns the 3-char left label for grid row i.
//   Row 0 (Sun), 2 (Tue), 4 (Thu), 6 (Sat) → blank
//   Row 1 (Mon), 3 (Wed), 5 (Fri) → labeled, matching GitHub
func weekdayLabel(i int) string {
	switch i {
	case 1:
		return "M  "
	case 3:
		return "W  "
	case 5:
		return "F  "
	default:
		return "   "
	}
}

// monthLabelRow builds month abbreviations aligned to week columns, matching
// GitHub's rolling-year layout (e.g. jun … may when today is in May).
func monthLabelRow(gridStart time.Time) string {
	const colWidth = 2
	buf := make([]byte, 52*colWidth)
	for i := range buf {
		buf[i] = ' '
	}

	prevMonth := 0
	for col := 0; col < 52; col++ {
		weekStart := gridStart.AddDate(0, 0, col*7)
		month := int(weekStart.Month())
		if month == prevMonth {
			continue
		}
		label := strings.ToLower(weekStart.Format("Jan"))
		pos := col * colWidth
		for i := 0; i < len(label) && pos+i < len(buf); i++ {
			buf[pos+i] = label[i]
		}
		prevMonth = month
	}
	return string(buf)
}

// commaInt formats a non-negative int with thousand separators (1247 → "1,247").
func commaInt(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var out strings.Builder
	pre := len(s) % 3
	if pre > 0 {
		out.WriteString(s[:pre])
		if len(s) > pre {
			out.WriteString(",")
		}
	}
	for i := pre; i < len(s); i += 3 {
		out.WriteString(s[i : i+3])
		if i+3 < len(s) {
			out.WriteString(",")
		}
	}
	return out.String()
}

// Indent returns a left-padded version of `text` indented by `cols`.
// Useful inside multi-line blocks that should sit at +NestedIndent.
func Indent(text string, cols int) string {
	prefix := strings.Repeat(" ", cols)
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}

// fmtRightAligned returns "left<spaces>right" padded to the given column width.
// Helper used by ProjectRow etc.
func fmtRightAligned(left, right string, width int) string {
	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(right)
	gap := width - leftLen - rightLen
	if gap < 2 {
		gap = 2
	}
	return left + strings.Repeat(" ", gap) + right
}

// _ keeps fmt imported even if unused in trimmed builds.
var _ = fmt.Sprintf
