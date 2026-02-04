# API Testing Plan

## Current State Assessment

### Test Files
**Zero test files exist.** The `api/` directory has no `*_test.go` files.

### CI/CD Configuration
**Exists but non-functional.** `.github/workflows/ci.yml` is configured but will pass vacuously since there are no tests:
- `api-test` job runs `go test -v -race -coverprofile=coverage.out ./...`
- Postgres service container is set up correctly
- Codecov integration is configured

### Code Structure Summary

| Layer | Package | Files | Testability |
|-------|---------|-------|-------------|
| Handler | `internal/handler/` | 6 files | High - needs mocking |
| Service | `internal/service/` | 3 files | High - interface-based deps |
| Repository | `internal/repository/postgres/` | 3 files | Medium - DB-dependent |
| Middleware | `internal/middleware/` | 2 files | High - pure HTTP |
| Domain | `internal/domain/` | 4 files | Highest - pure logic |
| Validator | `internal/validator/` | 1 file | High - pure validation |
| Storage | `internal/storage/` | 2 files | Medium - file-based |
| HTTP Utils | `pkg/httputil/` | 1 file | High - pure functions |

---

## Testing Strategy

### Priority Order

1. **Domain layer** (pure unit tests, no mocks needed)
2. **Service layer** (mock repositories)
3. **Handler layer** (mock services, use httptest)
4. **Middleware** (httptest with mock handlers)
5. **Repository layer** (integration tests with test DB)

### Test Types Needed

| Type | Purpose | Where |
|------|---------|-------|
| Unit | Business logic, validation, error mapping | Domain, Service, Handler |
| Integration | Database queries work correctly | Repository |
| HTTP | Endpoint contracts, error responses | Handler (via httptest) |

---

## Phase 1: Foundation (Priority: P0)

### 1.1 Domain Tests
**File:** `api/internal/domain/user_test.go`

| Test Case | Input | Expected |
|-----------|-------|----------|
| `TestUser_CanCreateCV_BelowLimit` | count=0, limit=1 | true |
| `TestUser_CanCreateCV_AtLimit` | count=1, limit=1 | false |
| `TestUser_CanUploadAsset_HasSpace` | used=0, limit=5MB, size=1MB | true |
| `TestUser_CanUploadAsset_NoSpace` | used=5MB, limit=5MB, size=1MB | false |
| `TestUser_RemainingStorage` | used=3MB, limit=5MB | 2MB |

**File:** `api/internal/domain/cv_test.go`

| Test Case | Input | Expected |
|-----------|-------|----------|
| `TestCV_ParseContent_Valid` | valid JSON | CVContent struct |
| `TestCV_ParseContent_Invalid` | malformed JSON | error |
| `TestCV_SetContent` | CVContent struct | marshaled JSON |

**File:** `api/internal/domain/errors_test.go`

| Test Case | Input | Expected |
|-----------|-------|----------|
| `TestValidationError_Error` | field="email", msg="invalid" | "invalid" |
| `TestValidationError_Unwrap` | ValidationError | ErrValidation |

### 1.2 Validator Tests
**File:** `api/internal/validator/cv_schema_test.go`

| Test Case | Input | Expected |
|-----------|-------|----------|
| `TestCVValidator_Validate_Valid` | minimal valid CV JSON | nil |
| `TestCVValidator_Validate_MissingRequired` | missing templateId | ValidationError |
| `TestCVValidator_Validate_InvalidJSON` | `{broken` | ValidationError |
| `TestValidateTemplateID_Valid` | "classic" | nil |
| `TestValidateTemplateID_Invalid` | "unknown" | ValidationError |
| `TestValidateLocale_Valid` | "en" | nil |
| `TestValidateLocale_Invalid` | "fr" | ValidationError |

### 1.3 HTTP Utils Tests
**File:** `api/pkg/httputil/response_test.go`

| Test Case | Input | Expected |
|-----------|-------|----------|
| `TestJSON_Success` | status=200, data | wrapped JSON response |
| `TestError_NotFound` | status=404, msg | error JSON with code |
| `TestHandleError_ErrNotFound` | domain.ErrNotFound | 404 response |
| `TestHandleError_ErrForbidden` | domain.ErrForbidden | 403 response |
| `TestHandleError_ValidationError` | ValidationError | 400 with field |
| `TestDecodeJSON_Valid` | valid JSON body | decoded struct |
| `TestDecodeJSON_Invalid` | malformed JSON | error |
| `TestDecodeJSON_UnknownFields` | extra fields | error (DisallowUnknownFields) |

