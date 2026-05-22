package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/sanikasurose/surose-os/internal/content"
)

type ProjectDetailModel struct {
	project    content.Project
	rendered   string // Glamour-rendered Markdown
	lines      []string
	scrollTop  int
}

func NewProjectDetailModel(p content.Project) ProjectDetailModel {
	m := ProjectDetailModel{project: p}
	m.rendered = renderMarkdown(p.LongDesc)
	m.lines = strings.Split(m.rendered, "\n")
	return m
}

func renderMarkdown(src string) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return src
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
	sb.WriteString(ScreenPad.Render(ScreenTitle.Render(p.Title)))
	sb.WriteString("\n")

	meta := fmt.Sprintf("%s  %s", p.Year, formatTags(p.Tags))
	sb.WriteString(ScreenPad.Render(MetadataLabel.Render(meta)))

	if p.Event != "" {
		sb.WriteString("\n")
		loc := p.Event
		if p.Location != "" {
			loc += " · " + p.Location
		}
		sb.WriteString(ScreenPad.Render(MetadataLabel.Render(loc)))
	}

	sb.WriteString("\n\n")

	visible := m.lines
	if m.scrollTop < len(m.lines) {
		visible = m.lines[m.scrollTop:]
	}
	for _, l := range visible {
		sb.WriteString(l + "\n")
	}

	if len(p.Links) > 0 {
		sb.WriteString("\n")
		for _, link := range p.Links {
			sb.WriteString("  ")
			sb.WriteString(AccentText.Render(link.Label+": "+link.URL))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(KeyHint.Render("  j/k  scroll    esc  back"))

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
