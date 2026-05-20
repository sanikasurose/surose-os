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

	for i, p := range content.Projects {
		tags := formatTags(p.Tags)
		year := MetadataLabel.Render(p.Year)

		titleLine := fmt.Sprintf("%s  %s", ProjectTitle.Render(p.Title), year)
		descLine := MetadataLabel.Render(p.ShortDesc)
		tagsLine := tags

		block := titleLine + "\n" + descLine + "\n" + tagsLine

		if i == m.cursor {
			sb.WriteString(ActiveBox.Render(block))
		} else {
			sb.WriteString(ContentBox.Render(block))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(KeyHint.Render("  j/k  navigate    enter  open    esc  back"))

	return sb.String()
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
