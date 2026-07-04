# Surose OS — Content Reference

This file contains all portfolio content for Sanika Surose.
The coding agent reads this file before writing or updating `internal/content/data.go`.
No content should be hardcoded in UI files — everything lives in data.go, sourced from here.

---

## Projects

All projects have a Category field: `hackathon` or `personal`.
Projects must appear in the exact order listed below within each category.

---

### Hackathon Projects
Display order: CLAIRITY → SeamSecure → Piratech → VeriHire

---

#### CLAIRITY
- **Slug:** `clairity`
- **Year:** 2026
- **Category:** hackathon
- **Event:** DeltaHacks 12
- **Location:** McMaster University, Hamilton, ON
- **Tags:** Python, Swift, LLMs, OCR, iOS
- **Short description:** AI-powered iOS app that translates complex medical and institutional language into clear, actionable next steps.
- **GitHub:** github.com/syaanmerchant/clAIrity
- **Long description (Markdown):**
```markdown
## CLAIRITY

Built at **DeltaHacks 12** at McMaster University.

An AI-powered iOS application that transforms complex medical, insurance, and institutional
language into clear, structured, actionable next steps — designed for patients navigating
high-stakes interactions like ER discharge instructions or insurance documents.

### What I built
- LLM-based reasoning and structured output generation for clinical document parsing
- OCR-driven extraction from PDFs, images, and typed text inputs
- Clinical NLP for identifying medications, diagnoses, timelines, and follow-up requirements
- Privacy-conscious AI orchestration — no personal data stored server-side
- Structured after-care guidance with conditional actions and safety triggers

### Stack
Python · Swift · LLMs · OCR · Clinical NLP · iOS
```

---

#### SeamSecure
- **Slug:** `seamsecure`
- **Year:** 2026
- **Category:** hackathon
- **Event:** ElleHacks 2026
- **Location:** York University, Toronto, ON
- **Tags:** Python, FastAPI, React, TypeScript, Google Gemini, OAuth
- **Short description:** AI-powered platform that detects conversation hijacking and email thread compromise using hybrid rule-based and LLM analysis.
- **GitHub:** github.com/SeamSecure/SeamSecure
- **Long description (Markdown):**
```markdown
## SeamSecure

Built at **ElleHacks 2026** at York University.

A full-stack cybersecurity web application that detects conversation hijacking and email
thread compromise. Most phishing filters catch obvious spam — SeamSecure catches attacks
that happen inside ongoing conversations between people who already trust each other.

### What I built
- AI-powered thread analysis using Google Gemini for intent drift detection and stylometric anomaly identification
- Rule-based heuristic engine for urgency escalation, credential harvesting, suspicious URLs, and domain impersonation
- Scores each thread 0–100% risk and categorizes as Safe, Suspicious, or Dangerous
- Google OAuth 2.0 with read-only Gmail API access for direct inbox thread analysis
- Privacy-first: all email content analyzed in real-time with no server-side storage

### Stack
Python · FastAPI · React · TypeScript · Tailwind CSS · Google Gemini · OAuth 2.0
```

---

#### Piratech
- **Slug:** `piratech`
- **Year:** 2026
- **Category:** hackathon
- **Event:** Hacktech 2026
- **Location:** California Institute of Technology, Los Angeles, CA
- **Tags:** Python, FastAPI, WebSockets, Semgrep, LLMs, React, Vite
- **Short description:** AI-powered security analysis platform for automated vulnerability detection and real-time code analysis.
- **GitHub:** github.com/bebopkenny/hacktech2026part2
- **Long description (Markdown):**
```markdown
## Piratech

Built at **Hacktech 2026** at the California Institute of Technology.

A full-stack security analysis platform that detects exploitable vulnerabilities in real
codebases. Most scanners flood developers with cryptic rule IDs — Piratech reads code the
way a senior security engineer would and explains exactly how each vulnerability can be exploited.

### What I built
- Two-layer pipeline: Semgrep static analysis for candidate detection + K2-Think-V2 reasoning model for cross-file taint tracing and exploitability assessment
- Traces taint paths across files to show exactly how attacker-controlled input reaches a sensitive sink
- Detects authentication gaps and missing access controls at the route and middleware level
- Filters false positives by verifying whether sanitization or ORM parameterization already covers each finding
- GitHub webhook integration for automatic re-scans on every push
- Deployed to Vultr from hour two of the build for live end-to-end pipeline testing

### Stack
Python · FastAPI · WebSockets · Semgrep · React · Vite · Backboard · Vultr
```

