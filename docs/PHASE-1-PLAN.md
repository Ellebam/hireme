# HireMe — Phase 1: Foundation

**Duration:** Weeks 1-3  
**Goal:** Working backend with data persistence, basic auth flow, and test coverage.

---

## Week 1: Project Bootstrap

### Milestone: Development environment fully operational

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 1.1 | Initialize Git repository | Repo created with .gitignore, initial commit |
| 1.2 | Create project structure | All directories created per architecture doc |
| 1.3 | Set up Go module | go.mod with dependencies, go mod tidy passes |
| 1.4 | Set up Next.js project | Next.js 14 scaffolded, TypeScript configured |
| 1.5 | Configure Taskfile | All dev commands working |
| 1.6 | Docker Compose for infra | PostgreSQL + Gotenberg start with `task infra:up` |
| 1.7 | Environment configuration | .env.example complete, config loading in Go |
| 1.8 | CI Pipeline (basic) | GitHub Actions: lint + test on PR |

### Detailed Tasks

#### 1.1 Initialize Git Repository
```bash
cd hireme
git init
git add .
git commit -m "Initial project structure"
```

#### 1.2 Create Project Structure
All core directories from architecture doc. Placeholder files where needed.

#### 1.3 Set Up Go Module
```bash
cd api
go mod init github.com/yourusername/hireme/api
go mod tidy
```

Install core dependencies:
- chi (router)
- pgx (PostgreSQL driver)
- godotenv (env loading)
- testify (testing)
- uuid (ID generation)

#### 1.4 Set Up Next.js Project
```bash
cd web
npx create-next-app@latest . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
```

Add dependencies:
- shadcn/ui (init + button, card, input)
- zustand (state management)
- next-intl (i18n)
- dnd-kit (drag and drop)

#### 1.5 Configure Taskfile
Verify all commands:
- `task infra:up` / `task infra:down`
- `task api:dev` / `task api:test`
- `task web:dev` / `task web:test`
- `task db:migrate`

#### 1.6 Docker Compose for Infrastructure
Test sequence:
1. `task infra:up`
2. Verify PostgreSQL accepts connections
3. Verify Gotenberg responds to health check

#### 1.7 Environment Configuration
Create `api/internal/config/config.go`:
- Load from .env.local
- Validate required vars
- Provide defaults for optional vars

#### 1.8 CI Pipeline
Create `.github/workflows/ci.yml`:
- Trigger on PR to main
- Run Go lint (golangci-lint)
- Run Go tests
- Run Next.js lint
- Run Next.js type check

### Week 1 Exit Criteria
- [ ] `git clone && ./scripts/setup-dev.sh` succeeds on fresh machine
- [ ] `task infra:up` starts PostgreSQL and Gotenberg
- [ ] `task api:dev` starts Go server (even if empty)
- [ ] `task web:dev` starts Next.js dev server
- [ ] GitHub Actions CI passes

---

## Week 2: Data Layer

### Milestone: Database schema implemented with migrations and type-safe queries

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 2.1 | Write initial migration | All tables created per schema design |
| 2.2 | Configure sqlc | sqlc.yaml configured, generates Go code |
| 2.3 | Write sqlc queries | All CRUD operations for users, cvs, assets |
| 2.4 | Implement repository layer | Interfaces + PostgreSQL implementations |
| 2.5 | Implement domain models | Pure Go structs, no external deps |
| 2.6 | Write repository tests | 80%+ coverage on repository layer |
| 2.7 | Seed script for dev user | Dev user created for auth bypass |

### Detailed Tasks

#### 2.1 Write Initial Migration
File: `api/db/migrations/000001_init.up.sql`
- users table
- cvs table (with JSONB content)
- assets table
- export_jobs table
- audit_log table
- Indexes
- RLS policies (not enabled yet)

#### 2.2 Configure sqlc
File: `api/db/sqlc.yaml`
- PostgreSQL engine
- pgx/v5 driver
- Custom type mappings (uuid, timestamptz, jsonb)

#### 2.3 Write sqlc Queries
Files in `api/db/queries/`:
- `users.sql`: GetByID, GetByExternalID, Create, Update, Delete
- `cvs.sql`: GetByID, GetByUserID, Create, Update, Delete, Count
- `assets.sql`: GetByID, Create, Delete, GetTotalSize

Run `task api:sqlc` to generate code.

#### 2.4 Implement Repository Layer

```go
// api/internal/repository/repository.go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*domain.User, error)
    GetByExternalID(ctx context.Context, provider, externalID string) (*domain.User, error)
    Create(ctx context.Context, user *domain.User) error
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}

type CVRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*domain.CV, error)
    GetByUserID(ctx context.Context, userID string) (*domain.CV, error)
    Create(ctx context.Context, cv *domain.CV) error
    Update(ctx context.Context, cv *domain.CV) error
    Delete(ctx context.Context, id uuid.UUID) error
    CountByUserID(ctx context.Context, userID string) (int, error)
}

type AssetRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
    Create(ctx context.Context, asset *domain.Asset) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetTotalSizeByUserID(ctx context.Context, userID string) (int64, error)
}
```

