package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sanikasurose/surose-os/internal/guestbook"
)

// GuestbookModel — the leave-a-note screen.
//
// Interaction model:
//   - Input is focused for the entire time the screen is on. Any printable
//     keypress (including j/k) goes into the textinput; ↑/↓ scroll history.
//   - Enter posts the message (when non-empty), reloads the list, scrolls to top.
//   - Esc with a non-empty input cancels (clears) that draft.
//   - Esc with an empty input exits to home.
//   - `:` is NOT intercepted by the global hisanika> prompt on this screen,
//     so the user can type colons in their message. The prompt is unreachable
//     from guestbook — that's the trade.
type GuestbookModel struct {
	store    *guestbook.Store
	input    textinput.Model
	handle   string              // "visitor_NNN" for this session
	messages []guestbook.Message // up to 50, newest first
	scroll   int                 // top-line offset into the wrapped body
	status   statusKind          // last action status — shown briefly under input
	statusMsg string
	exit     bool                // signal to RootModel: pop back to home
}

type statusKind int

const (
	statusNone statusKind = iota
	statusPosted
	statusCancelled
	statusError
)

// NewGuestbookModel builds an unloaded model. Call Refresh() when you switch
// to this screen so the list reflects fresh DB state.
func NewGuestbookModel(store *guestbook.Store, handle string) GuestbookModel {
	ti := textinput.New()
	ti.CharLimit = guestbook.MaxBodyChars
	ti.Placeholder = "leave a note for the next visitor…"
	ti.Prompt = ""
	ti.Focus()
	// Use the design palette for cursor + placeholder.
	ti.PromptStyle = PromptLabel
	ti.TextStyle = BodyText
	ti.PlaceholderStyle = GhostText

	return GuestbookModel{
		store:  store,
		input:  ti,
		handle: handle,
	}
}

// Refresh reloads the messages list from the store. Call when entering the
// screen (so a visitor sees what was posted while they were on other screens
// in this same session — and the most recent state on first visit).
func (m GuestbookModel) Refresh() GuestbookModel {
	if m.store == nil {
		return m
	}
	msgs, err := m.store.Recent(50)
	if err == nil {
		m.messages = msgs
	}
	m.scroll = 0
	return m
}

// ShouldExit / ResetExit form the signal the RootModel watches to pop back
// to the home screen when Esc is hit on an empty input.
func (m GuestbookModel) ShouldExit() bool       { return m.exit }
func (m GuestbookModel) ResetExit() GuestbookModel { m.exit = false; return m }

// Init returns the textinput cursor-blink command. RootModel forwards this
// when it switches to the guestbook screen so the cursor blinks immediately.
func (m GuestbookModel) Init() tea.Cmd { return textinput.Blink }

func (m GuestbookModel) Update(msg tea.Msg) (GuestbookModel, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "enter":
			return m.submit()
		case "esc":
			if strings.TrimSpace(m.input.Value()) != "" {
				m.input.SetValue("")
				m.status = statusCancelled
				m.statusMsg = "cancelled"
				return m, nil
			}
			m.exit = true
			return m, nil
		case "up":
			if m.scroll > 0 {
				m.scroll--
			}
			return m, nil
		case "down":
			m.scroll++
			return m, nil
		case "pgup":
			m.scroll -= 10
			if m.scroll < 0 {
				m.scroll = 0
			}
			return m, nil
		case "pgdown":
			m.scroll += 10
			return m, nil
		}
	}

	// Everything else (printable keys, cursor blink ticks, paste msgs etc)
	// goes to the textinput.
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	// Clear "posted"/"cancelled" status the moment the user starts typing.
	if m.status != statusNone && m.input.Value() != "" {
		m.status = statusNone
		m.statusMsg = ""
	}
	return m, cmd
}

func (m GuestbookModel) submit() (GuestbookModel, tea.Cmd) {
	body := strings.TrimSpace(m.input.Value())
	if body == "" {
		// Quiet no-op — nothing to post.
		return m, nil
	}
	if m.store == nil {
		m.status = statusError
		m.statusMsg = "guestbook unavailable"
		return m, nil
	}
	err := m.store.Insert(m.handle, body)
	switch {
	case err == nil:
		m.input.SetValue("")
		m = m.Refresh()
		m.status = statusPosted
		m.statusMsg = "✓ posted as " + m.handle
	case errors.Is(err, guestbook.ErrTooLong):
		m.status = statusError
		m.statusMsg = fmt.Sprintf("too long — %d/%d chars", utf8.RuneCountInString(body), guestbook.MaxBodyChars)
	case errors.Is(err, guestbook.ErrEmpty):
		// shouldn't happen because of TrimSpace check, but handle defensively
		return m, nil
	default:
		m.status = statusError
		m.statusMsg = "couldn't post — try again in a moment"
	}
	return m, nil
}

