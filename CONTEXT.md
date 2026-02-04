# HireMe — Project Context

## Current State

**Status:** API backend functional, Frontend MVP complete

### What's Working

**Backend:**
- ✅ Infrastructure: PostgreSQL + Gotenberg via Docker
- ✅ Database: Migrations applied, dev data seeded
- ✅ API endpoints: Health, Users, CVs (full CRUD), Assets
- ✅ Auth bypass: Dev mode uses `dev-user-001`
- ✅ Repository layer: sqlc queries wired to domain types
- ✅ Asset management: Upload, retrieve, delete with local storage
- ✅ Storage abstraction: LocalStorage implemented, R2 placeholder ready

**Frontend:**
- ✅ Three-column editor layout (palette | preview | properties)
- ✅ Dashboard page with CV list
- ✅ 6 section editors (Personal, Summary, Experience, Education, Skills, Languages)
- ✅ Drag & drop section reordering (dnd-kit)
- ✅ Live A4 CV preview with zoom (50-200%)
- ✅ Undo/redo with 50-step history
- ✅ Auto-save with 2s debounce
- ✅ Export modal (PDF, DOCX, JSON)
- ✅ Keyboard shortcuts (Ctrl+Z/Y, Ctrl+/-/0)
- ✅ 64 passing tests (Vitest)

### What's Not Working Yet
- ❌ Export (`POST /api/v1/export/{format}`) — returns 501, needs Gotenberg integration
- ❌ R2 cloud storage — placeholder only, needs AWS SDK implementation

### Frontend Architecture
See `web/FRONTEND_MVP_PLAN.md` for full details.

### To-Do (MVP API)

- [ ] **Finishing touches on MVP API**
  - Gotenberg integration for export endpoint
  - Error response consistency audit
  - Input validation coverage

- [ ] **Testing** → See [TESTING_PLAN.md](./TESTING_PLAN.md)
  - Phase 1: Domain, Validator, HTTP Utils (P0)
  - Phase 2: Service layer with mocks (P0)
  - Phase 3: Handler layer (P1)
  - Phase 4: Auth middleware (P1)
  - Phase 5: Repository integration (P2)
  - Phase 6: Storage tests (P2)

- [ ] **CI/CD improvements**
  - Add `go vet` to lint job
  - Run migrations before integration tests
  - Coverage threshold (60% minimum)

### API Endpoints (Verified Working)

```
GET  /health                    → {"status": "healthy"}
GET  /ready                     → {"status": "ready", "services": {...}}
GET  /api/v1/users/me           → Current user profile
PATCH /api/v1/users/me          → Update user profile
GET  /api/v1/cv                 → User's active CV
POST /api/v1/cv                 → Create new CV
PUT  /api/v1/cv/{id}            → Update CV
DELETE /api/v1/cv/{id}          → Delete CV
POST /api/v1/assets             → Upload asset (image)
GET  /api/v1/assets/{id}        → Get asset metadata (or file with Accept: image/*)
DELETE /api/v1/assets/{id}      → Delete asset
```

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.22+, Chi router, sqlc |
| Frontend | Next.js 14, React 18, Tailwind, shadcn/ui |
| Database | PostgreSQL 16 (JSONB for CV content) |
| Export | Gotenberg (HTML → PDF/DOCX) |
| Storage | Local filesystem (dev), Cloudflare R2 (prod) |

## Project Structure

```
hireme/
├── api/                    # Go backend
│   ├── cmd/server/         # Entry point
│   ├── db/                 # Migrations, queries, seeds
│   │   ├── migrations/     # SQL migrations
│   │   ├── queries/        # sqlc query definitions
│   │   └── seed/           # dev_seed.sql
│   └── internal/
│       ├── config/         # Environment loading
│       ├── domain/         # Domain types (User, CV, Asset)
│       ├── handler/        # HTTP handlers
│       ├── middleware/     # Auth, logging
│       ├── repository/     # Data access (postgres/)
│       ├── service/        # Business logic
│       └── validator/      # JSON schema validation
├── web/                    # Next.js frontend
│   └── src/
│       ├── app/            # Next.js pages (dashboard, editor)
│       ├── components/
│       │   ├── ui/         # shadcn/ui components
│       │   ├── layout/     # AppShell, Header
│       │   └── editor/     # CV editor components
│       │       ├── editors/  # Section editors
│       │       └── modals/   # Export, Delete modals
│       ├── stores/         # Zustand (editor + UI)
│       ├── hooks/          # useAutoSave, useKeyboardShortcuts
│       ├── lib/
│       │   ├── api/        # Typed HTTP client
│       │   └── dnd/        # dnd-kit configuration
│       └── types/          # TypeScript types
├── docker/                 # Docker configs
├── schemas/                # Shared JSON schemas
└── scripts/                # Dev utilities
```

## Architecture

```
HTTP Request → Handler → Service → Repository → PostgreSQL
                 ↓
              Middleware (Auth, Logging)
```

- **Handlers**: Parse HTTP, validate input, call services
- **Services**: Business logic, orchestration
- **Repositories**: Data access via sqlc-generated code
- **Domain**: Pure types, no external dependencies

## Key Patterns

### Repository Implementation
```go
// pgx.ErrNoRows → domain.ErrNotFound
user, err := r.q.GetUserByID(ctx, id)
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domain.ErrNotFound
}
```

### Auth Bypass (Development)
Set `AUTH_BYPASS_ENABLED=true` in `.env` — all requests use `dev-user-001`

### sqlc Workflow
1. Write SQL in `api/db/queries/*.sql`
2. Run `task api:sqlc`
3. Use generated `queries.Queries` in repositories

## Common Commands

```bash
# Infrastructure
task infra:up          # Start PostgreSQL + Gotenberg
task infra:down        # Stop infrastructure

# Development
task api:dev           # Run API with hot reload
task web:dev           # Run Next.js dev server

# Database
task db:migrate        # Run migrations
task db:seed           # Seed dev user + sample CV
task db:psql           # Open psql shell

# Code Generation
task api:sqlc          # Generate sqlc code
```

## Development Setup

```bash
task setup             # Full setup (infra, migrate, seed, deps)
# OR manually:
task infra:up
task db:migrate
task db:seed
task api:dev
```

## Testing the API

```bash
# Health check
curl http://localhost:8080/health

# Get current user (auth bypass must be enabled)
curl http://localhost:8080/api/v1/users/me

# Get CV
curl http://localhost:8080/api/v1/cv
```
