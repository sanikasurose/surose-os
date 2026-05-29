# Surose OS — Planning Document
### PRD + TRD Reference · Last updated: May 2026

---

## Part 1 — Product Requirements (PRD)

### Problem Statement

Most software engineering portfolios are either a generic website built from a template, or a GitHub profile with a few pinned repos. Neither communicates anything interesting about how the person actually thinks. For an engineer who is early in their career, the portfolio has to do extra work — it has to show personality, technical range, and the willingness to take on a challenge, all at once.

The secondary problem: personal websites are forgettable. Hundreds of portfolios exist that look identical. Standing out requires doing something unexpected.

### Solution

An SSH-accessible interactive portfolio. When someone runs `ssh enter@sanikasurose.com`, they get a fully navigable terminal application instead of a shell. The experience is immediate, distinctive, and impossible to replicate with a template. The fact that it exists is itself a signal — it communicates that the builder thinks in systems, not surfaces.

The visual identity is the differentiator within that format. Most SSH portfolios look like system logs. This one should look like it was designed by someone who cares about every pixel.

A **web version** with the same visual identity is built separately for non-technical visitors and will share the `sanikasurose.com` domain (HTTPS on 443, SSH on 22).

### Target Audience

**Primary:**
- Engineers at companies Sanika wants to work at — people who will appreciate the technical choice and recognize the stack
- Startup founders and technical hiring managers — people who value initiative and unconventional thinking over conventional credentials

**Secondary:**
- Other students and early-career engineers who find it via social sharing or GitHub
- Recruiters at larger tech companies — the "wow, that's different" factor is enough to make the portfolio memorable even if they don't fully understand the stack

**Not the audience (for the SSH app):**
- Non-technical users, family, general public — use the web version (separate project)

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

**Shipped (Phases 1–3 — complete):**
- SSH server (local dev on `:2222`; production target `ssh enter@sanikasurose.com`)
- Boot sequence: cinematic name reveal
- Screens: Home, Projects (menu + list), Project Detail, Experience, About, Contact, Guestbook
- Keyboard navigation: arrows, j/k, Enter, Esc, q
- `hisanika>` command prompt with full command set (including `guestbook`, `open <slug>`)
- Full visual design system: custom palette, selection styling, Glamour theme
- Real content: all projects, experience, bio, contact (see **docs/CONTENT.md**)
- Visitor counter and guestbook (SQLite)
- GitHub contribution grid on Home (optional `GITHUB_TOKEN`)

**Phase 4 (remaining — deploy only):**
- Production deployment on DigitalOcean
- Live at `ssh enter@sanikasurose.com` with no `-p` port flag
- DNS: `sanikasurose.com` A record → Droplet IP

**Deferred / not planned for this repo:**
- Smooth screen transitions
- Window resize polish beyond storing dimensions
- Boot sequence skip on keypress
- `?` global help overlay (per-screen `HelpHint` footers instead)
- Windows / limited-terminal disclaimer at boot
- Easter eggs (hidden commands, konami code equivalent)
- ASCII art portrait on About (`chafa` or equivalent)
- Animated cursor in the global prompt

**Nice-to-have (implemented):**
- Visitor counter ✅

**Out of scope for this Go project:**
- Authentication of any kind
- User accounts or persistent sessions beyond guestbook messages
- Web interface (separate project; same domain when live)
- Mobile or Windows-native apps (Windows Terminal works for SSH)

### User Stories

- As an engineer who received this link, I want to SSH in and immediately understand what I'm looking at, so I can decide whether to explore further.
- As a recruiter, I want to find contact information quickly, so I can reach out without spending time navigating.
- As a technical hiring manager, I want to see the actual projects with context, so I can evaluate the depth of work.
- As a fellow student who found this on GitHub, I want to understand how it was built, so I can learn from it or build my own.

### Timeline

| Phase | Goal | Status |
|---|---|---|
| Phase 1 | Working SSH server + navigation skeleton | ✅ Complete |
| Phase 2 | Real content + full command system + Experience | ✅ Complete |
| Phase 3 | Full visual identity | ✅ Complete |
| Phase 4 | Production deployment | 🔄 Remaining |

**Order of launch:** Web portfolio (separate repo) in progress first; SSH production deploy after both are ready. Domain `sanikasurose.com` required for production entry command and public website.

---

## Part 2 — Technical Requirements (TRD)

### Tech Stack Rationale

| Technology | Why |
|---|---|
| **Go** | Compiled to a single binary, trivial to deploy, excellent concurrency model for handling multiple SSH sessions. Also: learning Go through a real project is the goal. |
| **Wish** | Charmbracelet's SSH server middleware for Bubbletea. Handles the SSH layer, key management, and session lifecycle. No need to build a raw SSH server. |
| **Bubbletea** | Elm-architecture TUI framework. Clean separation of model/update/view makes the codebase navigable and testable. The right tool for a multi-screen stateful TUI. |
| **Lipgloss** | Declarative terminal styling. Lets us build a proper design system (named styles, consistent tokens) rather than inline ANSI codes everywhere. |
| **Glamour** | Markdown rendering in the terminal. Custom JSON theme in Project Detail. |
| **SQLite** | Guestbook messages only; content stays in `data.go`. |
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
    │  ─ one Bubbletea program per session
    ▼
