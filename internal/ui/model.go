package ui

import (
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
	screenExperience
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
	experience    ExperienceModel
	about         AboutModel
	contact       ContactModel

	promptFocused bool
	promptInput   string
	promptOutput  string
}

func NewRootModel() RootModel {
	return RootModel{
		current:    screenBoot,
		boot:       NewBootModel(),
		home:       NewHomeModel(),
		projects:   NewProjectsModel(),
		experience: NewExperienceModel(),
		about:      NewAboutModel(),
		contact:    NewContactModel(),
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
		if m.promptFocused {
			return m.updatePrompt(msg)
		}

		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == ":" {
			m.promptFocused = true
			m.promptInput = ""
			m.promptOutput = ""
			return m, nil
		}

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
			m.current = screenExperience
		case "3":
			m.current = screenAbout
		case "4":
			m.current = screenContact
		case "enter", " ":
			switch m.home.Selected() {
			case 0:
				m.current = screenProjects
			case 1:
				m.current = screenExperience
			case 2:
				m.current = screenAbout
			case 3:
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

	case screenExperience:
		switch msg.String() {
		case "q", "esc":
			m.current = screenHome
		case "down", "j":
			m.experience = m.experience.ScrollDown()
		case "up", "k":
			m.experience = m.experience.ScrollUp()
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
		raw := strings.TrimSpace(m.promptInput)
		m.promptInput = ""
		m, cmd := m.execPromptCmd(raw)
		return m, cmd

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

// execPromptCmd handles all hisanika> commands. Returns an updated model and
// an optional tea.Cmd (nil unless quitting).
func (m RootModel) execPromptCmd(input string) (RootModel, tea.Cmd) {
	if input == "" {
		m.promptOutput = ""
		return m, nil
	}

	// Handle "open <slug>" separately.
	if strings.HasPrefix(input, "open") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			m.promptOutput = "usage: open <slug>"
			return m, nil
		}
		slug := parts[1]
		for _, p := range content.Projects {
			if p.Slug == slug {
				m.projectDetail = NewProjectDetailModel(p)
				m.current = screenProjectDetail
				m.promptFocused = false
				m.promptOutput = ""
				return m, nil
			}
		}
		m.promptOutput = "project not found. type projects to browse."
		return m, nil
	}

	switch input {
	case "help":
		m.promptOutput = "commands: help · projects · experience · about · contact · open <slug> · clear · quit"
	case "projects":
		m.current = screenProjects
		m.promptFocused = false
		m.promptOutput = ""
	case "experience":
		m.current = screenExperience
		m.promptFocused = false
		m.promptOutput = ""
	case "about":
		m.current = screenAbout
		m.promptFocused = false
		m.promptOutput = ""
	case "contact":
		m.current = screenContact
		m.promptFocused = false
		m.promptOutput = ""
	case "clear":
		m.promptOutput = ""
	case "quit", "exit":
		m.promptOutput = "goodbye."
		return m, tea.Quit
	default:
		m.promptOutput = "unknown command. type help for a list."
	}

	return m, nil
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
	case screenExperience:
		body = m.experience.View(m.width, m.height)
	case screenAbout:
		body = m.about.View(m.width, m.height)
	case screenContact:
		body = m.contact.View(m.width)
	}

	return body + "\n" + m.renderPrompt()
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
		sb.WriteString(AccentText.Render("█"))
	} else {
		sb.WriteString(GhostText.Render("press : to type a command"))
	}

	if m.promptOutput != "" {
		sb.WriteString("\n  ")
		sb.WriteString(PromptOutput.Render(m.promptOutput))
	}

	return sb.String()
}
