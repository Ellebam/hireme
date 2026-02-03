# Work Log

Chronological record of development activity.

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
