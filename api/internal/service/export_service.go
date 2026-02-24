package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/export"
	"github.com/ellebam/hireme/api/internal/repository"
)

// HTMLRenderer renders CV content to an HTML string.
type HTMLRenderer interface {
	Render(content domain.CVContent) (string, error)
}

// ExportService handles document export operations.
type ExportService struct {
	cvRepo   repository.CVRepository
	renderer HTMLRenderer
	pdf      export.PDFConverter
	docx     export.DOCXGenerator
}

// NewExportService creates a new ExportService.
func NewExportService(cvRepo repository.CVRepository, renderer HTMLRenderer, pdf export.PDFConverter, docx export.DOCXGenerator) *ExportService {
	return &ExportService{
		cvRepo:   cvRepo,
		renderer: renderer,
		pdf:      pdf,
		docx:     docx,
	}
}

// fetchCVContent fetches a specific CV by ID, verifies ownership, and parses its structured content.
func (s *ExportService) fetchCVContent(ctx context.Context, cvID uuid.UUID, userID string) (*domain.CVContent, error) {
	cv, err := s.cvRepo.GetByID(ctx, cvID)
	if err != nil {
		return nil, fmt.Errorf("fetching CV: %w", err)
	}

	if cv.UserID != userID {
		return nil, domain.ErrForbidden
	}

	content, err := cv.ParseContent()
	if err != nil {
		return nil, fmt.Errorf("parsing CV content: %w", err)
	}

	return content, nil
}

// renderHTML fetches a specific CV, parses its content, and renders it to HTML.
func (s *ExportService) renderHTML(ctx context.Context, cvID uuid.UUID, userID string) (string, error) {
	content, err := s.fetchCVContent(ctx, cvID, userID)
	if err != nil {
		return "", err
	}

	html, err := s.renderer.Render(*content)
	if err != nil {
		return "", fmt.Errorf("rendering HTML: %w", err)
	}

	return html, nil
}

// ExportPDF generates a PDF for a specific CV.
func (s *ExportService) ExportPDF(ctx context.Context, cvID uuid.UUID, userID string) ([]byte, error) {
	html, err := s.renderHTML(ctx, cvID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdf.ConvertHTMLToPDF(ctx, html)
	if err != nil {
		return nil, fmt.Errorf("converting to PDF: %w", err)
	}

	return pdf, nil
}

// ExportDOCX generates a DOCX for a specific CV.
func (s *ExportService) ExportDOCX(ctx context.Context, cvID uuid.UUID, userID string) ([]byte, error) {
	content, err := s.fetchCVContent(ctx, cvID, userID)
	if err != nil {
		return nil, err
	}

	docxBytes, err := s.docx.Generate(*content)
	if err != nil {
		return nil, fmt.Errorf("generating DOCX: %w", err)
	}

	return docxBytes, nil
}
