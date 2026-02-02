# HireMe — Project Context

## Project Overview

HireMe is an open-source CV/resume generation web application. It allows professionals to build, manage, and export CVs using a schema-driven approach.

**Philosophy:** "Pragmatic Quality" — high modularity and best practices without over-engineering.

## Tech Stack

| Component | Technology | Notes |
|-----------|------------|-------|
| Language | Go 1.22+, TypeScript | Backend + Frontend |
| Backend | Chi router, sqlc | Go REST API |
| Frontend | Next.js 14 (App Router), React 18, Tailwind CSS, shadcn/ui | |
| Database | PostgreSQL 16 | JSONB for CV content |
| Export | Gotenberg | HTML → PDF/DOCX |
| Storage | Local filesystem (dev), Cloudflare R2 (prod) | |

## Project Structure

```
hireme/
├── api/          # Go backend (Chi + sqlc)
├── web/          # Next.js frontend
├── docker/       # Container configs
├── schemas/      # Shared JSON schemas
├── scripts/      # Dev utilities
└── docs/         # Documentation
```

## Conventions

### Code Style
- Go: `gofmt`, `golangci-lint`
- Frontend: ESLint, Prettier
- Naming: Go standard, TypeScript camelCase

### Git
- Branch naming: `feature/`, `fix/`, `chore/`
- Commit format: `[TAG] description`
- Tags: FEAT, FIX, REFACTOR, DOCS, TEST, CHORE

## Common Commands

All commands via Taskfile:

```bash
# Infrastructure
task infra:up      # Start PostgreSQL + Gotenberg
task infra:down    # Stop infrastructure

# Development
task api:dev       # Run Go API with hot reload
task web:dev       # Run Next.js dev server

# Database
task db:migrate    # Run database migrations
task db:reset      # Reset database

# Quality
task test          # Run all tests
task lint          # Run all linters
```

## Architecture Principles

1. **Separation of Concerns**: Handlers → Services → Repositories
2. **Domain-First**: Pure domain models with no external dependencies
3. **Interface-Based**: All major components use interfaces for testability
4. **Configuration-Driven**: Feature flags and limits via environment variables

## Key Patterns

### Backend (Go)
- Handlers receive HTTP requests, validate input, call services
- Services contain business logic, orchestrate repositories
- Repositories handle data access (sqlc generated)
- Middleware handles cross-cutting concerns (auth, logging)

### Frontend (Next.js)
- App Router with locale-based routing
- Server Components for static/marketing pages
- Client Components for interactive editor
- Zustand for client-side state management

## Authentication

**Development Mode:** Auth bypass enabled via `AUTH_BYPASS_ENABLED=true`
- Uses a seeded dev user (id: `dev-user-001`)
- No OIDC provider needed for local development

**Production:** Google OIDC via go-oidc library

## Database

PostgreSQL with:
- JSONB for flexible CV content storage
- Row Level Security for tenant isolation
- pgcrypto for field-level encryption

Migrations in `api/db/migrations/`, managed by golang-migrate.

## Testing Strategy

- Go: `go test` with testify assertions
- Frontend: Vitest + React Testing Library
- E2E: Playwright (post-MVP)

## Current Focus

Check `docs/ROADMAP.md` for current development phase and priorities.

## Directory-Specific Context

Each major directory has its own `CLAUDE.md` with specific guidance:
- `api/CLAUDE.md` — Backend development context
- `web/CLAUDE.md` — Frontend development context

## Common Tasks

### Adding a New API Endpoint
1. Define types in `api/internal/domain/`
2. Add sqlc query in `api/db/queries/`
3. Run `task api:sqlc` to generate code
4. Implement repository method in `api/internal/repository/postgres/`
5. Implement service method in `api/internal/service/`
6. Add handler in `api/internal/handler/`
7. Register route in `api/cmd/server/main.go`
8. Write tests

### Adding a New CV Section Type
1. Update `schemas/cv-schema.json`
2. Run `task generate:types` to update TypeScript types
3. Add section editor component in `web/src/components/editor/sections/`
4. Update `SectionList.tsx` to render new type
5. Update preview templates in `api/templates/`

### Modifying Database Schema
1. Create new migration: `task db:migration:create name=description`
2. Write up/down SQL in `api/db/migrations/`
3. Update sqlc queries if needed
4. Run `task db:migrate`
5. Regenerate sqlc: `task api:sqlc`