// View — full render. Lays out:
//   1. Screen header (▌ guestbook                    hisanika ╱ guestbook · N)
//   2. "you are visitor_NNN" handle hint
//   3. Wrapped, scrollable message list
//   4. "── leave a note" subsection + input + char counter + status
//   5. HelpHint footer
//
// Width / height come from RootModel.View().
func (m GuestbookModel) View(width, height int) string {
	var sb strings.Builder

	count := len(m.messages)
	bc := fmt.Sprintf("hisanika ╱ guestbook  ·  %d messages", count)
	sb.WriteString(ScreenHeader("guestbook", bc, width))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+2))
	sb.WriteString(MetadataLabel.Render("you are "))
	sb.WriteString(AccentText.Render(m.handle))
	sb.WriteString("\n\n")

	// Build the full unscrolled body — list of styled lines.
	bodyW := width - ScreenLeftPad - 4
	if bodyW < 40 {
		bodyW = 40
	}
	lines := m.renderList(bodyW)

	// Reserve space at the bottom for the input area (subsection + spacer +
	// input + spacer + status + hint) — about 6 lines.
	reserved := 8
	available := height - 6 - reserved // 6 = header + handle + padding
	if available < 6 {
		available = 6
	}

	if m.scroll > len(lines)-available {
		m.scroll = len(lines) - available
	}
	if m.scroll < 0 {
		m.scroll = 0
	}
	end := m.scroll + available
	if end > len(lines) {
		end = len(lines)
	}
	for _, l := range lines[m.scroll:end] {
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	// If list is shorter than the available area, pad with blank lines so
	// the input always sits at the same vertical position.
	if end-m.scroll < available {
		for i := 0; i < available-(end-m.scroll); i++ {
			sb.WriteString("\n")
		}
	}

	// ── input area ──────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad))
	sb.WriteString(SectionLabel("leave a note"))
	sb.WriteString("\n\n")

	// Caret + textinput.
	n := utf8.RuneCountInString(m.input.Value())
	counter := fmt.Sprintf("%d/%d", n, guestbook.MaxBodyChars)

	// Resize textinput to fill the remaining width, minus caret + counter + spacers.
	inputW := width - ScreenLeftPad - 4 - len(counter) - 4
	if inputW < 20 {
		inputW = 20
	}
	m.input.Width = inputW

	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+1))
	sb.WriteString(AccentText.Render(GlyphCaret))
	sb.WriteString("  ")
	sb.WriteString(m.input.View())
	sb.WriteString("    ")
	if n >= guestbook.MaxBodyChars-10 {
		sb.WriteString(AccentText.Render(counter))
	} else {
		sb.WriteString(GhostText.Render(counter))
	}
	sb.WriteString("\n\n")

	// Status line (one row reserved even when empty so layout doesn't jump).
	sb.WriteString(strings.Repeat(" ", ScreenLeftPad+3))
	switch m.status {
	case statusPosted:
		sb.WriteString(SuccessText.Render(m.statusMsg))
	case statusError:
		sb.WriteString(AccentText.Render(m.statusMsg))
	case statusCancelled:
		sb.WriteString(GhostText.Render(m.statusMsg))
	default:
		sb.WriteString(" ")
	}
	sb.WriteString("\n")

	sb.WriteString(HelpHint("↵ post", "esc cancel / back", "↑/↓ scroll", "pgup/pgdown"))

	return sb.String()
}

// renderList builds the styled, soft-wrapped list of all messages.
func (m GuestbookModel) renderList(bodyW int) []string {
	pad := strings.Repeat(" ", ScreenLeftPad+2)
	var out []string
	if len(m.messages) == 0 {
		out = append(out, pad+GhostText.Render("no messages yet — be the first."))
		return out
	}
	for _, mm := range m.messages {
		ts := formatGuestbookTime(mm.CreatedAt)
		head := AccentText.Render(mm.Handle) +
			"   " + GhostText.Render("·") + "   " +
			MetadataLabel.Render(ts)
		out = append(out, pad+head)
		wrapped := softwrap(mm.Body, bodyW-2)
		for _, w := range wrapped {
			out = append(out, pad+BodyText.Render(w))
		}
		out = append(out, "") // blank between entries
	}
	return out
}

// guestbookTZ is the display timezone for message timestamps. Stored values
// remain UTC in SQLite; only rendering converts. America/Toronto handles
// EST/EDT automatically so deploy location doesn't matter.
var guestbookTZ = func() *time.Location {
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		return time.UTC
	}
	return loc
}()

func formatGuestbookTime(t time.Time) string {
	local := t.In(guestbookTZ)
	return local.Format("2006-01-02 15:04 ") + local.Format("MST")
}

// softwrap splits s into lines of at most w runes, preferring to break at
// the rightmost space within a 15-rune look-back window. Pure rune-aware —
// safe with emoji and non-ASCII.
func softwrap(s string, w int) []string {
	if w <= 0 {
		return []string{s}
	}
	runes := []rune(s)
	var out []string
	for len(runes) > w {
		breakAt := w
		for i := w; i > w-15 && i > 0; i-- {
			if runes[i] == ' ' {
				breakAt = i
				break
			}
		}
		out = append(out, string(runes[:breakAt]))
		runes = runes[breakAt:]
		if len(runes) > 0 && runes[0] == ' ' {
			runes = runes[1:]
		}
	}
	if len(runes) > 0 {
		out = append(out, string(runes))
	}
	return out
}
