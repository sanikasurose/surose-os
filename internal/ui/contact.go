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

	pad := strings.Repeat(" ", ScreenLeftPad)

	sb.WriteString(ScreenHeader("contact", "hisanika ╱ contact", width))
	sb.WriteString("\n")

	sb.WriteString(pad + MetadataLabel.Render("the best way to reach me is email."))
	sb.WriteString("\n\n")

	// Contact rows — label + value, with fixed-width label column.
	type row struct{ label, value string }
	rows := []row{
		{"email", c.Email},
		{"linkedin", c.LinkedIn},
		{"github", c.GitHub},
	}
	// Label column width = longest label + 2 spaces padding.
	labelW := 0
	for _, r := range rows {
		if len(r.label) > labelW {
			labelW = len(r.label)
		}
	}
	for _, r := range rows {
		label := MetadataLabel.Render(r.label)
		spacing := strings.Repeat(" ", labelW-len(r.label)+3)
		sb.WriteString(pad + label + spacing + BodyText.Render(r.value) + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(HelpHint("esc back"))

	return sb.String()
}
