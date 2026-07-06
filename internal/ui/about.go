package ui

import (
	"strings"
)

type AboutModel struct {
	scrollTop int
}

func NewAboutModel() AboutModel { return AboutModel{} }

func (m AboutModel) View(width, height int) string {
	var lines []string
	add := func(s string) { lines = append(lines, s) }

	pad := strings.Repeat(" ", ScreenLeftPad)
	indent := strings.Repeat(" ", ScreenLeftPad+2)

	// Header — kept outside the scrollable `lines` slice below so it stays
	// pinned to the top of the screen instead of scrolling away.
	header := ScreenHeader("about", "hisanika ╱ about", width)

	// Name + subtitle.
	add(pad + ScreenTitle.Render("Sanika Surose"))
	add(pad + MetadataLabel.Render("software engineer · mcmaster university, class of 2029"))
	add("")

	// ── now ───────────────────────────────────────────────────────────
	add(pad + SectionLabel("now"))
	add("")

	nowLines := []string{
		"i'm a Software Engineer Intern at TranQuility Inc., an AI startup, where I",
		"build backend APIs, RAG data pipelines, and full-stack features using",
		"FastAPI, React, and Azure. My work has included shipping a production RAG",
		"pipeline (Python + Qdrant) spanning multiple document formats, cutting AI",
		"feature costs significantly through model and token optimization,",
		"resolving a production outage caused by async I/O bottlenecks, and",
		"migrating our deployment workflow from a multi-step manual process to a",
		"streamlined Azure Container Registry pipeline.",
		"",
		"i'm also an Undergraduate Research Assistant at McMaster University under",
		"Dr. Charles Welch, working on NLP and LLM evaluation. I've built a crash-",
		"safe, resumable inference pipeline serving Qwen3-8B via vLLM on a 4×H100",
		"GPU cluster, and designed zero-shot and few-shot emotion classification",
		"experiments on the GoEmotions dataset, benchmarking model performance",
		"across prompting strategies.",
	}
	for _, l := range nowLines {
		if l == "" {
			add("")
			continue
		}
		add(indent + Quote(l))
	}
	add("")

	// ── how ───────────────────────────────────────────────────────────
	add(pad + SectionLabel("how"))
	add("")

	add(indent + Quote("third-year software engineering co-op student. i work across"))
	add(indent + Quote("the stack — backends, distributed systems, LLM tooling, and"))
	add(indent + Quote("occasionally things like this."))
	add("")

	// ── elsewhere ─────────────────────────────────────────────────────
	add(pad + SectionLabel("elsewhere"))
	add("")

	elsewhereLines := []string{
		"sponsorship exec for DeltaHacks, McMaster's largest hackathon with 600+",
		"attendees. embedded subteam for McMaster Aerial Robotics and Drone Club,",
		"competing in the AEAC Student UAS Competition. building on the software",
		"subteam for McMaster Design League, McMaster's largest CAD-a-thon. outreach",
		"for McMaster Advanced Space Systems.",
	}
	for _, l := range elsewhereLines {
		add(indent + Quote(l))
	}
	add("")

	// ── education ─────────────────────────────────────────────────────
	add(pad + SectionLabel("education"))
	add("")

	school := "McMaster University"
	degree := "B.Eng. software engineering (co-op) · 2029"
	quotePrefixWidth := 3 // GlyphQuoteBar + two spaces
	gap := width - ScreenLeftPad - 2 - quotePrefixWidth - len(school) - len(degree)
	if gap < 2 {
		gap = 2
	}
	add(indent + QuoteBar.Render(GlyphQuoteBar) + "  " +
		BodyText.Render(school) +
		strings.Repeat(" ", gap) +
		MetadataLabel.Render(degree))
	add("")
	add(indent + MetadataLabel.Render("built this in Go, which i'd never used before. that was the point."))
	add("")

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
