package service

import (
	"context"
	"fmt"

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
}

// NewExportService creates a new ExportService.
func NewExportService(cvRepo repository.CVRepository, renderer HTMLRenderer, pdf export.PDFConverter) *ExportService {
	return &ExportService{
		cvRepo:   cvRepo,
		renderer: renderer,
		pdf:      pdf,
	}
}

// ExportPDF generates a PDF for the user's active CV.
func (s *ExportService) ExportPDF(ctx context.Context, userID string) ([]byte, error) {
	cv, err := s.cvRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching CV: %w", err)
	}

	content, err := cv.ParseContent()
	if err != nil {
		return nil, fmt.Errorf("parsing CV content: %w", err)
	}

	html, err := s.renderer.Render(*content)
	if err != nil {
		return nil, fmt.Errorf("rendering HTML: %w", err)
	}

	pdf, err := s.pdf.ConvertHTMLToPDF(ctx, html)
	if err != nil {
		return nil, fmt.Errorf("converting to PDF: %w", err)
	}

	return pdf, nil
}
