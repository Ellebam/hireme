# HireMe — Project Context

## Current State

**Status:** MVP functional — API backend + Frontend editor working end-to-end

### What's Working

**Backend:**
- ✅ Infrastructure: PostgreSQL + Gotenberg via Docker
- ✅ Database: Migrations applied, dev data seeded
- ✅ API endpoints: Health, Users, CVs (full CRUD), Assets
- ✅ Auth bypass: Dev mode uses `dev-user-001`
- ✅ Repository layer: sqlc queries wired to domain types
- ✅ Asset management: Upload, retrieve, delete with local storage
- ✅ Storage abstraction: LocalStorage implemented, R2 placeholder ready
- ✅ API response format: `{ data: ... }` wrapper with proper error handling
- ✅ Go tests: Domain, service, handler, middleware, repository, storage layers tested
- ✅ CI/CD: GitHub Actions with lint, test, build for both API and web

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
- ✅ CV templates system (starter + blank templates)
- ✅ Logging system with configurable levels
- ✅ 81 passing tests (Vitest)
- ✅ Build verified, full stack integration tested

### What's Not Working Yet
- ❌ Export (`POST /api/v1/export/{format}`) — returns 501, needs Gotenberg integration
- ❌ R2 cloud storage — placeholder only, needs AWS SDK implementation

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.22+, Chi router, sqlc |
| Frontend | Next.js 14, React 18, Tailwind, shadcn/ui, Zustand, dnd-kit |
| Database | PostgreSQL 16 (JSONB for CV content) |
| Export | Gotenberg (HTML → PDF/DOCX) |
| Storage | Local filesystem (dev), Cloudflare R2 (prod) |
| Testing | Vitest (frontend), Go testing (backend - planned) |

---

## Frontend Architecture

### Editor Layout
```
+------------------------------------------------------------------+
|                         Top Toolbar                               |
|  [Undo] [Redo] | Zoom [50-200%] | [Export v]                     |
+------------------------------------------------------------------+
|           |                              |                        |
|  Section  |      CV Preview              |    Properties          |
|  Palette  |      (Live A4 Render)        |    Panel               |
|  (240px)  |      (Flex)                  |    (320px)             |
|           |                              |                        |
+------------------------------------------------------------------+
```

### State Management (Zustand)
- **EditorStore**: CV data, sections, undo/redo history, save status
- **UIStore**: Sidebar toggles, modals, preview scale

### Key Frontend Files
```
web/src/
├── stores/
│   ├── editor-store.ts    # CV content + history
│   └── ui-store.ts        # UI state
├── lib/
│   ├── api/client.ts      # Typed API client (unwraps responses)
│   ├── logger.ts          # Configurable logging
│   └── templates/         # CV templates (starter, blank)
├── components/editor/
│   ├── EditorLayout.tsx   # Main layout
│   ├── CVPreview.tsx      # Live preview
│   ├── SectionPalette.tsx # Add sections
│   ├── PropertiesPanel.tsx# Edit section content
│   └── editors/           # Section-specific editors
└── hooks/
    ├── useAutoSave.ts     # 2s debounced save
    └── useKeyboardShortcuts.ts
```

### Frontend Debugging
```js
// Enable debug logs in browser console:
localStorage.setItem('LOG_LEVEL', 'debug');
location.reload();
```

---

## Backend Architecture

```
HTTP Request → Handler → Service → Repository → PostgreSQL
                 ↓
              Middleware (Auth, Logging, CORS)
```

- **Handlers**: Parse HTTP, validate input, call services
- **Services**: Business logic, orchestration
- **Repositories**: Data access via sqlc-generated code
- **Domain**: Pure types, no external dependencies

### API Response Format
All responses wrapped in standard format:
```json
// Success
{ "data": { ... } }

// Error
{ "error": { "code": "...", "message": "..." } }
```

### Key Backend Files
```
api/
├── cmd/server/main.go     # Entry point, router setup
├── internal/
│   ├── handler/           # HTTP handlers
│   ├── service/           # Business logic
│   ├── repository/        # Data access
│   ├── domain/            # Types + errors
│   └── middleware/        # Auth, CORS
└── db/
    ├── migrations/        # SQL migrations
    ├── queries/           # sqlc queries
    └── seed/              # Dev seed data
```

---

## Project Structure

```
hireme/
├── api/                    # Go backend
├── web/                    # Next.js frontend
├── docker/                 # Docker configs
├── schemas/                # Shared JSON schemas
├── scripts/                # Dev utilities
├── CONTEXT.md              # This file
├── WORKLOG.md              # Session history
└── CLAUDE.md               # AI assistant instructions
```

