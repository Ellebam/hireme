package handler

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
)

func TestUploadAsset_Success(t *testing.T) {
	userID := "user-123"
	testAsset := createTestAsset(userID)

	mockAssetSvc := &MockAssetService{
		UploadFunc: func(ctx context.Context, uid string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error) {
			if uid != userID {
				t.Errorf("Upload called with wrong userID: got %s, want %s", uid, userID)
			}
			if header.Filename != "photo.jpg" {
				t.Errorf("Upload called with wrong filename: got %s, want photo.jpg", header.Filename)
			}
			return testAsset, nil
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	// Create multipart request with a file
	fileContent := []byte("fake image content")
	req, err := createMultipartRequest("/assets", userID, "file", "photo.jpg", fileContent)
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}
	rr := httptest.NewRecorder()

	h.UploadAsset(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error in response: %v", resp.Error)
	}

	assetResp, err := parseAssetResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse asset response: %v", err)
	}

	if assetResp.ID != testAsset.ID.String() {
		t.Errorf("expected asset ID %s, got %s", testAsset.ID.String(), assetResp.ID)
	}
	if assetResp.Filename != testAsset.Filename {
		t.Errorf("expected filename %s, got %s", testAsset.Filename, assetResp.Filename)
	}
	if assetResp.OriginalFilename != testAsset.OriginalFilename {
		t.Errorf("expected original filename %s, got %s", testAsset.OriginalFilename, assetResp.OriginalFilename)
	}
	if assetResp.MimeType != testAsset.MimeType {
		t.Errorf("expected mime type %s, got %s", testAsset.MimeType, assetResp.MimeType)
	}
	if assetResp.SizeBytes != testAsset.SizeBytes {
		t.Errorf("expected size %d, got %d", testAsset.SizeBytes, assetResp.SizeBytes)
	}
}

func TestUploadAsset_NoFile(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, nil, &MockAssetService{})

	// Create request without multipart form
	req := newAuthenticatedRequest("POST", "/assets", userID, nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	rr := httptest.NewRecorder()

	h.UploadAsset(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}
}