PostgreSQL implementations in `api/internal/repository/postgres/`.

#### 2.5 Implement Domain Models

```go
// api/internal/domain/user.go
type User struct {
    ID                string
    ExternalID        string
    Provider          string
    Email             string
    EmailVerified     bool
    DisplayName       string
    Tier              string
    CVLimit           int
    StorageLimitBytes int64
    StorageUsedBytes  int64
    Locale            string
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

// api/internal/domain/cv.go
type CV struct {
    ID            uuid.UUID
    UserID        string
    Title         string
    SchemaVersion string
    Content       json.RawMessage
    IsActive      bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

#### 2.6 Write Repository Tests
Use testcontainers-go for integration tests:
```go
func TestUserRepository_Create(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(t) // Uses testcontainers
    repo := postgres.NewUserRepository(db)
    
    user := &domain.User{
        ID:         "test-user-1",
        ExternalID: "ext-123",
        Provider:   "google",
        Email:      "test@example.com",
    }
    
    err := repo.Create(ctx, user)
    assert.NoError(t, err)
    
    retrieved, err := repo.GetByID(ctx, user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, retrieved.Email)
}
```

#### 2.7 Seed Script for Dev User
Update `scripts/seed-dev-user.sh` to insert user matching `AUTH_BYPASS_USER_ID`.

### Week 2 Exit Criteria
- [ ] `task db:migrate` applies all migrations
- [ ] `task api:sqlc` generates query code without errors
- [ ] All repository interfaces have PostgreSQL implementations
- [ ] Repository tests pass with 80%+ coverage
- [ ] `task db:seed` creates dev user

---

## Week 3: Core API

### Milestone: REST API with authentication and basic CRUD operations

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 3.1 | Set up Chi router | Router configured with middleware |
| 3.2 | Implement middleware | Auth, logging, request ID, error handling |
| 3.3 | Implement auth bypass | Dev mode skips OIDC, uses hardcoded user |
| 3.4 | Implement user handlers | GET /users/me works |
| 3.5 | Implement CV handlers | Full CRUD for CVs |
| 3.6 | Implement JSON Schema validation | CV content validated against schema |
| 3.7 | Write handler tests | 70%+ coverage on handlers |
| 3.8 | API documentation | OpenAPI spec (basic) |

### Detailed Tasks

#### 3.1 Set Up Chi Router

```go
// api/cmd/server/main.go
func main() {
    cfg := config.Load()
    db := database.Connect(cfg.DatabaseURL)
    
    // Repositories
    userRepo := postgres.NewUserRepository(db)
    cvRepo := postgres.NewCVRepository(db)
    
    // Services
    userSvc := service.NewUserService(userRepo)
    cvSvc := service.NewCVService(cvRepo, userRepo)
    
    // Handlers
    h := handler.New(userSvc, cvSvc)
    
    // Router
    r := chi.NewRouter()
    
    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{...}))
    
    // Health endpoints (no auth)
    r.Get("/health", h.Health)
    r.Get("/ready", h.Ready)
    
    // API routes
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(middleware.Auth(cfg))
        
        r.Get("/users/me", h.GetCurrentUser)
        r.Patch("/users/me", h.UpdateCurrentUser)
        
        r.Get("/cv", h.GetCV)
        r.Post("/cv", h.CreateCV)
        r.Put("/cv/{id}", h.UpdateCV)
        r.Delete("/cv/{id}", h.DeleteCV)
    })
    
    slog.Info("starting server", "port", cfg.ServerPort)
    http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r)
}
```

#### 3.2 Implement Middleware

**Request ID:**
```go
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Logger:**
```go
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        ww := &responseWriter{ResponseWriter: w, status: 200}
        
        next.ServeHTTP(ww, r)
        
        slog.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", ww.status,
            "duration", time.Since(start),
            "request_id", GetRequestID(r.Context()),
        )
    })
}
```

#### 3.3 Implement Auth Bypass

