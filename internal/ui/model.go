package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sanikasurose/surose-os/internal/content"
	"github.com/sanikasurose/surose-os/internal/guestbook"
	"github.com/sanikasurose/surose-os/internal/stats"
)

type screen int

const (
	screenBoot screen = iota
	screenHome
	screenProjectsMenu
	screenProjects
	screenProjectDetail
	screenExperience
	screenAbout
	screenContact
	screenGuestbook
)

// RootModel is the top-level Bubbletea model. It owns the screen router,
// window dimensions, the hisanika> prompt, and the guestbook screen.
type RootModel struct {
	current       screen
	width, height int

	boot          BootModel
	home          HomeModel
	projectsMenu  ProjectsMenuModel
	projects      ProjectsModel
	projectDetail ProjectDetailModel
	experience    ExperienceModel
	about         AboutModel
	contact       ContactModel
	guestbook     GuestbookModel

	promptFocused bool
	promptInput   string
	promptOutput  string
}

// NewRootModel takes the per-session stats snapshot AND the shared guestbook
// store. Both come from the handler (see handler.go MakeHandler).
func NewRootModel(snap stats.Snapshot, gb *guestbook.Store) RootModel {
	handle := visitorHandle(snap.Visitor)
	return RootModel{
		current:      screenBoot,
		boot:         NewBootModel(),
		home:         NewHomeModel(snap),
		projectsMenu: NewProjectsMenuModel(),
		projects:     NewProjectsModel("hackathon"),
		experience:   NewExperienceModel(),
		about:        NewAboutModel(),
		contact:      NewContactModel(),
		guestbook:    NewGuestbookModel(gb, handle),
	}
}

// visitorHandle formats "visitor_NNN", zero-padded to 3 digits.
// 4-digit (and higher) numbers grow naturally.
func visitorHandle(n int) string {
	if n <= 0 {
		return "visitor_???"
	}
	return "visitor_" + zeroPad(n, 3)
}

func zeroPad(n, width int) string {
	s := ""
	for v := n; v > 0; v /= 10 {
		s = string(rune('0'+v%10)) + s
	}
	for len(s) < width {
		s = "0" + s
	}
	return s
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
			// Kick off the contribution-grid draw-in animation as soon as
			// home becomes visible — only when GitHub data loaded successfully.
			cmds := []tea.Cmd{cmd}
			if m.home.StatsAvailable() {
				cmds = append(cmds, m.home.StartGridAnimation())
			}
			return m, tea.Batch(cmds...)
		}
		return m, cmd

	case GridTickMsg:
		// Forward all grid-animation ticks to the home model.
		var hcmd tea.Cmd
		m.home, hcmd = m.home.Update(msg)
		return m, hcmd

	case tea.KeyMsg:
		if m.promptFocused {
			return m.updatePrompt(msg)
		}

		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Reserve `:` for the global prompt EXCEPT on the guestbook screen,
		// where the user is typing into a textinput and colons are legitimate.
		if msg.String() == ":" && m.current != screenGuestbook {
			m.promptFocused = true
			m.promptInput = ""
			m.promptOutput = ""
			return m, nil
		}

		// Guestbook handles its own keys via its embedded textinput.
		if m.current == screenGuestbook {
			var cmd tea.Cmd
			m.guestbook, cmd = m.guestbook.Update(msg)
			if m.guestbook.ShouldExit() {
				m.guestbook = m.guestbook.ResetExit()
				m.current = screenHome
				return m, nil
			}
			return m, cmd
		}

		return m.updateScreen(msg)
	}

	// Forward unhandled messages (textinput cursor blink, paste, etc.) to
	// the active screen if it cares. Today, only the guestbook does.
	if m.current == screenGuestbook {
		var cmd tea.Cmd
		m.guestbook, cmd = m.guestbook.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m RootModel) updateScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.current {

	case screenHome:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h", "up", "k":
			m.home = m.home.CursorUp()
		case "right", "l", "down", "j":
			m.home = m.home.CursorDown()
		case "1":
			m.current = screenAbout
		case "2":
			m.current = screenExperience
		case "3":
			m.current = screenProjectsMenu
		case "4":
			return m.enterGuestbook()
		case "5":
			m.current = screenContact
		case "enter", " ":
			switch m.home.Selected() {
			case 0:
				m.current = screenAbout
			case 1:
				m.current = screenExperience
			case 2:
				m.current = screenProjectsMenu
			case 3:
				return m.enterGuestbook()
			case 4:
				m.current = screenContact
			}
		}

	case screenProjectsMenu:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			m.current = screenHome
		case "up", "k":
			m.projectsMenu = m.projectsMenu.CursorUp()
		case "down", "j":
			m.projectsMenu = m.projectsMenu.CursorDown()
		case "1":
			m.projects = NewProjectsModel("hackathon")
			m.current = screenProjects
		case "2":
			m.projects = NewProjectsModel("personal")
			m.current = screenProjects
		case "enter":
			m.projects = NewProjectsModel(m.projectsMenu.Category())
			m.current = screenProjects
		}

	case screenProjects:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			m.current = screenProjectsMenu
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
		case "q":
			return m, tea.Quit
		case "esc":
			m.current = screenProjects
		case "o":
			m.projectDetail, _ = m.projectDetail.CopyLink()
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

// enterGuestbook refreshes the message list from the store and switches screens.
// Returns the textinput cursor-blink command so the input cursor blinks
// immediately rather than waiting for the next tick.
func (m RootModel) enterGuestbook() (tea.Model, tea.Cmd) {
	m.guestbook = m.guestbook.Refresh()
	m.current = screenGuestbook
	return m, m.guestbook.Init()
}

func (m RootModel) updatePrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.promptFocused = false
		m.promptInput = ""
		m.promptOutput = ""
		return m, nil

	case "enter":
		raw := strings.TrimSpace(m.promptInput)
		m.promptInput = ""
		m, cmd := m.execPromptCmd(raw)
		return m, cmd

	case "backspace", "ctrl+h":
		if len(m.promptInput) > 0 {
			runes := []rune(m.promptInput)
			m.promptInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		s := msg.String()
		if len([]rune(s)) == 1 {
			r := []rune(s)[0]
			if r >= 32 && r != 127 {
				m.promptInput += s
				return m, nil
			}
		}
		if len(msg.Runes) > 0 {
			m.promptInput += string(msg.Runes)
			return m, nil
		}
	}

	return m, nil
}

func (m RootModel) execPromptCmd(input string) (RootModel, tea.Cmd) {
	if input == "" {
		m.promptOutput = ""
		return m, nil
	}

	if strings.HasPrefix(input, "open") {
		parts := strings.Fields(input)
		if len(parts) < 2 {
			m.promptOutput = "usage: open <project>"
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
		m.promptOutput = "commands: help · projects · experience · about · guestbook · contact · open <project> · clear · quit"
	case "projects":
		m.current = screenProjectsMenu
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
	case "guestbook":
		m.guestbook = m.guestbook.Refresh()
		m.current = screenGuestbook
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
	case screenProjectsMenu:
		body = m.projectsMenu.View(m.width)
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
	case screenGuestbook:
		body = m.guestbook.View(m.width, m.height)
	}

	// Guestbook handles its own bottom — don't double-render the hisanika> bar
	// while the user is typing into the textinput.
	if m.current == screenGuestbook {
		return body
	}

	return body + "\n" + m.renderPrompt()
}

// renderPrompt delegates to the styles.go PromptBar helper.
func (m RootModel) renderPrompt() string {
	state := "idle"
	if m.promptFocused {
		state = "typing"
	}
	return PromptBar(state, m.promptInput, m.promptOutput)
}
