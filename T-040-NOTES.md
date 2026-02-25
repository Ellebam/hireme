# T-040: CI/CD — Build & Push Docker Images to GHCR

## Summary
- **Goal:** GitHub Actions workflow that builds Docker images for API and web, and pushes them to GHCR on main branch push or release tags
- **Acceptance:** Pushing to main (or tagging a release) builds both images and pushes them to `ghcr.io/ellebam/hireme/hireme-api` and `ghcr.io/ellebam/hireme/hireme-web`
- **Branch:** `feat/t-040-cicd-docker-ghcr`

---

## Investigation Checklist
- [x] Existing CI/CD workflows (`.github/workflows/ci.yml`)
- [x] Existing Dockerfiles (`docker/Dockerfile.api`, `docker/Dockerfile.web`)
- [x] Production docker-compose (`docker/docker-compose.yml`) — image refs
- [x] Go build: version, CGO, embedded assets
- [x] Next.js build: standalone output, env vars
- [x] Taskfile docker tasks
- [x] `.dockerignore` files
- [x] GHCR repo path (`Ellebam/hireme`)

## Findings

### Existing CI Pipeline (`.github/workflows/ci.yml`)
- Runs on push to main + PRs
- Jobs: api-lint, api-test (with PostgreSQL service), api-build, web-lint, web-typecheck, web-test, web-build, validate-schema
- Go 1.24, Node 20
- All jobs run independently (no dependencies between them)
- No Docker build/push step

### Existing Dockerfiles
Both exist at `docker/Dockerfile.api` and `docker/Dockerfile.web`.

**Dockerfile.api (lines 1-53):**
- Multi-stage: `golang:1.22-alpine` → `alpine:3.19`
- `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` (static binary)
- Non-root `appuser`, healthcheck on `/health`
- Exposes 8080

**Dockerfile.web (lines 1-52):**
- Multi-stage: `node:20-alpine` (build) → `node:20-alpine` (runner)
- Standalone Next.js output (`COPY .next/standalone`)
- Non-root `nextjs:nodejs`, healthcheck on `/`
- Exposes 3000, `HOSTNAME=0.0.0.0`

### Production docker-compose.yml
Already references GHCR images:
```yaml
api:  image: ghcr.io/${GITHUB_REPO}/hireme-api:${VERSION:-latest}
web:  image: ghcr.io/${GITHUB_REPO}/hireme-web:${VERSION:-latest}
```

### Go Embedded Assets
Templates and schema are embedded via `go:embed` — no runtime file dependencies:
- `api/internal/export/renderer.go:15` — `//go:embed templates/*.tmpl`
- `api/internal/validator/cv_schema.go:14` — `//go:embed schema/cv-schema.json`

### Taskfile Docker Tasks
```bash
task docker:build:api  # docker build -f docker/Dockerfile.api -t hireme-api:local .
task docker:build:web  # docker build -f docker/Dockerfile.web -t hireme-web:local .
```
Build context is `.` (project root) for both.

---

## Gaps Found

### BLOCKER — Dockerfile.api: Go version mismatch
- `Dockerfile.api` line 2: `FROM golang:1.22-alpine`
- `api/go.mod` line 3: `go 1.24.0` (toolchain `go1.24.7`)
- Build will fail because Go 1.22 can't compile code requiring Go 1.24
- **Fix:** Update to `golang:1.24-alpine`

### BLOCKER — Dockerfile.api: Build context mismatch
- Dockerfile assumes build context is `api/` (e.g., `COPY go.mod go.sum ./`)
- Taskfile passes `.` (project root) as context: `docker build -f docker/Dockerfile.api .`
- `go.mod` is at `api/go.mod`, not at project root
- **Fix:** Either change Dockerfiles to use `api/` prefixed paths, OR change build context to `api/` and `web/` respectively. Since all assets are embedded (no cross-directory deps), using subdirectory contexts is cleanest.

### BLOCKER — Dockerfile.api: Phantom templates COPY
- Line 37: `COPY --from=builder /app/templates /app/templates`
- Templates live at `api/internal/export/templates/` and are embedded via `go:embed`
- This path doesn't exist at the API root level; it will fail the build
- **Fix:** Remove this line — templates are in the binary

