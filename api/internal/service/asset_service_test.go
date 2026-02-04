package service

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/config"
	"github.com/ellebam/hireme/api/internal/domain"
)

// testConfig returns a config for testing
func testConfig() *config.Config {
	return &config.Config{
		Storage: config.StorageConfig{
			Backend:   "local",
			LocalPath: "/tmp/test-uploads",
		},
		Limits: config.LimitsConfig{
			MaxAssetSizeBytes: 2 * 1024 * 1024, // 2MB
		},
	}
}

// createTestPNGImage creates a simple PNG image for testing
func createTestPNGImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with a simple color
	for y := range height {
		for x := range width {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// mockMultipartFile wraps a byte slice to implement multipart.File
type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

// createMockFileHeader creates a mock multipart.FileHeader for testing
func createMockFileHeader(filename, contentType string, size int64) *multipart.FileHeader {
	header := &multipart.FileHeader{
		Filename: filename,
		Size:     size,
		Header:   make(textproto.MIMEHeader),
	}
	header.Header.Set("Content-Type", contentType)
	return header
}

func TestAssetService_Upload_Success(t *testing.T) {
	// Setup
	userID := "user-123"
	imageData := createTestPNGImage(100, 100)

	user := &domain.User{
		ID:                userID,
		StorageLimitBytes: 10 * 1024 * 1024, // 10MB
		StorageUsedBytes:  0,
	}

	var createdAsset *domain.Asset
	storagePutCalled := false

	mockAssetRepo := &MockAssetRepository{
		GetByChecksumFunc: func(ctx context.Context, uid, checksum string) (*domain.Asset, error) {
			return nil, domain.ErrNotFound // No duplicate
		},
		CreateFunc: func(ctx context.Context, asset *domain.Asset) error {
			createdAsset = asset
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateStorageUsedFunc: func(ctx context.Context, id string, bytes int64) error {
			return nil
		},
	}

	mockStorage := &MockStorage{
		PutFunc: func(ctx context.Context, key string, reader io.Reader) (string, error) {
			storagePutCalled = true
			return key, nil
		},
	}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Create mock file
	file := &mockMultipartFile{bytes.NewReader(imageData)}
	header := createMockFileHeader("test.png", "image/png", int64(len(imageData)))

	// Execute
	asset, err := svc.Upload(context.Background(), userID, file, header)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if asset == nil {
		t.Fatal("expected asset, got nil")
	}
	if !storagePutCalled {
		t.Error("expected storage Put to be called")
	}
	if createdAsset == nil {
		t.Fatal("expected Create to be called")
	}
	if createdAsset.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, createdAsset.UserID)
	}
	if createdAsset.MimeType != "image/png" {
		t.Errorf("expected MIME type 'image/png', got '%s'", createdAsset.MimeType)
	}
	if createdAsset.OriginalFilename != "test.png" {
		t.Errorf("expected original filename 'test.png', got '%s'", createdAsset.OriginalFilename)
	}
	if createdAsset.SizeBytes != int64(len(imageData)) {
		t.Errorf("expected size %d, got %d", len(imageData), createdAsset.SizeBytes)
	}
	if createdAsset.Width == nil || *createdAsset.Width != 100 {
		t.Error("expected width to be 100")
	}
	if createdAsset.Height == nil || *createdAsset.Height != 100 {
		t.Error("expected height to be 100")
	}
}

func TestAssetService_Upload_FileTooLarge(t *testing.T) {
	// Setup
	userID := "user-123"

	user := &domain.User{
		ID:                userID,
		StorageLimitBytes: 10 * 1024 * 1024,
		StorageUsedBytes:  0,
	}

	mockAssetRepo := &MockAssetRepository{}
	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	cfg.Limits.MaxAssetSizeBytes = 1024 // 1KB limit
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Create a file that exceeds the limit
	largeData := make([]byte, 2048) // 2KB
	file := &mockMultipartFile{bytes.NewReader(largeData)}
	header := createMockFileHeader("large.png", "image/png", int64(len(largeData)))

	// Execute
	asset, err := svc.Upload(context.Background(), userID, file, header)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrFileTooLarge) {
		t.Errorf("expected ErrFileTooLarge, got %v", err)
	}
	if asset != nil {
		t.Error("expected nil asset, got non-nil")
	}
}

