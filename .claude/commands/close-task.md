# Close Task

You are acting as **@pm** closing out a completed task. The tech lead has approved the PR. Your job is to update all tracking files, clean up temporary artifacts, commit, and push.

## Input

The user will provide a task ID (e.g. `T-001`). If not provided, infer it from the current branch name (e.g. `feat/t-001-*` → `T-001`).

## Process

### Phase 1: Identify the Task

1. **Read `CONTEXT.md`** — Find the task in the Active or Backlog table.
2. **Confirm the task ID, title, and current branch.** If the task isn't found in Active/Backlog, stop and ask the user.

### Phase 2: Update CONTEXT.md

Make all of the following edits to `CONTEXT.md`:

1. **Move to Done** — Remove the task row from Active or Backlog. Add a row to the Done table with the task ID, title, and branch name.
2. **Remove task details** — Delete the task's entry from the `## Task Details` section.
3. **Update dependencies** — For any task in Backlog that listed this task in "Blocked By":
   - Remove this task ID from their "Blocked By" column.
   - If "Blocked By" becomes empty, set it to `—`.
4. **Update dependency graph** — Remove this task from the ASCII dependency graph. Simplify any chains that no longer need it.
5. **Update status line** — Recalculate and update the task counts in `## Current State` → `**Status:**` line (total remaining, unblocked count).
6. **Update "Remaining Work"** — If the task resolves an item listed under `### Remaining Work`, remove that line.

### Phase 3: Clean Up Artifacts

1. **Delete notes file** — If `<TASK-ID>-NOTES.md` exists in the project root, delete it.
2. **Check for other temporary files** — Look for any other files specific to this task that should be removed.

### Phase 4: Commit & Push

1. **Stage all changes** — `git add` the modified `CONTEXT.md` and any deleted files.
2. **Commit** with message format:
   ```
   [DOCS] Close <TASK-ID>, update task board and dependencies
   ```
3. **Push** the current branch to remote.

### Phase 5: Report

Present a summary to the user:
- Task closed: ID + title
- Tasks unblocked (if any)
- Files changed
- Branch pushed

## Rules

- **CONTEXT.md is the source of truth.** All task state changes happen there.
- **Don't modify source code.** This command only touches tracking/docs files.
- **Be thorough with dependencies.** Missing an unblock can stall the next task.
- **Verify before pushing.** Run `git status` after staging to confirm only expected files are included.
- **If the task is not found**, stop and report — don't guess.

## Arguments

$TASK_DESCRIPTION
