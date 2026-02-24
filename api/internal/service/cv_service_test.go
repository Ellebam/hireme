package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/validator"
)

// createTestCVValidator creates a real validator for tests
func createTestCVValidator(t *testing.T) *validator.CVValidator {
	v, err := validator.NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create CV validator: %v", err)
	}
	return v
}

// validCVContent returns valid CV content for testing
func validCVContent() json.RawMessage {
	return json.RawMessage(`{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"locale": "en",
		"sections": [
			{
				"id": "personal-1",
				"type": "personal",
				"order": 1,
				"visible": true,
				"content": {
					"firstName": "John",
					"lastName": "Doe",
					"email": "john@example.com"
				}
			}
		]
	}`)
}

// invalidCVContent returns invalid CV content for testing
func invalidCVContent() json.RawMessage {
	return json.RawMessage(`{
		"invalid": "content"
	}`)
}

func TestCVService_GetByID_Success(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	expectedCV := &domain.CV{
		ID:            cvID,
		UserID:        userID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				return expectedCV, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.GetByID(context.Background(), cvID, userID)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cv == nil {
		t.Fatal("expected CV, got nil")
	}
	if cv.ID != cvID {
		t.Errorf("expected CV ID %s, got %s", cvID, cv.ID)
	}
	if cv.Title != "My CV" {
		t.Errorf("expected title 'My CV', got '%s'", cv.Title)
	}
}

func TestCVService_GetByID_NotFound(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.GetByID(context.Background(), cvID, userID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if cv != nil {
		t.Error("expected nil CV, got non-nil")
	}
}

func TestCVService_GetByID_WrongUser(t *testing.T) {
	// Setup
	cvID := uuid.New()
	ownerID := "owner-user"
	requestingUserID := "different-user"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        ownerID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				return existingCV, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.GetByID(context.Background(), cvID, requestingUserID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if cv != nil {
		t.Error("expected nil CV, got non-nil")
	}
}

func TestCVService_Create_Success(t *testing.T) {
	// Setup
	userID := "user-123"

	user := &domain.User{
		ID:       userID,
		CVLimit:  5,
		Tier:     domain.TierPro,
	}

	var createdCV *domain.CV

	mockCVRepo := &MockCVRepository{
		CountByUserIDFunc: func(ctx context.Context, uid string) (int, error) {
			return 1, nil // User has 1 CV, limit is 5
		},
		CreateFunc: func(ctx context.Context, cv *domain.CV) error {
			createdCV = cv
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	cvValidator := createTestCVValidator(t)
	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.Create(context.Background(), userID, "New CV", validCVContent())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cv == nil {
		t.Fatal("expected CV, got nil")
	}
	if createdCV == nil {
		t.Fatal("expected Create to be called")
	}
	if createdCV.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, createdCV.UserID)
	}
	if createdCV.Title != "New CV" {
		t.Errorf("expected title 'New CV', got '%s'", createdCV.Title)
	}
	if createdCV.SchemaVersion != "1.0.0" {
		t.Errorf("expected schema version '1.0.0', got '%s'", createdCV.SchemaVersion)
	}
	if !createdCV.IsActive {
		t.Error("expected CV to be active")
	}
}

func TestCVService_Create_LimitReached(t *testing.T) {
	// Setup
	userID := "user-123"

	user := &domain.User{
		ID:      userID,
		CVLimit: 1,
		Tier:    domain.TierFree,
	}

	mockCVRepo := &MockCVRepository{
		CountByUserIDFunc: func(ctx context.Context, uid string) (int, error) {
			return 1, nil // User has 1 CV, limit is 1
		},
		CreateFunc: func(ctx context.Context, cv *domain.CV) error {
			t.Error("Create should not be called when limit is reached")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	cvValidator := createTestCVValidator(t)
	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.Create(context.Background(), userID, "New CV", validCVContent())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrCVLimitReached) {
		t.Errorf("expected ErrCVLimitReached, got %v", err)
	}
	if cv != nil {
		t.Error("expected nil CV, got non-nil")
	}
}

func TestCVService_Create_InvalidSchema(t *testing.T) {
	// Setup
	userID := "user-123"

	user := &domain.User{
		ID:      userID,
		CVLimit: 5,
		Tier:    domain.TierPro,
	}

	mockCVRepo := &MockCVRepository{
		CountByUserIDFunc: func(ctx context.Context, uid string) (int, error) {
			return 0, nil
		},
		CreateFunc: func(ctx context.Context, cv *domain.CV) error {
			t.Error("Create should not be called for invalid content")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id == userID {
				return user, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	cvValidator := createTestCVValidator(t)
	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.Create(context.Background(), userID, "New CV", invalidCVContent())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ValidationError, got %v", err)
	}
	if cv != nil {
		t.Error("expected nil CV, got non-nil")
	}
}

func TestCVService_Update_TitleOnly(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        userID,
		Title:         "Old Title",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	var updatedCV *domain.CV

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				cvCopy := *existingCV
				return &cvCopy, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, cv *domain.CV) error {
			updatedCV = cv
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	newTitle := "New Title"
	cv, err := svc.Update(context.Background(), cvID, userID, &newTitle, nil)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cv == nil {
		t.Fatal("expected CV, got nil")
	}
	if cv.Title != "New Title" {
		t.Errorf("expected title 'New Title', got '%s'", cv.Title)
	}
	if updatedCV == nil {
		t.Fatal("expected Update to be called")
	}
	if updatedCV.Title != "New Title" {
		t.Errorf("expected updated title 'New Title', got '%s'", updatedCV.Title)
	}
}

func TestCVService_Update_ContentOnly(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        userID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	newContent := json.RawMessage(`{
		"schemaVersion": "1.0.0",
		"templateId": "modern",
		"locale": "de",
		"sections": [
			{
				"id": "personal-1",
				"type": "personal",
				"order": 1,
				"visible": true,
				"content": {
					"firstName": "Jane",
					"lastName": "Smith",
					"email": "jane@example.com"
				}
			}
		]
	}`)

	var updatedCV *domain.CV

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				cvCopy := *existingCV
				return &cvCopy, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, cv *domain.CV) error {
			updatedCV = cv
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	cv, err := svc.Update(context.Background(), cvID, userID, nil, &newContent)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cv == nil {
		t.Fatal("expected CV, got nil")
	}
	if updatedCV == nil {
		t.Fatal("expected Update to be called")
	}

	// Verify content was updated
	var content map[string]any
	if err := json.Unmarshal(updatedCV.Content, &content); err != nil {
		t.Fatalf("failed to unmarshal content: %v", err)
	}
	if content["templateId"] != "modern" {
		t.Errorf("expected templateId 'modern', got '%v'", content["templateId"])
	}
}

func TestCVService_Update_WrongUser(t *testing.T) {
	// Setup
	cvID := uuid.New()
	ownerID := "owner-user"
	requestingUserID := "different-user"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        ownerID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				return existingCV, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, cv *domain.CV) error {
			t.Error("Update should not be called for wrong user")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	newTitle := "Hacked Title"
	cv, err := svc.Update(context.Background(), cvID, requestingUserID, &newTitle, nil)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if cv != nil {
		t.Error("expected nil CV, got non-nil")
	}
}

func TestCVService_Update_InvalidContent(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        userID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				cvCopy := *existingCV
				return &cvCopy, nil
			}
			return nil, domain.ErrNotFound
		},
		UpdateFunc: func(ctx context.Context, cv *domain.CV) error {
			t.Error("Update should not be called for invalid content")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	invalidContent := invalidCVContent()
	cv, err := svc.Update(context.Background(), cvID, userID, nil, &invalidContent)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ValidationError, got %v", err)
	}
	if cv != nil {
		t.Error("expected nil CV, got non-nil")
	}
}

func TestCVService_Delete_Success(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        userID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	deleteCalled := false

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				return existingCV, nil
			}
			return nil, domain.ErrNotFound
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			deleteCalled = true
			if id != cvID {
				t.Errorf("expected delete ID %s, got %s", cvID, id)
			}
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	err := svc.Delete(context.Background(), cvID, userID)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !deleteCalled {
		t.Error("expected Delete to be called")
	}
}

func TestCVService_Delete_WrongUser(t *testing.T) {
	// Setup
	cvID := uuid.New()
	ownerID := "owner-user"
	requestingUserID := "different-user"

	existingCV := &domain.CV{
		ID:            cvID,
		UserID:        ownerID,
		Title:         "My CV",
		SchemaVersion: "1.0.0",
		Content:       validCVContent(),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			if id == cvID {
				return existingCV, nil
			}
			return nil, domain.ErrNotFound
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			t.Error("Delete should not be called for wrong user")
			return nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	err := svc.Delete(context.Background(), cvID, requestingUserID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestCVService_Delete_NotFound(t *testing.T) {
	// Setup
	cvID := uuid.New()
	userID := "user-123"

	mockCVRepo := &MockCVRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
			return nil, domain.ErrNotFound
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	// Execute
	err := svc.Delete(context.Background(), cvID, userID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCVService_ListByUserID_Success(t *testing.T) {
	userID := "user-123"

	cv1 := &domain.CV{
		ID:     uuid.New(),
		UserID: userID,
		Title:  "CV One",
	}
	cv2 := &domain.CV{
		ID:     uuid.New(),
		UserID: userID,
		Title:  "CV Two",
	}

	mockCVRepo := &MockCVRepository{
		ListByUserIDFunc: func(ctx context.Context, uid string) ([]*domain.CV, error) {
			if uid != userID {
				t.Errorf("expected userID %s, got %s", userID, uid)
			}
			return []*domain.CV{cv1, cv2}, nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	cvs, err := svc.ListByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cvs) != 2 {
		t.Fatalf("expected 2 CVs, got %d", len(cvs))
	}
	if cvs[0].Title != "CV One" {
		t.Errorf("expected first CV title 'CV One', got '%s'", cvs[0].Title)
	}
	if cvs[1].Title != "CV Two" {
		t.Errorf("expected second CV title 'CV Two', got '%s'", cvs[1].Title)
	}
}

func TestCVService_ListByUserID_Empty(t *testing.T) {
	userID := "user-123"

	mockCVRepo := &MockCVRepository{
		ListByUserIDFunc: func(ctx context.Context, uid string) ([]*domain.CV, error) {
			return []*domain.CV{}, nil
		},
	}

	mockUserRepo := &MockUserRepository{}
	cvValidator := createTestCVValidator(t)

	svc := NewCVService(mockCVRepo, mockUserRepo, cvValidator)

	cvs, err := svc.ListByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cvs) != 0 {
		t.Errorf("expected 0 CVs, got %d", len(cvs))
	}
}
