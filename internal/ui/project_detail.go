package ui

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/aymanbagabas/go-osc52/v2"
	"github.com/charmbracelet/glamour"
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
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(glamourStyleJSON),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		// Fall back to the standard dark theme if the JSON fails to load —
		// the body will still render, just not on-palette.
		r, err = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
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
	var sb strings.Builder
	p := m.project

	// Top: back-hint row.
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
	sb.WriteString(GhostText.Render("╴╴ back"))
	sb.WriteString("\n\n")

	// Title row with accent bar.
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
	sb.WriteString(AccentBar.Render(GlyphAccentBar))
	sb.WriteString(" ")
	sb.WriteString(ProjectTitle.Render(p.Title))
	sb.WriteString("\n")

	// Meta row: year · category · event · location.
	meta := p.Year
	if p.Event != "" {
		meta = fmt.Sprintf("%s · %s", meta, p.Event)
	}
	if p.Location != "" {
		meta = fmt.Sprintf("%s · %s", meta, p.Location)
	}
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
	sb.WriteString(MetadataLabel.Render(meta))
	sb.WriteString("\n\n")

	// Tag row — bracketed annotation style.
	if len(p.Tags) > 0 {
		sb.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
		sb.WriteString(Tags(p.Tags))
		sb.WriteString("\n\n")
	}

	// Glamour-rendered markdown body (scrollable).
	visible := m.lines
	if m.scrollTop < len(m.lines) {
		visible = m.lines[m.scrollTop:]
	}
	for _, l := range visible {
		sb.WriteString(l)
		sb.WriteString("\n")
	}

	// Link footer.
	if len(p.Links) > 0 {
		sb.WriteString("\n")
		sb.WriteString(SectionLabel("links"))
		sb.WriteString("\n\n")
		for _, link := range p.Links {
			sb.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
			sb.WriteString(AccentText.Render(GlyphArrow))
			sb.WriteString("  ")
			sb.WriteString(MetadataLabel.Render(link.Label))
			sb.WriteString("   ")
			sb.WriteString(BodyText.Render(normalizeURL(link.URL)))
			sb.WriteString("\n")
		}
	}

	if m.statusMsg != "" {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
		if strings.HasPrefix(m.statusMsg, "copied") {
			sb.WriteString(SuccessText.Render(m.statusMsg))
		} else {
			sb.WriteString(MetadataLabel.Render(m.statusMsg))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(HelpHint("j/k scroll", "o copy link", "esc back"))

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
