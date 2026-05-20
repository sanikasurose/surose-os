package ui

import (
	"strings"

	"github.com/sanikasurose/surose-os/internal/content"
)

type AboutModel struct {
	scrollTop int
}

func NewAboutModel() AboutModel { return AboutModel{} }

func (m AboutModel) View(width, height int) string {
	var sb strings.Builder

	sb.WriteString(ScreenPad.Render(ScreenTitle.Render("about")))
	sb.WriteString("\n\n")

	lines := strings.Split(content.AboutText, "\n")
	visible := lines
	if m.scrollTop < len(lines) {
		visible = lines[m.scrollTop:]
	}
	for _, l := range visible {
		sb.WriteString("  ")
		sb.WriteString(BodyText.Render(l))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(KeyHint.Render("  j/k  scroll    esc  back"))

	return sb.String()
}

func (m AboutModel) ScrollDown() AboutModel {
	m.scrollTop++
	return m
}

func (m AboutModel) ScrollUp() AboutModel {
	if m.scrollTop > 0 {
		m.scrollTop--
	}
	return m
}
