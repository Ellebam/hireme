# T-016: Fix save bug after template switch

## Summary
- **Goal:** Identify and fix the root cause of saves breaking after template switching
- **Acceptance:** Saves succeed reliably after switching templates; no schema validation mismatches between frontend and backend
- **Branch:** `fix/t-016-save-after-template-switch`

---

## Investigation Checklist
- [x] Trace `updateTemplateId` in editor store
- [x] Trace `useAutoSave` hook trigger logic
- [x] Trace `saveImmediately` → API client → backend handler
- [x] Trace backend validation chain (handler → service → validator → schema)
- [x] Compare frontend `CVContent` type with backend JSON schema
- [x] Check for `additionalProperties: false` violations
- [x] Check template ID enum mismatch (frontend vs backend)
- [x] Check `pushHistory` interaction with auto-save
- [x] Check seed data vs starterTemplate content structure
- [x] Review existing validator tests

## Findings

### 1. Auto-save mechanism: CORRECT
- `updateTemplateId` (editor-store.ts:507-514): Sets `cvContent.templateId`, `isDirty = true`, calls `pushHistory()`
- `useAutoSave` effect (useAutoSave.ts:55-59): Watches `isDirty`, `cv`, `cvContent` — triggers `debouncedSave()` on change
- `saveImmediately` (useAutoSave.ts:15-36): Uses `useEditorStore.getState()` for fresh state, sends `{ content: currentContent }` to API
- Dependencies are stable (Zustand actions), no stale closure issues
- **Verdict:** The auto-save correctly fires after template switch

### 2. Template ID mismatch: BUG FOUND
- **Frontend** (`types/cv.ts:12`): `TEMPLATE_IDS = ['classic', 'modern', 'visionary', 'blank']` — 4 values
- **Backend schema** (`cv-schema.json:15`): `enum: ["classic", "modern", "visionary"]` — 3 values
- **Backend test** (`cv_schema_test.go:176`): Explicitly tests `"blank"` as INVALID
- **`blankTemplate()`** (`templates/blank.ts:14`): Creates content with `templateId: 'blank'`
- **Toolbar** (`EditorToolbar.tsx:29-33`): Only shows Classic/Modern/Visionary (not Blank)
- **Impact:** Creating a CV from `blankTemplate()` would fail backend validation. Currently latent because `EditorPage` uses `starterTemplate()`, but `'blank'` exists in the type system and registry.

### 3. Language entry `id` field: BUG FOUND (PRIMARY ROOT CAUSE)
- **Backend schema** (`cv-schema.json:309-321`): `languagesContent` items have properties: `language`, `proficiency`, `certification`. **No `id` property.** + `additionalProperties: false`
- **Frontend type** (`types/cv.ts:157-162`): `LanguageEntry` has `id: string` (required)
- **LanguagesEditor** (`LanguagesEditor.tsx:30-31`): Generates `id: generateId()` for new entries
- **starterTemplate** (`templates/starter.ts:168-169`): Language entries include `id: generateId()`

**This is the primary root cause.** When the frontend sends content with language entries that have `id` fields, the backend rejects it because `additionalProperties: false` disallows the `id` property.

**Why this manifests as "saves break after template switch":**
1. Seed CV loads from DB (language entries have NO `id` — inserted directly via SQL)
2. User makes various edits. If they add a new language entry, it gets an `id` field
3. User switches template → auto-save fires → sends ALL content to backend
4. Backend validation rejects the content because language entry has `id` field
5. Save fails → user sees error

**OR for new CVs:**
1. `EditorPage` creates CV using `starterTemplate()` which includes language entries with `id`
2. `api.cv.create()` sends content to backend → backend validates → REJECTS
3. CV creation fails → user sees error immediately

### 4. Seed data structure
- `dev_seed.sql` language entries: `{"language": "German", "proficiency": "native"}` — **NO `id` field**
- This is correct per schema but diverges from the TypeScript type
- Seed data bypasses API validation (direct SQL INSERT)

### 5. Save status race condition (LOW)
- `markSaved()` sets `isDirty = false` after API call completes
- If user edits between `getState()` read and `markSaved()`, edit appears saved but wasn't
- Next edit triggers another save, so data isn't permanently lost
- Not specific to template switch — general auto-save issue

## Gaps Found

### BLOCKER
1. **Language entry `id` not in schema** — Backend schema doesn't allow `id` on language entries, but frontend always generates it. Saves fail for any CV with new language entries.
2. **`blank` template ID not in schema** — Frontend type allows `'blank'`, backend rejects it. `blankTemplate()` creates invalid content.

