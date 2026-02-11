# Project Configuration

## Quick Commands

| Command | Action |
|---------|--------|
| `checkpoint` | Update CONTEXT.md + commit |
| `@engineer` | Switch to implementation mode |
| `@architect` | Switch to design/planning mode |
| `/investigate <task>` | PM-led investigation: trace code, document findings, prepare handoff |
| `/plan <task>` | Architect-led planning: research, decisions, implementation plan appended to notes |
| `/qa-plan <task>` | QA test planning: audit coverage gaps, append test recommendations to notes |
| `/close-task <task>` | PM-led task close: update board, clean up notes, commit + push |
| `/local-qa` | Full local QA gate: unit tests, type check, build, E2E browser testing |
| `/ship` | Idempotent commit + push + PR creation (safe to re-run) |

## Agent System

| Agent | Purpose |
|-------|---------|
| Engineer | Write code, debug, implement |
| Architect | Design decisions, structure |
| QA | Code review, testing, quality |
| DevOps | CI/CD, deployment, infrastructure |
| Security | Security review, vulnerabilities |
| PM | Planning, prioritization |

Agent files in `.claude/agents/` for detailed processes.

## Key Files

| File | Purpose |
|------|---------|
| `CONTEXT.md` | Tech stack, architecture, task board, current state |
| `<TASK-ID>-NOTES.md` | Per-task investigation notes (temporary, created by `/investigate`) |
| `.env` | Environment variables (DATABASE_URL, AUTH_BYPASS_ENABLED) |

## Workflow

### Starting a Session
1. Check `CONTEXT.md` → "Current State" and "Task Board" sections
2. Pick a task or resume an active one

### Working on a Task
1. `/investigate <task>` — PM traces code, documents findings, prepares handoff
2. `/plan <task>` — Architect researches, makes decisions, writes implementation plan
3. `/qa-plan <task>` — QA audits coverage gaps, appends test recommendations
4. `@engineer` — Implements code + tests per plan
5. `/local-qa` — Full local QA gate before commit (tests, types, build, E2E)
6. QA review — tech lead + CodeRabbit review the PR
7. `/close-task <task>` — PM moves task to Done, cleans up

### Ending a Session
Say `checkpoint` to:
1. Update `CONTEXT.md` current state and task board
2. Commit with proper tag

### Commit Format
```
[TAG] Short description

Tags: FEAT, FIX, REFACTOR, DOCS, TEST, CHORE
```

## Common Commands

```bash
task infra:up      # Start PostgreSQL + Gotenberg
task api:dev       # Run Go API (hot reload)
task web:dev       # Run Next.js
task db:migrate    # Run migrations
task db:seed       # Seed dev data (user + CV)
task api:sqlc      # Regenerate sqlc code
```

## Dev Environment Notes

- **`next build` kills `task web:dev`** — Running `npx next build` while the dev server is running will stop it. If you need both, run E2E tests first, then the build, or restart the dev server after.
- **Template selector in editor** — It's a native `<select>` (combobox). Use `form_input` tool to change it; coordinate clicks may not open native dropdowns.

## GitHub CLI (gh) Reference

```bash
# Dependabot alerts (all, with details)
gh api repos/{owner}/{repo}/dependabot/alerts --jq '.[] | {number, state, severity: .security_advisory.severity, package: .security_vulnerability.package.name, ecosystem: .security_vulnerability.package.ecosystem, summary: .security_advisory.summary, patched: .security_vulnerability.first_patched_version.identifier}'

# Dependabot alerts (open only)
gh api repos/{owner}/{repo}/dependabot/alerts?state=open

# Dismiss a Dependabot alert
gh api -X PATCH repos/{owner}/{repo}/dependabot/alerts/{alert_number} -f state=dismissed -f dismissed_reason=tolerable_risk -f dismissed_comment="reason"

# PR operations
gh pr create --title "title" --body "body"
gh pr view <number>
gh pr checks <number>
gh api repos/{owner}/{repo}/pulls/{number}/comments
```