---

## Phase 2: Service Layer (Priority: P0)

### 2.1 Test Infrastructure
**File:** `api/internal/service/mocks_test.go`

Create mock implementations:
- `MockUserRepository`
- `MockCVRepository`
- `MockAssetRepository`
- `MockStorage`

### 2.2 User Service Tests
**File:** `api/internal/service/user_service_test.go`

| Test Case | Setup | Expected |
|-----------|-------|----------|
| `TestUserService_GetByID_Success` | mock returns user | user, nil |
| `TestUserService_GetByID_NotFound` | mock returns ErrNotFound | nil, ErrNotFound |
| `TestUserService_Update_DisplayName` | mock returns user | updated displayName |
| `TestUserService_Update_InvalidLocale` | locale="fr" | ValidationError |
| `TestUserService_GetOrCreate_Exists` | user exists | existing user |
| `TestUserService_GetOrCreate_Creates` | user doesn't exist | new user created |

### 2.3 CV Service Tests
**File:** `api/internal/service/cv_service_test.go`

| Test Case | Setup | Expected |
|-----------|-------|----------|
| `TestCVService_GetByID_Success` | mock returns CV, userID matches | CV, nil |
| `TestCVService_GetByID_NotFound` | mock returns ErrNotFound | nil, ErrNotFound |
| `TestCVService_GetByID_WrongUser` | CV belongs to different user | nil, ErrForbidden |
| `TestCVService_Create_Success` | under limit, valid content | CV, nil |
| `TestCVService_Create_LimitReached` | at CV limit | nil, ErrCVLimitReached |
| `TestCVService_Create_InvalidSchema` | invalid JSON schema | nil, ValidationError |
| `TestCVService_Update_TitleOnly` | only title provided | updated title |
| `TestCVService_Update_ContentOnly` | only content provided | updated content |
| `TestCVService_Update_WrongUser` | CV belongs to different user | ErrForbidden |
| `TestCVService_Delete_Success` | CV exists, user owns it | nil |
| `TestCVService_Delete_WrongUser` | CV belongs to different user | ErrForbidden |

### 2.4 Asset Service Tests
**File:** `api/internal/service/asset_service_test.go`

| Test Case | Setup | Expected |
|-----------|-------|----------|
| `TestAssetService_Upload_Success` | valid image, under limits | Asset, nil |
| `TestAssetService_Upload_FileTooLarge` | exceeds max size | nil, ErrFileTooLarge |
| `TestAssetService_Upload_StorageFull` | user at storage limit | nil, ErrStorageLimitReached |
| `TestAssetService_Upload_InvalidType` | non-image MIME type | nil, ErrInvalidFileType |
| `TestAssetService_Upload_Duplicate` | same checksum exists | existing Asset |
| `TestAssetService_GetByID_Success` | exists, user owns it | Asset, nil |
| `TestAssetService_GetByID_WrongUser` | belongs to different user | nil, ErrForbidden |
| `TestAssetService_Delete_Success` | exists, user owns it | nil |
| `TestAssetService_Delete_WrongUser` | belongs to different user | ErrForbidden |

---

## Phase 3: Handler Layer (Priority: P1)

### 3.1 Test Infrastructure
**File:** `api/internal/handler/testutil_test.go`

Create test helpers:
- `MockUserService`
- `MockCVService`
- `MockAssetService`
- `newTestHandler()` factory
- `newAuthenticatedRequest()` helper (sets userID in context)

### 3.2 User Handler Tests
**File:** `api/internal/handler/user_test.go`

| Test Case | Request | Expected |
|-----------|---------|----------|
| `TestGetCurrentUser_Success` | GET /users/me | 200, user JSON |
| `TestGetCurrentUser_NotFound` | service returns NotFound | 404 |
| `TestUpdateCurrentUser_Success` | PATCH /users/me | 200, updated user |
| `TestUpdateCurrentUser_InvalidBody` | malformed JSON | 400 |

### 3.3 CV Handler Tests
**File:** `api/internal/handler/cv_test.go`

