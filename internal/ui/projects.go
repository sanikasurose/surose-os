package ui

import (
	"fmt"
	"strings"

	"github.com/sanikasurose/surose-os/internal/content"
)

type ProjectsModel struct {
	cursor int
}

func NewProjectsModel() ProjectsModel { return ProjectsModel{} }

func (m ProjectsModel) View(width int) string {
	var sb strings.Builder

	sb.WriteString(ScreenPad.Render(ScreenTitle.Render("projects")))
	sb.WriteString("\n\n")

	// Hackathons group.
	sb.WriteString("  " + MetadataLabel.Render("hackathons") + "\n")
	sb.WriteString("  " + GhostText.Render(strings.Repeat("─", 38)) + "\n")
	for i, p := range content.Projects {
		if p.Category == "hackathon" {
			sb.WriteString(m.renderRow(p, i))
		}
	}

	sb.WriteString("\n")

	// Personal projects group.
	sb.WriteString("  " + MetadataLabel.Render("personal projects") + "\n")
	sb.WriteString("  " + GhostText.Render(strings.Repeat("─", 38)) + "\n")
	for i, p := range content.Projects {
		if p.Category == "personal" {
			sb.WriteString(m.renderRow(p, i))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(KeyHint.Render("  j/k  navigate    enter  open    esc  back"))

	return sb.String()
}

func (m ProjectsModel) renderRow(p content.Project, idx int) string {
	tags := formatTags(p.Tags)
	year := MetadataLabel.Render(p.Year)
	titleLine := fmt.Sprintf("  %s  %s", ProjectTitle.Render(p.Title), year)
	descLine := "  " + MetadataLabel.Render(p.ShortDesc)
	tagsLine := "  " + tags
	block := titleLine + "\n" + descLine + "\n" + tagsLine + "\n"

	if idx == m.cursor {
		return AccentText.Render("▸") + block[1:]
	}
	return block
}

func (m ProjectsModel) CursorUp() ProjectsModel {
	if m.cursor > 0 {
		m.cursor--
	}
	return m
}

func (m ProjectsModel) CursorDown() ProjectsModel {
	if m.cursor < len(content.Projects)-1 {
		m.cursor++
	}
	return m
}

func (m ProjectsModel) Selected() int { return m.cursor }

func formatTags(tags []string) string {
	var parts []string
	for _, t := range tags {
		parts = append(parts, TagPill.Render(t))
	}
	return strings.Join(parts, " ")
}
