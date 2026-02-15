package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
)

// mockCVRepo implements repository.CVRepository for testing
type mockCVRepo struct {
	GetByUserIDFunc func(ctx context.Context, userID string) (*domain.CV, error)
}

func (m *mockCVRepo) GetByUserID(ctx context.Context, userID string) (*domain.CV, error) {
	return m.GetByUserIDFunc(ctx, userID)
}

func (m *mockCVRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.CV, error) {
	return nil, nil
}
func (m *mockCVRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.CV, error) {
	return nil, nil
}
func (m *mockCVRepo) CountByUserID(ctx context.Context, userID string) (int, error) { return 0, nil }
func (m *mockCVRepo) Create(ctx context.Context, cv *domain.CV) error              { return nil }
func (m *mockCVRepo) Update(ctx context.Context, cv *domain.CV) error              { return nil }
func (m *mockCVRepo) Delete(ctx context.Context, id uuid.UUID) error               { return nil }

// mockRenderer implements HTMLRenderer for testing
type mockRenderer struct {
	RenderFunc func(content domain.CVContent) (string, error)
}

func (m *mockRenderer) Render(content domain.CVContent) (string, error) {
	return m.RenderFunc(content)
}

// mockPDFConverter implements export.PDFConverter for testing
type mockPDFConverter struct {
	ConvertFunc func(ctx context.Context, html string) ([]byte, error)
}

func (m *mockPDFConverter) ConvertHTMLToPDF(ctx context.Context, html string) ([]byte, error) {
	return m.ConvertFunc(ctx, html)
}

func exportTestCVContent() json.RawMessage {
	return json.RawMessage(`{"schemaVersion":"1.0.0","templateId":"classic","sections":[]}`)
}

func TestExportPDF_Success(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 test")

	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), UserID: userID, Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{
		RenderFunc: func(content domain.CVContent) (string, error) {
			return "<html>rendered</html>", nil
		},
	}
	pdf := &mockPDFConverter{
		ConvertFunc: func(ctx context.Context, html string) ([]byte, error) {
			if html != "<html>rendered</html>" {
				t.Errorf("expected rendered HTML, got %q", html)
			}
			return expectedPDF, nil
		},
	}

	svc := NewExportService(repo, renderer, pdf, nil)
	result, err := svc.ExportPDF(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != string(expectedPDF) {
		t.Errorf("expected %q, got %q", expectedPDF, result)
	}
}