| Test Case | Request | Expected |
|-----------|---------|----------|
| `TestGetCV_Success` | GET /cv | 200, CV JSON |
| `TestGetCV_NotFound` | service returns NotFound | 404 |
| `TestCreateCV_Success` | POST /cv with valid body | 201, CV JSON |
| `TestCreateCV_MissingTitle` | POST /cv without title | 400, validation error |
| `TestCreateCV_MissingContent` | POST /cv without content | 400, validation error |
| `TestCreateCV_InvalidBody` | malformed JSON | 400 |
| `TestCreateCV_LimitReached` | service returns CVLimitReached | 403 |
| `TestUpdateCV_Success` | PUT /cv/{id} | 200, updated CV |
| `TestUpdateCV_InvalidID` | PUT /cv/not-a-uuid | 400 |
| `TestUpdateCV_NotFound` | service returns NotFound | 404 |
| `TestUpdateCV_Forbidden` | service returns Forbidden | 403 |
| `TestDeleteCV_Success` | DELETE /cv/{id} | 204 |
| `TestDeleteCV_InvalidID` | DELETE /cv/not-a-uuid | 400 |
| `TestDeleteCV_NotFound` | service returns NotFound | 404 |

### 3.4 Asset Handler Tests
**File:** `api/internal/handler/asset_test.go`