---

#### VeriHire
- **Slug:** `verihire`
- **Year:** 2026
- **Category:** hackathon
- **Event:** Midnight Hackathon 2026
- **Location:** Online
- **Tags:** Python, FastAPI, Next.js, TypeScript, PostgreSQL, Docker, Google Gemini, Zero-Knowledge Proofs
- **Short description:** Privacy-first hiring verification platform — candidates prove qualifications without exposing their documents.
- **GitHub:** github.com/sanikasurose/midnight-hackathon
- **Long description (Markdown):**
```markdown
## VeriHire

Built at **Midnight Hackathon 2026** (online).

A privacy-first AI-powered hiring verification platform. Traditional hiring requires
candidates to overshare sensitive documents. VeriHire lets candidates prove their
qualifications using zero-knowledge proofs — employers verify claims without ever
seeing the underlying documents.

### What I built
- AI-powered claim extraction using Google Gemini: parses resumes for education, GPA, skills, certifications
- Zero-knowledge proof generation via the Midnight blockchain ecosystem
- Full privacy-preserving verification workflow — employers receive only proof results
- FastAPI backend with SQLAlchemy + PostgreSQL, Next.js + TypeScript frontend
- Dockerized local development environment with docker-compose

### My role
Backend Lead — owned the API layer, document processing pipeline, database schema, and ZK proof integration.

### Stack
Python · FastAPI · SQLAlchemy · PostgreSQL · Next.js · TypeScript · Tailwind CSS · Google Gemini · Midnight / ZK Proofs · Docker
```

---

### Personal Projects
Display order: typing-test → outreach-research-scraper → neurosnake-rl → sanika-ui → throttle → pulp → committr → surose-os

---

#### Typing Test
- **Slug:** `typing-test`
- **Year:** 2024
- **Category:** personal
- **Tags:** Python, CustomTkinter, GUI
- **Short description:** Feature-rich graphical typing trainer with real-time analytics, session history, and per-character feedback.
- **GitHub:** github.com/sanikasurose/typing-test
- **Long description (Markdown):**
```markdown
## Typing Test

A fully-featured graphical typing trainer built with Python and CustomTkinter.
Migrated from a terminal-based v1 into a multi-screen GUI with persistent session tracking.

### Features
- Real-time per-character accuracy highlighting (green correct, red incorrect)
- Live WPM calculation and accuracy percentage
- Per-character mistake analytics — identifies your most frequently mistyped keys
- Average keystroke delay tracking for consistency analysis
- Persistent session history across runs (sessions.json)
- Performance profile classification: Balanced / Fast but inaccurate / Accurate but slow / Needs consistency
- Multi-screen flow: start → typing → results → progress summary

### Stack
Python · CustomTkinter · Session analytics · Modular GUI architecture
```

---

#### Outreach Research Scraper
- **Slug:** `outreach-research-scraper`
- **Year:** 2025
- **Category:** personal
- **Tags:** Python, Web Scraping, NLP, CLI
- **Short description:** CLI tool that crawls university faculty pages and produces clean summaries of professor research interests and undergraduate opportunities.
- **GitHub:** github.com/sanikasurose/academic-research-scraper
- **Long description (Markdown):**
```markdown
## Outreach Research Scraper

A Python CLI tool for academic research discovery. Crawls university faculty pages
and produces clean, human-readable summaries of professor research interests,
recent publications, and undergraduate research opportunities.

### Features
- Multi-page crawling per professor profile
- Structured ProfessorProfile schema
- Dual output: Markdown summaries for human review + JSON for further processing
- Fully unit-tested formatting layer

### Usage
```bash
python main.py input/professor_urls.txt output/
```

### Stack
Python · Web scraping · NLP · CLI · Unit testing
```

---

