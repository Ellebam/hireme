package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourusername/hireme/api/internal/domain"
	"github.com/yourusername/hireme/api/internal/repository"
	"github.com/yourusername/hireme/api/internal/repository/postgres/queries"
)

// UserRepository implements repository.UserRepository using PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepository{
		db: db,
		q:  queries.New(db),
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return userToDomain(user), nil
}

func (r *UserRepository) GetByExternalID(ctx context.Context, provider, externalID string) (*domain.User, error) {
	user, err := r.q.GetUserByExternalID(ctx, queries.GetUserByExternalIDParams{
		Provider:   provider,
		ExternalID: externalID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return userToDomain(user), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return userToDomain(user), nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	created, err := r.q.CreateUser(ctx, queries.CreateUserParams{
		ID:                user.ID,
		ExternalID:        user.ExternalID,
		Provider:          user.Provider,
		Email:             user.Email,
		EmailVerified:     &user.EmailVerified,
		DisplayName:       strPtr(user.DisplayName),
		Tier:              strPtr(user.Tier),
		CvLimit:           int32Ptr(int32(user.CVLimit)),
		StorageLimitBytes: &user.StorageLimitBytes,
		Locale:            strPtr(user.Locale),
	})
	if err != nil {
		return err
	}

	// Update the user with the created values
	*user = *userToDomain(created)
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	updated, err := r.q.UpdateUser(ctx, queries.UpdateUserParams{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: &user.EmailVerified,
		DisplayName:   strPtr(user.DisplayName),
		Locale:        strPtr(user.Locale),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}

	// Update the user with the updated values
	*user = *userToDomain(updated)
	return nil
}

func (r *UserRepository) UpdateStorageUsed(ctx context.Context, id string, bytes int64) error {
	return r.q.UpdateUserStorageUsed(ctx, queries.UpdateUserStorageUsedParams{
		ID:               id,
		StorageUsedBytes: &bytes,
	})
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.q.DeleteUser(ctx, id)
}

// userToDomain converts a sqlc User to a domain User
func userToDomain(u queries.User) *domain.User {
	user := &domain.User{
		ID:         u.ID,
		ExternalID: u.ExternalID,
		Provider:   u.Provider,
		Email:      u.Email,
	}

	if u.EmailVerified != nil {
		user.EmailVerified = *u.EmailVerified
	}
	if u.DisplayName != nil {
		user.DisplayName = *u.DisplayName
	}
	if u.Tier != nil {
		user.Tier = *u.Tier
	}
	if u.CvLimit != nil {
		user.CVLimit = int(*u.CvLimit)
	}
	if u.StorageLimitBytes != nil {
		user.StorageLimitBytes = *u.StorageLimitBytes
	}
	if u.StorageUsedBytes != nil {
		user.StorageUsedBytes = *u.StorageUsedBytes
	}
	if u.Locale != nil {
		user.Locale = *u.Locale
	}
	if u.CreatedAt.Valid {
		user.CreatedAt = u.CreatedAt.Time
	}
	if u.UpdatedAt.Valid {
		user.UpdatedAt = u.UpdatedAt.Time
	}

	return user
}

// Helper functions for pointer conversions
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
