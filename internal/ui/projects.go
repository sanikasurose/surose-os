package ui

import (
	"fmt"
	"strings"

	"github.com/sanikasurose/surose-os/internal/content"
)

// ── Projects menu (hackathons / personal projects chooser) ─────────────────

type ProjectsMenuModel struct {
	cursor int
}

var projectsMenuItems = []struct {
	key      string
	label    string
	category string
}{
	{"1", "hackathons", "hackathon"},
	{"2", "personal projects", "personal"},
}

func NewProjectsMenuModel() ProjectsMenuModel { return ProjectsMenuModel{} }

func (m ProjectsMenuModel) View(width int) string {
	var sb strings.Builder
	pad := strings.Repeat(" ", ScreenLeftPad)

	sb.WriteString(strings.Repeat("\n", ScreenTopPad))
	sb.WriteString(ScreenHeader("projects", "hisanika ╱ projects", width))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
	sb.WriteString(MetadataLabel.Render("two groups — pick one"))
	sb.WriteString("\n\n")

	for i, item := range projectsMenuItems {
		count := countCategory(item.category)
		countStr := fmt.Sprintf("%d entries", count)

		// Surface fill = width - ScreenLeftPad - 2.
		// SelectedRow.Render(" " + inner + " ") adds 2 chars, so len(inner) = surfaceCols - 2.
		surfaceCols := width - ScreenLeftPad - 2
		innerLen := surfaceCols - 2

		keyLabel := "[" + item.key + "] " + item.label
		gap := innerLen - len(keyLabel) - len(countStr)
		if gap < 1 {
			gap = 1
		}
		inner := keyLabel + strings.Repeat(" ", gap) + countStr

		if i == m.cursor {
			sb.WriteString(SelectedItemLine(inner, width))
		} else {
			// Align with caret+space variant: ScreenLeftPad + 2 leading chars.
			rowPad := strings.Repeat(" ", ScreenLeftPad+2)
			rowLeft := MetadataLabel.Render("["+item.key+"]") + " " + BodyText.Render(item.label)
			rowRight := MetadataLabel.Render(countStr)
			normalGap := innerLen - len("["+item.key+"]") - 1 - len(item.label) - len(countStr)
			if normalGap < 1 {
				normalGap = 1
			}
			sb.WriteString(rowPad + rowLeft + strings.Repeat(" ", normalGap) + rowRight)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(pad + GhostText.Render(strings.Repeat("─", width-ScreenLeftPad*2)))
	sb.WriteString("\n")
	sb.WriteString(HelpHint("1/2 select", "↵ open", "esc back"))

	return sb.String()
}

func countCategory(cat string) int {
	n := 0
	for _, p := range content.Projects {
		if p.Category == cat {
			n++
		}
	}
	return n
}

func (m ProjectsMenuModel) CursorUp() ProjectsMenuModel {
	if m.cursor > 0 {
		m.cursor--
	}
	return m
}

func (m ProjectsMenuModel) CursorDown() ProjectsMenuModel {
	if m.cursor < len(projectsMenuItems)-1 {
		m.cursor++
	}
	return m
}

// Category returns the content category string for the selected menu item.
func (m ProjectsMenuModel) Category() string {
	return projectsMenuItems[m.cursor].category
}

// ── Projects list (one category at a time) ────────────────────────────────

type ProjectsModel struct {
	category string
	indices  []int // positions in content.Projects matching this category
	cursor   int
}

// NewProjectsModel builds a filtered list for the given category ("hackathon" or "personal").
func NewProjectsModel(category string) ProjectsModel {
	var idx []int
	for i, p := range content.Projects {
		if p.Category == category {
			idx = append(idx, i)
		}
	}
	return ProjectsModel{category: category, indices: idx}
}

func (m ProjectsModel) View(width int) string {
	var sb strings.Builder

	catLabel := m.category
	if catLabel == "personal" {
		catLabel = "personal projects"
	}
	breadcrumb := "hisanika ╱ projects ╱ " + catLabel
	sb.WriteString(strings.Repeat("\n", ScreenTopPad))
	sb.WriteString(ScreenHeader("projects ╴ "+catLabel, breadcrumb, width))
	sb.WriteString("\n")

	pad := strings.Repeat(" ", ScreenLeftPad)
	indent := strings.Repeat(" ", ScreenLeftPad+3)

	// Count line.
	sb.WriteString(pad + MetadataLabel.Render(fmt.Sprintf("showing %d", len(m.indices))))
	sb.WriteString("\n\n")

	for i, idx := range m.indices {
		p := content.Projects[idx]
		isSel := i == m.cursor

		// Title + year row.
		if isSel {
			// inner fills surfaceCols - 2 (the Render adds leading+trailing space).
			surfaceCols := width - ScreenLeftPad - 2
			innerLen := surfaceCols - 2
			gap := innerLen - len(p.Title) - len(p.Year)
			if gap < 1 {
				gap = 1
			}
			inner := p.Title + strings.Repeat(" ", gap) + p.Year
			sb.WriteString(SelectedItemLine(inner, width))
		} else {
			rowPad := strings.Repeat(" ", ScreenLeftPad+2)
			gap := width - ScreenLeftPad - 2 - len(p.Title) - len(p.Year)
			if gap < 1 {
				gap = 1
			}
			sb.WriteString(rowPad +
				ProjectTitle.Render(p.Title) +
				strings.Repeat(" ", gap) +
				MetadataLabel.Render(p.Year))
		}
		sb.WriteString("\n")

		// Description — split into at most two lines at a word boundary.
		line1, line2 := splitDesc(p.ShortDesc, width-len(indent)-2)
		sb.WriteString(indent + MetadataLabel.Render(line1) + "\n")
		if line2 != "" {
			sb.WriteString(indent + MetadataLabel.Render(line2) + "\n")
		}

		// Tags — bracketed annotation style via Tags() helper.
		sb.WriteString(indent + Tags(p.Tags) + "\n")

		if i < len(m.indices)-1 {
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(HelpHint("j/k scroll", "↵ open", "o copy link on detail", "esc back"))

	return sb.String()
}

// splitDesc splits s into at most two lines at a word boundary near maxW.
func splitDesc(s string, maxW int) (string, string) {
	if len(s) <= maxW {
		return s, ""
	}
	cut := maxW
	for i := maxW; i > maxW-20 && i > 0; i-- {
		if s[i] == ' ' {
			cut = i
			break
		}
	}
	return strings.TrimRight(s[:cut], " "), strings.TrimLeft(s[cut:], " ")
}

func (m ProjectsModel) CursorUp() ProjectsModel {
	if m.cursor > 0 {
		m.cursor--
	}
	return m
}

func (m ProjectsModel) CursorDown() ProjectsModel {
	if m.cursor < len(m.indices)-1 {
		m.cursor++
	}
	return m
}

// Selected returns the index into content.Projects for the highlighted item.
func (m ProjectsModel) Selected() int {
	if len(m.indices) == 0 {
		return 0
	}
	return m.indices[m.cursor]
}
