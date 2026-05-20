package content

// Project holds all data for a single portfolio project.
type Project struct {
	Title     string
	Slug      string
	Year      string
	Tags      []string
	ShortDesc string
	LongDesc  string // Markdown, rendered by Glamour in Phase 2
	Links     []Link
	Event     string
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

var Projects = []Project{
	{
		Title:     "VeriHire",
		Slug:      "verihire",
		Year:      "2026",
		Tags:      []string{"Python", "FastAPI", "PostgreSQL", "Docker", "LLMs"},
		ShortDesc: "Privacy-preserving AI hiring platform built at a hackathon.",
		Event:     "Midnight Hackathon 2026 — Mississauga, ON",
		LongDesc: `## VeriHire

A privacy-first AI-powered hiring platform built at **Midnight Hackathon 2026**.
Designed to process sensitive candidate information without exposing underlying personal data.

### What I built
- Scalable backend APIs and document-processing workflows using FastAPI, PostgreSQL, and Docker
- LLM-powered candidate evaluation, document parsing, and verification assistance
- Midnight blockchain concepts for privacy-preserving credential validation and controlled data sharing

### Stack
FastAPI · PostgreSQL · Docker · LangChain · Prompt Engineering · Cursor · Claude Code
`,
	},
	{
		Title:     "Piratech",
		Slug:      "piratech",
		Year:      "2026",
		Tags:      []string{"Python", "FastAPI", "WebSockets", "LLMs", "Security"},
		ShortDesc: "AI-powered security analysis platform for automated vulnerability detection.",
		Event:     "Hacktech 2026 — California Institute of Technology, Los Angeles, CA",
		LongDesc: `## Piratech

An AI-powered security analysis platform built at **Hacktech 2026** at Caltech.
Performs automated vulnerability detection and real-time code analysis.

### What I built
- Backend services and async processing workflows using FastAPI and WebSockets
- Static analysis tools + LLM-based reasoning to identify risks and generate remediation suggestions
- Real-time frontend dashboards for interactive security monitoring

### Stack
FastAPI · WebSockets · LLM reasoning workflows · Cursor · Claude Code · GitHub Copilot
`,
	},
	{
		Title:     "Committr",
		Slug:      "committr",
		Year:      "2026",
		Tags:      []string{"Java", "Spring Boot", "PostgreSQL", "Docker", "CI/CD"},
		ShortDesc: "Full-stack developer analytics platform for tracking GitHub activity.",
		LongDesc: `## Committr

A full-stack developer analytics platform for tracking GitHub activity,
engineering productivity metrics, and repository insights.

### What I built
- Scalable REST APIs and relational database models using Spring Boot and PostgreSQL
- Automated ETL-style processing pipelines and real-time dashboards
- Docker-based deployment with CI/CD integration and modular backend architecture

### Stack
Java · Spring Boot · PostgreSQL · Docker · CI/CD · Jenkins
`,
	},
	{
		Title:     "Surose OS",
		Slug:      "surose-os",
		Year:      "2025",
		Tags:      []string{"Go", "Wish", "Bubbletea", "Lipgloss", "SSH"},
		ShortDesc: "The SSH portfolio TUI you're currently inside.",
		Links:     []Link{{Label: "GitHub", URL: "github.com/sanikasurose/surose-os"}},
		LongDesc: `## Surose OS

The thing you're using right now.

Built to prove a point: terminal interfaces don't have to look like 1985.
Also built because I wanted to learn Go, and the best way I know how to learn
something is to build something real with it.

### Stack
Go · Wish · Bubbletea · Lipgloss · Glamour

### Hosting
DigitalOcean Droplet · systemd · sanikasurose.com
`,
	},
}

var Contact = ContactInfo{
	Email:    "sanikasurose@gmail.com",
	LinkedIn: "linkedin.com/in/sanikasurose",
	GitHub:   "github.com/sanikasurose",
	Note:     "Email is the best way to reach me.",
}

const AboutText = `Sanika Surose
Software Engineer · McMaster University, Class of 2029

Currently a Software Engineer Intern at Tranquility Inc., building backend APIs,
automation workflows, and full-stack features with FastAPI and React/Next.js.
Also an Undergraduate Research Assistant at McMaster, fine-tuning transformer-based
NLP models and building reproducible ML pipelines under Dr. Charles Welch.

Second-year Software Engineering co-op student. I work across the stack —
backends, distributed systems, LLM tooling, and occasionally things like this.

Built this in Go, which I'd never used before. That was the point.

McMaster University
B.Eng. Software Engineering (Co-op) · Expected 2029`
