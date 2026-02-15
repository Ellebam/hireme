package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/config"
	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/middleware"
	"github.com/ellebam/hireme/api/pkg/httputil"
)

// UserServiceInterface defines the user service methods used by handlers
type UserServiceInterface interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, displayName, locale *string) (*domain.User, error)
}

// CVServiceInterface defines the CV service methods used by handlers
type CVServiceInterface interface {
	GetByUserID(ctx context.Context, userID string) (*domain.CV, error)
	Create(ctx context.Context, userID, title string, content json.RawMessage) (*domain.CV, error)
	Update(ctx context.Context, id uuid.UUID, userID string, title *string, content *json.RawMessage) (*domain.CV, error)
	Delete(ctx context.Context, id uuid.UUID, userID string) error
}

// AssetServiceInterface defines the asset service methods used by handlers
type AssetServiceInterface interface {
	Upload(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error)
	GetByID(ctx context.Context, id uuid.UUID, userID string) (*domain.Asset, error)
	GetFileContent(ctx context.Context, asset *domain.Asset) ([]byte, error)
	Delete(ctx context.Context, id uuid.UUID, userID string) error
}

// ExportServiceInterface defines the export service methods used by handlers
type ExportServiceInterface interface {
	ExportPDF(ctx context.Context, userID string) ([]byte, error)
	ExportDOCX(ctx context.Context, userID string) ([]byte, error)
}

// MockUserService is a mock implementation of UserServiceInterface
type MockUserService struct {
	GetByIDFunc func(ctx context.Context, id string) (*domain.User, error)
	UpdateFunc  func(ctx context.Context, id string, displayName, locale *string) (*domain.User, error)
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockUserService) Update(ctx context.Context, id string, displayName, locale *string) (*domain.User, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, displayName, locale)
	}
	return nil, domain.ErrNotFound
}

// MockCVService is a mock implementation of CVServiceInterface
type MockCVService struct {
	GetByUserIDFunc func(ctx context.Context, userID string) (*domain.CV, error)
	CreateFunc      func(ctx context.Context, userID, title string, content json.RawMessage) (*domain.CV, error)
	UpdateFunc      func(ctx context.Context, id uuid.UUID, userID string, title *string, content *json.RawMessage) (*domain.CV, error)
	DeleteFunc      func(ctx context.Context, id uuid.UUID, userID string) error
}

func (m *MockCVService) GetByUserID(ctx context.Context, userID string) (*domain.CV, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}
	return nil, domain.ErrNotFound
}

func (m *MockCVService) Create(ctx context.Context, userID, title string, content json.RawMessage) (*domain.CV, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, userID, title, content)
	}
	return nil, domain.ErrInternal
}

func (m *MockCVService) Update(ctx context.Context, id uuid.UUID, userID string, title *string, content *json.RawMessage) (*domain.CV, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, userID, title, content)
	}
	return nil, domain.ErrNotFound
}

func (m *MockCVService) Delete(ctx context.Context, id uuid.UUID, userID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, userID)
	}
	return nil
}

// MockAssetService is a mock implementation of AssetServiceInterface
type MockAssetService struct {
	UploadFunc         func(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error)
	GetByIDFunc        func(ctx context.Context, id uuid.UUID, userID string) (*domain.Asset, error)
	GetFileContentFunc func(ctx context.Context, asset *domain.Asset) ([]byte, error)
	DeleteFunc         func(ctx context.Context, id uuid.UUID, userID string) error
}

func (m *MockAssetService) Upload(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(ctx, userID, file, header)
	}
	return nil, nil
}

func (m *MockAssetService) GetByID(ctx context.Context, id uuid.UUID, userID string) (*domain.Asset, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id, userID)
	}
	return nil, domain.ErrNotFound
}

func (m *MockAssetService) GetFileContent(ctx context.Context, asset *domain.Asset) ([]byte, error) {
	if m.GetFileContentFunc != nil {
		return m.GetFileContentFunc(ctx, asset)
	}
	return nil, domain.ErrNotFound
}

