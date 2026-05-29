# Surose OS — Navigation Spec

This file contains the complete navigation specification for Surose OS.
The coding agent reads this file before implementing screen routing, keybindings, or the hisanika> command prompt.

**Status:** Phases 1–3 implemented. Phase 4 is production deploy only.

---

## Screen Transition Map

```
[Boot] ──→ [Home]
             │
             ├──[1] / Enter──→ [About]              Esc──→ [Home]
             │
             ├──[2] / Enter──→ [Experience]         Esc──→ [Home]
             │
             ├──[3] / Enter──→ [Projects Menu]
             │                     │
             │              [1|2] / Enter──→ [Projects List]
             │                     │
             │              Enter──→ [Project Detail]
             │                     │
             │                  Esc──→ [Projects List]
             │                     │
             │                  Esc──→ [Projects Menu]
             │
             ├──[4] / Enter──→ [Guestbook]           Esc (empty input)──→ [Home]
             │
             └──[5] / Enter──→ [Contact]            Esc──→ [Home]
```

Esc always goes back exactly one level. It never exits the app.
`q` always quits the app from any screen.

On **Guestbook**, `:` does not open the global `hisanika>` prompt — input goes to the message field.

---

## Home Screen

| Key | Action |
|---|---|
| `1` | Go to About |
| `2` | Go to Experience |
| `3` | Go to Projects menu |
| `4` | Go to Guestbook |
| `5` | Go to Contact |
| `↑` / `k` | Move menu selection up |
| `↓` / `j` | Move menu selection down |
| `Enter` | Open selected menu item |
| `:` | Open `hisanika>` command prompt |
| `q` / `Ctrl+C` | Quit |

Per-screen key hints appear in the footer (`HelpHint`); there is no global `?` overlay.

---

## Projects Menu Screen

| Key | Action |
|---|---|
| `1` | Open hackathons list |
| `2` | Open personal projects list |
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Enter` | Open selected category list |
| `Esc` | Return to Home |
| `q` | Quit |
| `:` | Open `hisanika>` prompt |

---

## Projects List Screen

Filtered single-column list for the chosen category (`hackathon` or `personal`).

| Key | Action |
|---|---|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Enter` | Open selected project detail |
| `Esc` | Return to Projects menu |
| `q` | Quit |
| `:` | Open `hisanika>` prompt |

---

## Project Detail Screen

| Key | Action |
|---|---|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `o` | Copy first link to clipboard (OSC 52); status line confirms URL |
| `Esc` | Return to Projects list |
| `q` | Quit |
| `:` | Open `hisanika>` prompt |

---

## Experience Screen

| Key | Action |
|---|---|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `Esc` | Return to Home |
| `q` | Quit |
| `:` | Open `hisanika>` prompt |

---

## About Screen

| Key | Action |
|---|---|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `Esc` | Return to Home |
| `q` | Quit |
| `:` | Open `hisanika>` prompt |

---

## Contact Screen

| Key | Action |
|---|---|
| `Esc` | Return to Home |
| `q` | Quit |
| `:` | Open `hisanika>` prompt |

---

## Guestbook Screen

| Key | Action |
|---|---|
| Type | Enter message (when input focused) |
| `Enter` | Post message |
| `Esc` | Back to Home if input empty; otherwise clear/focus behavior per implementation |
| `↑` / `k` | Scroll message list up |
| `↓` / `j` | Scroll message list down |
| `pgup` / `pgdown` | Scroll message list |
| `q` | Quit |

No `:` prompt on this screen.

---

## Global Keybindings

| Key | Action |
|---|---|
| `q` | Quit application (all screens) |
| `Ctrl+C` | Quit application |
| `:` | Open `hisanika>` prompt (all screens except Guestbook) |

---

## hisanika> Command Prompt

Available on all screens except Guestbook. Activated by pressing `:` (colon).
The prompt label `hisanika>` renders in AccentPrimary (#C8847A).
User input renders in TextPrimary (#E8E4DC).
Command output renders in TextSecondary (#7A7A82).
Unknown command response: `unknown command. type help for a list.`

### Command List

| Command | Action |
|---|---|
| `help` | Display all available commands |
| `projects` | Navigate to Projects menu |
| `experience` | Navigate to Experience screen |
| `about` | Navigate to About screen |
| `guestbook` | Navigate to Guestbook screen |
| `contact` | Navigate to Contact screen |
| `open <slug>` | Open a project detail screen directly by slug |
| `clear` | Clear prompt output |
| `quit` | Quit application |
| `exit` | Quit application |

### Valid Project Slugs

| Slug | Project |
|---|---|
| `clairity` | CLAIRITY |
| `seamsecure` | SeamSecure |
| `piratech` | Piratech |
| `verihire` | VeriHire |
| `typing-test` | Typing Test |
| `outreach-research-scraper` | Outreach Research Scraper |
| `neurosnake-rl` | neurosnake-rl |
| `sanika-ui` | Sanika-UI |
| `throttle` | Throttle |
| `pulp` | Pulp |
| `committr` | Committr |
| `airport-baggage` | International Airport Baggage |
| `surose-os` | Surose OS |

### Edge Cases

- `open` with no argument → `usage: open <project>`
- `open <invalid-slug>` → `project not found. type projects to browse.`
- Any unrecognized command → `unknown command. type help for a list.`
- Empty input (Enter with no text) → no output, prompt clears
