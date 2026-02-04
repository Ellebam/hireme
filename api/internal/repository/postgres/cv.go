package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/repository"
	"github.com/ellebam/hireme/api/internal/repository/postgres/queries"
)

// CVRepository implements repository.CVRepository using PostgreSQL
type CVRepository struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

// NewCVRepository creates a new CVRepository
func NewCVRepository(db *pgxpool.Pool) repository.CVRepository {
	return &CVRepository{
		db: db,
		q:  queries.New(db),
	}
}

func (r *CVRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
	cv, err := r.q.GetCVByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return cvToDomain(cv), nil
}

func (r *CVRepository) GetByUserID(ctx context.Context, userID string) (*domain.CV, error) {
	cv, err := r.q.GetCVByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return cvToDomain(cv), nil
}

func (r *CVRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.CV, error) {
	cvs, err := r.q.ListCVsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.CV, len(cvs))
	for i, cv := range cvs {
		result[i] = cvToDomain(cv)
	}
	return result, nil
}

func (r *CVRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	count, err := r.q.CountCVsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *CVRepository) Create(ctx context.Context, cv *domain.CV) error {
	created, err := r.q.CreateCV(ctx, queries.CreateCVParams{
		UserID:        cv.UserID,
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
	})
	if err != nil {
		return err
	}

	// Update the cv with the created values
	*cv = *cvToDomain(created)
	return nil
}

func (r *CVRepository) Update(ctx context.Context, cv *domain.CV) error {
	updated, err := r.q.UpdateCV(ctx, queries.UpdateCVParams{
		ID:      cv.ID,
		Title:   cv.Title,
		Content: cv.Content,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}

	// Update the cv with the updated values
	*cv = *cvToDomain(updated)
	return nil
}

func (r *CVRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteCV(ctx, id)
}

// cvToDomain converts a sqlc Cv to a domain CV
func cvToDomain(c queries.Cv) *domain.CV {
	cv := &domain.CV{
		ID:            c.ID,
		UserID:        c.UserID,
		Title:         c.Title,
		SchemaVersion: c.SchemaVersion,
		Content:       c.Content,
	}

	if c.IsActive != nil {
		cv.IsActive = *c.IsActive
	}
	if c.CreatedAt.Valid {
		cv.CreatedAt = c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		cv.UpdatedAt = c.UpdatedAt.Time
	}

	return cv
}
