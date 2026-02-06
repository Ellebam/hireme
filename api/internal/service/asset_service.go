package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/config"
	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/repository"
	"github.com/ellebam/hireme/api/internal/storage"
)

// AssetService handles asset-related business logic
type AssetService struct {
	assetRepo repository.AssetRepository
	userRepo  repository.UserRepository
	storage   storage.Storage
	config    *config.Config
}

// NewAssetService creates a new AssetService
func NewAssetService(
	assetRepo repository.AssetRepository,
	userRepo repository.UserRepository,
	store storage.Storage,
	cfg *config.Config,
) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		userRepo:  userRepo,
		storage:   store,
		config:    cfg,
	}
}

// Upload handles file upload
func (s *AssetService) Upload(
	ctx context.Context,
	userID string,
	file multipart.File,
	header *multipart.FileHeader,
) (*domain.Asset, error) {
	// 1. Get user to check limits
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Check file size
	if header.Size > s.config.Limits.MaxAssetSizeBytes {
		return nil, domain.ErrFileTooLarge
	}

	// 3. Check storage limit
	if !user.CanUploadAsset(header.Size) {
		return nil, domain.ErrStorageLimitReached
	}

	// 4. Check MIME type
	mimeType := header.Header.Get("Content-Type")
	if !domain.IsAllowedImageType(mimeType) {
		return nil, domain.ErrInvalidFileType
	}

	// 5. Read file content for checksum and dimensions
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// 6. Calculate SHA-256 checksum
	checksum := calculateChecksum(fileContent)

	// 7. Check for duplicate (same user, same checksum)
	existing, err := s.assetRepo.GetByChecksum(ctx, userID, checksum)
	if err == nil && existing != nil {
		// Return existing asset instead of uploading duplicate
		return existing, nil
	}

	// 8. Get image dimensions
	width, height := getImageDimensions(fileContent)

	// 9. Generate unique filename
	ext := domain.GetExtension(mimeType)
	filename := uuid.New().String() + ext

	// Storage path: userID/year-month/filename
	storagePath := generateStoragePath(userID, filename)

	// 10. Store file
	reader := bytes.NewReader(fileContent)
	_, err = s.storage.Put(ctx, storagePath, reader)
	if err != nil {
		return nil, fmt.Errorf("storing file: %w", err)
	}

	// 11. Create asset record
	asset := &domain.Asset{
		ID:               uuid.New(),
		UserID:           userID,
		Filename:         filename,
		OriginalFilename: header.Filename,
		MimeType:         mimeType,
		SizeBytes:        header.Size,
		StoragePath:      storagePath,
		StorageBackend:   s.config.Storage.Backend,
		Checksum:         checksum,
		Width:            width,
		Height:           height,
		CreatedAt:        time.Now(),
	}

	if err := s.assetRepo.Create(ctx, asset); err != nil {
		// Cleanup: delete from storage on DB failure
		_ = s.storage.Delete(ctx, storagePath)
		return nil, fmt.Errorf("creating asset record: %w", err)
	}

	// 12. Update user storage usage
	if err := s.userRepo.UpdateStorageUsed(ctx, userID, header.Size); err != nil {
		slog.Warn("failed to update user storage usage", "userID", userID, "error", err)
	}

	return asset, nil
}

// GetByID retrieves an asset by ID, verifying ownership
func (s *AssetService) GetByID(ctx context.Context, id uuid.UUID, userID string) (*domain.Asset, error) {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if asset.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return asset, nil
}

// GetFileContent retrieves the actual file content
func (s *AssetService) GetFileContent(ctx context.Context, asset *domain.Asset) ([]byte, error) {
	reader, err := s.storage.Get(ctx, asset.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("retrieving file: %w", err)
	}
	defer func() { _ = reader.Close() }()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading file content: %w", err)
	}

	return content, nil
}

// Delete deletes an asset
func (s *AssetService) Delete(ctx context.Context, id uuid.UUID, userID string) error {
	// 1. Get asset to verify ownership and get storage path
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Verify ownership
	if asset.UserID != userID {
		return domain.ErrForbidden
	}

	// 3. Delete from storage backend
	if err := s.storage.Delete(ctx, asset.StoragePath); err != nil {
		// Log but continue - orphaned files can be cleaned up later
		slog.Warn("failed to delete file from storage", "path", asset.StoragePath, "error", err)
	}

	// 4. Delete database record
	if err := s.assetRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 5. Update user storage usage (decrement)
	return s.userRepo.UpdateStorageUsed(ctx, userID, -asset.SizeBytes)
}

// calculateChecksum computes SHA-256 hash of content
func calculateChecksum(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// getImageDimensions extracts width/height from image data
func getImageDimensions(content []byte) (*int, *int) {
	reader := bytes.NewReader(content)
	cfg, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, nil
	}
	width := cfg.Width
	height := cfg.Height
	return &width, &height
}

// generateStoragePath creates a storage path: userID/YYYY-MM/filename
func generateStoragePath(userID, filename string) string {
	now := time.Now()
	return filepath.Join(userID, now.Format("2006-01"), filename)
}
