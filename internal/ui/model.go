package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sanikasurose/surose-os/internal/content"
)

type screen int

const (
	screenBoot screen = iota
	screenHome
	screenProjects
	screenProjectDetail
	screenAbout
	screenContact
)

// RootModel is the top-level Bubbletea model. It owns the screen router,
// window dimensions, and the hisanika> prompt.
type RootModel struct {
	current       screen
	width, height int

	boot          BootModel
	home          HomeModel
	projects      ProjectsModel
	projectDetail ProjectDetailModel
	about         AboutModel
	contact       ContactModel

	promptFocused bool
	promptInput   string
	promptOutput  string
}

func NewRootModel() RootModel {
	return RootModel{
		current:  screenBoot,
		boot:     NewBootModel(),
		home:     NewHomeModel(),
		projects: NewProjectsModel(),
		about:    NewAboutModel(),
		contact:  NewContactModel(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.boot.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case bootTickMsg:
		newBoot, cmd := m.boot.Update(msg)
		m.boot = newBoot
		if m.boot.Done() {
			m.current = screenHome
		}
		return m, cmd

	case tea.KeyMsg:
		// Prompt is focused — route all keys there.
		if m.promptFocused {
			return m.updatePrompt(msg)
		}

		// Global quit.
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// `:` opens the prompt from any screen.
		if msg.String() == ":" {
			m.promptFocused = true
			m.promptInput = ""
			m.promptOutput = ""
			return m, nil
		}

		// Delegate to active screen.
		return m.updateScreen(msg)
	}

	return m, nil
}

func (m RootModel) updateScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.current {

	case screenHome:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.home = m.home.CursorUp()
		case "down", "j":
			m.home = m.home.CursorDown()
		case "1":
			m.current = screenProjects
		case "2":
			m.current = screenAbout
		case "3":
			m.current = screenContact
		case "enter", " ":
			switch m.home.Selected() {
			case 0:
				m.current = screenProjects
			case 1:
				m.current = screenAbout
			case 2:
				m.current = screenContact
			}
		}

	case screenProjects:
		switch msg.String() {
		case "q", "esc":
			m.current = screenHome
		case "up", "k":
			m.projects = m.projects.CursorUp()
		case "down", "j":
			m.projects = m.projects.CursorDown()
		case "enter":
			p := content.Projects[m.projects.Selected()]
			m.projectDetail = NewProjectDetailModel(p)
			m.current = screenProjectDetail
		}

	case screenProjectDetail:
		switch msg.String() {
		case "q", "esc":
			m.current = screenProjects
		case "down", "j":
			m.projectDetail = m.projectDetail.ScrollDown()
		case "up", "k":
			m.projectDetail = m.projectDetail.ScrollUp()
		}

	case screenAbout:
		switch msg.String() {
		case "q", "esc":
			m.current = screenHome
		case "down", "j":
			m.about = m.about.ScrollDown()
		case "up", "k":
			m.about = m.about.ScrollUp()
		}

	case screenContact:
		switch msg.String() {
		case "q", "esc":
			m.current = screenHome
		}
	}

	return m, nil
}

func (m RootModel) updatePrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.promptFocused = false
		m.promptInput = ""
		m.promptOutput = ""
		return m, nil

	case tea.KeyEnter:
		cmd := strings.TrimSpace(m.promptInput)
		m.promptOutput = m.execPromptCmd(cmd)
		m.promptInput = ""
		if cmd == "quit" || cmd == "exit" {
			return m, tea.Quit
		}
		return m, nil

	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.promptInput) > 0 {
			m.promptInput = m.promptInput[:len(m.promptInput)-1]
		}
		return m, nil

	case tea.KeyRunes:
		m.promptInput += string(msg.Runes)
		return m, nil
	}

	return m, nil
}

// execPromptCmd handles whitelisted commands. Phase 1: only help and quit.
func (m RootModel) execPromptCmd(cmd string) string {
	switch cmd {
	case "help":
		return "commands: help · projects · about · contact · open <slug> · clear · quit"
	case "quit", "exit":
		return "goodbye."
	case "projects":
		return "(navigation coming in phase 2)"
	case "about":
		return "(navigation coming in phase 2)"
	case "contact":
		return "(navigation coming in phase 2)"
	case "clear":
		return ""
	case "":
		return ""
	default:
		return fmt.Sprintf("unknown command. type help for a list.")
	}
}

func (m RootModel) View() string {
	if m.current == screenBoot {
		return m.boot.View()
	}

	var body string
	switch m.current {
	case screenHome:
		body = m.home.View(m.width)
	case screenProjects:
		body = m.projects.View(m.width)
	case screenProjectDetail:
		body = m.projectDetail.View(m.width, m.height)
	case screenAbout:
		body = m.about.View(m.width, m.height)
	case screenContact:
		body = m.contact.View(m.width)
	}

	prompt := m.renderPrompt()

	return body + "\n" + prompt
}

func (m RootModel) renderPrompt() string {
	var sb strings.Builder
	sb.WriteString("  ")
	sb.WriteString(GhostText.Render(strings.Repeat("─", 40)))
	sb.WriteString("\n  ")
	sb.WriteString(PromptLabel.Render("hisanika>"))
	sb.WriteString(" ")

	if m.promptFocused {
		sb.WriteString(PromptInput.Render(m.promptInput))
		sb.WriteString(AccentText.Render("█")) // block cursor
	} else {
		sb.WriteString(GhostText.Render("press : to type a command"))
	}

	if m.promptOutput != "" {
		sb.WriteString("\n  ")
		sb.WriteString(PromptOutput.Render(m.promptOutput))
	}

	return sb.String()
}