func TestExportPDF_CVNotFound(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return nil, domain.ErrNotFound
		},
	}
	renderer := &mockRenderer{RenderFunc: func(domain.CVContent) (string, error) { return "", nil }}
	pdf := &mockPDFConverter{ConvertFunc: func(context.Context, string) ([]byte, error) { return nil, nil }}

	svc := NewExportService(repo, renderer, pdf, nil)
	_, err := svc.ExportPDF(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestExportPDF_RenderError(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{
		RenderFunc: func(domain.CVContent) (string, error) {
			return "", errors.New("template error")
		},
	}
	pdf := &mockPDFConverter{ConvertFunc: func(context.Context, string) ([]byte, error) { return nil, nil }}

	svc := NewExportService(repo, renderer, pdf, nil)
	_, err := svc.ExportPDF(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "rendering HTML: template error" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestExportPDF_ConvertError(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{
		RenderFunc: func(domain.CVContent) (string, error) { return "<html></html>", nil },
	}
	pdf := &mockPDFConverter{
		ConvertFunc: func(context.Context, string) ([]byte, error) {
			return nil, errors.New("gotenberg timeout")
		},
	}

	svc := NewExportService(repo, renderer, pdf, nil)
	_, err := svc.ExportPDF(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "converting to PDF: gotenberg timeout" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestExportPDF_ParseContentError(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: json.RawMessage("not json")}, nil
		},
	}
	renderer := &mockRenderer{RenderFunc: func(domain.CVContent) (string, error) { return "", nil }}
	pdf := &mockPDFConverter{ConvertFunc: func(context.Context, string) ([]byte, error) { return nil, nil }}

	svc := NewExportService(repo, renderer, pdf, nil)
	_, err := svc.ExportPDF(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error from invalid JSON, got nil")
	}
}

func TestExportPDF_EmptyHTML(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{
		RenderFunc: func(domain.CVContent) (string, error) { return "", nil },
	}
	pdf := &mockPDFConverter{
		ConvertFunc: func(ctx context.Context, html string) ([]byte, error) {
			return []byte("%PDF-empty"), nil
		},
	}

	svc := NewExportService(repo, renderer, pdf, nil)
	result, err := svc.ExportPDF(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "%PDF-empty" {
		t.Errorf("expected PDF bytes, got %q", result)
	}
}

// mockDOCXGenerator implements export.DOCXGenerator for testing
type mockDOCXGenerator struct {
	GenerateFunc func(content domain.CVContent) ([]byte, error)
}

func (m *mockDOCXGenerator) Generate(content domain.CVContent) ([]byte, error) {
	return m.GenerateFunc(content)
}

func TestExportDOCX_Success(t *testing.T) {
	expectedDOCX := []byte("PK fake docx")

	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), UserID: userID, Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{
		RenderFunc: func(domain.CVContent) (string, error) {
			t.Fatal("renderer should not be called for DOCX")
			return "", nil
		},
	}
	docx := &mockDOCXGenerator{
		GenerateFunc: func(content domain.CVContent) ([]byte, error) {
			if content.SchemaVersion != "1.0.0" {
				t.Errorf("expected schemaVersion 1.0.0, got %q", content.SchemaVersion)
			}
			if content.TemplateID != "classic" {
				t.Errorf("expected templateId classic, got %q", content.TemplateID)
			}
			return expectedDOCX, nil
		},
	}

	svc := NewExportService(repo, renderer, nil, docx)
	result, err := svc.ExportDOCX(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != string(expectedDOCX) {
		t.Errorf("expected %q, got %q", expectedDOCX, result)
	}
}

func TestExportDOCX_CVNotFound(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return nil, domain.ErrNotFound
		},
	}
	renderer := &mockRenderer{RenderFunc: func(domain.CVContent) (string, error) { return "", nil }}
	docx := &mockDOCXGenerator{GenerateFunc: func(domain.CVContent) ([]byte, error) { return nil, nil }}

	svc := NewExportService(repo, renderer, nil, docx)
	_, err := svc.ExportDOCX(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestExportDOCX_GenerateError(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{RenderFunc: func(domain.CVContent) (string, error) { return "", nil }}
	docx := &mockDOCXGenerator{
		GenerateFunc: func(domain.CVContent) ([]byte, error) {
			return nil, errors.New("docx generation failed")
		},
	}

	svc := NewExportService(repo, renderer, nil, docx)
	_, err := svc.ExportDOCX(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "generating DOCX: docx generation failed" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestExportDOCX_ParseContentError(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: json.RawMessage("not json")}, nil
		},
	}
	renderer := &mockRenderer{RenderFunc: func(domain.CVContent) (string, error) { return "", nil }}
	docx := &mockDOCXGenerator{GenerateFunc: func(domain.CVContent) ([]byte, error) { return nil, nil }}

	svc := NewExportService(repo, renderer, nil, docx)
	_, err := svc.ExportDOCX(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error from invalid JSON, got nil")
	}
}

func TestExportDOCX_EmptySections(t *testing.T) {
	repo := &mockCVRepo{
		GetByUserIDFunc: func(ctx context.Context, userID string) (*domain.CV, error) {
			return &domain.CV{ID: uuid.New(), Content: exportTestCVContent()}, nil
		},
	}
	renderer := &mockRenderer{
		RenderFunc: func(domain.CVContent) (string, error) {
			t.Fatal("renderer should not be called for DOCX")
			return "", nil
		},
	}
	docx := &mockDOCXGenerator{
		GenerateFunc: func(content domain.CVContent) ([]byte, error) {
			return []byte("PK-empty"), nil
		},
	}

	svc := NewExportService(repo, renderer, nil, docx)
	result, err := svc.ExportDOCX(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "PK-empty" {
		t.Errorf("expected DOCX bytes, got %q", result)
	}
}
