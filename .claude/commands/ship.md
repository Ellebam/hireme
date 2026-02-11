# Ship

You are acting as **@engineer** shipping changes: commit, push, and open (or update) a PR. This command is **idempotent** — running it multiple times is safe. It will never create duplicate PRs or duplicate commits.

## Input

The user may optionally provide a commit message hint or PR description context. If not provided, infer from the changes.

## Process

### Phase 1: Assess State

Run all of these in parallel:

1. **`git status`** — Check for staged/unstaged/untracked changes (never use `-uall`).
2. **`git diff --stat`** — See what files changed.
3. **`git log --oneline -5`** — See recent commit messages for style.
4. **`git branch --show-current`** — Get current branch name.
5. **`gh pr list --head $(git branch --show-current) --json number,url,state --jq '.[0]'`** — Check if a PR already exists for this branch.

**Guard rails:**
- If on `main` or `master`, **stop and warn the user**. Do not commit or push to the default branch.
- If there are no changes to commit AND no unpushed commits, report "Nothing to ship" and stop.

### Phase 2: Commit (if there are uncommitted changes)

1. **Review the diff** — `git diff` (unstaged) and `git diff --cached` (staged) to understand all changes.
2. **Group related changes** — If the changes span logically distinct concerns (e.g., a fix + a new command + docs), create **separate commits** for each group. If they're all one concern, use a single commit.
3. **Stage files explicitly** — `git add <specific files>` for each commit. Never use `git add -A` or `git add .`.
4. **Write commit messages** following project format:
   ```
   [TAG] Short description

   Optional body with details.
   ```
   Tags: FEAT, FIX, REFACTOR, DOCS, TEST, CHORE
5. **Never add Co-Authored-By lines.**

If there are no uncommitted changes, skip to Phase 3.

### Phase 3: Push

1. **Check if branch has a remote tracking branch:**
   - If yes: `git push`
   - If no: `git push -u origin <branch-name>`
2. **Verify push succeeded.**

### Phase 4: PR (create or skip)

**Check the result from Phase 1, step 5** (the `gh pr list` check):

- **If a PR already exists** for this branch:
  - Report the existing PR URL.
  - Do NOT create a new PR.
  - The push in Phase 3 already updated the PR.

- **If no PR exists**, create one:
  1. **Determine the base branch** — Use `main` unless the user specifies otherwise.
  2. **Analyze ALL commits** on the branch vs base — `git log main..HEAD --oneline` and `git diff main...HEAD --stat`.
  3. **Draft PR title** — Under 70 chars, descriptive of the full branch scope.
  4. **Draft PR body** using this format:
     ```
     gh pr create --title "title" --body "$(cat <<'PREOF'
     ## Summary
     <1-3 bullet points covering ALL commits on the branch>

     ## Test plan
     <Bulleted checklist — what was tested and how>
     PREOF
     )"
     ```
  5. **Report the new PR URL.**

### Phase 5: Report

Present a concise summary:

```
## Shipped

**Branch:** <branch-name>
**Commits:** <number of new commits made>
**PR:** <URL> (new | existing — updated)
**Files changed:** <count>
```

## Rules

- **Idempotent.** Running this command twice with no new changes does nothing harmful. It reports "Nothing to ship" or "PR already exists, branch is up to date."
- **Never push to main/master.** Always stop and warn.
- **Never force push.** Use `git push`, never `git push --force`.
- **Never amend commits.** Always create new commits.
- **Never skip hooks.** No `--no-verify`.
- **Never add Co-Authored-By lines.**
- **Explicit staging only.** No `git add -A` or `git add .` — always name specific files.
- **One PR per branch.** If a PR exists, the push updates it. No duplicate PRs.
- **Don't include secrets.** Skip `.env`, credentials, tokens. Warn if they appear in the diff.
- **Separate concerns.** If changes are logically distinct, make separate commits rather than one monolithic commit.

## Arguments

$TASK_DESCRIPTION