func (m *MockAssetService) Delete(ctx context.Context, id uuid.UUID, userID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, userID)
	}
	return nil
}

// MockExportService is a mock implementation of ExportServiceInterface
type MockExportService struct {
	ExportPDFFunc  func(ctx context.Context, userID string) ([]byte, error)
	ExportDOCXFunc func(ctx context.Context, userID string) ([]byte, error)
}

func (m *MockExportService) ExportPDF(ctx context.Context, userID string) ([]byte, error) {
	if m.ExportPDFFunc != nil {
		return m.ExportPDFFunc(ctx, userID)
	}
	return nil, domain.ErrNotFound
}

func (m *MockExportService) ExportDOCX(ctx context.Context, userID string) ([]byte, error) {
	if m.ExportDOCXFunc != nil {
		return m.ExportDOCXFunc(ctx, userID)
	}
	return nil, domain.ErrNotFound
}

// TestHandler is a handler wrapper that uses interfaces for testing
type TestHandler struct {
	userService   UserServiceInterface
	cvService     CVServiceInterface
	assetService  AssetServiceInterface
	exportService ExportServiceInterface
	config        *config.Config
}

// NewTestHandler creates a handler with mock services for testing
func NewTestHandler(userSvc UserServiceInterface, cvSvc CVServiceInterface, assetSvc AssetServiceInterface, exportSvc ExportServiceInterface) *TestHandler {
	return &TestHandler{
		userService:   userSvc,
		cvService:     cvSvc,
		assetService:  assetSvc,
		exportService: exportSvc,
		config:        &config.Config{Features: config.FeaturesConfig{EnableExportPDF: true, EnableExportDOCX: true}},
	}
}

// GetCurrentUser returns the authenticated user's profile
func (h *TestHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := UserResponse{
		ID:                user.ID,
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		Tier:              user.Tier,
		CVLimit:           user.CVLimit,
		StorageLimitBytes: user.StorageLimitBytes,
		StorageUsedBytes:  user.StorageUsedBytes,
		Locale:            user.Locale,
	}

	httputil.JSON(w, http.StatusOK, response)
}

// UpdateCurrentUser updates the authenticated user's profile
func (h *TestHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	var req UpdateUserRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Update(ctx, userID, req.DisplayName, req.Locale)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := UserResponse{
		ID:                user.ID,
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		Tier:              user.Tier,
		CVLimit:           user.CVLimit,
		StorageLimitBytes: user.StorageLimitBytes,
		StorageUsedBytes:  user.StorageUsedBytes,
		Locale:            user.Locale,
	}

	httputil.JSON(w, http.StatusOK, response)
}

// GetCV returns the authenticated user's active CV
func (h *TestHandler) GetCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	cv, err := h.cvService.GetByUserID(ctx, userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := CVResponse{
		ID:            cv.ID.String(),
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
		CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.JSON(w, http.StatusOK, response)
}

// CreateCV creates a new CV for the authenticated user
func (h *TestHandler) CreateCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	var req CreateCVRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" {
		httputil.ValidationError(w, "title", "title is required")
		return
	}
	if len(req.Content) == 0 {
		httputil.ValidationError(w, "content", "content is required")
		return
	}

	cv, err := h.cvService.Create(ctx, userID, req.Title, req.Content)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := CVResponse{
		ID:            cv.ID.String(),
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
		CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.Created(w, response)
}

// UpdateCV updates an existing CV
func (h *TestHandler) UpdateCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse CV ID from URL
	idParam := chi.URLParam(r, "id")
	cvID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid CV ID")
		return
	}

	var req UpdateCVRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cv, err := h.cvService.Update(ctx, cvID, userID, req.Title, req.Content)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := CVResponse{
		ID:            cv.ID.String(),
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
		CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.JSON(w, http.StatusOK, response)
}