func TestAssetService_Upload_StorageFull(t *testing.T) {
	// Setup
	userID := "user-123"
	imageData := createTestPNGImage(100, 100)

	user := &domain.User{
		ID:                userID,
		StorageLimitBytes: 1024, // 1KB limit
		StorageUsedBytes:  1000, // Already using most of it
	}

	mockAssetRepo := &MockAssetRepository{}
	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Create mock file (larger than remaining storage)
	file := &mockMultipartFile{bytes.NewReader(imageData)}
	header := createMockFileHeader("test.png", "image/png", int64(len(imageData)))

	// Execute
	asset, err := svc.Upload(context.Background(), userID, file, header)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrStorageLimitReached) {
		t.Errorf("expected ErrStorageLimitReached, got %v", err)
	}
	if asset != nil {
		t.Error("expected nil asset, got non-nil")
	}
}

func TestAssetService_Upload_InvalidType(t *testing.T) {
	// Setup
	userID := "user-123"

	user := &domain.User{
		ID:                userID,
		StorageLimitBytes: 10 * 1024 * 1024,
		StorageUsedBytes:  0,
	}

	mockAssetRepo := &MockAssetRepository{}
	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Create a text file (invalid type)
	textData := []byte("This is a text file, not an image")
	file := &mockMultipartFile{bytes.NewReader(textData)}
	header := createMockFileHeader("document.txt", "text/plain", int64(len(textData)))

	// Execute
	asset, err := svc.Upload(context.Background(), userID, file, header)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidFileType) {
		t.Errorf("expected ErrInvalidFileType, got %v", err)
	}
	if asset != nil {
		t.Error("expected nil asset, got non-nil")
	}
}

