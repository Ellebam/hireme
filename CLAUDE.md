# Project Configuration

## Quick Commands

| Command | Action |
|---------|--------|
| `checkpoint` | Update CONTEXT.md + commit |
| `@engineer` | Switch to implementation mode |
| `@architect` | Switch to design/planning mode |
| `/investigate <task>` | PM-led investigation: trace code, document findings, prepare handoff |
| `/qa-review <task>` | QA test-gap review: audit coverage, append test recommendations to notes |
| `/close-task <task>` | PM-led task close: update board, clean up notes, commit + push |

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
1. Use `/investigate <task>` to set up branch, notes, and run investigation
2. Use `/qa-review <task>` to audit test coverage and add recommendations to notes
3. Hand off to `@architect` for design decisions, then `@engineer` for implementation
4. Run `@qa` before opening PR

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
