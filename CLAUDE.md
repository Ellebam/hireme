# Project Configuration

## Agent System

Invoke specialized agents with @mentions or let context determine the mode.

| Agent | Invoke | Purpose |
|-------|--------|---------|
| Engineer | `@engineer` | Write code, debug, implement |
| Architect | `@architect` | Design decisions, structure |
| QA | `@qa` | Code review, testing, quality |
| DevOps | `@devops` | CI/CD, deployment, infrastructure |
| Security | `@security` | Security review, vulnerabilities |
| PM | `@pm` | Planning, prioritization |

**Usage:**
- Explicit: `@engineer fix the login bug`
- Implicit: Agent auto-selects based on task
- Chained: Agents reference each other for handoffs

Read agent files in `.claude/agents/` for detailed processes.

## Context

**Project specifics are in `CONTEXT.md`** - tech stack, structure, conventions.

Always check CONTEXT.md when:
- Starting work in this project
- Making decisions about patterns to follow
- Running project-specific commands

## Workflow

### Starting a Session
1. Check `WORKLOG.md` for where you left off
2. Check `CONTEXT.md` for project specifics
3. Pick up the next task

### Ending a Session
1. Update `WORKLOG.md` with what was done
2. Commit work (WIP commits are fine)
3. Note blockers or next steps

### Commit Format
```
[TAG] Short description

Optional longer explanation.

Tags: FEAT, FIX, REFACTOR, DOCS, TEST, CHORE
```

## Quick Reference

**List agents**: Check the table above or read `.claude/agents/`
**Project context**: See `CONTEXT.md`
**Work history**: See `WORKLOG.md`
