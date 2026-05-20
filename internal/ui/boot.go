package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	bootTickMs  = 70 * time.Millisecond
	bootPauseTk = 5  // ticks of blank before typing starts
	bootNameTk  = 13 // "Sanika Surose" = 13 chars
	bootPause2  = 6  // ticks after name complete before subtitle
	bootSubtTk  = 1  // subtitle appears in one step
	bootPause3  = 8  // ticks of pause after subtitle
	bootDivTk   = 1  // divider appears
	bootPause4  = 5  // ticks before done
)

var bootName = "Sanika Surose"
var bootSubtitle = "software engineer · mcmaster university"

// totalBootSteps is the sum of all phase step counts.
var totalBootSteps = bootPauseTk + bootNameTk + bootPause2 + bootSubtTk + bootPause3 + bootDivTk + bootPause4

type bootTickMsg struct{}

type BootModel struct {
	step int
	done bool
}

func NewBootModel() BootModel { return BootModel{} }

func (m BootModel) Init() tea.Cmd {
	return tea.Tick(bootTickMs, func(time.Time) tea.Msg { return bootTickMsg{} })
}

func (m BootModel) Update(msg tea.Msg) (BootModel, tea.Cmd) {
	if _, ok := msg.(bootTickMsg); !ok {
		return m, nil
	}
	m.step++
	if m.step >= totalBootSteps {
		m.done = true
		return m, nil
	}
	return m, tea.Tick(bootTickMs, func(time.Time) tea.Msg { return bootTickMsg{} })
}

func (m BootModel) Done() bool { return m.done }

func (m BootModel) View() string {
	var sb strings.Builder

	// Phase boundaries (cumulative step counts).
	nameStart := bootPauseTk
	nameEnd := nameStart + bootNameTk
	subtitleStep := nameEnd + bootPause2
	dividerStep := subtitleStep + bootSubtTk + bootPause3

	sb.WriteString("\n\n\n")

	// Name — reveal one character per step during typing phase.
	if m.step >= nameStart {
		shown := m.step - nameStart
		if shown > len(bootName) {
			shown = len(bootName)
		}
		sb.WriteString("  ")
		sb.WriteString(BootName.Render(bootName[:shown]))
	}

	sb.WriteString("\n")

	// Subtitle.
	if m.step >= subtitleStep {
		sb.WriteString("  ")
		sb.WriteString(BootSubtitle.Render(bootSubtitle))
	}

	sb.WriteString("\n")

	// Divider.
	if m.step >= dividerStep {
		sb.WriteString("\n  ")
		sb.WriteString(GhostText.Render(strings.Repeat("─", 38)))
	}

	return sb.String()
}
