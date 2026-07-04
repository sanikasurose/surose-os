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

	// Header.
	add(ScreenHeader("about", "hisanika ╱ about", width))
	add("")

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

	add(indent + BodyText.Render("second-year software engineering co-op student. i work across"))
	add(indent + BodyText.Render("the stack — backends, distributed systems, LLM tooling, and"))
	add(indent + BodyText.Render("occasionally things like this."))
	add("")
	add(indent + MetadataLabel.Render("built this in Go, which i'd never used before. that was the point."))
	add("")

	// ── education ─────────────────────────────────────────────────────
	add(pad + SectionLabel("education"))
	add("")

	school := "McMaster University"
	degree := "B.Eng. software engineering (co-op) · 2029"
	gap := width - ScreenLeftPad - 2 - len(school) - len(degree)
	if gap < 2 {
		gap = 2
	}
	add(indent +
		BodyText.Render(school) +
		strings.Repeat(" ", gap) +
		MetadataLabel.Render(degree))
	add("")

	// Apply scroll.
	visible := lines
	if m.scrollTop < len(lines) {
		visible = lines[m.scrollTop:]
	}

	var sb strings.Builder
	for _, l := range visible {
		sb.WriteString(l + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString(HelpHint("j/k scroll", "esc back"))

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