[ Bubbletea Root Model ]
    │  ─ screen router + hisanika> prompt
    ▼
[ Screen Models ] ── [ Content Layer (data.go) ]
  boot.go                 Projects, Experience, About, Contact
  home.go                 (+ stats cache for grid)
  projects.go
  project_detail.go
  experience.go
  about.go
  contact.go
  guestbook.go
    │
    ├── [ styles.go ]     ← Lipgloss design system
    ├── [ guestbook/ ]    ← SQLite (shared across sessions)
    └── [ stats/ ]        ← visitors.count + GitHub API cache
```

Each SSH connection gets its own independent Bubbletea program instance. Guestbook and visitor counter are shared via filesystem stores under `SUROSE_DATA_DIR`.

### Data Model

Portfolio content: Go structs in `internal/content/data.go` (no migrations).

```go
type Project struct {
    Title       string
    Slug        string
    Year        string
    Category    string   // "hackathon" or "personal"
    Tags        []string
    ShortDesc   string
    LongDesc    string   // Markdown, rendered by Glamour
    Links       []Link
    Event       string
    Location    string
}

type ExperienceEntry struct {
    Kind     string   // "work" or "hackathon"
    Title    string
    Org      string
    // ...
}
```

Guestbook: SQLite at `{SUROSE_DATA_DIR}/guestbook.db`. Visitor counter: `{SUROSE_DATA_DIR}/visitors.count`.

To add or edit a project: edit `data.go`, recompile, redeploy.

### Screen State Machine

```
[Boot] ──→ [Home] ──→ [Projects Menu] ──→ [Projects List] ──→ [Project Detail]
             │              │                    │                  Esc
             ├── [Experience]                   Esc
             ├── [About]
             ├── [Guestbook]
             └── [Contact]
```

Esc goes back one level. `q` quits from any screen.

### Third-Party Integrations

- **Build time:** Go module dependencies only
- **Runtime (optional):** GitHub REST API for contribution calendar, streak, and stars on Home — requires `GITHUB_TOKEN` on the server

### Performance Targets

- Connection to first rendered frame: < 500ms
- Boot sequence duration: ~3 seconds (intentional, not a load time)
- Keypress response time: < 16ms (one frame at 60fps)
- Memory per session: < 10MB
- Concurrent sessions supported: 50+ (limited by Droplet RAM, not the app)

### Security Considerations

- No authentication by design — this is a public read-only portfolio (guestbook accepts short text only, validated and stored in SQLite)
- No user input is executed as code
- The `hisanika>` prompt is a command router, not a shell — only whitelisted commands are accepted
- SSH host key is generated once, stored on the server, not in the repo
- The `enter` system user on the Droplet has no shell (`/usr/sbin/nologin`), no sudo — only receives SSH connections into the app

### Infrastructure & Deployment

**Server:** DigitalOcean Droplet, Ubuntu 22.04 LTS, Basic $6/month plan (1 vCPU, 1GB RAM, 25GB SSD)

**Domain:** Register or configure `sanikasurose.com`. Same apex can serve SSH (port 22) and HTTPS web (port 443) on one VPS, or split web to a static host with DNS CNAME/`www` as needed.

**Build process:**
```bash
GOOS=linux GOARCH=amd64 go build -o surose-os ./cmd/surose-os
scp surose-os user@<droplet-ip>:/home/surose/surose-os/
```

**Port strategy:**
- Port 22: admin SSH access (for Sanika to manage the server)
- Port 2222: the Wish SSH server (the portfolio app)
- UFW NAT rule redirects incoming :22 connections for the `enter` user to :2222
- This gives visitors the clean `ssh enter@sanikasurose.com` entry with no port flag

**Process management:** systemd. The binary runs as a dedicated `surose` system user. On crash, systemd restarts it after 5 seconds.

**Makefile targets (planned for Phase 4):**
```makefile
make build     # compile for Linux AMD64
make run       # run locally for development
make deploy    # build + scp to Droplet + restart service (to add)
make logs      # tail journalctl output from Droplet (to add)
make ssh       # SSH into the Droplet as admin (to add)
```

**DNS:** A record on `sanikasurose.com` pointing to the Droplet's public IP for SSH. Web uses 443 on the same host or a separate hosting provider.

### Windows Compatibility

Windows Terminal (Windows 10+) supports SSH natively and renders the app correctly. No special disclaimer or detection is implemented; recommend Windows Terminal, iTerm2, or WezTerm in README if needed.
