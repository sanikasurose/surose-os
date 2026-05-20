package ui

import (
	"fmt"
	"strings"

	"github.com/sanikasurose/surose-os/internal/content"
)

type ProjectDetailModel struct {
	project    content.Project
	scrollTop  int
	lineHeight int
}

func NewProjectDetailModel(p content.Project) ProjectDetailModel {
	return ProjectDetailModel{project: p}
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
		sb.WriteString(ScreenPad.Render(MetadataLabel.Render(p.Event)))
	}

	sb.WriteString("\n\n")

	// Phase 1: render long desc as plain text. Glamour in Phase 2.
	lines := strings.Split(p.LongDesc, "\n")
	visible := lines
	if m.scrollTop < len(lines) {
		visible = lines[m.scrollTop:]
	}
	for _, l := range visible {
		sb.WriteString("  ")
		sb.WriteString(BodyText.Render(l))
		sb.WriteString("\n")
	}

	if len(p.Links) > 0 {
		sb.WriteString("\n")
		for _, link := range p.Links {
			sb.WriteString("  ")
			sb.WriteString(AccentText.Render(link.Label + ": " + link.URL))
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
