package service

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	GetByIDFunc            func(ctx context.Context, id string) (*domain.User, error)
	GetByExternalIDFunc    func(ctx context.Context, provider, externalID string) (*domain.User, error)
	GetByEmailFunc         func(ctx context.Context, email string) (*domain.User, error)
	CreateFunc             func(ctx context.Context, user *domain.User) error
	UpdateFunc             func(ctx context.Context, user *domain.User) error
	UpdateStorageUsedFunc  func(ctx context.Context, id string, bytes int64) error
	DeleteFunc             func(ctx context.Context, id string) error
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockUserRepository) GetByExternalID(ctx context.Context, provider, externalID string) (*domain.User, error) {
	if m.GetByExternalIDFunc != nil {
		return m.GetByExternalIDFunc(ctx, provider, externalID)
	}
	return nil, domain.ErrNotFound
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, domain.ErrNotFound
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) UpdateStorageUsed(ctx context.Context, id string, bytes int64) error {
	if m.UpdateStorageUsedFunc != nil {
		return m.UpdateStorageUsedFunc(ctx, id, bytes)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockCVRepository is a mock implementation of repository.CVRepository
type MockCVRepository struct {
	GetByIDFunc       func(ctx context.Context, id uuid.UUID) (*domain.CV, error)
	GetByUserIDFunc   func(ctx context.Context, userID string) (*domain.CV, error)
	ListByUserIDFunc  func(ctx context.Context, userID string) ([]*domain.CV, error)
	CountByUserIDFunc func(ctx context.Context, userID string) (int, error)
	CreateFunc        func(ctx context.Context, cv *domain.CV) error
	UpdateFunc        func(ctx context.Context, cv *domain.CV) error
	DeleteFunc        func(ctx context.Context, id uuid.UUID) error
}

func (m *MockCVRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockCVRepository) GetByUserID(ctx context.Context, userID string) (*domain.CV, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}
	return nil, domain.ErrNotFound
}

func (m *MockCVRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.CV, error) {
	if m.ListByUserIDFunc != nil {
		return m.ListByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockCVRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	if m.CountByUserIDFunc != nil {
		return m.CountByUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockCVRepository) Create(ctx context.Context, cv *domain.CV) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, cv)
	}
	return nil
}

func (m *MockCVRepository) Update(ctx context.Context, cv *domain.CV) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, cv)
	}
	return nil
}

func (m *MockCVRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockAssetRepository is a mock implementation of repository.AssetRepository
type MockAssetRepository struct {
	GetByIDFunc            func(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	GetByChecksumFunc      func(ctx context.Context, userID, checksum string) (*domain.Asset, error)
	ListByUserIDFunc       func(ctx context.Context, userID string) ([]*domain.Asset, error)
	GetTotalSizeByUserIDFunc func(ctx context.Context, userID string) (int64, error)
	CreateFunc             func(ctx context.Context, asset *domain.Asset) error
	DeleteFunc             func(ctx context.Context, id uuid.UUID) error
}

func (m *MockAssetRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockAssetRepository) GetByChecksum(ctx context.Context, userID, checksum string) (*domain.Asset, error) {
	if m.GetByChecksumFunc != nil {
		return m.GetByChecksumFunc(ctx, userID, checksum)
	}
	return nil, domain.ErrNotFound
}

func (m *MockAssetRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Asset, error) {
	if m.ListByUserIDFunc != nil {
		return m.ListByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockAssetRepository) GetTotalSizeByUserID(ctx context.Context, userID string) (int64, error) {
	if m.GetTotalSizeByUserIDFunc != nil {
		return m.GetTotalSizeByUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockAssetRepository) Create(ctx context.Context, asset *domain.Asset) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, asset)
	}
	return nil
}

func (m *MockAssetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockStorage is a mock implementation of storage.Storage
type MockStorage struct {
	PutFunc    func(ctx context.Context, key string, reader io.Reader) (string, error)
	GetFunc    func(ctx context.Context, path string) (io.ReadCloser, error)
	DeleteFunc func(ctx context.Context, path string) error
	ExistsFunc func(ctx context.Context, path string) (bool, error)
}

func (m *MockStorage) Put(ctx context.Context, key string, reader io.Reader) (string, error) {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, key, reader)
	}
	return key, nil
}

func (m *MockStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, path)
	}
	return nil, domain.ErrNotFound
}

func (m *MockStorage) Delete(ctx context.Context, path string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, path)
	}
	return nil
}

func (m *MockStorage) Exists(ctx context.Context, path string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, path)
	}
	return false, nil
}

// MockCVValidator is a mock implementation for testing
type MockCVValidator struct {
	ValidateFunc func(content []byte) error
}

func (m *MockCVValidator) Validate(content []byte) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(content)
	}
	return nil
}
