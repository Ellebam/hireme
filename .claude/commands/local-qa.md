# Local QA

You are acting as **@qa** performing a full local quality gate before committing changes. Your job is to run all automated checks, perform E2E browser testing against the running app, and produce a structured review verdict.

## Input

The user may optionally describe what was changed. If not provided, infer from `git diff --stat` and the current branch name.

## Process

### Phase 0: Chrome Extension Preflight

Before anything else, verify the Chrome browser extension is connected — E2E tests depend on it.

1. **Call `tabs_context_mcp`** with `createIfEmpty: true`.
2. **If it returns an error** (e.g. "No Chrome extension connected"):
   - **STOP immediately.** Do NOT proceed to any other phase.
   - Tell the user: "Chrome extension is not connected. E2E browser tests cannot run. Please connect the Claude in Chrome extension and re-run `/local-qa`."
   - Exit the command.
3. **If it succeeds** — Log the tab group info and continue to Phase 1.

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

**Test sequence using Chrome browser extension** (`mcp__claude-in-chrome__*` tools):

1. **Get tab context** — `tabs_context_mcp` to get or create a tab
2. **Dashboard page** (`http://localhost:3000/`):
   - Page loads without crash
   - CV card(s) render with correct data (title, template, section count)
   - Navigation links present (Dashboard, Editor)
   - Check console for errors (`read_console_messages` with `onlyErrors: true`)
   - Note: Dark Reader extension hydration mismatches are expected, not our code
3. **Editor page** (click Edit on a CV card):
   - Page loads with CV preview showing all sections
   - Left sidebar shows section list and CV structure
   - Toolbar renders (template selector, zoom, undo/redo, export)
   - Click a section in preview — Properties panel opens on the right
   - Template switching — Change template via `form_input` on the combobox, verify re-render
   - Auto-save — Check for "Saved" indicator after changes
   - Check console for errors
4. **Navigation** — Go back to Dashboard, verify the persisted change is reflected
5. **API verification** — Confirm data round-trips (changes made in editor appear on dashboard)

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
| Console errors | PASS/FAIL | ... |

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

## Rules

- **Run ALL checks.** Don't skip phases even if early results look good.
- **Distinguish new issues from pre-existing.** Don't blame the current change for old warnings.
- **Dark Reader hydration mismatches are expected** — The browser extension injects `data-darkreader-*` attributes that cause SSR/client mismatch. Always note these as "browser extension, not our code."
- **Don't modify code.** This command only tests and reports. If you find issues, report them — fixing is for @engineer.
- **Be honest about failures.** A FAIL verdict is valuable — it prevents broken code from being committed.
- **Always clean up.** If you started dev servers, kill them in Phase 5. Never leave orphaned processes.

## Arguments

$TASK_DESCRIPTION