```go
func Auth(cfg *config.Config) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Dev bypass
            if cfg.AuthBypassEnabled {
                ctx := context.WithValue(r.Context(), UserIDKey, cfg.AuthBypassUserID)
                next.ServeHTTP(w, r.WithContext(ctx))
                return
            }
            
            // Real auth: validate JWT
            token := extractBearerToken(r)
            if token == "" {
                httputil.Error(w, http.StatusUnauthorized, "missing token")
                return
            }
            
            claims, err := validateJWT(token, cfg.JWTSecret)
            if err != nil {
                httputil.Error(w, http.StatusUnauthorized, "invalid token")
                return
            }
            
            ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

#### 3.4-3.5 Implement Handlers

```go
// api/internal/handler/cv.go
func (h *Handler) GetCV(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := middleware.GetUserID(ctx)
    
    cv, err := h.cvService.GetByUserID(ctx, userID)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            httputil.Error(w, http.StatusNotFound, "cv not found")
            return
        }
        httputil.Error(w, http.StatusInternalServerError, "internal error")
        return
    }
    
    httputil.JSON(w, http.StatusOK, cv)
}

func (h *Handler) CreateCV(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := middleware.GetUserID(ctx)
    
    var req CreateCVRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    // Validate CV content against JSON Schema
    if err := h.validator.ValidateCV(req.Content); err != nil {
        httputil.Error(w, http.StatusBadRequest, err.Error())
        return
    }
    
    cv, err := h.cvService.Create(ctx, userID, req.Title, req.Content)
    if err != nil {
        if errors.Is(err, domain.ErrCVLimitReached) {
            httputil.Error(w, http.StatusForbidden, "cv limit reached")
            return
        }
        httputil.Error(w, http.StatusInternalServerError, "internal error")
        return
    }
    
    httputil.JSON(w, http.StatusCreated, cv)
}
```

#### 3.6 Implement JSON Schema Validation

```go
// api/internal/validator/cv_schema.go
type CVValidator struct {
    schema *jsonschema.Schema
}

func NewCVValidator(schemaPath string) (*CVValidator, error) {
    data, err := os.ReadFile(schemaPath)
    if err != nil {
        return nil, fmt.Errorf("reading schema: %w", err)
    }
    
    compiler := jsonschema.NewCompiler()
    if err := compiler.AddResource("cv-schema.json", bytes.NewReader(data)); err != nil {
        return nil, fmt.Errorf("adding schema: %w", err)
    }
    
    schema, err := compiler.Compile("cv-schema.json")
    if err != nil {
        return nil, fmt.Errorf("compiling schema: %w", err)
    }
    
    return &CVValidator{schema: schema}, nil
}

func (v *CVValidator) ValidateCV(content json.RawMessage) error {
    var data interface{}
    if err := json.Unmarshal(content, &data); err != nil {
        return fmt.Errorf("invalid JSON: %w", err)
    }
    
    if err := v.schema.Validate(data); err != nil {
        return fmt.Errorf("schema validation failed: %w", err)
    }
    
    return nil
}
```

#### 3.7 Write Handler Tests

```go
func TestHandler_CreateCV(t *testing.T) {
    // Setup
    mockCVService := mock.NewCVService()
    h := handler.New(nil, mockCVService, nil)
    
    body := `{"title": "My CV", "content": {...}}`
    req := httptest.NewRequest("POST", "/api/v1/cv", strings.NewReader(body))
    req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "test-user"))
    
    w := httptest.NewRecorder()
    
    // Act
    h.CreateCV(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

#### 3.8 API Documentation
Create `docs/API.md` with endpoint documentation.
Consider OpenAPI/Swagger spec for future.

### Week 3 Exit Criteria
- [ ] `task api:dev` starts server, responds to `/health`
- [ ] Auth bypass works: requests authenticated as dev user
- [ ] `POST /api/v1/cv` creates CV with valid schema
- [ ] `GET /api/v1/cv` returns user's CV
- [ ] `PUT /api/v1/cv/:id` updates CV (only owner)
- [ ] Invalid CV content rejected with clear error
- [ ] Handler tests pass with 70%+ coverage
- [ ] API responds with proper error format

---

## Phase 1 Completion Checklist

### Infrastructure
- [ ] PostgreSQL running via Docker Compose
- [ ] Gotenberg running via Docker Compose
- [ ] All Taskfile commands working
- [ ] CI pipeline passing

### Backend
- [ ] Configuration loading from environment
- [ ] Database migrations applied
- [ ] sqlc queries generated
- [ ] Repository layer implemented + tested
- [ ] Service layer with business logic
- [ ] Chi router with middleware
- [ ] Auth bypass for development
- [ ] CV CRUD endpoints working
- [ ] JSON Schema validation
- [ ] Structured logging

### Testing
- [ ] Repository tests: 80%+ coverage
- [ ] Handler tests: 70%+ coverage
- [ ] Integration test setup with testcontainers

### Documentation
- [ ] CLAUDE.md files in place
- [ ] API documentation
- [ ] README with setup instructions

---

## Next Phase Preview

**Phase 2: Asset Management & Export (Weeks 4-5)**
- Storage interface (local + R2)
- Image upload with validation
- Image processing (resize, EXIF strip)
- Gotenberg integration for PDF/DOCX
- Export job queue
