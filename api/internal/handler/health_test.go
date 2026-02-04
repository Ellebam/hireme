package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	h := NewTestHandler(nil, nil, nil)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	h.Health(rr, req)

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

	healthResp, err := parseHealthResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse health response: %v", err)
	}

	if healthResp.Status != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", healthResp.Status)
	}

	if healthResp.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", healthResp.Version)
	}
}

func TestReady_AllHealthy(t *testing.T) {
	h := NewTestHandler(nil, nil, nil)

	req := httptest.NewRequest("GET", "/ready", nil)
	rr := httptest.NewRecorder()

	h.Ready(rr, req)

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

	readyResp, err := parseReadyResponse(resp)
	if err != nil {
		t.Fatalf("failed to parse ready response: %v", err)
	}

	if readyResp.Status != "ready" {
		t.Errorf("expected status 'ready', got '%s'", readyResp.Status)
	}

	// Check services
	if readyResp.Services == nil {
		t.Fatal("expected services map to be non-nil")
	}

	if readyResp.Services["database"] != "ok" {
		t.Errorf("expected database status 'ok', got '%s'", readyResp.Services["database"])
	}

	if readyResp.Services["storage"] != "ok" {
		t.Errorf("expected storage status 'ok', got '%s'", readyResp.Services["storage"])
	}
}
