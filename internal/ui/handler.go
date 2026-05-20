package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
)

// Handler is the Wish bubbletea middleware handler.
// It is called once per SSH connection and returns a fresh Bubbletea model
// so that each session is completely independent.
func Handler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return NewRootModel(), []tea.ProgramOption{tea.WithAltScreen()}
}
