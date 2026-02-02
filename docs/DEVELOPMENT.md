# HireMe — Development Guide

Complete guide for setting up and working with the HireMe development environment.

---

## Prerequisites

### Required Tools

| Tool | Version | Installation |
|------|---------|--------------|
| Go | 1.22+ | https://go.dev/dl/ |
| Node.js | 20+ | https://nodejs.org/ |
| Docker | Latest | https://www.docker.com/ |
| Task | 3.x | https://taskfile.dev/installation/ |

### Recommended Tools

| Tool | Purpose |
|------|---------|
| Air | Go hot reload (auto-installed by `task api:dev`) |
| golangci-lint | Go linter (auto-installed by `task api:lint`) |
| sqlc | SQL code gen (auto-installed by `task api:sqlc`) |

### Verify Installation

```bash
go version      # Should be 1.22+
node --version  # Should be 20+
docker --version
task --version
```

---

## Initial Setup

### 1. Clone and Enter Repository

```bash
git clone https://github.com/yourusername/hireme.git
cd hireme
```

### 2. Run Setup Script

The setup script handles everything:

```bash
./scripts/setup-dev.sh
```

This will:
- Check dependencies
- Create `.env.local` from `.env.example`
- Create data directories
- Start PostgreSQL and Gotenberg
- Install Go and Node dependencies
- Run database migrations
- Seed the development user

### 3. Manual Setup (Alternative)

If you prefer manual setup:

```bash
# Create environment file
cp .env.example .env.local

# Create data directories
mkdir -p data/uploads data/postgres

# Start infrastructure
task infra:up

# Install dependencies
task api:deps
task web:deps

# Run migrations
task db:migrate

# Seed dev user
task db:seed
```

---

## Running the Application

### Start Infrastructure

```bash
task infra:up
```

This starts:
- **PostgreSQL** on `localhost:5432`
- **Gotenberg** on `localhost:3001`

### Start Development Servers

In separate terminals:

```bash
# Terminal 1: Go API (with hot reload)
task api:dev
# → http://localhost:8080

# Terminal 2: Next.js (with hot reload)
task web:dev
# → http://localhost:3000
```

### Stop Everything

```bash
task infra:down
```

---

## Common Development Tasks

### Database Operations

```bash
task db:psql           # Open PostgreSQL shell
task db:migrate        # Apply pending migrations
task db:migrate:down   # Rollback last migration
task db:reset          # Drop and recreate database (DESTRUCTIVE)

# Create new migration
task db:migration:create name=add_new_feature
```

### Code Generation

```bash
task api:sqlc          # Regenerate sqlc Go code
task generate:types    # Generate TypeScript types from JSON Schema
task generate          # Run all code generation
```

### Testing

```bash
task test              # Run all tests
task api:test          # Run Go tests only
task api:test:coverage # Go tests with coverage report
task web:test          # Run frontend tests
task web:test:watch    # Frontend tests in watch mode
```

### Linting & Formatting

```bash
task lint              # Run all linters
task fmt               # Format all code
task api:lint          # Go linter only
task web:lint          # ESLint only
task web:typecheck     # TypeScript type check
```

### Building

```bash
task api:build         # Build Go binary → api/bin/server
task web:build         # Build Next.js → web/.next
task docker:build      # Build Docker images locally
```

---

## Authentication

### Development Mode (Default)

Auth bypass is enabled by default in `.env.local`:

```env
AUTH_BYPASS_ENABLED=true
AUTH_BYPASS_USER_ID=dev-user-001
```

All API requests are automatically authenticated as the dev user. No login required.

### Testing with Real Auth

To test Google OIDC:

1. Create credentials at https://console.cloud.google.com/apis/credentials
2. Update `.env.local`:
   ```env
   AUTH_BYPASS_ENABLED=false
   GOOGLE_CLIENT_ID=your-client-id
   GOOGLE_CLIENT_SECRET=your-client-secret
   ```
3. Restart the API server

---

## Project Structure

```
hireme/
├── api/                      # Go Backend
│   ├── cmd/server/           # Entry point
│   ├── internal/
│   │   ├── config/           # Configuration
│   │   ├── domain/           # Domain models
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # HTTP middleware
│   │   ├── repository/       # Data access
│   │   ├── service/          # Business logic
│   │   ├── storage/          # File storage
│   │   └── validator/        # Input validation
│   ├── db/
│   │   ├── migrations/       # SQL migrations
│   │   └── queries/          # sqlc queries
│   └── templates/            # HTML export templates
│
├── web/                      # Next.js Frontend
│   ├── src/
│   │   ├── app/              # App Router pages
│   │   ├── components/       # React components
│   │   ├── lib/              # Utilities
│   │   ├── stores/           # Zustand stores
│   │   └── i18n/             # Translations
│   └── public/               # Static assets
│
├── docker/                   # Docker configs
├── schemas/                  # JSON schemas
├── scripts/                  # Dev scripts
└── docs/                     # Documentation
```

