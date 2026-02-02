package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/yourusername/hireme/api/internal/domain"
	"github.com/yourusername/hireme/api/internal/repository"
	"github.com/yourusername/hireme/api/internal/validator"
)

// CVService handles CV-related business logic
type CVService struct {
	cvRepo    repository.CVRepository
	userRepo  repository.UserRepository
	validator *validator.CVValidator
}

// NewCVService creates a new CVService
func NewCVService(
	cvRepo repository.CVRepository,
	userRepo repository.UserRepository,
	validator *validator.CVValidator,
) *CVService {
	return &CVService{
		cvRepo:    cvRepo,
		userRepo:  userRepo,
		validator: validator,
	}
}

// GetByID retrieves a CV by ID, verifying ownership
func (s *CVService) GetByID(ctx context.Context, id uuid.UUID, userID string) (*domain.CV, error) {
	cv, err := s.cvRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if cv.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return cv, nil
}

// GetByUserID retrieves the user's active CV
func (s *CVService) GetByUserID(ctx context.Context, userID string) (*domain.CV, error) {
	return s.cvRepo.GetByUserID(ctx, userID)
}

// Create creates a new CV for a user
func (s *CVService) Create(ctx context.Context, userID, title string, content json.RawMessage) (*domain.CV, error) {
	// Get user to check limits
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check CV limit
	count, err := s.cvRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.CanCreateCV(count) {
		return nil, domain.ErrCVLimitReached
	}

	// Validate content against JSON schema
	if err := s.validator.Validate(content); err != nil {
		return nil, err
	}

	// Create CV
	cv := &domain.CV{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         title,
		SchemaVersion: "1.0.0",
		Content:       content,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.cvRepo.Create(ctx, cv); err != nil {
		return nil, err
	}

	return cv, nil
}

// Update updates an existing CV
func (s *CVService) Update(ctx context.Context, id uuid.UUID, userID string, title *string, content *json.RawMessage) (*domain.CV, error) {
	// Get existing CV
	cv, err := s.cvRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if cv.UserID != userID {
		return nil, domain.ErrForbidden
	}

	// Update fields
	if title != nil {
		cv.Title = *title
	}
	if content != nil {
		// Validate new content
		if err := s.validator.Validate(*content); err != nil {
			return nil, err
		}
		cv.Content = *content
	}

	cv.UpdatedAt = time.Now()

	if err := s.cvRepo.Update(ctx, cv); err != nil {
		return nil, err
	}

	return cv, nil
}

// Delete deletes a CV
func (s *CVService) Delete(ctx context.Context, id uuid.UUID, userID string) error {
	// Get existing CV
	cv, err := s.cvRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if cv.UserID != userID {
		return domain.ErrForbidden
	}

	return s.cvRepo.Delete(ctx, id)
}
