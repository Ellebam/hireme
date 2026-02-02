package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourusername/hireme/api/internal/domain"
	"github.com/yourusername/hireme/api/internal/repository"
)

// UserRepository implements repository.UserRepository using PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	// TODO: Implement using sqlc generated queries
	return nil, domain.ErrNotFound
}

func (r *UserRepository) GetByExternalID(ctx context.Context, provider, externalID string) (*domain.User, error) {
	// TODO: Implement using sqlc generated queries
	return nil, domain.ErrNotFound
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	// TODO: Implement using sqlc generated queries
	return nil, domain.ErrNotFound
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	// TODO: Implement using sqlc generated queries
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	// TODO: Implement using sqlc generated queries
	return nil
}

func (r *UserRepository) UpdateStorageUsed(ctx context.Context, id string, bytes int64) error {
	// TODO: Implement using sqlc generated queries
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implement using sqlc generated queries
	return nil
}