### Medium — No `.dockerignore` files
- No `.dockerignore` exists anywhere in the project
- Build context includes `.git/`, `node_modules/`, `data/`, test files, etc.
- Makes Docker builds slow and images potentially larger
- **Fix:** Add `.dockerignore` files for `api/` and `web/` (or root level)

### Medium — Dockerfile.web: Build context mismatch
- Same issue as API: expects `package.json` at root of context
- `package.json` is at `web/package.json`
- **Fix:** Same approach — change build context to `web/`

---

## Architect Handoff

### Decisions Needed

**1. Build context strategy** (Quick Decision)
- **Option A (Recommended):** Change build context to subdirectories (`api/` and `web/`). Both Dockerfiles already expect this. No cross-directory deps since templates and schema are embedded.
- **Option B:** Rewrite Dockerfiles with project-root context (`COPY api/go.mod api/go.sum ./`). More complex, needed only if Dockerfiles need files from sibling directories.
- **Recommendation:** Option A — simpler, matches how Dockerfiles are already written.

**2. Workflow trigger strategy** (Quick Decision)
- **Option A (Recommended):** Trigger on push to main + version tags (`v*`). Main gets `latest` tag, version tags get `v1.0.0` tag. CI must pass first (use `needs:` to depend on existing CI jobs, or run after CI workflow).
- **Option B:** Trigger only on version tags. Simpler but no `latest` image on main.
- **Recommendation:** Option A — supports both development (pull `latest`) and releases (pull `v1.0.0`).

**3. CI dependency strategy** (Quick Decision)
- **Option A (Recommended):** Add Docker build/push jobs to existing `ci.yml` with `needs:` on all CI jobs. Single workflow, images only built when all tests pass.
- **Option B:** Separate `build-push.yml` workflow triggered by `workflow_run` completion of CI. More modular but more complex.
- **Recommendation:** Option A — simpler, single workflow, `needs:` ensures tests pass first.

**4. Multi-arch builds** (Quick Decision)
- **Option A (Recommended):** Single arch `linux/amd64` only. Covers all common deployment targets.
- **Option B:** Multi-arch `linux/amd64,linux/arm64`. Supports ARM servers (AWS Graviton, Apple Silicon). More build time.
- **Recommendation:** Option A for now — can add arm64 later if needed.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `docker/Dockerfile.api` | Fix Go version (1.22→1.24), remove phantom templates COPY | BLOCKER |
| `docker/Dockerfile.web` | No changes needed (if context changed to `web/`) | — |
| `.github/workflows/ci.yml` | Add docker build+push jobs with `needs:` on CI jobs | BLOCKER (new) |
| `Taskfile.yml` | Fix build context for `docker:build:api` and `docker:build:web` | Medium |
| `api/.dockerignore` | New file — exclude test files, docs, etc. | Medium |
| `web/.dockerignore` | New file — exclude test files, node_modules, etc. | Medium |

### Workflow Design (Sketch)

```yaml
# Added to .github/workflows/ci.yml
docker-build-push:
  name: Build & Push Docker Images
  runs-on: ubuntu-latest
  needs: [api-lint, api-test, api-build, web-lint, web-typecheck, web-test, web-build]
  if: github.ref == 'refs/heads/main' || startsWith(github.ref, 'refs/tags/v')
  permissions:
    contents: read
    packages: write
  steps:
    - uses: actions/checkout@v6
    - uses: docker/login-action@v3    # Login to GHCR
    - uses: docker/metadata-action@v5  # Generate tags (latest, version)
    - uses: docker/build-push-action@v6  # Build + push API
    - uses: docker/build-push-action@v6  # Build + push Web
```

### GHCR Image Paths
- API: `ghcr.io/ellebam/hireme/hireme-api`
- Web: `ghcr.io/ellebam/hireme/hireme-web`
- Tags: `latest` (main), `v1.0.0` (release tags), `sha-<short>` (commit SHA)

### Tests to Add/Update
- No unit tests needed — this is infrastructure
- **Verification:** After implementation, push to a branch and verify the workflow runs (can do dry-run with `push: false` first)
- The existing CI jobs already validate that builds succeed