---

## API Endpoints

```
GET  /health                    → {"data": {"status": "healthy"}}
GET  /ready                     → {"data": {"status": "ready"}}
GET  /api/v1/users/me           → Current user profile
PATCH /api/v1/users/me          → Update user profile
GET  /api/v1/cv                 → User's active CV
POST /api/v1/cv                 → Create new CV
PUT  /api/v1/cv/{id}            → Update CV
DELETE /api/v1/cv/{id}          → Delete CV
POST /api/v1/assets             → Upload asset (image)
GET  /api/v1/assets/{id}        → Get asset metadata
DELETE /api/v1/assets/{id}      → Delete asset
```

---

## Common Commands

```bash
# Full Development (recommended)
task dev               # Start infra + API + web together
task dev:stop          # Stop everything + kill dev ports
task dev:restart       # Full stop and restart

# Partial Development
task dev:api           # Start infra + API only
task dev:web           # Start infra + web only

# Infrastructure
task infra:up          # Start PostgreSQL + Gotenberg
task infra:down        # Stop infrastructure

# Database
task db:migrate        # Run migrations
task db:seed           # Seed dev user + sample CV
task db:psql           # Open psql shell

# Testing
npm test               # Frontend tests (in web/)
task api:test          # Backend tests (when implemented)

# Code Generation
task api:sqlc          # Generate sqlc code
```

---

## To-Do

### High Priority
- [ ] Gotenberg integration for PDF/DOCX export
- [ ] Production security hardening (see Security Recommendations)
- [ ] Branch protection rules for main branch

### Medium Priority
- [ ] Mobile responsive improvements
- [ ] Additional section types (certifications, projects)
- [x] ~~Upgrade Next.js to patch security vulnerabilities~~ (pinned to 14.2.35)
- [x] ~~CI/CD: Add coverage thresholds~~ (configured in vitest.config.ts)

### Low Priority
- [ ] R2 cloud storage implementation
- [ ] OAuth authentication (Google OIDC)
- [ ] Multiple CV support

### Frontend Testing (add alongside refactors/bugfixes)

Current state: 81 tests cover stores, API client, and templates. Zero component or interaction tests. The items below are ordered by impact — each one catches a different class of bug.

**P0 — Catch render/type breakage:**
- [ ] Smoke tests for all 6 section editors (Personal, Summary, Experience, Education, Skills, Languages) — render with mock data, verify key elements appear. Catches broken imports, missing props, and type mismatches after refactors.
- [ ] Smoke tests for page-level components (Dashboard, Editor) — render with mocked stores, verify basic structure loads.

**P1 — Catch interaction bugs:**
- [ ] `useAutoSave` hook test — verify debounce timing, dirty-state triggers, and cleanup on unmount. This hook has real async logic that can silently break.
- [ ] Editor interaction tests — add/edit/delete entries in ExperienceEditor and SkillsEditor (these have modal forms and the most complex user flows). Use `@testing-library/user-event`.
- [ ] Export API error path tests — verify abort timeout and network error map to `ApiError` (new code, currently 0% covered).

**P2 — Catch UX regressions:**
- [ ] Keyboard shortcut tests — verify Ctrl+Z/Y trigger undo/redo on the store.
- [ ] Drag-and-drop section reorder test — verify `SectionPalette` reorder callback updates store order.

**Testing notes:**
- Use `@testing-library/react` + `@testing-library/user-event` (already installed)
- Mock `useEditorStore`/`useUIStore` via Zustand's test pattern (see `src/test/setup.ts`)
- Component tests go next to their source: `ComponentName.test.tsx`
- Target: reach 35% statement coverage with P0+P1 items, which makes the threshold meaningful

---

## Key Patterns

### Auth Bypass (Development)
Set `AUTH_BYPASS_ENABLED=true` in `.env` — all requests use `dev-user-001`

### sqlc Workflow
1. Write SQL in `api/db/queries/*.sql`
2. Run `task api:sqlc`
3. Use generated `queries.Queries` in repositories

### Repository Error Mapping
```go
// pgx.ErrNoRows → domain.ErrNotFound
user, err := r.q.GetUserByID(ctx, id)
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domain.ErrNotFound
}
```

---

## Development Setup

```bash
task setup             # Full setup (infra, migrate, seed, deps)
task dev               # Start everything (API :8080, Web :3000)
```

## Testing

```bash
# Frontend
cd web && npm test

# API Health check
curl http://localhost:8080/health

# Get CV (auth bypass must be enabled)
curl http://localhost:8080/api/v1/cv
```
