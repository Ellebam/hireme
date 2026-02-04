package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ellebam/hireme/api/internal/middleware"
	"github.com/ellebam/hireme/api/pkg/httputil"
)

// AssetResponse represents an asset in API responses
type AssetResponse struct {
	ID               string `json:"id"`
	Filename         string `json:"filename"`
	OriginalFilename string `json:"originalFilename"`
	MimeType         string `json:"mimeType"`
	SizeBytes        int64  `json:"sizeBytes"`
	Width            *int   `json:"width,omitempty"`
	Height           *int   `json:"height,omitempty"`
	URL              string `json:"url"`
	CreatedAt        string `json:"createdAt"`
}

// UploadAsset handles file uploads
func (h *Handler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputil.Error(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	// Upload via service
	asset, err := h.assetService.Upload(ctx, userID, file, header)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := AssetResponse{
		ID:               asset.ID.String(),
		Filename:         asset.Filename,
		OriginalFilename: asset.OriginalFilename,
		MimeType:         asset.MimeType,
		SizeBytes:        asset.SizeBytes,
		Width:            asset.Width,
		Height:           asset.Height,
		URL:              "/api/v1/assets/" + asset.ID.String(), // TODO: generate proper URL
		CreatedAt:        asset.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	httputil.Created(w, response)
}

// GetAsset retrieves an asset (serves the file or returns metadata)
func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse asset ID from URL
	idParam := chi.URLParam(r, "id")
	assetID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	// Get asset metadata
	asset, err := h.assetService.GetByID(ctx, assetID, userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Check if client wants file content (explicit image accept header) or metadata (default)
	// Only serve file when explicitly requesting the image MIME type or image/*
	accept := r.Header.Get("Accept")
	wantsFile := accept == asset.MimeType || accept == "image/*"

	if !wantsFile {
		// Default: return JSON metadata
		response := AssetResponse{
			ID:               asset.ID.String(),
			Filename:         asset.Filename,
			OriginalFilename: asset.OriginalFilename,
			MimeType:         asset.MimeType,
			SizeBytes:        asset.SizeBytes,
			Width:            asset.Width,
			Height:           asset.Height,
			URL:              "/api/v1/assets/" + asset.ID.String(),
			CreatedAt:        asset.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		httputil.JSON(w, http.StatusOK, response)
		return
	}

	// Serve the file (only when explicitly requesting image content)
	fileContent, err := h.assetService.GetFileContent(ctx, asset)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", asset.MimeType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+asset.Filename+"\"")
	w.Write(fileContent)
}

// DeleteAsset deletes an asset
func (h *Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	// Parse asset ID from URL
	idParam := chi.URLParam(r, "id")
	assetID, err := uuid.Parse(idParam)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	if err := h.assetService.Delete(ctx, assetID, userID); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.NoContent(w)
}
