package ui

import (
	"strings"

	"github.com/sanikasurose/surose-os/internal/content"
)

type ContactModel struct{}

func NewContactModel() ContactModel { return ContactModel{} }

func (m ContactModel) View(width int) string {
	var sb strings.Builder
	c := content.Contact

	sb.WriteString(ScreenPad.Render(ScreenTitle.Render("contact")))
	sb.WriteString("\n\n")

	sb.WriteString("  " + MetadataLabel.Render("email     ") + BodyText.Render(c.Email) + "\n")
	sb.WriteString("  " + MetadataLabel.Render("linkedin  ") + BodyText.Render(c.LinkedIn) + "\n")
	sb.WriteString("  " + MetadataLabel.Render("github    ") + BodyText.Render(c.GitHub) + "\n")
	sb.WriteString("\n")
	sb.WriteString("  " + MetadataLabel.Render(c.Note) + "\n")
	sb.WriteString("\n")
	sb.WriteString(KeyHint.Render("  esc  back"))

	return sb.String()
}
