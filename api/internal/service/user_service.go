package service

import (
	"context"

	"github.com/yourusername/hireme/api/internal/domain"
	"github.com/yourusername/hireme/api/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByExternalID retrieves a user by provider and external ID
func (s *UserService) GetByExternalID(ctx context.Context, provider, externalID string) (*domain.User, error) {
	return s.repo.GetByExternalID(ctx, provider, externalID)
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	return s.repo.Create(ctx, user)
}

// Update updates a user's profile
func (s *UserService) Update(ctx context.Context, id string, displayName, locale *string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if displayName != nil {
		user.DisplayName = *displayName
	}
	if locale != nil {
		// Validate locale
		if *locale != domain.LocaleEN && *locale != domain.LocaleDE {
			return nil, domain.NewValidationError("locale", "invalid locale, must be 'en' or 'de'")
		}
		user.Locale = *locale
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetOrCreate finds a user by external ID or creates one if not found
func (s *UserService) GetOrCreate(ctx context.Context, provider, externalID, email, displayName string) (*domain.User, error) {
	// Try to find existing user
	user, err := s.repo.GetByExternalID(ctx, provider, externalID)
	if err == nil {
		return user, nil
	}

	// Create new user if not found
	if err == domain.ErrNotFound {
		user = &domain.User{
			ID:                externalID, // Use external ID as internal ID for simplicity
			ExternalID:        externalID,
			Provider:          provider,
			Email:             email,
			EmailVerified:     true,
			DisplayName:       displayName,
			Tier:              domain.TierFree,
			CVLimit:           1,
			StorageLimitBytes: 5 * 1024 * 1024, // 5MB default
			Locale:            domain.LocaleEN,
		}

		if err := s.repo.Create(ctx, user); err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, err
}