### Medium
3. **Schema and TypeScript type drift** — No automated check that `CVContent` type matches `cv-schema.json`. Schema is source of truth but frontend can silently diverge.
4. **No test for "save after template switch"** — Zero test coverage for this flow in either frontend or backend tests.

### Low
5. **Auto-save race condition** — `markSaved()` can clear `isDirty` for edits made during an in-flight save. Data not lost (next save catches it), but status indicator can be briefly misleading.

## Architect Handoff

### Decisions Needed

**Quick Decision 1: Fix language entry `id` in schema or frontend?**
- **Option A (Recommended):** Add `id` property to language entry schema (matches all other entry types: experience, education, certifications, projects all have `id`)
- **Option B:** Remove `id` from `LanguageEntry` TypeScript type and editor — breaks consistency with other entry types
- **Recommendation:** Option A. Every other entry type has `id`. Languages should too. It was clearly an oversight.

**Quick Decision 2: What to do with `blank` template?**
- **Option A (Recommended):** Remove `'blank'` from `TEMPLATE_IDS`, make `blankTemplate()` use `templateId: 'classic'` (blank refers to content, not visual style)
- **Option B:** Add `'blank'` to backend schema enum (but there's no corresponding visual template component)
- **Recommendation:** Option A. The `blank` template is a *content* template (empty sections), not a *visual* template. It should render with a valid visual template ID. `'classic'` is the sensible default.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `schemas/cv-schema.json` | Add `id` property to languagesContent entry items | #1 |
| `api/internal/validator/schema/cv-schema.json` | Same (embedded copy, must match) | #1 |
| `api/internal/validator/cv_schema_test.go` | Add test: language entries with `id` pass validation | #1 |
| `web/src/types/cv.ts` | Remove `'blank'` from `TEMPLATE_IDS` | #2 |
| `web/src/lib/templates/blank.ts` | Change `templateId: 'blank'` → `templateId: 'classic'` | #2 |
| `web/src/stores/__tests__/editor-store.test.ts` | Add test: save payload after template switch is schema-valid | #4 |

### Tests to Add/Update
- **Backend:** Validator test with language entry containing `id` (should pass after fix)
- **Backend:** Validator test with full starterTemplate content (regression)
- **Frontend:** Editor store test: `updateTemplateId` produces schema-valid content
- **Frontend:** Auto-save test: save fires after template switch with correct payload

### Constraints
- Both schema copies (`schemas/` and `api/internal/validator/schema/`) MUST stay in sync
- < 300 lines changed
- No migration needed — DB stores raw JSONB, schema only validates on write

### Recommended Next Agent
**@architect** for quick decision confirmation, then **@engineer** for implementation. The fix is small (XS) — add `id` to schema, fix `blank` template, add tests.

---

## Architect Plan

### Decisions

**Decision 1: Fix language entry `id` — schema or frontend?**
- Option A: Add `id` to language entry schema — Pros: matches all 5 other entry types (experience, education, certifications, projects, awards all have `id`); no frontend changes; consistent pattern. Cons: none meaningful.
- Option B: Remove `id` from frontend `LanguageEntry` type — Pros: no schema change. Cons: breaks consistency; requires removing `id` from `LanguagesEditor`, `starterTemplate`, and test mocks; makes language entries the only ones without stable IDs (breaks React `key` patterns).
- **Choice:** A — **Rationale:** This was a simple oversight. Every other entry type has `id: { type: "string", format: "uuid" }` + required. Languages should match.

**Decision 2: What to do with `blank` template ID?**
- Option A: Make `blankTemplate()` use `templateId: 'classic'` and remove `'blank'` from `TEMPLATE_IDS` — Pros: `blank` is a *content* template (empty sections), not a *visual* template; no backend schema change needed; no dead-code visual template component. Cons: none.
- Option B: Add `'blank'` to backend schema and create a visual template component — Pros: consistent naming. Cons: `blank` doesn't describe a visual style; requires a new React component that would be identical to `classic`; adds maintenance burden.
- **Choice:** A — **Rationale:** The template registry concept (`starter`/`blank`) is about *content presets*, while `templateId` in the schema is about *visual rendering*. Conflating them was the mistake. `blankTemplate()` should create empty sections rendered with the `classic` visual template.

### Code Adjustments

| # | File | Change | Why |
|---|------|--------|-----|
| 1 | `schemas/cv-schema.json` (lines 310-320) | Add `"id": { "type": "string", "format": "uuid" }` to language entry properties; add `"id"` to `required` array | Fix BLOCKER #1 — match all other entry types |
| 2 | `api/internal/validator/schema/cv-schema.json` (lines 310-320) | Exact same change as #1 | Embedded copy must stay in sync |
| 3 | `api/internal/validator/cv_schema_test.go` | Add test: language section with populated entries (including `id`) passes validation | Regression test for the fix |
| 4 | `web/src/types/cv.ts` (line 12) | Change `['classic', 'modern', 'visionary', 'blank']` → `['classic', 'modern', 'visionary']` | Fix BLOCKER #2 — remove invalid template ID from type system |
| 5 | `web/src/lib/templates/blank.ts` (line 14) | Change `templateId: 'blank'` → `templateId: 'classic'` | Fix BLOCKER #2 — blank content template uses valid visual template |
| 6 | `api/db/seed/dev_seed.sql` (lines 176-177) | Add `"id"` UUID fields to language entries | Keep seed data consistent with updated schema |

### Files NOT Changing
- `web/src/stores/editor-store.ts` — `updateTemplateId` works correctly
- `web/src/hooks/useAutoSave.ts` — auto-save mechanism works correctly
- `web/src/components/editor/EditorToolbar.tsx` — `TEMPLATE_OPTIONS` already only shows valid 3 templates
- `web/src/components/editor/editors/LanguagesEditor.tsx` — already generates `id` correctly
- `web/src/lib/templates/starter.ts` — language entries already have `id`, now schema accepts them
- `web/src/lib/templates/registry.ts` — `blank` template name/description are fine (it's the content `templateId` that was wrong)
- `api/internal/handler/cv.go` — no handler changes needed
- `api/internal/service/cv_service.go` — no service changes needed

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Schema copies drift out of sync | Low | High (silent validation mismatch) | Step 1 of execution: edit both files, diff to confirm identical |
| Existing DB data has language entries without `id` | Medium | Low | Schema only validates on API write; existing data untouched. Seed data updated for new installs. |
| Removing `'blank'` from TypeScript type causes compile error | Low | Low | Only used in `blankTemplate()` which is updated in same PR; `EditorToolbar` already doesn't reference `'blank'` |

### Engineer Execution Steps

1. **Edit both schema files** — Add `id` to `languagesContent` entry items in:
   - `schemas/cv-schema.json` (lines 311-319)
   - `api/internal/validator/schema/cv-schema.json` (lines 311-319)

   Add after line 311 (`"properties": {`):
   ```json
   "id": { "type": "string", "format": "uuid" },
   ```
   Change `"required"` (line 319) from:
   ```json
   "required": ["language", "proficiency"],
   ```
   to:
   ```json
   "required": ["id", "language", "proficiency"],
   ```

2. **Diff the two schema files** — Verify they are byte-identical:
   ```bash
   diff schemas/cv-schema.json api/internal/validator/schema/cv-schema.json
   ```

3. **Add backend validator test** — In `api/internal/validator/cv_schema_test.go`, add a test `TestCVValidator_Validate_PopulatedLanguagesSection` with a language section containing entries with `id`, `language`, and `proficiency` fields. Follow the pattern from `TestCVValidator_Validate_PopulatedCertificationsSection`.

4. **Run backend tests** — Verify the new test passes and existing tests still pass:
   ```bash
   cd api && go test ./internal/validator/...
   ```

5. **Fix frontend `TEMPLATE_IDS`** — In `web/src/types/cv.ts` line 12, remove `'blank'`:
   ```typescript
   export const TEMPLATE_IDS = ['classic', 'modern', 'visionary'] as const;
   ```

6. **Fix blank template** — In `web/src/lib/templates/blank.ts` line 14, change:
   ```typescript
   templateId: 'classic',
   ```

7. **Update seed data** — In `api/db/seed/dev_seed.sql`, add `id` fields to language entries:
   ```json
   {"id": "lang-001", "language": "German", "proficiency": "native"},
   {"id": "lang-002", "language": "English", "proficiency": "fluent"}
   ```
   Note: seed bypasses validation (direct INSERT), but IDs should be present for consistency. These don't need to be UUIDs since the seed data uses short IDs elsewhere (e.g., `"exp-001"`).

8. **Run frontend tests** — Verify nothing breaks:
   ```bash
   cd web && npm test
   ```

9. **Run frontend type check** — Verify no TypeScript errors:
   ```bash
   cd web && npx tsc --noEmit
   ```

### Verification Gates
1. `cd api && go test ./internal/validator/...` — all pass (including new language test)
2. `cd web && npm test` — all 167+ tests pass
3. `cd web && npx tsc --noEmit` — no type errors
4. `diff schemas/cv-schema.json api/internal/validator/schema/cv-schema.json` — no diff
5. `cd web && npx next build` — build succeeds

---

## Test Plan

### Existing Coverage Reviewed

| Test File | Tests | Lines | Relevant Coverage |
|-----------|-------|-------|-------------------|
| `api/internal/validator/cv_schema_test.go` | 15 | 403 | Valid/invalid template IDs, missing required fields, empty/populated certifications sections, locales. **No language section tests.** |
| `web/src/lib/templates/templates.test.ts` | 13 | 133 | Starter/blank template structure, section types, registry. **Line 69 asserts `templateId === 'blank'` — will break after fix.** |
| `web/src/stores/editor-store.test.ts` | 43 | 531 | setCV, sections CRUD, undo/redo, save state, updateTemplateId (3 tests: updates ID, no-op without CV, pushes history). **No test that content structure is preserved after template switch.** |
| `web/src/hooks/useAutoSave.test.ts` | 10 | 286 | saveNow registration, immediate save, concurrent prevention, error handling, debounce. **No test for payload structure after template switch.** |
| `web/src/test/mocks/cv.ts` | — | 337 | `mockLanguagesContent` entries include `id` field (matches frontend type, currently rejected by backend schema). After fix, mock data will be valid against both. |

### QA-Recommended Additional Tests

| # | Test | Layer | Priority | Description |
|---|------|-------|----------|-------------|
| 1 | Populated language section with `id` passes validation | Backend validator | High | In `cv_schema_test.go`, add `TestCVValidator_Validate_PopulatedLanguagesSection` following the pattern of `TestCVValidator_Validate_PopulatedCertificationsSection` (line 311). Include a language entry with `id` (UUID format), `language`, `proficiency`, and optional `certification`. Assert validation passes. This is the core regression test for the fix. |
| 2 | Language entry without `id` fails validation | Backend validator | High | In `cv_schema_test.go`, add a test that submits a language section where entries have `language` + `proficiency` but NO `id` field. Since `id` is being added to `required`, this must fail with a `ValidationError`. Confirms the schema enforces the new requirement. |
| 3 | blankTemplate uses valid visual templateId | Frontend templates | High | In `templates.test.ts`, update the existing assertion at line 69 from `toBe('blank')` to `toBe('classic')`. Additionally, add a new assertion that `TEMPLATE_IDS` (imported from `@/types/cv`) includes `template.templateId` — verifies the template produces a schema-valid templateId. |
| 4 | Full realistic CV payload validates (multi-section regression guard) | Backend validator | Medium | In `cv_schema_test.go`, add `TestCVValidator_Validate_FullRealisticCV` with a CV containing personal, summary, experience, education, skills, languages (with `id`), certifications, and projects sections — all with populated entries. Serves as a drift guard: if any section schema changes, this catches it. Follow existing section patterns for entry structure. |
| 5 | updateTemplateId preserves all section content including language `id`s | Frontend editor store | Medium | In `editor-store.test.ts`, in the `updateTemplateId` describe block (line 507), add a test that: loads `mockCV` (which has language entries with `id`), calls `updateTemplateId('classic')`, verifies the language section's entries still have their `id` fields and all other properties intact. Proves template switch doesn't strip or mutate content. |
| 6 | Auto-save payload after template switch includes correct content | Frontend useAutoSave | Medium | In `useAutoSave.test.ts`, in the `saveNow (immediate save)` describe block (line 104), add a test that: loads `mockCV`, calls `updateTemplateId('visionary')` + `markDirty()`, calls `saveNow()`, then asserts `mockUpdate` was called with `{ content: <full CVContent> }` where `templateId` is `'visionary'` and language entries have `id` fields. End-to-end payload verification for the bug scenario. |
| 7 | Language entry with extra property rejected | Backend validator | Medium | In `cv_schema_test.go`, add a test with a language entry that has `id`, `language`, `proficiency`, plus an unknown property (e.g., `"extra": "field"`). Verifies `additionalProperties: false` still works after adding `id` to the schema. |
