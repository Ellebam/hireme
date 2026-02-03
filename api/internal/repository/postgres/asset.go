package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourusername/hireme/api/internal/domain"
	"github.com/yourusername/hireme/api/internal/repository"
	"github.com/yourusername/hireme/api/internal/repository/postgres/queries"
)

// AssetRepository implements repository.AssetRepository using PostgreSQL
type AssetRepository struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

// NewAssetRepository creates a new AssetRepository
func NewAssetRepository(db *pgxpool.Pool) repository.AssetRepository {
	return &AssetRepository{
		db: db,
		q:  queries.New(db),
	}
}

func (r *AssetRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	asset, err := r.q.GetAssetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return assetToDomain(asset), nil
}

func (r *AssetRepository) GetByChecksum(ctx context.Context, userID, checksum string) (*domain.Asset, error) {
	asset, err := r.q.GetAssetByChecksum(ctx, queries.GetAssetByChecksumParams{
		UserID:   userID,
		Checksum: checksum,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return assetToDomain(asset), nil
}

func (r *AssetRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Asset, error) {
	assets, err := r.q.ListAssetsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Asset, len(assets))
	for i, asset := range assets {
		result[i] = assetToDomain(asset)
	}
	return result, nil
}

func (r *AssetRepository) GetTotalSizeByUserID(ctx context.Context, userID string) (int64, error) {
	return r.q.GetTotalAssetSizeByUserID(ctx, userID)
}

func (r *AssetRepository) Create(ctx context.Context, asset *domain.Asset) error {
	// Convert metadata to bytes
	var metadata []byte
	if asset.Metadata != nil {
		metadata = asset.Metadata
	} else {
		metadata = []byte("{}")
	}

	// Convert dimensions
	var width, height *int32
	if asset.Width != nil {
		w := int32(*asset.Width)
		width = &w
	}
	if asset.Height != nil {
		h := int32(*asset.Height)
		height = &h
	}

	created, err := r.q.CreateAsset(ctx, queries.CreateAssetParams{
		UserID:           asset.UserID,
		Filename:         asset.Filename,
		OriginalFilename: asset.OriginalFilename,
		MimeType:         asset.MimeType,
		SizeBytes:        int32(asset.SizeBytes),
		StoragePath:      asset.StoragePath,
		StorageBackend:   asset.StorageBackend,
		Checksum:         asset.Checksum,
		Width:            width,
		Height:           height,
		Metadata:         metadata,
	})
	if err != nil {
		return err
	}

	// Update the asset with the created values
	*asset = *assetToDomain(created)
	return nil
}

func (r *AssetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteAsset(ctx, id)
}

// assetToDomain converts a sqlc Asset to a domain Asset
func assetToDomain(a queries.Asset) *domain.Asset {
	asset := &domain.Asset{
		ID:               a.ID,
		UserID:           a.UserID,
		Filename:         a.Filename,
		OriginalFilename: a.OriginalFilename,
		MimeType:         a.MimeType,
		SizeBytes:        int64(a.SizeBytes),
		StoragePath:      a.StoragePath,
		StorageBackend:   a.StorageBackend,
		Checksum:         a.Checksum,
	}

	if a.Width != nil {
		w := int(*a.Width)
		asset.Width = &w
	}
	if a.Height != nil {
		h := int(*a.Height)
		asset.Height = &h
	}
	if a.Metadata != nil {
		asset.Metadata = json.RawMessage(a.Metadata)
	}
	if a.CreatedAt.Valid {
		asset.CreatedAt = a.CreatedAt.Time
	}

	return asset
}
