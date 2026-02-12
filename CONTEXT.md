# HireMe — Project Context

## Current State

**Status:** Task-based workflow active — T-001, T-006, T-007 done, 12 tasks remaining, 9 unblocked

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
- ✅ Auto-save with 2s debounce + manual save button
- ✅ Export modal (PDF, DOCX, JSON)
- ✅ Keyboard shortcuts (Ctrl+Z/Y, Ctrl+/-/0)
- ✅ CV templates system (starter + blank templates)
- ✅ Logging system with configurable levels
- ✅ Save error feedback — tooltip on error icon, click-to-retry
- ✅ Dashboard is root page (`/`), `/dashboard` redirects
- ✅ Edit button links to `/editor` (not `/editor/${cv.id}`)
- ✅ Burger menu (MoreVertical) has DropdownMenu with Edit/Delete
- ✅ Delete confirmation dialog on CV card
- ✅ Back-to-dashboard link in editor toolbar
- ✅ shadcn/ui dropdown-menu component
- ✅ 127 passing tests (Vitest)
- ✅ Build verified, full stack integration tested

### Remaining Work
- ❌ Export (`POST /api/v1/export/{format}`) — returns 501, needs Gotenberg integration → T-002, T-003, T-004, T-005
- ❌ R2 cloud storage — placeholder only, needs AWS SDK implementation → T-013

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.22+, Chi router, sqlc |
| Frontend | Next.js 15, React 19, Tailwind, shadcn/ui, Zustand, dnd-kit |
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
│   ├── EditorLayout.tsx   # Main layout (responsive breakpoints)
│   ├── CVPreview.tsx      # Live preview (template dispatch)
│   ├── SectionPalette.tsx # Add sections
│   ├── PropertiesPanel.tsx# Edit section content
│   └── editors/           # Section-specific editors
├── components/templates/
│   ├── ClassicTemplate.tsx# Traditional CV layout
│   ├── ModernTemplate.tsx # Contemporary with timeline
│   ├── VisionaryTemplate.tsx # Two-column with sidebar
│   └── index.ts           # Template exports
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
├── WORKLOG.md              # Trimmed — historical logs in git history
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

## Task Board

### Active
| ID | Task | Branch | Blocked By | Status |
|----|------|--------|------------|--------|
| — | — | — | — | — |

### Backlog
| ID | Task | Blocked By | Size |
|----|------|------------|------|
| T-002 | Create HTML generation for CV export | — | S |
| T-003 | Implement PDF export via Gotenberg | T-002 | S |
| T-004 | Implement DOCX export via Gotenberg | T-002 | S |
| T-005 | Wire frontend export modal to real API | T-003, T-004 | S |
| T-008 | P1 interaction tests — editors | — | S |
| T-009 | P1 interaction tests — export error paths | — | XS |
| T-010 | P2 UX regression tests | — | S |
| T-011 | Add Certifications section | — | S |
| T-012 | Add Projects section | — | S |
| T-013 | R2 cloud storage implementation | — | M |
| T-014 | OAuth authentication (Google OIDC) | — | M |
| T-015 | Multiple CV support | — | M |

### Done
| ID | Task | PR |
|----|------|----|
| T-001 | Verify template switching persists to backend | feat/t-001-template-persist |
| T-006 | P0 smoke tests — section editors | test/t-006-007-smoke-tests |
| T-007 | P0 smoke tests — page components | test/t-006-007-smoke-tests |

---

## Task Details

**T-002: Create HTML generation for CV export** — S
- Go service that takes CV content + templateId and produces a self-contained HTML string
- One function per template (classic, modern, visionary) with inline CSS
- Shared foundation for both PDF and DOCX export
- Tests: unit tests generating HTML from sample CV data for each template
- Acceptance: Unit tests pass, HTML renders correctly in browser

**T-003: Implement PDF export via Gotenberg** — S
- Wire `POST /api/v1/export/pdf` to call Gotenberg's chromium HTML-to-PDF endpoint
- Use HTML from T-002, POST to `http://gotenberg:3000/forms/chromium/convert/html`
- Return PDF binary with correct Content-Type
- Tests: handler test (mock Gotenberg), integration test (real Gotenberg)
- Acceptance: `curl` export endpoint returns valid PDF