func TestUploadAsset_TooLarge(t *testing.T) {
	userID := "user-123"

	mockAssetSvc := &MockAssetService{
		UploadFunc: func(ctx context.Context, uid string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error) {
			return nil, domain.ErrFileTooLarge
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	fileContent := []byte("large file content")
	req, err := createMultipartRequest("/assets", userID, "file", "large.jpg", fileContent)
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}
	rr := httptest.NewRecorder()

	h.UploadAsset(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}

	if resp.Error.Code != "file_too_large" {
		t.Errorf("expected error code 'file_too_large', got '%s'", resp.Error.Code)
	}
}

func TestUploadAsset_InvalidType(t *testing.T) {
	userID := "user-123"

	mockAssetSvc := &MockAssetService{
		UploadFunc: func(ctx context.Context, uid string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error) {
			return nil, domain.ErrInvalidFileType
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	fileContent := []byte("text file content")
	req, err := createMultipartRequest("/assets", userID, "file", "document.txt", fileContent)
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}
	rr := httptest.NewRecorder()

	h.UploadAsset(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}

	if resp.Error.Code != "invalid_file_type" {
		t.Errorf("expected error code 'invalid_file_type', got '%s'", resp.Error.Code)
	}
}

func TestUploadAsset_StorageLimitReached(t *testing.T) {
	userID := "user-123"

	mockAssetSvc := &MockAssetService{
		UploadFunc: func(ctx context.Context, uid string, file multipart.File, header *multipart.FileHeader) (*domain.Asset, error) {
			return nil, domain.ErrStorageLimitReached
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	fileContent := []byte("image content")
	req, err := createMultipartRequest("/assets", userID, "file", "photo.jpg", fileContent)
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}
	rr := httptest.NewRecorder()

	h.UploadAsset(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}

	if resp.Error.Code != "storage_limit_reached" {
		t.Errorf("expected error code 'storage_limit_reached', got '%s'", resp.Error.Code)
	}
}

func TestGetAsset_Metadata(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()
	testAsset := createTestAsset(userID)
	testAsset.ID = assetID

	mockAssetSvc := &MockAssetService{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID, uid string) (*domain.Asset, error) {
			if id != assetID {
				t.Errorf("GetByID called with wrong ID: got %s, want %s", id, assetID)
			}
			if uid != userID {
				t.Errorf("GetByID called with wrong userID: got %s, want %s", uid, userID)
			}
			return testAsset, nil
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("GET", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	rr := httptest.NewRecorder()

	h.GetAsset(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error in response: %v", resp.Error)
	}

	assetResp, err := parseAssetResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse asset response: %v", err)
	}

	if assetResp.ID != testAsset.ID.String() {
		t.Errorf("expected asset ID %s, got %s", testAsset.ID.String(), assetResp.ID)
	}
	if assetResp.MimeType != testAsset.MimeType {
		t.Errorf("expected mime type %s, got %s", testAsset.MimeType, assetResp.MimeType)
	}
}

func TestGetAsset_FileContent(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()
	testAsset := createTestAsset(userID)
	testAsset.ID = assetID
	fileContent := []byte("binary image data")

	mockAssetSvc := &MockAssetService{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID, uid string) (*domain.Asset, error) {
			return testAsset, nil
		},
		GetFileContentFunc: func(ctx context.Context, asset *domain.Asset) ([]byte, error) {
			if asset.ID != assetID {
				t.Errorf("GetFileContent called with wrong asset ID: got %s, want %s", asset.ID, assetID)
			}
			return fileContent, nil
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("GET", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	req.Header.Set("Accept", "image/*") // Request file content
	rr := httptest.NewRecorder()

	h.GetAsset(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check content type header
	contentType := rr.Header().Get("Content-Type")
	if contentType != testAsset.MimeType {
		t.Errorf("expected content type %s, got %s", testAsset.MimeType, contentType)
	}

	// Check content disposition header
	contentDisposition := rr.Header().Get("Content-Disposition")
	expectedDisposition := "inline; filename=\"" + testAsset.Filename + "\""
	if contentDisposition != expectedDisposition {
		t.Errorf("expected content disposition %s, got %s", expectedDisposition, contentDisposition)
	}

	// Check body is the file content
	if rr.Body.String() != string(fileContent) {
		t.Errorf("expected body to be file content")
	}
}

func TestGetAsset_FileContentWithExactMimeType(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()
	testAsset := createTestAsset(userID)
	testAsset.ID = assetID
	fileContent := []byte("binary image data")

	mockAssetSvc := &MockAssetService{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID, uid string) (*domain.Asset, error) {
			return testAsset, nil
		},
		GetFileContentFunc: func(ctx context.Context, asset *domain.Asset) ([]byte, error) {
			return fileContent, nil
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("GET", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	req.Header.Set("Accept", "image/jpeg") // Request file content with exact MIME type
	rr := httptest.NewRecorder()

	h.GetAsset(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check content type header
	contentType := rr.Header().Get("Content-Type")
	if contentType != testAsset.MimeType {
		t.Errorf("expected content type %s, got %s", testAsset.MimeType, contentType)
	}
}

func TestGetAsset_NotFound(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()

	mockAssetSvc := &MockAssetService{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID, uid string) (*domain.Asset, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("GET", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	rr := httptest.NewRecorder()

	h.GetAsset(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestGetAsset_InvalidID(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, nil, &MockAssetService{})

	req := newAuthenticatedRequestWithParams("GET", "/assets/not-a-uuid", userID, nil, map[string]string{"id": "not-a-uuid"})
	rr := httptest.NewRecorder()

	h.GetAsset(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetAsset_Forbidden(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()

	mockAssetSvc := &MockAssetService{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID, uid string) (*domain.Asset, error) {
			return nil, domain.ErrForbidden
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("GET", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	rr := httptest.NewRecorder()

	h.GetAsset(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}

func TestDeleteAsset_Success(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()
	deleteCalled := false

	mockAssetSvc := &MockAssetService{
		DeleteFunc: func(ctx context.Context, id uuid.UUID, uid string) error {
			deleteCalled = true
			if id != assetID {
				t.Errorf("Delete called with wrong ID: got %s, want %s", id, assetID)
			}
			if uid != userID {
				t.Errorf("Delete called with wrong userID: got %s, want %s", uid, userID)
			}
			return nil
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("DELETE", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	rr := httptest.NewRecorder()

	h.DeleteAsset(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	if !deleteCalled {
		t.Error("expected Delete to be called")
	}
}

func TestDeleteAsset_NotFound(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()

	mockAssetSvc := &MockAssetService{
		DeleteFunc: func(ctx context.Context, id uuid.UUID, uid string) error {
			return domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("DELETE", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	rr := httptest.NewRecorder()

	h.DeleteAsset(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestDeleteAsset_InvalidID(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, nil, &MockAssetService{})

	req := newAuthenticatedRequestWithParams("DELETE", "/assets/not-a-uuid", userID, nil, map[string]string{"id": "not-a-uuid"})
	rr := httptest.NewRecorder()

	h.DeleteAsset(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDeleteAsset_Forbidden(t *testing.T) {
	userID := "user-123"
	assetID := uuid.New()

	mockAssetSvc := &MockAssetService{
		DeleteFunc: func(ctx context.Context, id uuid.UUID, uid string) error {
			return domain.ErrForbidden
		},
	}

	h := NewTestHandler(nil, nil, mockAssetSvc)

	req := newAuthenticatedRequestWithParams("DELETE", "/assets/"+assetID.String(), userID, nil, map[string]string{"id": assetID.String()})
	rr := httptest.NewRecorder()

	h.DeleteAsset(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}
