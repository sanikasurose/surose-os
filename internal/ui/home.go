package ui

import (
	"fmt"
	"strings"
)

type HomeModel struct {
	cursor int
}

var homeItems = []struct {
	key   string
	label string
}{
	{"1", "projects"},
	{"2", "experience"},
	{"3", "about"},
	{"4", "contact"},
}

func NewHomeModel() HomeModel { return HomeModel{} }

func (m HomeModel) View(width int) string {
	var sb strings.Builder

	sb.WriteString(ScreenPad.Render(
		ScreenTitle.Render("Sanika Surose") + "\n" +
			MetadataLabel.Render("software engineer · mcmaster university"),
	))

	sb.WriteString("\n\n")

	for i, item := range homeItems {
		hint := KeyHint.Render(fmt.Sprintf("  [%s]", item.key))
		label := item.label
		var line string
		if i == m.cursor {
			line = MenuItemSelected.Render(fmt.Sprintf("%s  %s", hint, label))
		} else {
			line = MenuItemNormal.Render(fmt.Sprintf("%s  %s", hint, label))
		}
		sb.WriteString(line + "\n")
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

// Selected returns the Screen the cursor is on (0=Projects, 1=About, 2=Contact).
func (m HomeModel) Selected() int { return m.cursor }
