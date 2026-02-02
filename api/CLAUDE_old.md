# HireMe API — Backend AI Context

## Overview

Go backend service using Chi router with a clean architecture pattern.

**Port:** 8080 (configurable via `SERVER_PORT`)

## Directory Structure

```
api/
├── cmd/server/main.go       # Entry point, DI, router setup
├── internal/
│   ├── config/              # Environment-based configuration
│   ├── domain/              # Pure domain models (NO external deps)
│   ├── handler/             # HTTP handlers (thin, validation only)
│   ├── middleware/          # Auth, logging, request ID
│   ├── repository/          # Data access interfaces + implementations
│   ├── service/             # Business logic layer
│   ├── storage/             # File storage abstraction
│   └── validator/           # Input validation, JSON Schema
├── pkg/                     # Shared utilities (can be imported)
├── db/
│   ├── migrations/          # SQL migration files
│   ├── queries/             # sqlc query definitions
│   └── sqlc.yaml            # sqlc configuration
└── templates/               # HTML templates for PDF export
```

## Code Conventions

### Handlers
```go
// Handlers are thin — validate input, call service, return response
func (h *Handler) CreateCV(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := middleware.GetUserID(ctx) // From auth middleware
    
    var req CreateCVRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    if err := h.validator.ValidateCV(req.Content); err != nil {
        httputil.Error(w, http.StatusBadRequest, err.Error())
        return
    }
    
    cv, err := h.cvService.Create(ctx, userID, req)
    if err != nil {
        // Service returns domain errors, handler maps to HTTP
        httputil.HandleError(w, err)
        return
    }
    
    httputil.JSON(w, http.StatusCreated, cv)
}
```

### Services
```go
// Services contain business logic, orchestrate repositories
func (s *CVService) Create(ctx context.Context, userID string, req CreateCVRequest) (*domain.CV, error) {
    // Check user's CV limit
    count, err := s.cvRepo.CountByUser(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("checking cv count: %w", err)
    }
    
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("getting user: %w", err)
    }
    
    if count >= user.CVLimit {
        return nil, domain.ErrCVLimitReached
    }
    
    cv := &domain.CV{
        ID:        uuid.New().String(),
        UserID:    userID,
        Content:   req.Content,
        CreatedAt: time.Now(),
    }
    
    if err := s.cvRepo.Create(ctx, cv); err != nil {
        return nil, fmt.Errorf("creating cv: %w", err)
    }
    
    return cv, nil
}
```

### Repositories
```go
// Repositories are interfaces — implementations in postgres/ or mock/
type CVRepository interface {
    Create(ctx context.Context, cv *domain.CV) error
    GetByID(ctx context.Context, id string) (*domain.CV, error)
    GetByUserID(ctx context.Context, userID string) (*domain.CV, error)
    Update(ctx context.Context, cv *domain.CV) error
    Delete(ctx context.Context, id string) error
    CountByUser(ctx context.Context, userID string) (int, error)
}
```

## Error Handling

Domain errors are defined in `internal/domain/errors.go`:

```go
var (
    ErrNotFound       = errors.New("not found")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
    ErrCVLimitReached = errors.New("cv limit reached")
    ErrInvalidInput   = errors.New("invalid input")
)
```

Handlers map these to HTTP status codes in `pkg/httputil/response.go`.

## Configuration

All config via environment variables. See `internal/config/config.go`:

```go
type Config struct {
    Environment       string // development, staging, production
    ServerPort        int
    DatabaseURL       string
    AuthBypassEnabled bool   // DEV ONLY: skip auth
    AuthBypassUserID  string // DEV ONLY: hardcoded user
    StorageBackend    string // "local" or "r2"
    StorageLocalPath  string // Path for local storage
    GotenbergURL      string // Export service URL
}
```

## Testing

### Unit Tests
```bash
task api:test              # Run all tests
task api:test:coverage     # With coverage report
```

### Test Patterns
```go
func TestCVService_Create(t *testing.T) {
    // Arrange
    mockRepo := mock.NewCVRepository()
    mockUserRepo := mock.NewUserRepository()
    mockUserRepo.Users["user-1"] = &domain.User{ID: "user-1", CVLimit: 1}
    
    svc := service.NewCVService(mockRepo, mockUserRepo)
    
    // Act
    cv, err := svc.Create(context.Background(), "user-1", CreateCVRequest{...})
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, cv.ID)
}
```

## Database

### sqlc Workflow
1. Write queries in `db/queries/*.sql`
2. Run `task api:sqlc` to generate Go code
3. Generated code appears in `internal/repository/postgres/queries/`

### Migration Commands
```bash
task db:migrate              # Apply migrations
task db:migrate:down         # Rollback one migration
task db:migration:create name=add_feature  # Create new migration
```

## API Routes

```
GET    /health                 # Health check
GET    /ready                  # Readiness check

POST   /api/v1/auth/callback   # OIDC callback
DELETE /api/v1/auth/logout     # Logout

GET    /api/v1/users/me        # Get current user
PATCH  /api/v1/users/me        # Update current user

GET    /api/v1/cv              # Get user's CV
POST   /api/v1/cv              # Create CV
PUT    /api/v1/cv/:id          # Update CV
DELETE /api/v1/cv/:id          # Delete CV

POST   /api/v1/assets          # Upload asset
GET    /api/v1/assets/:id      # Get asset
DELETE /api/v1/assets/:id      # Delete asset

POST   /api/v1/export/:format  # Trigger export (pdf, docx, json, yaml)
GET    /api/v1/export/:id      # Get export status/download
```

## Common Tasks

### Adding a New Endpoint
1. Add route in `cmd/server/main.go`
2. Create handler method in appropriate `internal/handler/*.go`
3. Add service method if business logic needed
4. Add repository method if data access needed
5. Write tests

### Adding a Database Field
1. Create migration: `task db:migration:create name=add_field_to_table`
2. Write ALTER TABLE in up.sql, reverse in down.sql
3. Update sqlc queries if the field should be selected/inserted
4. Run `task db:migrate && task api:sqlc`
5. Update domain model
6. Update service/handler as needed

## Logging

Structured logging with slog:

```go
slog.Info("cv created",
    "user_id", userID,
    "cv_id", cv.ID,
    "request_id", middleware.GetRequestID(ctx),
)
```

## Development Bypass

When `AUTH_BYPASS_ENABLED=true`:
- All requests are authenticated as `AUTH_BYPASS_USER_ID`
- No JWT validation
- OIDC endpoints still work but aren't required

Enable in `.env.local`:
```
AUTH_BYPASS_ENABLED=true
AUTH_BYPASS_USER_ID=dev-user-001
```