func TestAssetService_Upload_Duplicate(t *testing.T) {
	// Setup
	userID := "user-123"
	imageData := createTestPNGImage(100, 100)

	user := &domain.User{
		ID:                userID,
		StorageLimitBytes: 10 * 1024 * 1024,
		StorageUsedBytes:  0,
	}

	existingAsset := &domain.Asset{
		ID:               uuid.New(),
		UserID:           userID,
		Filename:         "existing.png",
		OriginalFilename: "original.png",
		MimeType:         "image/png",
		SizeBytes:        int64(len(imageData)),
		StoragePath:      "user-123/2024-01/existing.png",
		StorageBackend:   "local",
		Checksum:         "some-checksum", // Will be matched
		CreatedAt:        time.Now(),
	}

	mockAssetRepo := &MockAssetRepository{
		GetByChecksumFunc: func(ctx context.Context, uid, checksum string) (*domain.Asset, error) {
			// Return existing asset for any checksum
			return existingAsset, nil
		},
		CreateFunc: func(ctx context.Context, asset *domain.Asset) error {
			t.Error("Create should not be called for duplicate")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	mockStorage := &MockStorage{
		PutFunc: func(ctx context.Context, key string, reader io.Reader) (string, error) {
			t.Error("Storage Put should not be called for duplicate")
			return "", nil
		},
	}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Create mock file
	file := &mockMultipartFile{bytes.NewReader(imageData)}
	header := createMockFileHeader("test.png", "image/png", int64(len(imageData)))

	// Execute
	asset, err := svc.Upload(context.Background(), userID, file, header)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if asset == nil {
		t.Fatal("expected asset, got nil")
	}
	if asset.ID != existingAsset.ID {
		t.Errorf("expected existing asset ID %s, got %s", existingAsset.ID, asset.ID)
	}
}

func TestAssetService_GetByID_Success(t *testing.T) {
	// Setup
	assetID := uuid.New()
	userID := "user-123"

	expectedAsset := &domain.Asset{
		ID:               assetID,
		UserID:           userID,
		Filename:         "image.png",
		OriginalFilename: "my-image.png",
		MimeType:         "image/png",
		SizeBytes:        1024,
		StoragePath:      "user-123/2024-01/image.png",
		StorageBackend:   "local",
		CreatedAt:        time.Now(),
	}

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			if id == assetID {
				return expectedAsset, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute
	asset, err := svc.GetByID(context.Background(), assetID, userID)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if asset == nil {
		t.Fatal("expected asset, got nil")
	}
	if asset.ID != assetID {
		t.Errorf("expected asset ID %s, got %s", assetID, asset.ID)
	}
}

func TestAssetService_GetByID_NotFound(t *testing.T) {
	// Setup
	assetID := uuid.New()
	userID := "user-123"

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute
	asset, err := svc.GetByID(context.Background(), assetID, userID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if asset != nil {
		t.Error("expected nil asset, got non-nil")
	}
}

func TestAssetService_GetByID_WrongUser(t *testing.T) {
	// Setup
	assetID := uuid.New()
	ownerID := "owner-user"
	requestingUserID := "different-user"

	existingAsset := &domain.Asset{
		ID:               assetID,
		UserID:           ownerID,
		Filename:         "image.png",
		OriginalFilename: "my-image.png",
		MimeType:         "image/png",
		SizeBytes:        1024,
		StoragePath:      "owner-user/2024-01/image.png",
		StorageBackend:   "local",
		CreatedAt:        time.Now(),
	}

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			if id == assetID {
				return existingAsset, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute
	asset, err := svc.GetByID(context.Background(), assetID, requestingUserID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if asset != nil {
		t.Error("expected nil asset, got non-nil")
	}
}

func TestAssetService_Delete_Success(t *testing.T) {
	// Setup
	assetID := uuid.New()
	userID := "user-123"

	existingAsset := &domain.Asset{
		ID:               assetID,
		UserID:           userID,
		Filename:         "image.png",
		OriginalFilename: "my-image.png",
		MimeType:         "image/png",
		SizeBytes:        1024,
		StoragePath:      "user-123/2024-01/image.png",
		StorageBackend:   "local",
		CreatedAt:        time.Now(),
	}

	storageDeleteCalled := false
	repoDeleteCalled := false

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			if id == assetID {
				return existingAsset, nil
			}
			return nil, domain.ErrNotFound
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			repoDeleteCalled = true
			if id != assetID {
				t.Errorf("expected delete ID %s, got %s", assetID, id)
			}
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		UpdateStorageUsedFunc: func(ctx context.Context, id string, bytes int64) error {
			if id != userID {
				t.Errorf("expected user ID %s, got %s", userID, id)
			}
			if bytes != -1024 {
				t.Errorf("expected bytes -1024, got %d", bytes)
			}
			return nil
		},
	}

	mockStorage := &MockStorage{
		DeleteFunc: func(ctx context.Context, path string) error {
			storageDeleteCalled = true
			if path != existingAsset.StoragePath {
				t.Errorf("expected path %s, got %s", existingAsset.StoragePath, path)
			}
			return nil
		},
	}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute
	err := svc.Delete(context.Background(), assetID, userID)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !storageDeleteCalled {
		t.Error("expected storage Delete to be called")
	}
	if !repoDeleteCalled {
		t.Error("expected repo Delete to be called")
	}
}

func TestAssetService_Delete_WrongUser(t *testing.T) {
	// Setup
	assetID := uuid.New()
	ownerID := "owner-user"
	requestingUserID := "different-user"

	existingAsset := &domain.Asset{
		ID:               assetID,
		UserID:           ownerID,
		Filename:         "image.png",
		OriginalFilename: "my-image.png",
		MimeType:         "image/png",
		SizeBytes:        1024,
		StoragePath:      "owner-user/2024-01/image.png",
		StorageBackend:   "local",
		CreatedAt:        time.Now(),
	}

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			if id == assetID {
				return existingAsset, nil
			}
			return nil, domain.ErrNotFound
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			t.Error("Delete should not be called for wrong user")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	mockStorage := &MockStorage{
		DeleteFunc: func(ctx context.Context, path string) error {
			t.Error("Storage Delete should not be called for wrong user")
			return nil
		},
	}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute
	err := svc.Delete(context.Background(), assetID, requestingUserID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestAssetService_Delete_NotFound(t *testing.T) {
	// Setup
	assetID := uuid.New()
	userID := "user-123"

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	mockStorage := &MockStorage{}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute
	err := svc.Delete(context.Background(), assetID, userID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAssetService_Delete_StorageErrorContinues(t *testing.T) {
	// Setup
	assetID := uuid.New()
	userID := "user-123"

	existingAsset := &domain.Asset{
		ID:               assetID,
		UserID:           userID,
		Filename:         "image.png",
		OriginalFilename: "my-image.png",
		MimeType:         "image/png",
		SizeBytes:        1024,
		StoragePath:      "user-123/2024-01/image.png",
		StorageBackend:   "local",
		CreatedAt:        time.Now(),
	}

	repoDeleteCalled := false

	mockAssetRepo := &MockAssetRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
			if id == assetID {
				return existingAsset, nil
			}
			return nil, domain.ErrNotFound
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			repoDeleteCalled = true
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		UpdateStorageUsedFunc: func(ctx context.Context, id string, bytes int64) error {
			return nil
		},
	}

	mockStorage := &MockStorage{
		DeleteFunc: func(ctx context.Context, path string) error {
			// Simulate storage error
			return errors.New("storage error")
		},
	}

	cfg := testConfig()
	svc := NewAssetService(mockAssetRepo, mockUserRepo, mockStorage, cfg)

	// Execute - should succeed even if storage delete fails
	err := svc.Delete(context.Background(), assetID, userID)

	// Assert - delete should still succeed
	if err != nil {
		t.Fatalf("expected no error despite storage failure, got %v", err)
	}
	if !repoDeleteCalled {
		t.Error("expected repo Delete to be called even after storage error")
	}
}
