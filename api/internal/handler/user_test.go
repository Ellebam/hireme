package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellebam/hireme/api/internal/domain"
)

func TestGetCurrentUser_Success(t *testing.T) {
	userID := "user-123"
	testUser := createTestUser(userID)

	mockUserSvc := &MockUserService{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			if id != userID {
				t.Errorf("GetByID called with wrong ID: got %s, want %s", id, userID)
			}
			return testUser, nil
		},
	}

	h := NewTestHandler(mockUserSvc, nil, nil)
	req := newAuthenticatedRequest("GET", "/users/me", userID, nil)
	rr := httptest.NewRecorder()

	h.GetCurrentUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error in response: %v", resp.Error)
	}

	userResp, err := parseUserResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse user response: %v", err)
	}

	if userResp.ID != testUser.ID {
		t.Errorf("expected user ID %s, got %s", testUser.ID, userResp.ID)
	}
	if userResp.Email != testUser.Email {
		t.Errorf("expected email %s, got %s", testUser.Email, userResp.Email)
	}
	if userResp.DisplayName != testUser.DisplayName {
		t.Errorf("expected display name %s, got %s", testUser.DisplayName, userResp.DisplayName)
	}
	if userResp.Tier != testUser.Tier {
		t.Errorf("expected tier %s, got %s", testUser.Tier, userResp.Tier)
	}
	if userResp.CVLimit != testUser.CVLimit {
		t.Errorf("expected CV limit %d, got %d", testUser.CVLimit, userResp.CVLimit)
	}
	if userResp.StorageLimitBytes != testUser.StorageLimitBytes {
		t.Errorf("expected storage limit %d, got %d", testUser.StorageLimitBytes, userResp.StorageLimitBytes)
	}
	if userResp.StorageUsedBytes != testUser.StorageUsedBytes {
		t.Errorf("expected storage used %d, got %d", testUser.StorageUsedBytes, userResp.StorageUsedBytes)
	}
	if userResp.Locale != testUser.Locale {
		t.Errorf("expected locale %s, got %s", testUser.Locale, userResp.Locale)
	}
}

func TestGetCurrentUser_NotFound(t *testing.T) {
	userID := "nonexistent-user"

	mockUserSvc := &MockUserService{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(mockUserSvc, nil, nil)
	req := newAuthenticatedRequest("GET", "/users/me", userID, nil)
	rr := httptest.NewRecorder()

	h.GetCurrentUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}
}

func TestUpdateCurrentUser_Success(t *testing.T) {
	userID := "user-123"
	newDisplayName := "Updated Name"
	testUser := createTestUser(userID)
	testUser.DisplayName = newDisplayName

	mockUserSvc := &MockUserService{
		UpdateFunc: func(ctx context.Context, id string, displayName, locale *string) (*domain.User, error) {
			if id != userID {
				t.Errorf("Update called with wrong ID: got %s, want %s", id, userID)
			}
			if displayName == nil || *displayName != newDisplayName {
				t.Errorf("Update called with wrong displayName")
			}
			return testUser, nil
		},
	}

	h := NewTestHandler(mockUserSvc, nil, nil)

	body := jsonBody(UpdateUserRequest{
		DisplayName: &newDisplayName,
	})
	req := newAuthenticatedRequest("PATCH", "/users/me", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCurrentUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error in response: %v", resp.Error)
	}

	userResp, err := parseUserResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse user response: %v", err)
	}

	if userResp.DisplayName != newDisplayName {
		t.Errorf("expected display name %s, got %s", newDisplayName, userResp.DisplayName)
	}
}

func TestUpdateCurrentUser_UpdateLocale(t *testing.T) {
	userID := "user-123"
	newLocale := domain.LocaleDE
	testUser := createTestUser(userID)
	testUser.Locale = newLocale

	mockUserSvc := &MockUserService{
		UpdateFunc: func(ctx context.Context, id string, displayName, locale *string) (*domain.User, error) {
			if locale == nil || *locale != newLocale {
				t.Errorf("Update called with wrong locale: got %v, want %s", locale, newLocale)
			}
			return testUser, nil
		},
	}

	h := NewTestHandler(mockUserSvc, nil, nil)

	body := jsonBody(UpdateUserRequest{
		Locale: &newLocale,
	})
	req := newAuthenticatedRequest("PATCH", "/users/me", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCurrentUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	userResp, err := parseUserResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse user response: %v", err)
	}

	if userResp.Locale != newLocale {
		t.Errorf("expected locale %s, got %s", newLocale, userResp.Locale)
	}
}

func TestUpdateCurrentUser_InvalidBody(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(&MockUserService{}, nil, nil)

	// Send invalid JSON
	req := newAuthenticatedRequest("PATCH", "/users/me", userID, jsonBody("not valid json structure"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCurrentUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}
}

func TestUpdateCurrentUser_ValidationError(t *testing.T) {
	userID := "user-123"
	invalidLocale := "invalid"

	mockUserSvc := &MockUserService{
		UpdateFunc: func(ctx context.Context, id string, displayName, locale *string) (*domain.User, error) {
			return nil, domain.NewValidationError("locale", "invalid locale, must be 'en' or 'de'")
		},
	}

	h := NewTestHandler(mockUserSvc, nil, nil)

	body := jsonBody(UpdateUserRequest{
		Locale: &invalidLocale,
	})
	req := newAuthenticatedRequest("PATCH", "/users/me", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCurrentUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}

	if resp.Error.Field != "locale" {
		t.Errorf("expected error field 'locale', got '%s'", resp.Error.Field)
	}
}

func TestUpdateCurrentUser_NotFound(t *testing.T) {
	userID := "nonexistent-user"
	newDisplayName := "Updated Name"

	mockUserSvc := &MockUserService{
		UpdateFunc: func(ctx context.Context, id string, displayName, locale *string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(mockUserSvc, nil, nil)

	body := jsonBody(UpdateUserRequest{
		DisplayName: &newDisplayName,
	})
	req := newAuthenticatedRequest("PATCH", "/users/me", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCurrentUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