| Test Case | Request | Expected |
|-----------|---------|----------|
| `TestUploadAsset_Success` | POST multipart with image | 201, Asset JSON |
| `TestUploadAsset_NoFile` | POST without file | 400 |
| `TestUploadAsset_TooLarge` | service returns FileTooLarge | 400 |
| `TestUploadAsset_InvalidType` | service returns InvalidFileType | 400 |
| `TestGetAsset_Metadata` | GET /assets/{id} | 200, Asset JSON |
| `TestGetAsset_FileContent` | GET /assets/{id} Accept: image/* | 200, binary |
| `TestGetAsset_NotFound` | service returns NotFound | 404 |
| `TestDeleteAsset_Success` | DELETE /assets/{id} | 204 |
| `TestDeleteAsset_NotFound` | service returns NotFound | 404 |

### 3.5 Health Handler Tests
**File:** `api/internal/handler/health_test.go`

| Test Case | Request | Expected |
|-----------|---------|----------|
| `TestHealth` | GET /health | 200, {"status": "healthy"} |
| `TestReady_AllHealthy` | GET /ready | 200, {"status": "ready"} |

---

## Phase 4: Middleware Tests (Priority: P1)

### 4.1 Auth Middleware Tests
**File:** `api/internal/middleware/auth_test.go`

| Test Case | Setup | Expected |
|-----------|-------|----------|
| `TestAuth_BypassEnabled` | bypass=true | userID in context |
| `TestAuth_ValidJWT` | valid Bearer token | userID from sub claim |
| `TestAuth_MissingHeader` | no Authorization header | 401 |
| `TestAuth_InvalidFormat` | "Basic token" | 401 |
| `TestAuth_InvalidToken` | malformed JWT | 401 |
| `TestAuth_ExpiredToken` | expired JWT | 401 |
| `TestAuth_WrongSigningMethod` | RS256 instead of HS256 | 401 |
| `TestGetUserID_Present` | userID in context | userID |
| `TestGetUserID_Missing` | no userID in context | "" |
| `TestMustGetUserID_Present` | userID in context | userID |
| `TestMustGetUserID_Missing` | no userID in context | panic |

---

## Phase 5: Repository Integration Tests (Priority: P2)

### 5.1 Test Infrastructure
**File:** `api/internal/repository/postgres/testutil_test.go`

- Setup test database connection
- Run migrations in test setup
- Cleanup after each test
- Transaction-based isolation (rollback after each test)

### 5.2 User Repository Tests
**File:** `api/internal/repository/postgres/user_test.go`

| Test Case | Operation | Expected |
|-----------|-----------|----------|
| `TestUserRepo_Create_Success` | Create user | user in DB |
| `TestUserRepo_GetByID_Success` | Get existing user | user returned |
| `TestUserRepo_GetByID_NotFound` | Get non-existent ID | ErrNotFound |
| `TestUserRepo_GetByEmail_Success` | Get by email | user returned |
| `TestUserRepo_Update_Success` | Update user | changes persisted |
| `TestUserRepo_UpdateStorageUsed` | Increment storage | new value |

### 5.3 CV Repository Tests
**File:** `api/internal/repository/postgres/cv_test.go`

| Test Case | Operation | Expected |
|-----------|-----------|----------|
| `TestCVRepo_Create_Success` | Create CV | CV in DB |
| `TestCVRepo_GetByID_Success` | Get existing CV | CV returned |
| `TestCVRepo_GetByID_NotFound` | Get non-existent ID | ErrNotFound |
| `TestCVRepo_GetByUserID_Success` | Get user's CV | CV returned |
| `TestCVRepo_CountByUserID` | Count user's CVs | correct count |
| `TestCVRepo_Update_Success` | Update CV | changes persisted |
| `TestCVRepo_Delete_Success` | Delete CV | CV removed |

### 5.4 Asset Repository Tests
**File:** `api/internal/repository/postgres/asset_test.go`

| Test Case | Operation | Expected |
|-----------|-----------|----------|
| `TestAssetRepo_Create_Success` | Create asset | asset in DB |
| `TestAssetRepo_GetByID_Success` | Get existing | asset returned |
| `TestAssetRepo_GetByChecksum` | Find duplicate | asset returned |
| `TestAssetRepo_ListByUserID` | List user's assets | array of assets |
| `TestAssetRepo_GetTotalSizeByUserID` | Sum storage | correct total |
| `TestAssetRepo_Delete_Success` | Delete asset | asset removed |

---

## Phase 6: Storage Tests (Priority: P2)

### 6.1 LocalStorage Tests
**File:** `api/internal/storage/storage_test.go`

| Test Case | Operation | Expected |
|-----------|-----------|----------|
| `TestLocalStorage_Put_Success` | Store file | file on disk |
| `TestLocalStorage_Put_CreatesDir` | Store in new path | dir created |
| `TestLocalStorage_Get_Success` | Retrieve file | file content |
| `TestLocalStorage_Get_NotFound` | Non-existent file | error |
| `TestLocalStorage_Delete_Success` | Delete file | file removed |
| `TestLocalStorage_Delete_NotFound` | Delete non-existent | no error |
| `TestLocalStorage_Exists_True` | File exists | true |
| `TestLocalStorage_Exists_False` | File doesn't exist | false |

---

## CI/CD Improvements

### Current State
The existing `ci.yml` is mostly correct. Needed changes:

### Recommended Updates

1. **Add go vet check** to `api-lint` job:
```yaml
- name: Run go vet
  working-directory: api
  run: go vet ./...
```

2. **Add coverage threshold** (after tests exist):
```yaml
- name: Check coverage threshold
  run: |
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$coverage < 60" | bc -l) )); then
      echo "Coverage $coverage% is below 60% threshold"
      exit 1
    fi
```

3. **Run migrations before tests** in `api-test` job:
```yaml
- name: Run migrations
  working-directory: api
  env:
    DATABASE_URL: postgres://hireme:test@localhost:5432/hireme_test?sslmode=disable
  run: |
    go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    migrate -path db/migrations -database "$DATABASE_URL" up
```

---

## Test Infrastructure Needs

### 1. Mock Generator (Optional)
Consider adding `mockgen` for interface mocking:
```bash
go install go.uber.org/mock/mockgen@latest
```

### 2. Test Fixtures
Create `api/internal/testdata/` directory:
- `valid_cv.json` - minimal valid CV content
- `invalid_cv.json` - schema-violating CV
- `test_image.png` - 1x1 pixel test image

### 3. Test Helpers Package
**File:** `api/internal/testutil/testutil.go`
- `NewTestUser()` - creates domain.User with defaults
- `NewTestCV()` - creates domain.CV with defaults
- `NewTestAsset()` - creates domain.Asset with defaults

---

## Implementation Order

### Week 1: Foundation
- [ ] Domain layer tests (`domain/*_test.go`)
- [ ] Validator tests (`validator/cv_schema_test.go`)
- [ ] HTTP utils tests (`pkg/httputil/response_test.go`)

### Week 2: Services
- [ ] Create mock implementations
- [ ] User service tests
- [ ] CV service tests
- [ ] Asset service tests

### Week 3: Handlers + Middleware
- [ ] Create handler test utilities
- [ ] User handler tests
- [ ] CV handler tests
- [ ] Asset handler tests
- [ ] Auth middleware tests

### Week 4: Integration
- [ ] Repository integration tests (with test DB)
- [ ] Storage tests
- [ ] CI/CD improvements
- [ ] Coverage reporting

---

## Success Criteria

| Metric | Target |
|--------|--------|
| Unit test coverage | > 70% |
| Service layer coverage | > 80% |
| Handler layer coverage | > 75% |
| All CI jobs passing | Yes |
| Test execution time | < 30 seconds (unit), < 2 min (integration) |

---

## Out of Scope

- End-to-end API tests (curl-based)
- Performance/load testing
- Frontend tests
- Export endpoint tests (endpoint not implemented)
- R2 storage tests (not implemented)
