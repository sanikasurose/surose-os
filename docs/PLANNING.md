# Surose OS — Planning Document
### PRD + TRD Reference · Last updated: May 2025

---

## Part 1 — Product Requirements (PRD)

### Problem Statement

Most software engineering portfolios are either a generic website built from a template, or a GitHub profile with a few pinned repos. Neither communicates anything interesting about how the person actually thinks. For an engineer who is early in their career, the portfolio has to do extra work — it has to show personality, technical range, and the willingness to take on a challenge, all at once.

The secondary problem: personal websites are forgettable. Hundreds of portfolios exist that look identical. Standing out requires doing something unexpected.

### Solution

An SSH-accessible interactive portfolio. When someone runs `ssh enter@sanikasurose.com`, they get a fully navigable terminal application instead of a shell. The experience is immediate, distinctive, and impossible to replicate with a template. The fact that it exists is itself a signal — it communicates that the builder thinks in systems, not surfaces.

The visual identity is the differentiator within that format. Most SSH portfolios look like system logs. This one should look like it was designed by someone who cares about every pixel.

### Target Audience

**Primary:**
- Engineers at companies Sanika wants to work at — people who will appreciate the technical choice and recognize the stack
- Startup founders and technical hiring managers — people who value initiative and unconventional thinking over conventional credentials

**Secondary:**
- Other students and early-career engineers who find it via social sharing or GitHub
- Recruiters at larger tech companies — the "wow, that's different" factor is enough to make the portfolio memorable even if they don't fully understand the stack

**Not the audience (for this version):**
- Non-technical users, family, general public — a web version will be built separately for this group

### Core Value Proposition

The portfolio communicates three things simultaneously, without stating them explicitly:

1. **Technical range** — building in Go, using SSH as a UI transport, deploying on a VPS with systemd. This is not a beginner's project.
2. **Design sensibility** — the visual identity proves that Sanika thinks about UX and aesthetics, not just functionality.
3. **Initiative** — the project didn't exist as a template. It was designed from scratch and built in a language she learned in order to build it.

### Success Metrics

This is not a product with DAUs or conversion funnels. Success looks like:

- The entry command gets shared on Twitter/X or LinkedIn at least once organically
- At least one recruiter or founder mentions it during an interview conversation
- The GitHub repo gets starred by engineers outside Sanika's network
- The project is referenced in at least one job offer discussion

### Feature List

**MVP (Phase 1–3, must ship):**
- SSH server accessible at `ssh enter@sanikasurose.com`
- Boot sequence: cinematic name reveal
- Home, Projects List, Project Detail, About, Contact screens
- Keyboard navigation: arrows, j/k, Enter, Esc, q
- `hisanika>` command prompt with full command set
- Full visual design system: custom palette, rounded borders, color hierarchy
- Real content: all four projects, real bio, real contact info

**Phase 4 (should have):**
- Smooth screen transitions
- Window resize handling
- Production deployment on DigitalOcean
- Live at `ssh enter@sanikasurose.com`

