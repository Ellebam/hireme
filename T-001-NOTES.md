# T-001: Verify Template Switching Persists to Backend

## Task Summary
- **Size:** XS
- **Goal:** Changing template in editor survives save-reload cycle
- **Acceptance:** Change template, reload page, template stays
- **Branch:** `feat/t-001-template-persist`

---

## Investigation Checklist

### Frontend
- [x] `editor-store.ts` — Does `updateTemplateId()` mark CV as dirty? **YES**
- [x] `useAutoSave.ts` — Does save payload include `templateId`? **YES** (sends full `cvContent`)
- [x] `api/client.ts` — Does `updateCV()` send full content with `templateId`? **YES**
- [x] `types/cv.ts` — What are the allowed `TEMPLATE_IDS`? **`['classic', 'modern', 'visionary', 'blank']`**

### Backend
- [x] `PUT /api/v1/cv/{id}` handler — Does it persist JSONB content as-is? **YES** (opaque `json.RawMessage`)
- [x] CV JSON schema — Is `visionary` in allowed `templateId` values? **NO — only `["classic", "modern", "minimal"]`**
- [x] Backend validation — Does it validate `templateId` at all? **YES** — JSON schema validation in service layer

### Schema
- [x] `api/internal/validator/schema/cv-schema.json` — Enum: `["classic", "modern", "minimal"]` — **MISMATCH**

---

## Findings

### Frontend Findings

**All good.** The frontend flow works correctly end-to-end:

1. **`editor-store.ts:507-514`** — `updateTemplateId()`:
   - Sets `cvContent.templateId = templateId`
   - Sets `isDirty = true` (triggers auto-save)
   - Calls `pushHistory()` (enables undo)

2. **`useAutoSave.ts:28`** — Sends `{ content: currentContent }` which includes the full `CVContent` object with `templateId`

3. **`api/client.ts:254-256`** — `update()` sends PUT with full body to `/api/v1/cv/{id}`

4. **`types/cv.ts:12`** — `TEMPLATE_IDS = ['classic', 'modern', 'visionary', 'blank']` — includes `visionary`

5. **`EditorToolbar.tsx:29-33`** — Template selector dropdown offers: Classic, Modern, Visionary (no Blank)

6. **`editor-store.test.ts:486-509`** — 3 existing tests cover:
   - Updates templateId + marks dirty
   - No-op when no CV loaded
   - Pushes to history

### Backend Findings

**Content flows as opaque JSONB through every layer:**

| Layer | File | What Happens |
|-------|------|--------------|
| Handler | `handler/cv.go:99-133` | Receives `Content` as `*json.RawMessage`, passes to service |
| Service | `service/cv_service.go:114-120` | **Validates content against JSON schema**, then stores |
| Validator | `validator/cv_schema.go:44-60` | Unmarshals + validates against embedded schema |
| Repository | `repository/postgres/cv.go:89-105` | Stores as raw JSONB, no transformation |
| SQL | `queries/cvs.sql.go:179-208` | `UPDATE cvs SET content = COALESCE($3, content)` |

**Domain constants** (`domain/cv.go:70-74`):
```go
TemplateClassic = "classic"
TemplateModern  = "modern"
TemplateMinimal = "minimal"   // <-- NO "visionary"
```

**Standalone `ValidateTemplateID()`** (`validator/cv_schema.go:76-84`):
- Only allows classic, modern, minimal
- Not called in the Update flow (schema validation handles it)

**Seed data** (`db/seed/dev_seed.sql:51`): Uses `"modern"` — valid

### Schema Findings

**`api/internal/validator/schema/cv-schema.json:15`:**
```json
"templateId": {
  "type": "string",
  "enum": ["classic", "modern", "minimal"]
}
```

- `"visionary"` is **NOT** in the enum
- `"blank"` is **NOT** in the enum
- `"minimal"` **IS** in the enum but the frontend renamed it to `"visionary"` in Session 10

---

## Gaps Found

