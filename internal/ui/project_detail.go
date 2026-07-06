package ui

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/aymanbagabas/go-osc52/v2"
	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
	"github.com/sanikasurose/surose-os/internal/content"
)

//go:embed glamour/surose.json
var glamourStyleJSON []byte

type ProjectDetailModel struct {
	project   content.Project
	rendered  string // Glamour-rendered Markdown
	lines     []string
	scrollTop int
	statusMsg string // copy-link feedback shown above the help row
}

func NewProjectDetailModel(p content.Project) ProjectDetailModel {
	m := ProjectDetailModel{project: p}
	m.rendered = renderMarkdown(p.LongDesc)
	m.lines = strings.Split(m.rendered, "\n")
	return m
}

// renderMarkdown uses the embedded surose.json Glamour theme so the markdown
// body picks up the same palette as the rest of the app.
func renderMarkdown(src string) string {
	// Glamour defaults to termenv.TrueColor internally regardless of the
	// app's own profile (see main.go's lipgloss.SetColorProfile) — set it
	// explicitly here too, or this renderer alone reverts to truecolor.
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(glamourStyleJSON),
		glamour.WithColorProfile(termenv.ANSI256),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		// Fall back to the standard dark theme if the JSON fails to load —
		// the body will still render, just not on-palette.
		r, err = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithColorProfile(termenv.ANSI256),
			glamour.WithWordWrap(80),
		)
		if err != nil {
			return src
		}
	}
	out, err := r.Render(src)
	if err != nil {
		return src
	}
	return out
}

func (m ProjectDetailModel) View(width, height int) string {
	p := m.project

	var top strings.Builder

	// Top: back-hint row.
	top.WriteString("\n")
	top.WriteString(strings.Repeat(" ", ScreenLeftPad))
	top.WriteString(GhostText.Render("╴╴ back"))
	top.WriteString("\n\n")

	// Title row with accent bar.
	top.WriteString(strings.Repeat(" ", ScreenLeftPad))
	top.WriteString(AccentBar.Render(GlyphAccentBar))
	top.WriteString(" ")
	top.WriteString(ProjectTitle.Render(p.Title))
	top.WriteString("\n")

	// Meta row: year · category · event · location.
	meta := p.Year
	if p.Event != "" {
		meta = fmt.Sprintf("%s · %s", meta, p.Event)
	}
	if p.Location != "" {
		meta = fmt.Sprintf("%s · %s", meta, p.Location)
	}
	top.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
	top.WriteString(MetadataLabel.Render(meta))
	top.WriteString("\n\n")

	// Tag row — bracketed annotation style.
	if len(p.Tags) > 0 {
		top.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
		top.WriteString(Tags(p.Tags))
		top.WriteString("\n\n")
	}

	var bottom strings.Builder

	// Link footer.
	if len(p.Links) > 0 {
		bottom.WriteString("\n")
		bottom.WriteString(strings.Repeat(" ", ScreenLeftPad))
		bottom.WriteString(SectionLabel("links"))
		bottom.WriteString("\n\n")
		for _, link := range p.Links {
			bottom.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
			bottom.WriteString(AccentText.Render(GlyphArrow))
			bottom.WriteString("  ")
			bottom.WriteString(MetadataLabel.Render(link.Label))
			bottom.WriteString("   ")
			bottom.WriteString(BodyText.Render(normalizeURL(link.URL)))
			bottom.WriteString("\n")
		}
	}

	if m.statusMsg != "" {
		bottom.WriteString("\n")
		bottom.WriteString(strings.Repeat(" ", ScreenLeftPad))
		if strings.HasPrefix(m.statusMsg, "copied") {
			bottom.WriteString(SuccessText.Render(m.statusMsg))
		} else {
			bottom.WriteString(MetadataLabel.Render(m.statusMsg))
		}
		bottom.WriteString("\n")
	}

	bottom.WriteString("\n")
	bottom.WriteString(HelpHint("j/k scroll", "o copy link", "esc back"))

	topStr := top.String()
	bottomStr := bottom.String()

	// Bound the scrollable body to what actually fits: height, minus the
	// fixed rows above and below it, minus 1 for the hisanika> bar RootModel
	// appends after this View returns. Without this, printing more rows than
	// the terminal has causes the terminal itself to scroll — dragging the
	// (fixed) header above off-screen along with it.
	topLines := strings.Count(topStr, "\n")
	bottomLines := strings.Count(bottomStr, "\n") + 1
	available := height - topLines - bottomLines - 1
	if available < 3 {
		available = 3
	}

	if m.scrollTop > len(m.lines)-available {
		m.scrollTop = len(m.lines) - available
	}
	if m.scrollTop < 0 {
		m.scrollTop = 0
	}
	end := m.scrollTop + available
	if end > len(m.lines) {
		end = len(m.lines)
	}
	visible := m.lines[m.scrollTop:end]

	var sb strings.Builder
	sb.WriteString(topStr)
	for _, l := range visible {
		sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	sb.WriteString(bottomStr)

	return sb.String()
}

func (m ProjectDetailModel) ScrollDown() ProjectDetailModel {
	m.scrollTop++
	return m
}

func (m ProjectDetailModel) ScrollUp() ProjectDetailModel {
	if m.scrollTop > 0 {
		m.scrollTop--
	}
	return m
}

// normalizeURL ensures bare host/path links render and copy as https URLs.
func normalizeURL(raw string) string {
	if raw == "" {
		return raw
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	return "https://" + raw
}

// CopyLink copies the first project link to the visitor's clipboard via OSC 52
// and sets an in-app status message. Returns the normalized URL when copied.
func (m ProjectDetailModel) CopyLink() (ProjectDetailModel, string) {
	if len(m.project.Links) == 0 {
		m.statusMsg = "no link to copy for this project"
		return m, ""
	}
	url := normalizeURL(m.project.Links[0].URL)
	fmt.Fprint(os.Stderr, osc52.New(url))
	m.statusMsg = "copied — paste in browser: " + url
	return m, url
}
