package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
)

func TestGetCV_Success(t *testing.T) {
	userID := "user-123"
	testCV := createTestCV(userID)

	mockCVSvc := &MockCVService{
		GetByUserIDFunc: func(ctx context.Context, uid string) (*domain.CV, error) {
			if uid != userID {
				t.Errorf("GetByUserID called with wrong userID: got %s, want %s", uid, userID)
			}
			return testCV, nil
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)
	req := newAuthenticatedRequest("GET", "/cv", userID, nil)
	rr := httptest.NewRecorder()

	h.GetCV(rr, req)

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

	cvResp, err := parseCVResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse CV response: %v", err)
	}

	if cvResp.ID != testCV.ID.String() {
		t.Errorf("expected CV ID %s, got %s", testCV.ID.String(), cvResp.ID)
	}
	if cvResp.Title != testCV.Title {
		t.Errorf("expected title %s, got %s", testCV.Title, cvResp.Title)
	}
	if cvResp.SchemaVersion != testCV.SchemaVersion {
		t.Errorf("expected schema version %s, got %s", testCV.SchemaVersion, cvResp.SchemaVersion)
	}
}

func TestGetCV_NotFound(t *testing.T) {
	userID := "user-123"

	mockCVSvc := &MockCVService{
		GetByUserIDFunc: func(ctx context.Context, uid string) (*domain.CV, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)
	req := newAuthenticatedRequest("GET", "/cv", userID, nil)
	rr := httptest.NewRecorder()

	h.GetCV(rr, req)

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

func TestCreateCV_Success(t *testing.T) {
	userID := "user-123"
	title := "Software Engineer CV"
	content := json.RawMessage(`{"sections":[]}`)
	testCV := createTestCV(userID)
	testCV.Title = title
	testCV.Content = content

	mockCVSvc := &MockCVService{
		CreateFunc: func(ctx context.Context, uid, ttl string, cnt json.RawMessage) (*domain.CV, error) {
			if uid != userID {
				t.Errorf("Create called with wrong userID: got %s, want %s", uid, userID)
			}
			if ttl != title {
				t.Errorf("Create called with wrong title: got %s, want %s", ttl, title)
			}
			return testCV, nil
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	body := jsonBody(CreateCVRequest{
		Title:   title,
		Content: content,
	})
	req := newAuthenticatedRequest("POST", "/cv", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateCV(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("unexpected error in response: %v", resp.Error)
	}

	cvResp, err := parseCVResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse CV response: %v", err)
	}

	if cvResp.Title != title {
		t.Errorf("expected title %s, got %s", title, cvResp.Title)
	}
}

func TestCreateCV_MissingTitle(t *testing.T) {
	userID := "user-123"
	content := json.RawMessage(`{"sections":[]}`)

	h := NewTestHandler(nil, &MockCVService{}, nil)

	body := jsonBody(CreateCVRequest{
		Title:   "", // Missing title
		Content: content,
	})
	req := newAuthenticatedRequest("POST", "/cv", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateCV(rr, req)

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

	if resp.Error.Field != "title" {
		t.Errorf("expected error field 'title', got '%s'", resp.Error.Field)
	}
}

func TestCreateCV_MissingContent(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, &MockCVService{}, nil)

	// Send JSON without content field (it will be empty/nil)
	body := strings.NewReader(`{"title": "My CV"}`)
	req := newAuthenticatedRequest("POST", "/cv", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateCV(rr, req)

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

	if resp.Error.Field != "content" {
		t.Errorf("expected error field 'content', got '%s'", resp.Error.Field)
	}
}

func TestCreateCV_InvalidBody(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, &MockCVService{}, nil)

	// Send invalid JSON
	req := newAuthenticatedRequest("POST", "/cv", userID, jsonBody("not valid json structure"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateCV(rr, req)

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

func TestCreateCV_LimitReached(t *testing.T) {
	userID := "user-123"
	title := "Second CV"
	content := json.RawMessage(`{"sections":[]}`)

	mockCVSvc := &MockCVService{
		CreateFunc: func(ctx context.Context, uid, ttl string, cnt json.RawMessage) (*domain.CV, error) {
			return nil, domain.ErrCVLimitReached
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	body := jsonBody(CreateCVRequest{
		Title:   title,
		Content: content,
	})
	req := newAuthenticatedRequest("POST", "/cv", userID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateCV(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp, err := parseJSONResponse(rr.Body)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}

	if resp.Error.Code != "cv_limit_reached" {
		t.Errorf("expected error code 'cv_limit_reached', got '%s'", resp.Error.Code)
	}
}

func TestUpdateCV_Success(t *testing.T) {
	userID := "user-123"
	cvID := uuid.New()
	newTitle := "Updated CV Title"
	testCV := createTestCV(userID)
	testCV.ID = cvID
	testCV.Title = newTitle

	mockCVSvc := &MockCVService{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, uid string, title *string, content *json.RawMessage) (*domain.CV, error) {
			if id != cvID {
				t.Errorf("Update called with wrong ID: got %s, want %s", id, cvID)
			}
			if uid != userID {
				t.Errorf("Update called with wrong userID: got %s, want %s", uid, userID)
			}
			if title == nil || *title != newTitle {
				t.Errorf("Update called with wrong title")
			}
			return testCV, nil
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	body := jsonBody(UpdateCVRequest{
		Title: &newTitle,
	})
	req := newAuthenticatedRequestWithParams("PUT", "/cv/"+cvID.String(), userID, body, map[string]string{"id": cvID.String()})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCV(rr, req)

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

	cvResp, err := parseCVResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse CV response: %v", err)
	}

	if cvResp.Title != newTitle {
		t.Errorf("expected title %s, got %s", newTitle, cvResp.Title)
	}
}

func TestUpdateCV_InvalidID(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, &MockCVService{}, nil)

	body := jsonBody(UpdateCVRequest{})
	req := newAuthenticatedRequestWithParams("PUT", "/cv/not-a-uuid", userID, body, map[string]string{"id": "not-a-uuid"})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCV(rr, req)

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

func TestUpdateCV_NotFound(t *testing.T) {
	userID := "user-123"
	cvID := uuid.New()
	newTitle := "Updated Title"

	mockCVSvc := &MockCVService{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, uid string, title *string, content *json.RawMessage) (*domain.CV, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	body := jsonBody(UpdateCVRequest{
		Title: &newTitle,
	})
	req := newAuthenticatedRequestWithParams("PUT", "/cv/"+cvID.String(), userID, body, map[string]string{"id": cvID.String()})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCV(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUpdateCV_Forbidden(t *testing.T) {
	userID := "user-123"
	cvID := uuid.New()
	newTitle := "Updated Title"

	mockCVSvc := &MockCVService{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, uid string, title *string, content *json.RawMessage) (*domain.CV, error) {
			return nil, domain.ErrForbidden
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	body := jsonBody(UpdateCVRequest{
		Title: &newTitle,
	})
	req := newAuthenticatedRequestWithParams("PUT", "/cv/"+cvID.String(), userID, body, map[string]string{"id": cvID.String()})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateCV(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}

func TestDeleteCV_Success(t *testing.T) {
	userID := "user-123"
	cvID := uuid.New()
	deleteCalled := false

	mockCVSvc := &MockCVService{
		DeleteFunc: func(ctx context.Context, id uuid.UUID, uid string) error {
			deleteCalled = true
			if id != cvID {
				t.Errorf("Delete called with wrong ID: got %s, want %s", id, cvID)
			}
			if uid != userID {
				t.Errorf("Delete called with wrong userID: got %s, want %s", uid, userID)
			}
			return nil
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	req := newAuthenticatedRequestWithParams("DELETE", "/cv/"+cvID.String(), userID, nil, map[string]string{"id": cvID.String()})
	rr := httptest.NewRecorder()

	h.DeleteCV(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	if !deleteCalled {
		t.Error("expected Delete to be called")
	}
}

func TestDeleteCV_InvalidID(t *testing.T) {
	userID := "user-123"

	h := NewTestHandler(nil, &MockCVService{}, nil)

	req := newAuthenticatedRequestWithParams("DELETE", "/cv/not-a-uuid", userID, nil, map[string]string{"id": "not-a-uuid"})
	rr := httptest.NewRecorder()

	h.DeleteCV(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDeleteCV_NotFound(t *testing.T) {
	userID := "user-123"
	cvID := uuid.New()

	mockCVSvc := &MockCVService{
		DeleteFunc: func(ctx context.Context, id uuid.UUID, uid string) error {
			return domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, mockCVSvc, nil)

	req := newAuthenticatedRequestWithParams("DELETE", "/cv/"+cvID.String(), userID, nil, map[string]string{"id": cvID.String()})
	rr := httptest.NewRecorder()

	h.DeleteCV(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
