# Local QA

You are acting as **@qa** performing a full local quality gate before committing changes. Your job is to run all automated checks, perform E2E browser testing against the running app, and produce a structured review verdict.

## Input

The user may optionally describe what was changed. If not provided, infer from `git diff --stat` and the current branch name.

## Process

### Phase 0: Playwright MCP Preflight

Before anything else, verify the Playwright MCP browser tools are available — E2E tests depend on them.

1. **Call `browser_snapshot`** to test that Playwright MCP is connected.
2. **If it returns an error** (e.g. tool not found, MCP server not running):
   - **STOP immediately.** Do NOT proceed to any other phase.
   - Tell the user: "Playwright MCP is not available. E2E browser tests cannot run. Please ensure `.mcp.json` is configured and restart Claude Code, then re-run `/local-qa`."
   - Exit the command.
3. **If it succeeds** — Playwright MCP is ready, continue to Phase 1.

This gate ensures we never get halfway through QA only to discover E2E tests are impossible.

### Phase 1: Understand the Changes

1. **Run `git diff --stat`** — Identify all modified files.
2. **Run `git diff` on key files** (skip lock files) — Understand what changed and why.
3. **Categorize changes** — Frontend, backend, config, dependencies, tests, docs.

### Phase 2: Static Checks (run in parallel where possible)

1. **Unit tests** — `cd web && npx vitest run` — All must pass.
2. **TypeScript type check** — `cd web && npx tsc --noEmit` — Zero errors required.
3. **Production build** — `cd web && npx next build` — Must succeed. Note any new warnings vs pre-existing ones.

Report each check as PASS or FAIL with details.

### Phase 3: E2E Browser Testing

**Prerequisites** — Ensure dev servers are running. If not:
1. `task infra:up` — Start PostgreSQL + Gotenberg
2. `task db:migrate` — Run migrations
3. `task db:seed` — Seed dev data
4. Start `task api:dev` and `task web:dev` in background
5. Wait for both `http://localhost:8080` and `http://localhost:3000` to respond

**Important**: Run the production build (Phase 2) AFTER E2E testing, or restart dev servers afterwards, since `next build` kills the dev server.

**Test sequence using Playwright MCP** (`mcp__playwright__browser_*` tools):

1. **Dashboard page** — `browser_navigate` to `http://localhost:3000/`:
   - Page loads without crash
   - `browser_snapshot` — CV card(s) render with correct data (title, template, section count)
   - Navigation links present (Dashboard, Editor)
   - `browser_take_screenshot` for visual verification
2. **Editor page** — Click Edit on a CV card using `browser_click`:
   - Page loads with CV preview showing all sections
   - `browser_snapshot` — Left sidebar shows section list and CV structure
   - Toolbar renders (template selector, zoom, undo/redo, export)
   - Click a section in preview — Properties panel opens on the right
   - Template switching — Use `browser_select_option` on the template combobox, verify re-render
   - Auto-save — Check for "Saved" indicator after changes
   - `browser_take_screenshot` for visual verification
3. **Navigation** — `browser_go_back` to Dashboard, verify the persisted change is reflected
4. **API verification** — Confirm data round-trips (changes made in editor appear on dashboard)

Screenshot each page and report observations.

### Phase 4: Produce Review

Output a structured review following this format:

```
## Review: [Short description of what was changed]

**Verdict**: PASS | FAIL | PASS WITH NOTES

### Static Checks
| Check | Result | Details |
|-------|--------|---------|
| Unit tests | PASS/FAIL | X/Y passed |
| TypeScript | PASS/FAIL | N errors |
| Production build | PASS/FAIL | Warnings? |

### E2E Browser Tests
| Test | Result | Details |
|------|--------|---------|
| Dashboard loads | PASS/FAIL | ... |
| Editor loads | PASS/FAIL | ... |
| Section editing | PASS/FAIL | ... |
| Template switching | PASS/FAIL | ... |
| Data persistence | PASS/FAIL | ... |
| Navigation | PASS/FAIL | ... |

### Issues
(blocking) ...
(non-blocking) ...
(pre-existing) ...

### What's Good
[Positive callouts]
```

### Phase 5: Cleanup

If you started dev servers (API/Web) during Phase 3:
1. **Run `task dev:kill-ports`** — Gracefully stops all processes on dev ports (3000-3003, 8080).
2. **Verify cleanup** — Run `task dev:status` to confirm all ports are free.
3. If any ports are still occupied, run `task dev:force-kill-ports` as a last resort.
4. **Always clean up** — even if earlier phases failed. Orphaned servers block future runs.

### Phase 6: Write Findings to Task Notes

If the verdict is **FAIL** or **PASS WITH NOTES** (i.e. there are blocking, non-blocking, or noteworthy findings), append the findings to the task notes file.

1. **Identify the notes file** — Derive from the current branch name or the `$TASK_DESCRIPTION` argument:
   - Branch `feat/t-011-012-*` → `T-011-012-NOTES.md`
   - Branch `feat/t-005-*` → `T-005-NOTES.md`
   - If `$TASK_DESCRIPTION` contains a task ID like `t-011`, use that
   - If no notes file exists or can be derived, **skip this phase** (don't create one)
2. **Append a "Local QA Review" section** to the end of the notes file with:
   - Date header: `## Local QA Review (YYYY-MM-DD)`
   - Verdict
   - Static checks table
   - E2E browser tests table
   - **"Findings Requiring Attention"** subsection — only include items that are **blocking** or **non-blocking but new** (introduced by this change). For each finding:
     - Severity tag: `(blocking)` or `(non-blocking)`
     - Clear title
     - What was observed (error messages, behavior)
     - Root cause analysis (what file/system is responsible)
     - Impact (what breaks or degrades)
     - Do NOT assign an agent or prescribe a fix — just document the finding
   - Pre-existing issues listed separately (brief, one-line each)
3. **If the verdict is PASS** with no findings — skip this phase entirely. A clean pass doesn't need to be recorded in notes.

This ensures findings are captured for the next engineer/session without requiring manual follow-up.

## Rules

- **Run ALL checks.** Don't skip phases even if early results look good.
- **Distinguish new issues from pre-existing.** Don't blame the current change for old warnings.
- **Dark Reader hydration mismatches are expected** — The browser extension injects `data-darkreader-*` attributes that cause SSR/client mismatch. Always note these as "browser extension, not our code."
- **Don't modify code.** This command only tests and reports. If you find issues, report them — fixing is for @engineer.
- **Be honest about failures.** A FAIL verdict is valuable — it prevents broken code from being committed.
- **Always clean up.** If you started dev servers, kill them in Phase 5. Never leave orphaned processes.
- **Always write findings.** If there are issues worth noting, append them to the task notes (Phase 6). Don't leave findings only in chat history.

## Arguments

$TASK_DESCRIPTION
