package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/config"
	"github.com/ellebam/hireme/api/internal/domain"
)

func TestCreateExport_PDF_Success(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 test content")
	userID := "user-123"
	cvID := uuid.New()

	mockExportSvc := &MockExportService{
		ExportPDFFunc: func(ctx context.Context, id uuid.UUID, uid string) ([]byte, error) {
			if id != cvID {
				t.Errorf("expected cvID %s, got %s", cvID, id)
			}
			if uid != userID {
				t.Errorf("expected userID %q, got %q", userID, uid)
			}
			return expectedPDF, nil
		},
	}

	h := NewTestHandler(nil, nil, nil, mockExportSvc)

	req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/pdf", userID, nil, map[string]string{"id": cvID.String(), "format": "pdf"})
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

func TestCreateExport_InvalidCVID(t *testing.T) {
	h := NewTestHandler(nil, nil, nil, &MockExportService{})

	req := newAuthenticatedRequestWithParams("POST", "/cv/not-a-uuid/export/pdf", "user-1", nil, map[string]string{"id": "not-a-uuid", "format": "pdf"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateExport_InvalidFormat(t *testing.T) {
	cvID := uuid.New()
	formats := []string{"", "html", "exe", "PDF"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			h := NewTestHandler(nil, nil, nil, &MockExportService{})
			req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/"+format, "user-1", nil, map[string]string{"id": cvID.String(), "format": format})
			rr := httptest.NewRecorder()

			h.CreateExport(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("format %q: expected 400, got %d", format, rr.Code)
			}
		})
	}
}

func TestCreateExport_PDF_FeatureDisabled(t *testing.T) {
	cvID := uuid.New()
	h := NewTestHandler(nil, nil, nil, &MockExportService{})
	h.config = &config.Config{
		Features: config.FeaturesConfig{EnableExportPDF: false},
	}

	req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/pdf", "user-1", nil, map[string]string{"id": cvID.String(), "format": "pdf"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", rr.Code)
	}
}

func TestCreateExport_PDF_ServiceError(t *testing.T) {
	cvID := uuid.New()
	mockExportSvc := &MockExportService{
		ExportPDFFunc: func(ctx context.Context, id uuid.UUID, userID string) ([]byte, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, nil, nil, mockExportSvc)

	req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/pdf", "user-1", nil, map[string]string{"id": cvID.String(), "format": "pdf"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestCreateExport_DOCX_Success(t *testing.T) {
	expectedDOCX := []byte("PK fake docx content")
	userID := "user-123"
	cvID := uuid.New()

	mockExportSvc := &MockExportService{
		ExportDOCXFunc: func(ctx context.Context, id uuid.UUID, uid string) ([]byte, error) {
			if id != cvID {
				t.Errorf("expected cvID %s, got %s", cvID, id)
			}
			if uid != userID {
				t.Errorf("expected userID %q, got %q", userID, uid)
			}
			return expectedDOCX, nil
		},
	}

	h := NewTestHandler(nil, nil, nil, mockExportSvc)

	req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/docx", userID, nil, map[string]string{"id": cvID.String(), "format": "docx"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		t.Errorf("expected DOCX Content-Type, got %q", ct)
	}
	if cd := rr.Header().Get("Content-Disposition"); cd != `attachment; filename="export.docx"` {
		t.Errorf("expected Content-Disposition with export.docx, got %q", cd)
	}
	if rr.Body.String() != string(expectedDOCX) {
		t.Errorf("expected body %q, got %q", expectedDOCX, rr.Body.String())
	}
}

func TestCreateExport_DOCX_FeatureDisabled(t *testing.T) {
	cvID := uuid.New()
	h := NewTestHandler(nil, nil, nil, &MockExportService{})
	h.config = &config.Config{
		Features: config.FeaturesConfig{EnableExportDOCX: false},
	}

	req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/docx", "user-1", nil, map[string]string{"id": cvID.String(), "format": "docx"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", rr.Code)
	}
}

func TestCreateExport_DOCX_ServiceError(t *testing.T) {
	cvID := uuid.New()
	mockExportSvc := &MockExportService{
		ExportDOCXFunc: func(ctx context.Context, id uuid.UUID, userID string) ([]byte, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := NewTestHandler(nil, nil, nil, mockExportSvc)

	req := newAuthenticatedRequestWithParams("POST", "/cv/"+cvID.String()+"/export/docx", "user-1", nil, map[string]string{"id": cvID.String(), "format": "docx"})
	rr := httptest.NewRecorder()

	h.CreateExport(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