**T-004: Implement DOCX export via Gotenberg** — S
- Wire `POST /api/v1/export/docx` using Gotenberg's LibreOffice endpoint
- Same HTML input as PDF, different Gotenberg conversion
- Tests: handler test (mock Gotenberg), integration test (real Gotenberg)
- Acceptance: `curl` export endpoint returns valid DOCX

**T-005: Wire frontend export modal to real API** — S
- Connect existing ExportModal buttons to actual API calls
- Download file to browser (PDF, DOCX); JSON export already works
- Show loading state and error handling
- Tests: component test for export flow, API client test for export methods
- Acceptance: Click "Export PDF" in UI, browser downloads PDF

**T-008: P1 interaction tests — editors** — S
- ExperienceEditor + SkillsEditor interaction tests
- Add/edit/delete entries via modal forms using `@testing-library/user-event`
- Acceptance: Tests cover add, edit, delete flows

**T-009: P1 interaction tests — export error paths** — XS
- Verify abort timeout and network errors map to `ApiError`
- Acceptance: Error path tests passing

**T-010: P2 UX regression tests** — S
- Keyboard shortcuts: verify Ctrl+Z/Y trigger undo/redo
- DnD section reorder: verify store updates
- Acceptance: Tests for shortcuts and DnD passing

**T-011: Add Certifications section** — S
- `CertificationsEditor.tsx` (modal form for cert entries)
- Add to `PropertiesPanel` switch + all 3 templates
- Schema already defines `certificationsContent`
- Tests: editor smoke test, interaction test, template render test
- Acceptance: Can add/edit/delete certifications, visible in preview

**T-012: Add Projects section** — S
- `ProjectsEditor.tsx` (modal form for project entries)
- Add to `PropertiesPanel` switch + all 3 templates
- Schema already defines `projectsContent`
- Tests: editor smoke test, interaction test, template render test
- Acceptance: Can add/edit/delete projects, visible in preview

**T-013: R2 cloud storage implementation** — M
- Implement `R2Storage` methods using AWS SDK
- Tests: storage interface tests with mock S3 client

**T-014: OAuth authentication (Google OIDC)** — M (needs splitting when picked up)
- Full auth flow: login, callback, JWT, middleware
- Tests: middleware tests, auth flow tests

**T-015: Multiple CV support** — M (needs splitting when picked up)
- Dashboard shows list, editor route `/editor/[id]`, create new CV
- Tests: routing tests, dashboard list test

### Dependency Graph
```
T-002 → T-003 → T-005
     → T-004 → T-005

T-008 through T-012: ALL INDEPENDENT (any order)
T-013, T-014, T-015: Future backlog, independent
```

9 out of 12 remaining tasks have zero dependencies and can be picked up in any order.

---

## Task Workflow

Every task follows this lifecycle:

```
Backlog → Active (in-progress) → QA Review → Tech Lead Approval → Done
```

### Steps
1. **Investigate** — `/investigate <task>` — PM traces code, documents findings, prepares handoff in notes file
2. **Plan** — `/plan <task>` — Architect researches, makes decisions, writes implementation plan in notes file
3. **QA Plan** — `/qa-plan <task>` — QA audits coverage gaps, appends test recommendations to notes file
4. **Implement** — `@engineer` implements code + tests per the plan
5. **QA Review** — Tech lead + CodeRabbit review the PR
6. **Close** — `/close-task <task>` — PM moves task to Done, cleans up notes

### Task Size Rules
- **XS** (< 1 hour): Config changes, verifications, tiny fixes
- **S** (1-3 hours): Single feature, single editor, test suite
- **M** (3-6 hours): Feature with backend+frontend; must be split if larger

### Branch & PR Convention
- Branch: `feat/t-003-pdf-export`, `fix/t-001-template-persist`, `test/t-006-editor-smoke`
- PR title: `[T-003] Implement PDF export via Gotenberg`
- Target: < 300 lines changed per PR

---

## Stakeholder Feedback

New feedback: say `stakeholder feedback: [description]` → PM agent creates/updates tasks in the Task Board.

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