---

## Code Conventions

### Go

- Use `internal/` for non-public packages
- Repository pattern for data access
- Service layer for business logic
- Handlers are thin — validate, call service, respond
- Structured logging with `slog`
- Errors wrap context: `fmt.Errorf("doing thing: %w", err)`

### TypeScript/React

- Functional components only
- Use `"use client"` only when necessary
- Zustand for client state
- API client in `lib/api/`
- Types in `types/` directory
- Tailwind for styling, shadcn/ui for components

### File Naming

| Type | Convention | Example |
|------|------------|---------|
| Go files | snake_case | `user_service.go` |
| Go test files | snake_case + _test | `user_service_test.go` |
| React components | PascalCase | `CVEditor.tsx` |
| Utilities | camelCase | `utils.ts` |

---

## Debugging

### API Logs

```bash
# View real-time logs
task api:dev

# Logs include request ID, method, path, status, duration
# 2024-01-15 10:30:45 INFO request method=GET path=/api/v1/cv status=200 duration=12ms request_id=abc123
```

### Database Queries

```bash
# Connect to database
task db:psql

# View recent CVs
SELECT id, user_id, title, updated_at FROM cvs ORDER BY updated_at DESC LIMIT 10;

# Check user limits
SELECT id, email, tier, cv_limit, storage_used_bytes FROM users;
```

### Infrastructure

```bash
task infra:status   # Container status
task infra:logs     # Container logs
```

---

## Troubleshooting

### PostgreSQL won't start

```bash
# Check if port 5432 is in use
lsof -i :5432

# Remove stale data and restart
task infra:clean
task infra:up
```

### Migrations fail

```bash
# Check migration status
task db:migrate:status

# If stuck, force version
docker exec hireme-postgres psql -U hireme -d hireme -c "DELETE FROM schema_migrations;"
task db:migrate
```

### Go dependencies issues

```bash
cd api
go clean -modcache
go mod download
go mod tidy
```

### Node dependencies issues

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

### Hot reload not working

**Go (Air):**
```bash
# Reinstall air
go install github.com/air-verse/air@latest
```

**Next.js:**
```bash
# Clear cache
rm -rf web/.next
task web:dev
```

---

## IDE Setup

### VS Code

Recommended extensions:
- Go (golang.go)
- ESLint
- Prettier
- Tailwind CSS IntelliSense
- GitLens

Settings (`.vscode/settings.json`):
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

### Agentic AI IDE (Antigravity)

The project includes `CLAUDE.md` files optimized for AI assistance:
- `/CLAUDE.md` — Project overview and common tasks
- `/api/CLAUDE.md` — Backend-specific context
- `/web/CLAUDE.md` — Frontend-specific context

These provide context for AI agents working on the codebase.

---

## Testing Strategy

### Unit Tests

Test individual functions and methods in isolation.

```go
// api/internal/service/cv_service_test.go
func TestCVService_Create_RespectsLimit(t *testing.T) {
    mockRepo := mock.NewCVRepository()
    mockRepo.CountReturn = 1 // Already has 1 CV
    
    userRepo := mock.NewUserRepository()
    userRepo.GetReturn = &domain.User{CVLimit: 1}
    
    svc := service.NewCVService(mockRepo, userRepo)
    
    _, err := svc.Create(ctx, "user-1", "Test", content)
    assert.ErrorIs(t, err, domain.ErrCVLimitReached)
}
```

### Integration Tests

Test components working together with real database.

```go
func TestCVHandler_CreateCV_Integration(t *testing.T) {
    db := setupTestDB(t) // Uses testcontainers
    defer db.Close()
    
    app := setupApp(db)
    
    req := httptest.NewRequest("POST", "/api/v1/cv", body)
    req.Header.Set("Authorization", "Bearer test-token")
    
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    assert.Equal(t, 201, w.Code)
}
```

### E2E Tests (Post-MVP)

Full user journey tests with Playwright.

---

## Git Workflow

### Branch Naming

```
feature/add-export-pdf
fix/cv-validation-error
chore/update-dependencies
docs/api-documentation
```

### Commit Messages

Follow conventional commits:
```
feat: add PDF export endpoint
fix: correct CV schema validation for dates
docs: update API documentation
test: add integration tests for asset upload
chore: update Go dependencies
```

### Pull Request Process

1. Create feature branch from `main`
2. Make changes with atomic commits
3. Run `task lint && task test`
4. Push and create PR
5. CI must pass
6. Review and merge
