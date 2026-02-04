package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ellebam/hireme/api/internal/domain"
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

// CreateExport initiates a new export job
func (h *Handler) CreateExport(w http.ResponseWriter, r *http.Request) {
	// Get format from URL
	format := chi.URLParam(r, "format")

	// Validate format
	if !domain.IsValidExportFormat(format) {
		httputil.Error(w, http.StatusBadRequest, "invalid export format")
		return
	}

	// TODO: Implement export job creation
	// 1. Get user's CV
	// 2. Create export job record
	// 3. Queue for processing (or process synchronously for MVP)
	// 4. Return job status

	httputil.Error(w, http.StatusNotImplemented, "export not yet implemented")
}

// GetExport returns the status of an export job
func (h *Handler) GetExport(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement export status retrieval
	// 1. Get export job by ID
	// 2. Verify user owns the job
	// 3. Return status and download URL if complete

	httputil.Error(w, http.StatusNotImplemented, "export not yet implemented")
}
