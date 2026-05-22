# Surose OS — Navigation Spec

This file contains the complete navigation specification for Surose OS.
The coding agent reads this file before implementing screen routing, keybindings, or the hisanika> command prompt.

---

## Screen Transition Map

```
[Boot] ──→ [Home]
             │
             ├──[1] / Enter──→ [Projects]
             │                     │
             │              Enter──→ [Project Detail]
             │                     │
             │                  Esc──→ back to [Projects]
             │
             ├──[2] / Enter──→ [Experience]
             │                  Esc──→ back to [Home]
             │
             ├──[3] / Enter──→ [About]
             │                  Esc──→ back to [Home]
             │
             └──[4] / Enter──→ [Contact]
                                Esc──→ back to [Home]
```

Esc always goes back exactly one level. It never exits the app.
q always quits the app from any screen.

---

## Home Screen

| Key | Action |
|---|---|
| `1` | Go to Projects |
| `2` | Go to Experience |
| `3` | Go to About |
| `4` | Go to Contact |
| `↑` / `k` | Move menu selection up |
| `↓` / `j` | Move menu selection down |
| `Enter` | Open selected menu item |
| `?` | Show help overlay |
| `q` / `Ctrl+C` | Quit |

---

## Projects Screen

Projects are displayed in two labeled groups: `hackathons` (top) and `personal projects` (bottom).
Selection moves continuously across both groups — j/k navigates through all items in order.

| Key | Action |
|---|---|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Enter` | Open selected project detail |
| `Esc` | Return to Home |
| `q` | Quit |

---

## Project Detail Screen

| Key | Action |
|---|---|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `o` | Open first link (GitHub or demo) in browser, if terminal supports it |
| `Esc` | Return to Projects |
| `q` | Quit |

---

## Experience Screen

| Key | Action |
|---|---|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `Esc` | Return to Home |
| `q` | Quit |

---

## About Screen

| Key | Action |
|---|---|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `Esc` | Return to Home |
| `q` | Quit |

---

## Contact Screen

| Key | Action |
|---|---|
| `Esc` | Return to Home |
| `q` | Quit |

---

## Global Keybindings (all screens)

| Key | Action |
|---|---|
| `q` | Quit application |
| `Ctrl+C` | Quit application |
| `?` | Show help overlay |

---

## hisanika> Command Prompt

Available on all screens. Activated by pressing `:` (colon).
The prompt label `hisanika>` renders in AccentPrimary (#C8847A).
User input renders in TextPrimary (#E8E4DC).
Command output renders in TextSecondary (#7A7A82).
Unknown command response: `unknown command. type help for a list.`

### Command List

| Command | Action |
|---|---|
| `help` | Display all available commands |
| `projects` | Navigate to Projects screen |
| `experience` | Navigate to Experience screen |
| `about` | Navigate to About screen |
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

- `open` with no argument → `usage: open <slug>`
- `open <invalid-slug>` → `project not found. type projects to browse.`
- Any unrecognized command → `unknown command. type help for a list.`
- Empty input (Enter with no text) → no output, prompt clears