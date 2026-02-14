package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellebam/hireme/api/internal/config"
	"github.com/ellebam/hireme/api/internal/domain"
)

func TestCreateExport_PDF_Success(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 test content")
	userID := "user-123"

	mockExportSvc := &MockExportService{
		ExportPDFFunc: func(ctx context.Context, uid string) ([]byte, error) {
			if uid != userID {
				t.Errorf("expected userID %q, got %q", userID, uid)
			}
			return expectedPDF, nil
		},
	}

	h := NewTestHandler(nil, nil, nil, mockExportSvc)

	req := newAuthenticatedRequestWithParams("POST", "/export/pdf", userID, nil, map[string]string{"format": "pdf"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("expected Content-Type application/pdf, got %q", ct)
	}
	if cd := rr.Header().Get("Content-Disposition"); cd != `attachment; filename="export.pdf"` {
		t.Errorf("expected Content-Disposition attachment, got %q", cd)
	}
	if rr.Body.String() != string(expectedPDF) {
		t.Errorf("expected body %q, got %q", expectedPDF, rr.Body.String())
	}
}

func TestCreateExport_InvalidFormat(t *testing.T) {
	formats := []string{"", "html", "exe", "PDF"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			h := NewTestHandler(nil, nil, nil, nil)
			req := newAuthenticatedRequestWithParams("POST", "/export/"+format, "user-1", nil, map[string]string{"format": format})
			rr := httptest.NewRecorder()

			h.CreateExport(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("format %q: expected 400, got %d", format, rr.Code)
			}
		})
	}
}

func TestCreateExport_PDF_FeatureDisabled(t *testing.T) {
	h := NewTestHandler(nil, nil, nil, &MockExportService{})
	h.config = &config.Config{
		Features: config.FeaturesConfig{EnableExportPDF: false},
	}

	req := newAuthenticatedRequestWithParams("POST", "/export/pdf", "user-1", nil, map[string]string{"format": "pdf"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", rr.Code)
	}
}

func TestCreateExport_PDF_ServiceError(t *testing.T) {
	mockExportSvc := &MockExportService{
		ExportPDFFunc: func(ctx context.Context, userID string) ([]byte, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, nil, nil, mockExportSvc)

	req := newAuthenticatedRequestWithParams("POST", "/export/pdf", "user-1", nil, map[string]string{"format": "pdf"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestCreateExport_DOCX_NotImplemented(t *testing.T) {
	h := NewTestHandler(nil, nil, nil, nil)

	req := newAuthenticatedRequestWithParams("POST", "/export/docx", "user-1", nil, map[string]string{"format": "docx"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", rr.Code)
	}
}
