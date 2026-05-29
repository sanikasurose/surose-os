# Surose OS

An SSH portfolio. Connect and you get a terminal app — projects, experience, about, guestbook, contact — not a shell.

```
ssh enter@sanikasurose.com   # production (after Phase 4 deploy)
```

**Status:** Phases 1–3 complete. **Remaining:** production deploy (no `-p` flag). A browser version with the same visual identity is in progress as a separate project before launch.

## Built with

Go · [Wish](https://github.com/charmbracelet/wish) · [Bubbletea](https://github.com/charmbracelet/bubbletea) · [Lipgloss](https://github.com/charmbracelet/lipgloss) · [Glamour](https://github.com/charmbracelet/glamour)

## Run locally

```bash
git clone https://github.com/sanikasurose/surose-os.git
cd surose-os
make run
```

In another terminal:

```bash
ssh enter@localhost -p 2222
```

Requires Go. Wish generates the host key on first run (`.ssh/id_ed25519`).

Optional: create `.env` with `GITHUB_TOKEN` for the Home contribution grid (see `.env.example` if present).

To build a Linux binary:

```bash
make build
```

## Keyboard navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Select / open |
| `1`–`5` | Home → About, Experience, Projects, Guestbook, Contact |
| `Esc` | Back one screen |
| `q` | Quit |
| `:` | Open `hisanika>` command prompt |
| `Ctrl+C` | Quit |

**Home:** `1` about · `2` experience · `3` projects · `4` guestbook · `5` contact

**Projects:** menu (`1` hackathons · `2` personal) → list → detail. On detail, `o` copies the first link.

**Prompt** (`:` then type, Enter to run)

| Command | Action |
|---------|--------|
| `help` | List commands |
| `projects` | Projects menu |
| `experience` | Experience screen |
| `about` | About screen |
| `guestbook` | Guestbook screen |
| `contact` | Contact screen |
| `open <slug>` | Jump to project detail |
| `clear` | Clear prompt output |
| `quit` / `exit` | Disconnect |

## Docs

- [PLANNING.md](docs/PLANNING.md) — PRD/TRD
- [NAVIGATION.md](docs/NAVIGATION.md) — keys and commands
- [CONTENT.md](docs/CONTENT.md) — portfolio copy for `data.go`

## Web

A browser portfolio with the same visual identity is in progress (separate repo). Plan: finish web, then deploy both SSH and web on `sanikasurose.com`.
