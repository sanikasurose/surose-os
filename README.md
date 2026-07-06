# Surose OS

An SSH-accessible portfolio. `ssh` in and you get a navigable terminal application — projects, experience, about, guestbook, contact — instead of a shell.

Most engineering portfolios are a templated website or a GitHub profile with a few pinned repos. Neither says much about how the person actually builds things. This one is a small system: an SSH server, a hand-rolled TUI, and a visual identity designed to look nothing like a system log. The format is part of the pitch.

A browser version with the same visual identity is being built separately and will share this project's domain once both are ready.

## Try it

```
ssh enter@sanikasurose.com
```

Coming soon — this entry point goes live once the browser version above is ready to launch alongside it.

First time connecting, your terminal will ask you to confirm the server's identity (a one-time security check every SSH connection does) — type `yes` and you're in.

## Navigating

Every screen shares the same core keys:

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection / scroll up |
| `↓` / `j` | Move selection / scroll down |
| `Enter` | Select / open |
| `Esc` | Back one screen |
| `q` / `Ctrl+C` | Quit, from anywhere |
| `:` | Open the `hisanika>` command prompt |

**Home** — `1` about · `2` experience · `3` projects · `4` guestbook · `5` contact (or use the arrows and `Enter`).

**Projects** — a menu (`1` hackathons, `2` personal projects) leads to a filtered list, which leads to a project's detail page. On a detail page, `o` copies the project's first link to your clipboard.

**Guestbook** is the one exception: typing goes straight into the message box instead of triggering commands, so `:` doesn't open the prompt there, and `Esc` only backs out to Home once the input is empty.

### The `hisanika>` prompt

Press `:` on any screen except the guestbook to drop into a small command line. Type a command and hit `Enter`:

| Command | Does |
|---------|------|
| `help` | Lists available commands |
| `projects` | Jump to the Projects menu |
| `experience` | Jump to Experience |
| `about` | Jump to About |
| `guestbook` | Jump to Guestbook |
| `contact` | Jump to Contact |
| `open <slug>` | Open a project's detail page directly (e.g. `open surose-os`) |
| `clear` | Clear the prompt's output |
| `quit` / `exit` | Disconnect |

Anything else gets `unknown command. type help for a list.`

## Built with

Go · [Wish](https://github.com/charmbracelet/wish) · [Bubbletea](https://github.com/charmbracelet/bubbletea) · [Lipgloss](https://github.com/charmbracelet/lipgloss) · [Glamour](https://github.com/charmbracelet/glamour)