#### neurosnake-rl
- **Slug:** `neurosnake-rl`
- **Year:** 2026
- **Category:** personal
- **Tags:** Python, PyTorch, Deep Q-Learning, SQLite, CNN, Reinforcement Learning
- **Short description:** End-to-end reinforcement learning pipeline — CNN-based DQN agent trained to play Snake, with full experiment tracking infrastructure.
- **GitHub:** github.com/sanikasurose/neurosnake-rl
- **Long description (Markdown):**
```markdown
## neurosnake-rl

Started as a Snake game. Became a full reinforcement learning research pipeline.

Most RL tutorials stop at "the model trains." Real ML work requires reproducibility,
experiment tracking, and structured evaluation. This project builds all of that from scratch.

### The agent
- Convolutional Neural Network with frame stacking for temporal awareness
- Double DQN to reduce Q-value overestimation
- Target network stabilization + replay buffer training
- Shaped reward signal: Manhattan distance, food bonuses, collision penalties, step cost to discourage looping
- Device-aware: runs on MPS, CUDA, or CPU with no code changes

### The pipeline
- Every run logged to SQLite with full hyperparameter metadata
- Auto-generated reward curves and convergence detection
- Experiment leaderboard for comparing runs objectively
- Offline analytics CLI for hyperparameter breakdowns and checkpoint rankings
- Deterministic seeding across Python, NumPy, and PyTorch for reproducible results

### Stack
Python · PyTorch · NumPy · SQLite · Matplotlib · CNN · Double DQN
```

---

#### Sanika-UI
- **Slug:** `sanika-ui`
- **Year:** 2025
- **Category:** personal
- **Tags:** React, TypeScript, CSS Variables, Component Library, Design Systems
- **Short description:** Lightweight React UI component library exploring reusable component architecture and theme systems.
- **GitHub:** github.com/sanikasurose/sanika-ui
- **Long description (Markdown):**
```markdown
## Sanika-UI

A lightweight React component library built to understand how design systems work
from the inside — not just how to use them.

### Components
Button · Navbar · Modal · Tabs · Toast · ProjectCard · Container

### Theme system
Custom theme engine using React Context and CSS variables. Three themes (Classic Pink,
Dark Pink, Dev Mode) switchable via UI or keyboard shortcut (T). All components
update automatically — no per-component theme logic.

### Why I built it
Most engineers use Material UI or Chakra without understanding how they work internally.
Building one manually reveals the architecture behind reusable component systems —
how state, context, and CSS variables compose into something cohesive.

### Stack
React · TypeScript · CSS Variables · Context API
```

---

#### Throttle
- **Slug:** `throttle`
- **Year:** 2026
- **Category:** personal
- **Tags:** Java, Spring Boot, Redis, Docker, Maven, JUnit
- **Short description:** Production-grade API rate limiting service with sliding-window algorithm, Redis-backed atomic execution, and 87% test coverage.
- **GitHub:** github.com/sanikasurose/throttle
- **Long description (Markdown):**
```markdown
## Throttle

A production-grade API rate limiting service. Enforces per-user request quotas using a
sliding window algorithm backed by Redis, with atomic Lua script execution to prevent
race conditions under concurrent load.

Fixed-window rate limiters have a boundary exploit — a user can double their quota by
sending requests at the end of one window and the start of the next. Sliding window
eliminates this. Throttle implements it correctly.

### Features
- Sliding-window throttling via atomic Redis Lua script (evict → count → add in one operation)
- Header-based identity via X-User-Id
- /metrics endpoint with per-user allow/reject counters
- Spring Boot Actuator health endpoint
- Docker Compose for local development (Redis + service)
- 87.1% test coverage · CI passing · Quality gate passed

### Stack
Java 17 · Spring Boot 3.5.12 · Redis 7 · Docker · Maven · JUnit 5
```

---