### Constraints
- PR size: ~150 lines (small: workflow YAML + Dockerfile fixes + dockerignore)
- No secrets need to be configured — `GITHUB_TOKEN` is automatic for GHCR
- `docker-compose.yml` already uses correct GHCR paths — no changes needed there

### Recommended Next Agent
**@architect** — Quick decisions on the 4 items above, then **@engineer** to implement. All decisions have clear recommendations, so architect pass should be fast.

---

## Test Plan
- Verify Dockerfiles build locally after fixes (`task docker:build`)
- Verify workflow YAML syntax (can use `actionlint` or GitHub's validation)
- After merge to main, verify images appear at `ghcr.io/ellebam/hireme/`

---

## Architect Plan

### Decisions

**Decision 1: Build context strategy**
- Option A: Subdirectory contexts (`api/` and `web/`) — Pros: Dockerfiles already written for this, smaller context, no cross-dir deps (assets embedded). Cons: None.
- Option B: Project-root context with prefixed COPY paths — Pros: Single context. Cons: Requires rewriting all Dockerfiles, needed only if cross-dir deps exist (they don't).
- **Choice:** A — **Rationale:** Both Dockerfiles already expect subdirectory context. All templates and schemas are embedded via `go:embed`. Zero cross-directory dependencies.

**Decision 2: Workflow trigger strategy**
- Option A: Push to main + version tags (`v*`) — Pros: `latest` always available for dev, versioned tags for releases. Cons: Slightly more config.
- Option B: Version tags only — Pros: Simpler. Cons: No `latest` image, can't test deployment without tagging.
- **Choice:** A — **Rationale:** T-041 (first release) needs both `latest` for iterative testing and `v1.0.0` for the release tag.

**Decision 3: CI dependency strategy**
- Option A: Add Docker jobs to existing `ci.yml` with `needs:` — Pros: Single workflow, images only built when all tests pass, no `workflow_run` complexity. Cons: Longer YAML file.
- Option B: Separate `build-push.yml` with `workflow_run` trigger — Pros: Modular. Cons: `workflow_run` has quirks (doesn't trigger on PRs, delayed execution, harder to debug).
- **Choice:** A — **Rationale:** Simpler, more reliable, and the CI YAML is still manageable size. Both docker jobs depend on ALL CI jobs passing to keep `latest` tag consistent across API and web.

**Decision 4: Multi-arch builds**
- Option A: `linux/amd64` only — Pros: Fast builds, simple. Cons: No ARM support.
- Option B: `linux/amd64` + `linux/arm64` — Pros: ARM support. Cons: 2x build time, QEMU emulation for cross-arch.
- **Choice:** A — **Rationale:** No ARM deployment targets planned. Can add `linux/arm64` to the `platforms:` list later with a one-line change.

### Workflow Architecture

Two parallel Docker jobs added to `ci.yml`, both gated on all CI jobs:

```
ci.yml
├── api-lint ────┐
├── api-test ────┤
├── api-build ───┤
├── web-lint ────┼──→ docker-api (build + push hireme-api)
├── web-typecheck┤
├── web-test ────┤
├── web-build ───┼──→ docker-web (build + push hireme-web)
└── validate-schema
```

- Both docker jobs run in parallel after all CI jobs pass
- Only on `push` to `main` or version tags (`v*`) — skipped on PRs
- Each job: checkout → setup-buildx → login GHCR → metadata → build+push
- GHA cache (`type=gha`) for Docker layer caching

### Tag Strategy

| Trigger | Tags produced |
|---------|--------------|
| Push to `main` | `latest`, `sha-abc1234` |
| Tag `v1.0.0` | `v1.0.0`, `sha-abc1234` |

Uses `docker/metadata-action@v5` with:
```yaml
tags: |
  type=raw,value=latest,enable={{is_default_branch}}
  type=ref,event=tag
  type=sha,prefix=sha-
```

### GHCR Image Paths
- `ghcr.io/ellebam/hireme/hireme-api`
- `ghcr.io/ellebam/hireme/hireme-web`

(`github.repository` = `Ellebam/hireme`, GHCR lowercases automatically)

### Files Changing

| # | File | Change | Lines |
|---|------|--------|-------|
| 1 | `docker/Dockerfile.api` | `golang:1.22-alpine` → `golang:1.24-alpine`, remove phantom `COPY templates` line | ~3 |
| 2 | `.github/workflows/ci.yml` | Add `docker-api` and `docker-web` jobs (append after existing jobs), add `on.push.tags` trigger | ~80 |
| 3 | `Taskfile.yml` | Fix build context: `.` → `api` and `.` → `web` in docker build commands | ~2 |
| 4 | `api/.dockerignore` | New file — exclude `bin/`, coverage, air config | ~8 |
| 5 | `web/.dockerignore` | New file — exclude `node_modules/`, `.next/`, coverage, env files | ~10 |

### Files NOT Changing

| File | Reason |
|------|--------|
| `docker/Dockerfile.web` | Already correct for `web/` build context |
| `docker/docker-compose.yml` | Already references correct GHCR paths |
| `docker/docker-compose.infra.yml` | Dev-only, unrelated |

### Exact Code Changes

#### 1. `docker/Dockerfile.api` — 2 edits

**Edit A:** Line 2 — Update Go version
```diff
-FROM golang:1.22-alpine AS builder
+FROM golang:1.24-alpine AS builder
```

**Edit B:** Lines 36-37 — Remove phantom templates COPY
```diff
 # Copy binary from builder
 COPY --from=builder /app/server /app/server
-
-# Copy templates
-COPY --from=builder /app/templates /app/templates
```

#### 2. `.github/workflows/ci.yml` — 2 edits

**Edit A:** Add version tag trigger (line 5 area)
```diff
 on:
   push:
     branches: [main]
+    tags: ['v*']
   pull_request:
     branches: [main]
```

**Edit B:** Append two new jobs after `validate-schema` job (after line 221)
```yaml
  # ══════════════════════════════════════════════════════════════════════════
  # DOCKER BUILD & PUSH
  # ══════════════════════════════════════════════════════════════════════════
  docker-api:
    name: Docker API Image
    runs-on: ubuntu-latest
    needs: [api-lint, api-test, api-build, web-lint, web-typecheck, web-test, web-build, validate-schema]
    if: github.event_name == 'push'
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v6

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/${{ github.repository }}/hireme-api
          tags: |
            type=raw,value=latest,enable={{is_default_branch}}
            type=ref,event=tag
            type=sha,prefix=sha-

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: ./api
          file: docker/Dockerfile.api
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha,scope=api
          cache-to: type=gha,mode=max,scope=api

  docker-web:
    name: Docker Web Image
    runs-on: ubuntu-latest
    needs: [api-lint, api-test, api-build, web-lint, web-typecheck, web-test, web-build, validate-schema]
    if: github.event_name == 'push'
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v6

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/${{ github.repository }}/hireme-web
          tags: |
            type=raw,value=latest,enable={{is_default_branch}}
            type=ref,event=tag
            type=sha,prefix=sha-

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: ./web
          file: docker/Dockerfile.web
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha,scope=web
          cache-to: type=gha,mode=max,scope=web
```

Note: The `if: github.event_name == 'push'` condition is sufficient because the `on.push` trigger already limits to `main` branch and `v*` tags. PRs use `pull_request` event which won't match. No need for a more complex `if` expression.

#### 3. `Taskfile.yml` — 2 line changes

```diff
  docker:build:api:
    desc: Build API Docker image
    cmds:
-     - docker build -f docker/Dockerfile.api -t hireme-api:local .
+     - docker build -f docker/Dockerfile.api -t hireme-api:local api
```

```diff
  docker:build:web:
    desc: Build Web Docker image
    cmds:
-     - docker build -f docker/Dockerfile.web -t hireme-web:local .
+     - docker build -f docker/Dockerfile.web -t hireme-web:local web
```

#### 4. `api/.dockerignore` — New file

```
bin/
coverage.*
*.html
.air.toml
```

#### 5. `web/.dockerignore` — New file

```
node_modules/
.next/
coverage/
.env
.env.*
```

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Docker build fails in CI due to untested Dockerfile | Low | Medium | Dockerfiles build locally first (`task docker:build`); CI runs tests before docker |
| GHCR push permission denied | Low | Low | `GITHUB_TOKEN` + `packages: write` is the standard pattern; well-documented |
| GHA cache miss on first run (slow build) | Certain | Low | First build is uncached (~2-3 min); subsequent builds use GHA cache |
| `golang:1.24-alpine` image unavailable | Very Low | High | Verified it exists on Docker Hub with multiple patch versions |

### Consequences
- **Easier:** Deploying via docker-compose (just `docker compose pull && docker compose up -d`)
- **Easier:** T-041 (first release) can proceed once this merges
- **Unchanged:** Local development workflow, existing CI jobs

### Engineer Execution Steps

1. **Edit `docker/Dockerfile.api`** — Update Go version and remove phantom templates COPY (2 edits)
2. **Edit `.github/workflows/ci.yml`** — Add `tags: ['v*']` trigger + append `docker-api` and `docker-web` jobs
3. **Edit `Taskfile.yml`** — Fix build contexts (2 line changes)
4. **Create `api/.dockerignore`** — Exclude build artifacts
5. **Create `web/.dockerignore`** — Exclude node_modules, .next, env files
6. **Verify locally** — Run `task docker:build` to confirm both images build successfully
7. **Verify image size** — `docker images | grep hireme` — API should be ~20-30MB, web ~100-150MB

### Verification Gates

1. `task docker:build:api` — builds without errors
2. `task docker:build:web` — builds without errors
3. `docker images | grep hireme` — both images present with reasonable sizes
4. Existing CI still passes (no changes to existing jobs)
5. Workflow YAML is valid (GitHub will validate on push)

### QA-Recommended Additional Verification

This is a pure infrastructure task — no unit or integration tests apply. All verification is manual/build-based.

| # | Check | Layer | Priority | Description |
|---|-------|-------|----------|-------------|
| 1 | API container runtime smoke test | Docker | High | After `task docker:build:api`, run the container and verify it starts: `docker run --rm -d -p 8080:8080 -e AUTH_BYPASS_ENABLED=true -e DATABASE_URL=fake hireme-api:local`. It will fail to connect to DB, but the binary should start (check logs with `docker logs`). Verifies: correct binary path, permissions, non-root user works. Stop it after checking logs. |
| 2 | Web container runtime smoke test | Docker | High | After `task docker:build:web`, run: `docker run --rm -d -p 3000:3000 hireme-web:local`. Verify `curl -s http://localhost:3000` returns HTML (Next.js standalone should serve even without API). Verifies: standalone output copied correctly, `server.js` runs, non-root user works, port binding. |
| 3 | Workflow YAML lint | CI | High | Run `actionlint .github/workflows/ci.yml` locally (install via `brew install actionlint` if needed). Catches syntax errors, invalid `needs` references, and expression typos that GitHub would only surface after push. This is a required gate, not optional. |
| 4 | PR trigger regression check | CI | Medium | After adding `tags: ['v*']` trigger, manually verify the workflow YAML logic: (a) `on.push.branches` still lists only `[main]`, (b) `on.push.tags` lists only `['v*']`, (c) Docker jobs have `if: github.event_name == 'push'`. This ensures PRs don't trigger Docker builds. Verified by reading the final YAML — no runtime test needed. |
| 5 | `.dockerignore` effectiveness | Docker | Medium | After creating `.dockerignore` files, verify they work by checking Docker build context size. Run `docker build -f docker/Dockerfile.api api` and look for the `Sending build context to Docker daemon` line — API context should be <50MB (without `.dockerignore` it could include all Go test data). For web, context should be <10MB (excludes `node_modules/` and `.next/`). |
| 6 | Image size sanity check | Docker | Medium | After builds, run `docker images | grep hireme`. Expected ranges: API ~20-40MB (static Go binary + Alpine), Web ~100-200MB (Node.js + standalone Next.js). If either image is >500MB, something is wrong (likely `.dockerignore` not excluding `node_modules/` or build context issue). |

**Notes on existing Test Plan items:**
- "Verify workflow YAML syntax" is upgraded from optional to **required** (check #3 above)
- "After merge to main, verify images appear at GHCR" remains valid but is post-merge — can't gate on it pre-commit
- Local Docker build verification (checks #1, #2) adds runtime validation beyond just "build succeeds"
