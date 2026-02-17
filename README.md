<p align="center">
  <img src="docs/hireme1.png" alt="HireMe Logo" width="160" />
</p>

<h1 align="center">HireMe</h1>

<p align="center">
  <strong>Open-source, schema-driven CV builder for professionals.</strong><br />
  Build, customize, and export beautiful CVs with a modern drag-and-drop editor.
</p>

<p align="center">
  <a href="https://github.com/Ellebam/hireme/actions/workflows/ci.yml"><img src="https://github.com/Ellebam/hireme/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="MIT License" /></a>
</p>

<p align="center">
  <img src="docs/hireme2.png" alt="HireMe Hero" width="480" />
</p>

---

## Features

- **3 Professional Templates** — Classic, Modern, and Visionary layouts
- **Drag & Drop Editor** — Reorder sections intuitively with live preview
- **Live A4 Preview** — Real-time rendering with zoom (50–200%)
- **Multi-Format Export** — PDF, DOCX, and JSON
- **8 Section Types** — Personal, Summary, Experience, Education, Skills, Languages, Certifications, Projects
- **Undo / Redo** — 50-step history with keyboard shortcuts
- **Auto-Save** — 2-second debounced saves with manual override
- **Keyboard Shortcuts** — Ctrl+Z/Y, Ctrl+/-/0 for zoom
- **Privacy First** — Self-hosted, your data stays yours

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 15, React 19, TypeScript, Tailwind CSS, shadcn/ui, Zustand |
| Backend | Go 1.24, Chi router, sqlc |
| Database | PostgreSQL 16 (JSONB for CV content) |
| Export | Gotenberg (HTML → PDF), godocx (CV data → DOCX) |
| Testing | Vitest (170+ frontend tests), Go testing (backend) |
| CI/CD | GitHub Actions — lint, test, typecheck, build |

## Quick Start

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker](https://www.docker.com/)
- [Task](https://taskfile.dev/) (task runner)

### Setup

```bash
git clone https://github.com/Ellebam/hireme.git
cd hireme
task setup
```

### Development

```bash
# Start everything (PostgreSQL, Gotenberg, API, Web)
task dev
```

Or start services individually:

```bash
task infra:up     # PostgreSQL + Gotenberg
task api:dev      # API server (hot reload) → :8080
task web:dev      # Next.js dev server → :3000
```

Open [http://localhost:3000](http://localhost:3000) to use the app.

### Useful Commands

```bash
task                      # List all available tasks
task dev                  # Start full dev environment
task dev:stop             # Stop everything
task test                 # Run all tests
task lint                 # Run linters
task db:migrate           # Run database migrations
task db:seed              # Seed sample data
task db:psql              # Open PostgreSQL shell
```

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Next.js Frontend                   │
│  Zustand stores  ·  shadcn/ui  ·  dnd-kit  ·  A4   │
└──────────────────────────┬──────────────────────────┘
                           │ REST API
┌──────────────────────────▼──────────────────────────┐
│                    Go Backend                        │
│  Handler → Service → Repository → PostgreSQL         │
│              ↓                                       │
│         Middleware (Auth, CORS, Logging)              │
└──────────────────────────┬──────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
        ┌─────▼─────┐           ┌──────▼──────┐
        │ PostgreSQL │           │  Gotenberg   │
        │  (JSONB)   │           │  (PDF export)│
        └────────────┘           └─────────────┘
```

## Project Structure

```
hireme/
├── api/                    # Go backend
│   ├── cmd/server/         # Entry point
│   ├── internal/           # Handlers, services, repos, domain
│   └── db/                 # Migrations, queries, seeds
├── web/                    # Next.js frontend
│   └── src/
│       ├── components/     # Editor, templates, UI
│       ├── stores/         # Zustand (editor + UI state)
│       ├── hooks/          # Auto-save, keyboard shortcuts
│       └── lib/            # API client, templates, logger
├── schemas/                # Shared JSON schemas
├── docker/                 # Container configs
├── scripts/                # Dev utilities
└── Taskfile.yml            # Task runner config
```

## Configuration

Copy `.env.example` to `.env` and adjust:

```bash
cp .env.example .env
```

Key settings:

| Variable | Description | Default |
|----------|-------------|---------|
| `AUTH_BYPASS_ENABLED` | Skip auth in development | `true` |
| `DATABASE_URL` | PostgreSQL connection string | (see `.env.example`) |
| `STORAGE_BACKEND` | File storage backend | `local` |

## API Endpoints

```
GET    /health                → Health check
GET    /api/v1/users/me       → Current user
PATCH  /api/v1/users/me       → Update profile
GET    /api/v1/cv             → Get active CV
POST   /api/v1/cv             → Create CV
PUT    /api/v1/cv/{id}        → Update CV
DELETE /api/v1/cv/{id}        → Delete CV
POST   /api/v1/assets         → Upload asset
GET    /api/v1/assets/{id}    → Get asset
DELETE /api/v1/assets/{id}    → Delete asset
POST   /api/v1/export/{fmt}   → Export CV (pdf, docx)
```

## Contributing

Contributions are welcome! Please check the [issue templates](.github/ISSUE_TEMPLATE) before opening a new issue.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes
4. Push and open a Pull Request

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Gotenberg](https://gotenberg.dev/) — HTML to PDF conversion
- [shadcn/ui](https://ui.shadcn.com/) — UI component library
- [dnd-kit](https://dndkit.com/) — Drag and drop toolkit
- [godocx](https://github.com/nicholasgasior/godocx) — DOCX generation in Go