// DeleteCV deletes a CV
func (h *TestHandler) DeleteCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse CV ID from URL
	idParam := chi.URLParam(r, "id")
	cvID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid CV ID")
		return
	}

	if err := h.cvService.Delete(ctx, cvID, userID); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.NoContent(w)
}

// UploadAsset handles file uploads
func (h *TestHandler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputil.Error(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "file is required")
		return
	}
	defer func() { _ = file.Close() }()

	// Upload via service
	asset, err := h.assetService.Upload(ctx, userID, file, header)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := AssetResponse{
		ID:               asset.ID.String(),
		Filename:         asset.Filename,
		OriginalFilename: asset.OriginalFilename,
		MimeType:         asset.MimeType,
		SizeBytes:        asset.SizeBytes,
		Width:            asset.Width,
		Height:           asset.Height,
		URL:              "/api/v1/assets/" + asset.ID.String(),
		CreatedAt:        asset.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.Created(w, response)
}

// GetAsset retrieves an asset (serves the file or returns metadata)
func (h *TestHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse asset ID from URL
	idParam := chi.URLParam(r, "id")
	assetID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	// Get asset metadata
	asset, err := h.assetService.GetByID(ctx, assetID, userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Check if client wants file content (explicit image accept header) or metadata (default)
	accept := r.Header.Get("Accept")
	wantsFile := accept == asset.MimeType || accept == "image/*"

	if !wantsFile {
		// Default: return JSON metadata
		response := AssetResponse{
			ID:               asset.ID.String(),
			Filename:         asset.Filename,
			OriginalFilename: asset.OriginalFilename,
			MimeType:         asset.MimeType,
			SizeBytes:        asset.SizeBytes,
			Width:            asset.Width,
			Height:           asset.Height,
			URL:              "/api/v1/assets/" + asset.ID.String(),
			CreatedAt:        asset.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		httputil.JSON(w, http.StatusOK, response)
		return
	}

	// Serve the file (only when explicitly requesting image content)
	fileContent, err := h.assetService.GetFileContent(ctx, asset)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", asset.MimeType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+asset.Filename+"\"")
	_, _ = w.Write(fileContent)
}

// DeleteAsset deletes an asset
func (h *TestHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse asset ID from URL
	idParam := chi.URLParam(r, "id")
	assetID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	if err := h.assetService.Delete(ctx, assetID, userID); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.NoContent(w)
}

// Health returns the application health status (no service dependency)
func (h *TestHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	}
	httputil.JSON(w, http.StatusOK, response)
}

// Ready returns the application readiness status (no service dependency)
func (h *TestHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response := ReadyResponse{
		Status: "ready",
		Services: map[string]string{
			"database": "ok",
			"storage":  "ok",
		},
	}
	httputil.JSON(w, http.StatusOK, response)
}

