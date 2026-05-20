# Surose OS

An SSH portfolio. Connect and you get a terminal app — projects, about, contact — not a shell.

```
ssh enter@sanikasurose.com
```

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

Requires Go. Wish generates the host key on first run.

To build a binary:

```bash
make build
```

## Keyboard navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Select / open |
| `1` `2` `3` | Home → Projects, About, Contact |
| `Esc` / `q` | Back (sub-screens) or quit (home) |
| `:` | Open `hisanika>` command prompt |
| `Ctrl+C` | Quit |

**Prompt** (`:` then type)

| Command | Action |
|---------|--------|
| `help` | List commands |
| `quit` / `exit` | Disconnect |

## Web

A browser version with the same visual identity is in progress.
