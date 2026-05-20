package ui

import "github.com/charmbracelet/lipgloss"

// Color palette — single source of truth.
const (
	colorBackground    = "#1A1A1C"
	colorSurface       = "#222226"
	colorAccentPrimary = "#C8847A"
	colorTextPrimary   = "#E8E4DC"
	colorTextSecondary = "#7A7A82"
	colorTextGhost     = "#3D3D44"
	colorAccentSuccess = "#7A9E82"
)

// All Lipgloss styles live here. No lipgloss.NewStyle() calls anywhere else.
var (
	ScreenTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	MetadataLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	BodyText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	GhostText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextGhost))

	AccentText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary))

	SuccessText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentSuccess))

	// List item styles.
	NormalRow = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary)).
			PaddingLeft(1)

	SelectedRow = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary)).
			Background(lipgloss.Color(colorSurface)).
			PaddingLeft(1)

	// Containers.
	ContentBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorTextGhost)).
			Padding(0, 1)

	ActiveBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorAccentPrimary)).
			Padding(0, 1)

	// Screen-level padding applied to root container.
	ScreenPad = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingTop(1)

	// Boot sequence.
	BootName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	BootSubtitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	// Prompt bar.
	PromptLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary))

	PromptInput = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	PromptOutput = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	// Tags / pills.
	TagPill = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary)).
			Background(lipgloss.Color(colorSurface)).
			PaddingLeft(1).
			PaddingRight(1)

	ProjectTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	// Key hint shown in menus and nav bars.
	KeyHint = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextGhost))

	// Home menu item styles.
	MenuItemNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary)).
			PaddingLeft(2)

	MenuItemSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccentPrimary)).
			Background(lipgloss.Color(colorSurface)).
			PaddingLeft(2)
)