// CreateExport handles export requests in tests
func (h *TestHandler) CreateExport(w http.ResponseWriter, r *http.Request) {
	format := chi.URLParam(r, "format")

	if !domain.IsValidExportFormat(format) {
		httputil.Error(w, http.StatusBadRequest, "invalid export format")
		return
	}

	switch format {
	case domain.ExportFormatPDF:
		if !h.config.Features.EnableExportPDF {
			httputil.Error(w, http.StatusNotImplemented, "PDF export is not enabled")
			return
		}

		ctx := r.Context()
		userID := middleware.MustGetUserID(ctx)

		pdfBytes, err := h.exportService.ExportPDF(ctx, userID)
		if err != nil {
			httputil.HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="export.pdf"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdfBytes)
	case domain.ExportFormatDOCX:
		if !h.config.Features.EnableExportDOCX {
			httputil.Error(w, http.StatusNotImplemented, "DOCX export is not enabled")
			return
		}

		ctx := r.Context()
		userID := middleware.MustGetUserID(ctx)

		docxBytes, err := h.exportService.ExportDOCX(ctx, userID)
		if err != nil {
			httputil.HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		w.Header().Set("Content-Disposition", `attachment; filename="export.docx"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(docxBytes)
	default:
		httputil.Error(w, http.StatusNotImplemented, "export format not yet implemented")
	}
}

// newAuthenticatedRequest creates an HTTP request with userID set in context
func newAuthenticatedRequest(method, url string, userID string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, url, body)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

// newAuthenticatedRequestWithParams creates an HTTP request with userID and chi URL params
func newAuthenticatedRequestWithParams(method, url string, userID string, body io.Reader, params map[string]string) *http.Request {
	req := newAuthenticatedRequest(method, url, userID, body)

	// Set chi URL params using chi's context
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	return req.WithContext(ctx)
}

// parseJSONResponse parses an httputil.Response from the response body
func parseJSONResponse(body io.Reader) (*httputil.Response, error) {
	var resp httputil.Response
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// parseUserResponse extracts UserResponse from httputil.Response.Data
func parseUserResponse(resp *httputil.Response) (*UserResponse, error) {
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	var user UserResponse
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// parseCVResponse extracts CVResponse from httputil.Response.Data
func parseCVResponse(resp *httputil.Response) (*CVResponse, error) {
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	var cv CVResponse
	if err := json.Unmarshal(data, &cv); err != nil {
		return nil, err
	}
	return &cv, nil
}

// parseAssetResponse extracts AssetResponse from httputil.Response.Data
func parseAssetResponse(resp *httputil.Response) (*AssetResponse, error) {
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	var asset AssetResponse
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, err
	}
	return &asset, nil
}

// parseHealthResponse extracts HealthResponse from httputil.Response.Data
func parseHealthResponse(resp *httputil.Response) (*HealthResponse, error) {
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	var health HealthResponse
	if err := json.Unmarshal(data, &health); err != nil {
		return nil, err
	}
	return &health, nil
}

// parseReadyResponse extracts ReadyResponse from httputil.Response.Data
func parseReadyResponse(resp *httputil.Response) (*ReadyResponse, error) {
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	var ready ReadyResponse
	if err := json.Unmarshal(data, &ready); err != nil {
		return nil, err
	}
	return &ready, nil
}

// createTestUser creates a test user with default values
func createTestUser(id string) *domain.User {
	return &domain.User{
		ID:                id,
		ExternalID:        id,
		Provider:          domain.ProviderGoogle,
		Email:             "test@example.com",
		EmailVerified:     true,
		DisplayName:       "Test User",
		Tier:              domain.TierFree,
		CVLimit:           1,
		StorageLimitBytes: 5 * 1024 * 1024,
		StorageUsedBytes:  0,
		Locale:            domain.LocaleEN,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// createTestCV creates a test CV with default values
func createTestCV(userID string) *domain.CV {
	return &domain.CV{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       json.RawMessage(`{"sections":[]}`),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// createTestAsset creates a test asset with default values
func createTestAsset(userID string) *domain.Asset {
	width := 100
	height := 100
	return &domain.Asset{
		ID:               uuid.New(),
		UserID:           userID,
		Filename:         "test-file.jpg",
		OriginalFilename: "photo.jpg",
		MimeType:         "image/jpeg",
		SizeBytes:        1024,
		StoragePath:      userID + "/2024-01/test-file.jpg",
		StorageBackend:   domain.StorageBackendLocal,
		Checksum:         "abc123",
		Width:            &width,
		Height:           &height,
		CreatedAt:        time.Now(),
	}
}

// createMultipartRequest creates a multipart form request with a file
func createMultipartRequest(url, userID, fieldName, fileName string, fileContent []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileContent); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := newAuthenticatedRequest("POST", url, userID, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// jsonBody creates a JSON request body from a value
func jsonBody(v any) io.Reader {
	data, _ := json.Marshal(v)
	return strings.NewReader(string(data))
}