**Nice-to-have (post-MVP):**
- Boot sequence skip on keypress
- Easter eggs (hidden commands, konami code equivalent)
- Visitor counter (how many people have SSH'd in)
- ASCII art portrait on the About screen (using `chafa` or equivalent)
- Animated cursor in the prompt

**Out of scope for this project:**
- Authentication of any kind
- User accounts or persistent sessions
- Any web interface (separate project)
- Mobile or Windows-native support (Windows Terminal works; older setups get a disclaimer)
- Analytics beyond a simple visitor counter

### User Stories

- As an engineer who received this link, I want to SSH in and immediately understand what I'm looking at, so I can decide whether to explore further.
- As a recruiter, I want to find contact information quickly, so I can reach out without spending time navigating.
- As a technical hiring manager, I want to see the actual projects with context, so I can evaluate the depth of work.
- As a fellow student who found this on GitHub, I want to understand how it was built, so I can learn from it or build my own.

### Timeline

| Phase | Goal | Target |
|---|---|---|
| Phase 1 | Working SSH server + navigation skeleton | Week 1 |
| Phase 2 | Real content + full command system | Week 2 |
| Phase 3 | Full visual identity | Week 3 |
| Phase 4 | Polish + production deployment | Week 4 |

---

## Part 2 — Technical Requirements (TRD)

### Tech Stack Rationale

| Technology | Why |
|---|---|
| **Go** | Compiled to a single binary, trivial to deploy, excellent concurrency model for handling multiple SSH sessions. Also: learning Go through a real project is the goal. |
| **Wish** | Charmbracelet's SSH server middleware for Bubbletea. Handles the SSH layer, key management, and session lifecycle. No need to build a raw SSH server. |
| **Bubbletea** | Elm-architecture TUI framework. Clean separation of model/update/view makes the codebase navigable and testable. The right tool for a multi-screen stateful TUI. |
| **Lipgloss** | Declarative terminal styling. Lets us build a proper design system (named styles, consistent tokens) rather than inline ANSI codes everywhere. |
| **Glamour** | Markdown rendering in the terminal. Used exclusively in Project Detail views. Lets project content be written in Markdown without custom rendering logic. |
| **DigitalOcean** | Simple Linux VPS. No container orchestration needed for a single Go binary. $6/month. Easier port forwarding setup than Fly.io for raw TCP/SSH traffic. |

### Architecture Overview

```
Internet
    │
    ▼ :22 (TCP)
[ UFW NAT ]
    │
    ▼ :2222 (TCP)
[ Wish SSH Server ]
    │  ─ handles connection lifecycle
    │  ─ accepts all public keys + passwords
    │  ─ one Bubbletea program per session
    ▼
[ Bubbletea Root Model ]
    │  ─ holds current screen state
    │  ─ routes keypresses to active screen
    │  ─ manages screen transitions
    ▼
[ Screen Models ] ── [ Content Layer (data.go) ]
  boot.go                 Project structs
  home.go                 About string (Markdown)
  projects.go             Contact struct
  project_detail.go
  about.go
  contact.go
    │
    ▼
[ styles.go ]  ← single source of truth for all Lipgloss styles
```

Each SSH connection gets its own independent Bubbletea program instance. Sessions do not share state. This means multiple people can SSH in simultaneously without interference.

### Data Model

No database. All content is defined as Go structs in `internal/content/data.go`.

```go
type Project struct {
    Title       string
    Slug        string
    Year        string
    Tags        []string
    ShortDesc   string
    LongDesc    string   // Markdown, rendered by Glamour
    Links       []Link
    Event       string   // optional: hackathon name / location
}

type Link struct {
    Label string
    URL   string
}

type ContactInfo struct {
    Email    string
    LinkedIn string
    GitHub   string
    Note     string
}
```

To add or edit a project: edit `data.go`, recompile, redeploy. No migrations, no database.

### Screen State Machine

```
[Boot] ──→ [Home] ──→ [Projects List] ──→ [Project Detail]
                │                                │
                ├──→ [About]              Esc ───┘
                │
                └──→ [Contact]
```

The root model holds a `currentScreen` enum and delegates all `Update()` and `View()` calls to the active screen model. Screen transitions are triggered by returning a custom `Msg` type from a screen's `Update()`.

### Third-Party Integrations

None. No external APIs, no payment systems, no auth providers. The only external dependency is the Go module system at build time.

### Performance Targets

- Connection to first rendered frame: < 500ms
- Boot sequence duration: 3–4 seconds (intentional, not a load time)
- Keypress response time: < 16ms (one frame at 60fps)
- Memory per session: < 10MB
- Concurrent sessions supported: 50+ (limited by Droplet RAM, not the app)

### Security Considerations

- No authentication by design — this is a public read-only portfolio
- No user input is executed as code or written to disk
- The `hisanika>` prompt is a command router, not a shell — only whitelisted commands are accepted
- SSH host key is generated once, stored on the server, not in the repo
- The `enter` system user on the Droplet has no shell (`/usr/sbin/nologin`), no sudo, and no home directory write access — it only exists to receive SSH connections

### Infrastructure & Deployment

**Server:** DigitalOcean Droplet, Ubuntu 22.04 LTS, Basic $6/month plan (1 vCPU, 1GB RAM, 25GB SSD)

**Build process:**
```bash
GOOS=linux GOARCH=amd64 go build -o surose-os ./cmd/surose-os
scp surose-os user@<droplet-ip>:/home/surose/surose-os/
```

**Port strategy:**
- Port 22: admin SSH access (for Sanika to manage the server)
- Port 2222: the Wish SSH server (the portfolio app)
- UFW NAT rule redirects incoming :22 connections from the `enter` user to :2222
- This gives visitors the clean `ssh enter@sanikasurose.com` entry with no port flag

**Process management:** systemd. The binary runs as a dedicated `surose` system user. On crash, systemd restarts it after 5 seconds.

**Makefile targets:**
```makefile
make build     # compile for Linux AMD64
make run       # run locally for development
make deploy    # build + scp to Droplet + restart service
make logs      # tail journalctl output from Droplet
make ssh       # SSH into the Droplet as admin
```

**Domain:** A record on `sanikasurose.com` pointing to the Droplet's public IP. No CDN, no load balancer, no HTTPS (SSH doesn't use TLS).

### Windows Compatibility

Windows Terminal (Windows 10+) supports SSH natively and renders the app correctly. The app should detect limited terminal environments at connection time by checking:

- `$TERM` environment variable
- `$WT_SESSION` (set by Windows Terminal)
- Terminal dimensions (too small = prompt to resize)

If a limited environment is detected, display a one-time disclaimer at the top of the boot sequence:

```
note: some visual elements may render differently in older terminals.
      windows terminal, iterm2, and wezterm are fully supported.
```

This message uses `TextSecondary` color and disappears after the boot sequence completes.
