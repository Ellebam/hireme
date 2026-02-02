# HireMe

**Open-source, schema-driven CV generator for professionals.**

Build, manage, and export beautiful CVs with a modern drag-and-drop editor. Export to PDF, DOCX, or raw data formats for use with AI tools.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features

- 🎨 **3 Professional Templates** — Classic, Modern, Minimal
- 🖱️ **Drag & Drop Editor** — Intuitive section management
- 📄 **Multi-Format Export** — PDF, DOCX, JSON, YAML
- 🔒 **Privacy First** — Your data stays yours
- 🌍 **Multi-Language** — English and German (more coming)
- 📱 **Responsive** — Works on desktop, tablet, and mobile

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 14, React 18, TypeScript, Tailwind CSS, shadcn/ui |
| Backend | Go 1.22+, Chi router, sqlc |
| Database | PostgreSQL 16 |
| Export | Gotenberg (HTML → PDF/DOCX) |
| Storage | Local filesystem / Cloudflare R2 |

## Quick Start

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker](https://www.docker.com/)
- [Task](https://taskfile.dev/) (task runner)

### Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/hireme.git
cd hireme

# Run setup script (installs deps, starts DB, runs migrations)
./scripts/setup-dev.sh

# Or do it manually:
task setup
```

### Development

```bash
# Start infrastructure (PostgreSQL, Gotenberg)
task infra:up

# In terminal 1: Start API server (with hot reload)
task api:dev

# In terminal 2: Start web dev server
task web:dev
```

The app will be available at:
- **Web**: http://localhost:3000
- **API**: http://localhost:8080

### Useful Commands

```bash
task                    # Show all available tasks
task dev                # Start dev environment
task test               # Run all tests
task lint               # Run linters
task db:migrate         # Run database migrations
task db:psql            # Open PostgreSQL shell
task infra:logs         # View infrastructure logs
```

## Project Structure

```
hireme/
├── api/                # Go backend
├── web/                # Next.js frontend
├── docker/             # Container configs
├── schemas/            # JSON schemas (source of truth)
├── scripts/            # Dev utilities
├── docs/               # Documentation
└── Taskfile.yml        # Task runner config
```

See [CLAUDE.md](CLAUDE.md) for detailed architecture information.

## Configuration

Copy `.env.example` to `.env.local` and adjust as needed:

```bash
cp .env.example .env.local
```

Key settings:
- `AUTH_BYPASS_ENABLED=true` — Skip auth in development
- `STORAGE_BACKEND=local` — Use local filesystem for uploads
- `DATABASE_URL` — PostgreSQL connection string

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Development Guide](docs/DEVELOPMENT.md)
- [Deployment](docs/DEPLOYMENT.md)

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gotenberg](https://gotenberg.dev/) for document generation
- [shadcn/ui](https://ui.shadcn.com/) for UI components
- [dnd-kit](https://dndkit.com/) for drag-and-drop
