package ui

import (
	"strings"

	"github.com/sanikasurose/surose-os/internal/content"
)

type ExperienceModel struct {
	scrollTop int
}

func NewExperienceModel() ExperienceModel { return ExperienceModel{} }

func (m ExperienceModel) View(width, height int) string {
	var sb strings.Builder

	sb.WriteString(ScreenPad.Render(ScreenTitle.Render("experience")))
	sb.WriteString("\n\n")

	// Build all lines first, then apply scroll.
	var lines []string
	for i, e := range content.Experience {
		if i > 0 {
			lines = append(lines, "")
		}
		if e.Kind == "work" {
			lines = append(lines, "  "+ProjectTitle.Render(e.Title))
			lines = append(lines, "  "+MetadataLabel.Render(e.Org+" · "+e.Dates))
			if e.Location != "" {
				lines = append(lines, "  "+MetadataLabel.Render(e.Location))
			}
			for _, b := range e.Bullets {
				lines = append(lines, "  "+GhostText.Render("·")+" "+BodyText.Render(b))
			}
		} else {
			lines = append(lines, "  "+ProjectTitle.Render(e.Title))
			lines = append(lines, "  "+MetadataLabel.Render(e.Location+" · "+e.Dates))
			lines = append(lines, "  "+BodyText.Render(e.Brief))
			lines = append(lines, "  "+GhostText.Render("→ full project detail in projects"))
		}
	}

	visible := lines
	if m.scrollTop < len(lines) {
		visible = lines[m.scrollTop:]
	}
	for _, l := range visible {
		sb.WriteString(l + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(KeyHint.Render("  j/k  scroll    esc  back"))

	return sb.String()
}

func (m ExperienceModel) ScrollDown() ExperienceModel {
	m.scrollTop++
	return m
}

func (m ExperienceModel) ScrollUp() ExperienceModel {
	if m.scrollTop > 0 {
		m.scrollTop--
	}
	return m
}
