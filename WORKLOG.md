# Work Log

Chronological record of development activity.

---

## 2026-02-04 (Session 3)

### Session Focus
Frontend MVP architecture planning

### Completed
- [DOCS] Created comprehensive `web/FRONTEND_MVP_PLAN.md` with full architecture
- [DOCS] Designed three-column editor layout (palette | preview | properties)
- [DOCS] Defined component hierarchy and state management strategy
- [DOCS] Planned drag & drop implementation with dnd-kit
- [DOCS] Established testing strategy with Vitest
- [DOCS] Created 7-phase implementation roadmap

### Notes
**Key architectural decisions:**
- Three-column layout: Section palette (left), CV preview (center), Properties panel (right)
- Zustand for state: `EditorStore` (CV data, undo/redo) + `UIStore` (panels, modals)
- Auto-save with 2s debounce + optimistic updates
- dnd-kit for section and entry reordering
- Responsive: Desktop (3-col) → Tablet (2-col) → Mobile (single + FAB)

**MVP scope:**
- 6 section types: Personal, Summary, Experience, Education, Skills, Languages
- Single CV editing with live preview
- Drag & drop reordering
- Export to PDF/DOCX/JSON
- Basic i18n (en/de)

### Next
- Phase 1: TypeScript types from CV schema
- Phase 1: API client with typed methods
- Phase 1: Zustand stores (editor + UI)
- Phase 1: Vitest setup and store tests

---

## 2026-02-04 (Session 2)

### Session Focus
Testing strategy and planning for MVP API

### Completed
- [DOCS] Created comprehensive `TESTING_PLAN.md` with phased approach
- [DOCS] Updated `CONTEXT.md` with testing to-dos referencing plan
- [CHORE] Assessed current test coverage (zero tests exist)
- [CHORE] Identified CI/CD improvements needed

### Notes
**Current state:** Zero test files in `api/` directory. CI passes vacuously.

**Testing phases:**
| Phase | Focus | Priority |
|-------|-------|----------|
| 1 | Domain, Validator, HTTP Utils | P0 |
| 2 | Service layer (mocked repos) | P0 |
| 3 | Handler layer (httptest) | P1 |
| 4 | Auth middleware | P1 |
| 5 | Repository integration | P2 |
| 6 | Storage tests | P2 |

**CI/CD gaps:**
- Missing `go vet` in lint job
- No migration step before integration tests
- No coverage threshold enforcement

### Next
- Start Phase 1: Domain layer tests (`domain/*_test.go`)
- Create test infrastructure (mocks, fixtures)
- Implement service layer tests

---

## 2026-02-04

### Session Focus
Implement asset management (upload, retrieve, delete) with storage abstraction

### Completed
- [FEAT] Implemented `AssetService.Upload()` with checksum, deduplication, image dimensions
- [FEAT] Implemented `AssetService.GetFileContent()` for file retrieval
- [FEAT] Completed `AssetService.Delete()` with storage cleanup
- [FEAT] Created storage abstraction with `NewStorage()` factory function
- [FEAT] Added `R2Storage` placeholder with TODOs for future cloud storage
- [FIX] Fixed `GetAsset` handler to return JSON metadata by default
- [CHORE] Updated Go module path to `github.com/ellebam/hireme/api`
- [DOCS] Updated CONTEXT.md with asset endpoints

### Notes
**Storage architecture:**
- `Storage` interface: `Put`, `Get`, `Delete`, `Exists`
- `LocalStorage`: Fully implemented, stores in `./data/uploads/{userID}/{YYYY-MM}/`
- `R2Storage`: Placeholder with error returns, ready for AWS SDK integration

**Asset features:**
- SHA-256 checksum for deduplication (same user + same file = returns existing)
- Image dimension extraction via Go's `image` package
- User storage tracking (increments/decrements `storage_used_bytes`)

### Next
- Implement export endpoints (Gotenberg integration)
- Implement R2Storage when ready for production
- Start frontend

---

## 2026-02-03

### Session Focus
Wire up repository layer to make API functional

### Completed
- [FIX] `api/db/sqlc.yaml` — fixed paths and type override syntax for sqlc v1.30
- [FEAT] Implemented `UserRepository` with sqlc queries
- [FEAT] Implemented `CVRepository` with sqlc queries
- [FEAT] Implemented `AssetRepository` with sqlc queries
- [FIX] `Taskfile.yml` — db:seed now runs full seed SQL (user + CV)
- [DOCS] Updated `CLAUDE.md` with quick commands
- [DOCS] Updated `CONTEXT.md` with current state

### Notes
**Root cause of 404s:** Repositories were stubs returning `ErrNotFound`. Auth bypass worked but `GetByID()` always failed.

**Key mappings:**
- `pgx.ErrNoRows` → `domain.ErrNotFound`
- `pgtype.Timestamptz` → `time.Time` (check `.Valid`)
- Nullable pointers → domain fields with nil checks

### Next
- Implement asset upload (file storage)
- Implement export endpoints (Gotenberg)
- Start frontend

---

## 2026-01-29

### Session Focus
Initial Claude Code setup

### Completed
- [CHORE] Customized agent files in `.claude/agents/`
- [DOCS] Updated all agents with HireMe-specific context

### Notes
Project discovery — well-structured Go backend, sqlc for queries, Next.js 14 frontend planned.

---

<!-- TEMPLATE:
## YYYY-MM-DD

### Session Focus
[Goal]

### Completed
- [TAG] Item

### Notes
[Observations]

### Next
[Continue with...]
-->
