# Project Configuration

## Quick Commands

| Command | Action |
|---------|--------|
| `update worklog` | Update WORKLOG.md with session progress |
| `checkpoint` | Update worklog + CONTEXT.md + commit |
| `@engineer` | Switch to implementation mode |
| `@architect` | Switch to design/planning mode |

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
| `CONTEXT.md` | Tech stack, architecture, current state |
| `WORKLOG.md` | Session history, what was done, next steps |
| `.env` | Environment variables (DATABASE_URL, AUTH_BYPASS_ENABLED) |

## Workflow

### Starting a Session
1. Read last entry in `WORKLOG.md` for context
2. Check `CONTEXT.md` → "Current State" section

### Ending a Session
Say `checkpoint` or `update worklog` to:
1. Update `WORKLOG.md` with completed work
2. Update `CONTEXT.md` current state if needed
3. Commit with proper tag

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
