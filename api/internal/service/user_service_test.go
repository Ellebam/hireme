package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ellebam/hireme/api/internal/domain"
)

func TestUserService_GetByID_Success(t *testing.T) {
	// Setup
	expectedUser := &domain.User{
		ID:          "user-123",
		ExternalID:  "ext-123",
		Provider:    domain.ProviderGoogle,
		Email:       "test@example.com",
		DisplayName: "Test User",
		Tier:        domain.TierFree,
		Locale:      domain.LocaleEN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-123" {
				return expectedUser, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	user, err := svc.GetByID(context.Background(), "user-123")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID %s, got %s", expectedUser.ID, user.ID)
	}
	if user.Email != expectedUser.Email {
		t.Errorf("expected email %s, got %s", expectedUser.Email, user.Email)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	// Setup
	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	user, err := svc.GetByID(context.Background(), "nonexistent-user")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if user != nil {
		t.Error("expected nil user, got non-nil")
	}
}

func TestUserService_Update_DisplayName(t *testing.T) {
	// Setup
	existingUser := &domain.User{
		ID:          "user-123",
		ExternalID:  "ext-123",
		Provider:    domain.ProviderGoogle,
		Email:       "test@example.com",
		DisplayName: "Old Name",
		Tier:        domain.TierFree,
		Locale:      domain.LocaleEN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var updatedUser *domain.User

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-123" {
				// Return a copy to simulate DB behavior
				userCopy := *existingUser
				return &userCopy, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			updatedUser = user
			return nil
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	newName := "New Name"
	user, err := svc.Update(context.Background(), "user-123", &newName, nil)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.DisplayName != "New Name" {
		t.Errorf("expected display name 'New Name', got '%s'", user.DisplayName)
	}
	if updatedUser == nil {
		t.Fatal("expected Update to be called")
	}
	if updatedUser.DisplayName != "New Name" {
		t.Errorf("expected updated user display name 'New Name', got '%s'", updatedUser.DisplayName)
	}
}

func TestUserService_Update_InvalidLocale(t *testing.T) {
	// Setup
	existingUser := &domain.User{
		ID:          "user-123",
		ExternalID:  "ext-123",
		Provider:    domain.ProviderGoogle,
		Email:       "test@example.com",
		DisplayName: "Test User",
		Tier:        domain.TierFree,
		Locale:      domain.LocaleEN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-123" {
				userCopy := *existingUser
				return &userCopy, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			t.Error("Update should not be called for invalid locale")
			return nil
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	invalidLocale := "fr"
	user, err := svc.Update(context.Background(), "user-123", nil, &invalidLocale)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ValidationError, got %v", err)
	}
	if user != nil {
		t.Error("expected nil user, got non-nil")
	}

	// Check it's a ValidationError with the right field
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		if validationErr.Field != "locale" {
			t.Errorf("expected field 'locale', got '%s'", validationErr.Field)
		}
	} else {
		t.Error("expected error to be *domain.ValidationError")
	}
}

func TestUserService_Update_ValidLocale(t *testing.T) {
	// Setup
	existingUser := &domain.User{
		ID:          "user-123",
		ExternalID:  "ext-123",
		Provider:    domain.ProviderGoogle,
		Email:       "test@example.com",
		DisplayName: "Test User",
		Tier:        domain.TierFree,
		Locale:      domain.LocaleEN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-123" {
				userCopy := *existingUser
				return &userCopy, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	newLocale := domain.LocaleDE
	user, err := svc.Update(context.Background(), "user-123", nil, &newLocale)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Locale != domain.LocaleDE {
		t.Errorf("expected locale '%s', got '%s'", domain.LocaleDE, user.Locale)
	}
}

func TestUserService_GetOrCreate_Exists(t *testing.T) {
	// Setup
	existingUser := &domain.User{
		ID:          "user-123",
		ExternalID:  "ext-123",
		Provider:    domain.ProviderGoogle,
		Email:       "test@example.com",
		DisplayName: "Existing User",
		Tier:        domain.TierFree,
		Locale:      domain.LocaleEN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createCalled := false

	mockRepo := &MockUserRepository{
		GetByExternalIDFunc: func(ctx context.Context, provider, externalID string) (*domain.User, error) {
			if provider == domain.ProviderGoogle && externalID == "ext-123" {
				return existingUser, nil
			}
			return nil, domain.ErrNotFound
		},
		CreateFunc: func(ctx context.Context, user *domain.User) error {
			createCalled = true
			return nil
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	user, err := svc.GetOrCreate(context.Background(), domain.ProviderGoogle, "ext-123", "test@example.com", "New Name")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != existingUser.ID {
		t.Errorf("expected user ID %s, got %s", existingUser.ID, user.ID)
	}
	if user.DisplayName != "Existing User" {
		t.Errorf("expected existing display name 'Existing User', got '%s'", user.DisplayName)
	}
	if createCalled {
		t.Error("Create should not be called when user exists")
	}
}

func TestUserService_GetOrCreate_Creates(t *testing.T) {
	// Setup
	var createdUser *domain.User

	mockRepo := &MockUserRepository{
		GetByExternalIDFunc: func(ctx context.Context, provider, externalID string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
		CreateFunc: func(ctx context.Context, user *domain.User) error {
			createdUser = user
			return nil
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	user, err := svc.GetOrCreate(context.Background(), domain.ProviderGoogle, "new-ext-id", "new@example.com", "New User")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if createdUser == nil {
		t.Fatal("expected Create to be called")
	}

	// Verify the created user has correct defaults
	if createdUser.ExternalID != "new-ext-id" {
		t.Errorf("expected external ID 'new-ext-id', got '%s'", createdUser.ExternalID)
	}
	if createdUser.Provider != domain.ProviderGoogle {
		t.Errorf("expected provider '%s', got '%s'", domain.ProviderGoogle, createdUser.Provider)
	}
	if createdUser.Email != "new@example.com" {
		t.Errorf("expected email 'new@example.com', got '%s'", createdUser.Email)
	}
	if createdUser.DisplayName != "New User" {
		t.Errorf("expected display name 'New User', got '%s'", createdUser.DisplayName)
	}
	if createdUser.Tier != domain.TierFree {
		t.Errorf("expected tier '%s', got '%s'", domain.TierFree, createdUser.Tier)
	}
	if createdUser.Locale != domain.LocaleEN {
		t.Errorf("expected locale '%s', got '%s'", domain.LocaleEN, createdUser.Locale)
	}
	if createdUser.CVLimit != 1 {
		t.Errorf("expected CV limit 1, got %d", createdUser.CVLimit)
	}
	if createdUser.StorageLimitBytes != 5*1024*1024 {
		t.Errorf("expected storage limit %d, got %d", 5*1024*1024, createdUser.StorageLimitBytes)
	}
}

func TestUserService_GetOrCreate_CreateError(t *testing.T) {
	// Setup
	createError := errors.New("database error")

	mockRepo := &MockUserRepository{
		GetByExternalIDFunc: func(ctx context.Context, provider, externalID string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
		CreateFunc: func(ctx context.Context, user *domain.User) error {
			return createError
		},
	}

	svc := NewUserService(mockRepo)

	// Execute
	user, err := svc.GetOrCreate(context.Background(), domain.ProviderGoogle, "new-ext-id", "new@example.com", "New User")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, createError) {
		t.Errorf("expected create error, got %v", err)
	}
	if user != nil {
		t.Error("expected nil user, got non-nil")
	}
}
