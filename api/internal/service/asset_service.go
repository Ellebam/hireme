package service

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/yourusername/hireme/api/internal/config"
	"github.com/yourusername/hireme/api/internal/domain"
	"github.com/yourusername/hireme/api/internal/repository"
)

// AssetService handles asset-related business logic
type AssetService struct {
	assetRepo repository.AssetRepository
	userRepo  repository.UserRepository
	config    *config.Config
}

// NewAssetService creates a new AssetService
func NewAssetService(
	assetRepo repository.AssetRepository,
	userRepo repository.UserRepository,
	cfg *config.Config,
) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		userRepo:  userRepo,
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
	// Get user to check limits
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check file size
	if header.Size > s.config.Limits.MaxAssetSizeBytes {
		return nil, domain.ErrFileTooLarge
	}

	// Check storage limit
	if !user.CanUploadAsset(header.Size) {
		return nil, domain.ErrStorageLimitReached
	}

	// Check MIME type
	mimeType := header.Header.Get("Content-Type")
	if !domain.IsAllowedImageType(mimeType) {
		return nil, domain.ErrInvalidFileType
	}

	// TODO: Implement actual file storage
	// 1. Generate unique filename
	// 2. Calculate checksum
	// 3. Store file (local or R2)
	// 4. Process image (resize, strip EXIF)
	// 5. Create asset record
	// 6. Update user storage usage

	return nil, domain.ErrInternal // Placeholder
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
	// TODO: Implement file retrieval from storage backend
	return nil, domain.ErrInternal // Placeholder
}

// Delete deletes an asset
func (s *AssetService) Delete(ctx context.Context, id uuid.UUID, userID string) error {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if asset.UserID != userID {
		return domain.ErrForbidden
	}

	// TODO: Delete file from storage backend

	// Delete record
	if err := s.assetRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Update user storage usage
	return s.userRepo.UpdateStorageUsed(ctx, userID, -asset.SizeBytes)
}