#### Pulp
- **Slug:** `pulp`
- **Year:** 2026
- **Category:** personal
- **Tags:** Python, Claude Haiku, OCR, Tesseract, Docker, CLI, LLMs
- **Short description:** Local-first CLI that converts PDFs (text-layer or scanned) into clean, semantically structured Markdown for LLM/RAG ingestion.
- **GitHub:** github.com/sanikasurose/pulp
- **Long description (Markdown):**
```markdown
## Pulp

A local-first CLI that converts PDFs — text-layer or scanned — into clean,
semantically structured Markdown suitable for LLM and RAG ingestion.

Most PDF-to-text tools produce noisy output that breaks downstream pipelines.
Pulp cleans, structures, and optionally uses an LLM to semantically organize
the result. It never crashes on LLM failure — falls back to heuristic output gracefully.

### Features
- Handles both text-layer and scanned PDFs (OCR via Tesseract + pdf2image)
- Optional LLM structuring stage via Claude Haiku (Anthropic API)
- --diff flag for token counts and per-stage cleaning summary
- Multi-column layout heuristic (auto-detects two-column PDFs)
- Docker support: Poppler + Tesseract bundled, runs as non-root user
- 86% test coverage · 36 tests passing · 0 code smells

### Stack
Python 3.11 · Claude Haiku · Tesseract · Poppler · Docker · uv
```

---

#### Committr
- **Slug:** `committr`
- **Year:** 2026
- **Category:** personal
- **Tags:** Java, Spring Boot, PostgreSQL, Docker, CI/CD, Jenkins
- **Short description:** Full-stack developer analytics platform for tracking GitHub activity, engineering productivity metrics, and repository insights.
- **GitHub:** github.com/sanikasurose/committr
- **Long description (Markdown):**
```markdown
## Committr

A full-stack developer analytics platform for tracking GitHub activity,
engineering productivity metrics, and repository insights.

### What I built
- Scalable REST APIs and relational database models using Spring Boot and PostgreSQL
- Automated ETL-style processing pipelines and real-time analytics dashboards
- Docker-based deployment with CI/CD integration and modular backend architecture

### Stack
Java · Spring Boot · PostgreSQL · Docker · CI/CD · Jenkins
```

---

#### Surose OS
- **Slug:** `surose-os`
- **Year:** 2025
- **Category:** personal
- **Tags:** Go, Wish, Bubbletea, Lipgloss, Glamour, SSH
- **Short description:** The SSH portfolio TUI you're currently inside.
- **GitHub:** github.com/sanikasurose/surose-os
- **Long description (Markdown):**
```markdown
## Surose OS

The thing you're using right now.

Built to prove a point: terminal interfaces don't have to look like 1985.
Also built because I wanted to learn Go, and the best way I know how to learn
something is to build something real with it.

### Stack
Go · Wish · Bubbletea · Lipgloss · Glamour

### Hosting
DigitalOcean Droplet · systemd · sanikasurose.com

### Source
Open source. Link below.
```

---

### Course Project

#### International Airport — Smart Baggage Handling Mechanism
- **Slug:** `airport-baggage`
- **Year:** 2024
- **Category:** personal
- **Tags:** Python, OOP, CAD, Simulation, Hardware Integration
- **Course:** ENG 1P13 — Integrated Cornerstone Design Projects in Engineering, McMaster University
- **Short description:** Drawbridge-inspired mechanical ramp system for airport baggage transfers — simulation, control logic, and physical prototype.
- **Long description (Markdown):**
```markdown
## International Airport — Smart Baggage Handling

A team engineering design project from **ENG 1P13** at McMaster University.
Designed and prototyped a drawbridge-inspired mechanical ramp system to improve
the efficiency and safety of airport baggage transfers between platforms.

### What I worked on
- Python-based simulation and control logic for a rotary actuator and rack-and-pinion mechanism
- Object-oriented models to represent mechanical components and motion behavior
- Software validation through iterative testing and refinement
- Integrating software control with the physical prototype in a team environment

### Stack
Python · OOP · CAD · Actuator control · Motion simulation · Hardware integration
```

---

## Experience

### Work Roles
Display order: most recent first.

---

**Software Engineer Intern**
- Company: TranQuility Inc.
- Location: Toronto, ON (Remote)
- Dates: May 2026 – Present
- Bullets:
  - Developing backend APIs and automation workflows using FastAPI, PostgreSQL, and Redis
  - Building full-stack application features with React/Next.js in a fast-paced startup environment
  - Integrating external APIs and supporting deployment, debugging, monitoring, and system improvements
  - Contributing to AI-driven application workflows and real-time user operations
  - Collaborating in Agile development cycles with ownership over feature implementation and testing

---

