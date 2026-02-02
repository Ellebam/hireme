package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/yourusername/hireme/api/internal/domain"
)

// UserRepository defines data access methods for users
type UserRepository interface {
	// GetByID retrieves a user by their internal ID
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetByExternalID retrieves a user by provider and external ID
	GetByExternalID(ctx context.Context, provider, externalID string) (*domain.User, error)

	// GetByEmail retrieves a user by email address
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// UpdateStorageUsed updates the user's storage usage
	UpdateStorageUsed(ctx context.Context, id string, bytes int64) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error
}

// CVRepository defines data access methods for CVs
type CVRepository interface {
	// GetByID retrieves a CV by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CV, error)

	// GetByUserID retrieves the active CV for a user
	GetByUserID(ctx context.Context, userID string) (*domain.CV, error)

	// ListByUserID retrieves all CVs for a user
	ListByUserID(ctx context.Context, userID string) ([]*domain.CV, error)

	// CountByUserID returns the number of CVs for a user
	CountByUserID(ctx context.Context, userID string) (int, error)

	// Create creates a new CV
	Create(ctx context.Context, cv *domain.CV) error

	// Update updates an existing CV
	Update(ctx context.Context, cv *domain.CV) error

	// Delete deletes a CV by ID
	Delete(ctx context.Context, id uuid.UUID) error
}

// AssetRepository defines data access methods for assets
type AssetRepository interface {
	// GetByID retrieves an asset by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error)

	// GetByChecksum retrieves an asset by user and checksum (for deduplication)
	GetByChecksum(ctx context.Context, userID, checksum string) (*domain.Asset, error)

	// ListByUserID retrieves all assets for a user
	ListByUserID(ctx context.Context, userID string) ([]*domain.Asset, error)

	// GetTotalSizeByUserID returns total storage used by a user
	GetTotalSizeByUserID(ctx context.Context, userID string) (int64, error)

	// Create creates a new asset record
	Create(ctx context.Context, asset *domain.Asset) error

	// Delete deletes an asset by ID
	Delete(ctx context.Context, id uuid.UUID) error
}

// ExportJobRepository defines data access methods for export jobs
type ExportJobRepository interface {
	// GetByID retrieves an export job by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ExportJob, error)

	// ListPending retrieves pending export jobs
	ListPending(ctx context.Context, limit int) ([]*domain.ExportJob, error)

	// Create creates a new export job
	Create(ctx context.Context, job *domain.ExportJob) error

	// UpdateStatus updates the status of an export job
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, resultPath, errorMsg *string) error
}
