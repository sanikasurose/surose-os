package ui

import (
	"strings"
	"unicode"

	"github.com/sanikasurose/surose-os/internal/content"
)

type ExperienceModel struct {
	scrollTop int
}

func NewExperienceModel() ExperienceModel { return ExperienceModel{} }

func (m ExperienceModel) View(width, height int) string {
	var lines []string
	add := func(s string) { lines = append(lines, s) }

	indent := strings.Repeat(" ", ScreenLeftPad+2)

	// Header — kept outside the scrollable `lines` slice below so it stays
	// pinned to the top of the screen instead of scrolling away.
	header := ScreenHeader("experience", "hisanika ╱ experience", width)

	// Partition into work and hackathon groups.
	var workEntries []content.ExperienceEntry
	var hackEntries []content.ExperienceEntry
	for _, e := range content.Experience {
		if e.Kind == "work" {
			workEntries = append(workEntries, e)
		} else {
			hackEntries = append(hackEntries, e)
		}
	}

	// ── 2026 — present (work) ─────────────────────────────────────────
	add(TimelineNode("2026 — present", width))
	add("")

	for _, e := range workEntries {
		// Both work entries use accent color (idx=0 of the 5-stop ramp).
		titleStyle := TimelineTitle(0, 5)
		titleText := lc(e.Title)
		orgText := lc(e.Org)

		// Title left, org right-aligned.
		gap := width - ScreenLeftPad - 2 - len(titleText) - len(orgText)
		if gap < 2 {
			gap = 2
		}
		add(indent +
			titleStyle.Render(titleText) +
			strings.Repeat(" ", gap) +
			MetadataLabel.Render(orgText))

		// Location · dates on one line.
		locationDates := lc(e.Location) + " · " + lc(e.Dates)
		add(indent + MetadataLabel.Render(locationDates))

		// Bullets using Bullet() helper; lowercase first character only.
		for _, b := range e.Bullets {
			add(indent + Bullet(lcFirst(b)))
		}
		add("")
	}

	// ── 2026 — hackathons ─────────────────────────────────────────────
	add(TimelineNode("2026 — hackathons", width))
	add("")

	for i, e := range hackEntries {
		// i=0 → TimelineTitle(1,5) (#B47770), …, i=3 → TimelineTitle(4,5) (#7A7A82).
		titleStyle := TimelineTitle(i+1, 5)
		titleText := lc(e.Title)
		venueText := lc(e.Location) + " · " + lc(e.Dates)

		// Title left, venue right-aligned.
		gap := width - ScreenLeftPad - 2 - len(titleText) - len(venueText)
		if gap < 2 {
			gap = 2
		}
		add(indent +
			titleStyle.Render(titleText) +
			strings.Repeat(" ", gap) +
			MetadataLabel.Render(venueText))

		// One-sentence brief.
		add(indent + BodyText.Render(lcFirst(e.Brief)))

		// Ghost link hint — matches the mockup's "→ full project detail in projects".
		add(indent + GhostText.Render("→ full project detail in projects"))

		if i < len(hackEntries)-1 {
			add("")
		}
	}

	// Bound the scrollable body to what actually fits: height, minus the
	// header above and footer below, minus 2 for the hisanika> bar RootModel
	// appends after this View returns (a divider row plus the prompt row
	// itself — not just one line). Without this, printing more rows than
	// the terminal has causes the terminal itself to scroll — dragging the
	// (fixed) header above off-screen along with it.
	topStr := header + "\n\n"
	footerStr := "\n" + HelpHint("j/k scroll", "esc back")
	topLines := strings.Count(topStr, "\n")
	footerLines := strings.Count(footerStr, "\n") + 1
	available := height - topLines - footerLines - 2
	if available < 3 {
		available = 3
	}

	if m.scrollTop > len(lines)-available {
		m.scrollTop = len(lines) - available
	}
	if m.scrollTop < 0 {
		m.scrollTop = 0
	}
	end := m.scrollTop + available
	if end > len(lines) {
		end = len(lines)
	}
	visible := lines[m.scrollTop:end]

	var sb strings.Builder
	sb.WriteString(topStr)
	for _, l := range visible {
		sb.WriteString(l + "\n")
	}
	sb.WriteString(footerStr)

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

// lc lowercases a string for the experience screen's tone-matching typography.
func lc(s string) string { return strings.ToLower(s) }

// lcFirst lowercases only the first rune of s, preserving acronyms.
func lcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