**Undergraduate Research Assistant**
- Supervisor: Dr. Charles Welch
- Institution: McMaster University — Faculty of Engineering
- Location: Hamilton, ON
- Dates: May 2026 – Present
- Research area: Natural Language Processing (NLP) and Applied Machine Learning
- Bullets:
  - Fine-tuning and evaluating transformer-based NLP models on large-scale text datasets using Python and PyTorch
  - Building reproducible data-processing and model-training pipelines for experimentation, benchmarking, and evaluation workflows
  - Analyzing model outputs, accuracy metrics, and performance trends across multiple datasets and training configurations
  - Working with modern ML and LLM workflows including data preprocessing, tokenization, and inference evaluation
  - Documenting experimental methodologies, technical findings, and research workflows in collaborative academic environments

---

### Hackathon Entries
Display order: most recent first (same order as hackathon projects list above reversed — VeriHire first on experience page).

**IMPORTANT:** On the Experience screen, hackathon entries show ONLY the event name, location, date, project name, and a single brief description line. No bullet points. No full project description. End with a note pointing the user to the Projects screen for full detail.

Format for each hackathon entry on the Experience screen:
```
[Event Name]
[Location] · [Month Year]
Built [Project Name] — [one sentence description]
→ full project detail in projects
```

---

**Midnight Hackathon 2026**
- Project: VeriHire
- Location: Online
- Date: May 2026
- Brief: Built VeriHire — a privacy-first hiring platform using zero-knowledge proofs so candidates prove qualifications without exposing documents.

**Hacktech 2026**
- Project: Piratech
- Location: California Institute of Technology, Los Angeles, CA
- Date: April 2026
- Brief: Built Piratech — an AI-powered security analysis platform that explains exactly how vulnerabilities in real codebases can be exploited.

**ElleHacks 2026**
- Project: SeamSecure
- Location: York University, Toronto, ON
- Date: February 2026
- Brief: Built SeamSecure — a platform that detects conversation hijacking and email thread compromise using LLM-based intent drift analysis.

**DeltaHacks 12**
- Project: CLAIRITY
- Location: McMaster University, Hamilton, ON
- Date: January 2026
- Brief: Built CLAIRITY — an iOS app that translates complex medical and institutional documents into clear, actionable next steps.

---

## About

Content lives directly in `internal/ui/about.go`, organized into sections.

```
Sanika Surose
Software Engineer · McMaster University, Class of 2029

── now ──────────────────────────────────────────────
i'm a Software Engineer Intern at TranQuility Inc., an AI startup, where I
build backend APIs, RAG data pipelines, and full-stack features using
FastAPI, React, and Azure. My work has included shipping a production RAG
pipeline (Python + Qdrant) spanning multiple document formats, cutting AI
feature costs significantly through model and token optimization,
resolving a production outage caused by async I/O bottlenecks, and
migrating our deployment workflow from a multi-step manual process to a
streamlined Azure Container Registry pipeline.

i'm also an Undergraduate Research Assistant at McMaster University under
Dr. Charles Welch, working on NLP and LLM evaluation. I've built a crash-
safe, resumable inference pipeline serving Qwen3-8B via vLLM on a 4×H100
GPU cluster, and designed zero-shot and few-shot emotion classification
experiments on the GoEmotions dataset, benchmarking model performance
across prompting strategies.

── how ──────────────────────────────────────────────
third-year software engineering co-op student. i work across
the stack — backends, distributed systems, LLM tooling, and
occasionally things like this.

── elsewhere ────────────────────────────────────────
sponsorship exec for DeltaHacks, McMaster's largest hackathon with 600+
attendees. embedded subteam for McMaster Aerial Robotics and Drone Club,
competing in the AEAC Student UAS Competition. building on the software
subteam for McMaster Design League, McMaster's largest CAD-a-thon. outreach
for McMaster Advanced Space Systems.

── education ────────────────────────────────────────
McMaster University                B.Eng. software engineering (co-op) · 2029

built this in Go, which i'd never used before. that was the point.
```

Note: `now`, `how`, and `elsewhere` sections all use the same left accent-bar (`Quote`) styling for visual consistency.

---

## Contact

```
sanikasurose@gmail.com
linkedin.com/in/sanikasurose
github.com/sanikasurose

Email is the best way to reach me.
```