### GAP-1: Backend schema rejects `"visionary"` (BLOCKER)
- **File:** `api/internal/validator/schema/cv-schema.json`
- **Problem:** Enum is `["classic", "modern", "minimal"]` — frontend sends `"visionary"` which fails validation
- **Effect:** Saving a CV with visionary template returns **400 Bad Request**, auto-save silently fails
- **Root cause:** Session 10 renamed `minimal` → `visionary` in frontend `types/cv.ts` but never updated the backend schema

### GAP-2: Domain constants out of sync
- **File:** `api/internal/domain/cv.go:70-74`
- **Problem:** Constants define `TemplateMinimal = "minimal"` but frontend uses `"visionary"`
- **Effect:** Any backend code using these constants references a value the frontend no longer sends

### GAP-3: `ValidateTemplateID()` standalone function out of sync
- **File:** `api/internal/validator/cv_schema.go:76-84`
- **Problem:** Only allows classic, modern, minimal
- **Effect:** Low impact (not called in Update flow), but inconsistent

### GAP-4: `"blank"` template not in backend schema
- **File:** `api/internal/validator/schema/cv-schema.json`
- **Problem:** Frontend defines `"blank"` as a valid `TemplateId` but backend schema doesn't include it
- **Effect:** Creating a CV from the blank template would also fail backend validation
- **Note:** Lower priority — blank isn't in the toolbar selector, only used for "new from scratch" flow

---

## Architect Handoff

> **Instructions for @architect:** Use the findings above to produce a design decision (Quick Decision or ADR format per your process). The investigation is complete — no further codebase exploration needed. Focus on the decisions below and produce a concrete change list for @engineer.

### Decisions Needed

**Decision 1: Replace `"minimal"` with `"visionary"`, or support both?**
- Context: Session 10 renamed `minimal` → `visionary` on the frontend. No production data exists yet (only dev seed using `"modern"`).
- Option A: **Replace** `"minimal"` with `"visionary"` everywhere (clean break, no legacy baggage)
- Option B: **Add** `"visionary"` alongside `"minimal"` (backwards-compatible, but keeps a dead value)
- Consideration: There is no migration concern — the seed data uses `"modern"`, not `"minimal"`

**Decision 2: Should `"blank"` be a valid backend templateId?**
- Context: Frontend `TEMPLATE_IDS` includes `"blank"`, used for "start from scratch" template. Not in toolbar selector but part of the type system. Backend schema would reject it.
- Option A: **Add** `"blank"` to backend schema enum (full parity with frontend)
- Option B: **Don't add** — `"blank"` is a frontend-only concept, CVs should always resolve to a render template before saving
- Consideration: If blank is just an initial state that gets replaced when user picks a real template, Option B is cleaner

### Files in Scope

| File | What Needs to Change | Gap |
|------|---------------------|-----|
| `api/internal/validator/schema/cv-schema.json` | Update `templateId` enum | GAP-1 |
| `api/internal/domain/cv.go` | Update/replace `TemplateMinimal` constant | GAP-2 |
| `api/internal/validator/cv_schema.go` | Update `ValidateTemplateID()` | GAP-3 |
| `web/src/types/cv.ts` | Possibly remove `"blank"` if Decision 2 = Option B | GAP-4 |

### Tests to Add/Update

| Test | Purpose |
|------|---------|
| Backend validator test | Verify `"visionary"` passes JSON schema validation |
| Backend service test | Round-trip: save CV with `"visionary"` → retrieve → verify persisted |
| Backend `ValidateTemplateID()` test | Unit test for standalone function with new values |
| Frontend round-trip test (optional) | Mock save → mock load → verify templateId survives |

### Constraints
- PR target: < 300 lines changed
- Task size: XS (< 1 hour implementation after decisions are made)
- Branch: `feat/t-001-template-persist`
- Must not break existing tests (101 frontend, Go backend suite)

---

## Test Plan

- [x] Unit test: change template → verify dirty state — **ALREADY EXISTS** (editor-store.test.ts:486-509)
- [ ] Backend schema test: `"visionary"` is accepted by JSON schema validation
- [ ] Backend service test: save with `"visionary"` → retrieve → verify templateId persists
- [ ] Backend `ValidateTemplateID()` unit test
- [ ] Integration test (optional): full round-trip through API with curl/httptest
