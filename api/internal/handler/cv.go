package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yourusername/hireme/api/internal/middleware"
	"github.com/yourusername/hireme/api/pkg/httputil"
)

// CVResponse represents a CV in API responses
type CVResponse struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	SchemaVersion string          `json:"schemaVersion"`
	Content       json.RawMessage `json:"content"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
}

// CreateCVRequest represents a CV creation request
type CreateCVRequest struct {
	Title   string          `json:"title"`
	Content json.RawMessage `json:"content"`
}

// UpdateCVRequest represents a CV update request
type UpdateCVRequest struct {
	Title   *string          `json:"title,omitempty"`
	Content *json.RawMessage `json:"content,omitempty"`
}

// GetCV returns the authenticated user's active CV
func (h *Handler) GetCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	cv, err := h.cvService.GetByUserID(ctx, userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := CVResponse{
		ID:            cv.ID.String(),
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
		CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.JSON(w, http.StatusOK, response)
}

// CreateCV creates a new CV for the authenticated user
func (h *Handler) CreateCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	var req CreateCVRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" {
		httputil.ValidationError(w, "title", "title is required")
		return
	}
	if len(req.Content) == 0 {
		httputil.ValidationError(w, "content", "content is required")
		return
	}

	cv, err := h.cvService.Create(ctx, userID, req.Title, req.Content)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := CVResponse{
		ID:            cv.ID.String(),
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
		CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.Created(w, response)
}

// UpdateCV updates an existing CV
func (h *Handler) UpdateCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse CV ID from URL
	idParam := chi.URLParam(r, "id")
	cvID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid CV ID")
		return
	}

	var req UpdateCVRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cv, err := h.cvService.Update(ctx, cvID, userID, req.Title, req.Content)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := CVResponse{
		ID:            cv.ID.String(),
		Title:         cv.Title,
		SchemaVersion: cv.SchemaVersion,
		Content:       cv.Content,
		CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.JSON(w, http.StatusOK, response)
}

// DeleteCV deletes a CV
func (h *Handler) DeleteCV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse CV ID from URL
	idParam := chi.URLParam(r, "id")
	cvID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid CV ID")
		return
	}

	if err := h.cvService.Delete(ctx, cvID, userID); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.NoContent(w)
}
