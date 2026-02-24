package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/middleware"
	"github.com/ellebam/hireme/api/pkg/httputil"
)

// ExportResponse represents an export job in API responses
type ExportResponse struct {
	ID        string  `json:"id"`
	Format    string  `json:"format"`
	Status    string  `json:"status"`
	URL       *string `json:"url,omitempty"`
	Error     *string `json:"error,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

// CreateExport handles synchronous export requests for a specific CV
func (h *Handler) CreateExport(w http.ResponseWriter, r *http.Request) {
	// Parse CV ID from URL
	idParam := chi.URLParam(r, "id")
	cvID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid CV ID")
		return
	}

	format := chi.URLParam(r, "format")

	if !domain.IsValidExportFormat(format) {
		httputil.Error(w, http.StatusBadRequest, "invalid export format")
		return
	}

	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	switch format {
	case domain.ExportFormatPDF:
		if !h.config.Features.EnableExportPDF {
			httputil.Error(w, http.StatusNotImplemented, "PDF export is not enabled")
			return
		}

		pdfBytes, err := h.exportService.ExportPDF(ctx, cvID, userID)
		if err != nil {
			httputil.HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="export.pdf"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdfBytes)
	case domain.ExportFormatDOCX:
		if !h.config.Features.EnableExportDOCX {
			httputil.Error(w, http.StatusNotImplemented, "DOCX export is not enabled")
			return
		}

		docxBytes, err := h.exportService.ExportDOCX(ctx, cvID, userID)
		if err != nil {
			httputil.HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		w.Header().Set("Content-Disposition", `attachment; filename="export.docx"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(docxBytes)
	default:
		httputil.Error(w, http.StatusNotImplemented, "export format not yet implemented")
	}
}

// GetExport returns the status of an export job
func (h *Handler) GetExport(w http.ResponseWriter, r *http.Request) {
	httputil.Error(w, http.StatusNotImplemented, "export not yet implemented")
}
