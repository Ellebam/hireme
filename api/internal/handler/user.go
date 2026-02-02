package handler

import (
	"net/http"

	"github.com/yourusername/hireme/api/internal/middleware"
	"github.com/yourusername/hireme/api/pkg/httputil"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	DisplayName       string `json:"displayName"`
	Tier              string `json:"tier"`
	CVLimit           int    `json:"cvLimit"`
	StorageLimitBytes int64  `json:"storageLimitBytes"`
	StorageUsedBytes  int64  `json:"storageUsedBytes"`
	Locale            string `json:"locale"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	DisplayName *string `json:"displayName,omitempty"`
	Locale      *string `json:"locale,omitempty"`
}

// GetCurrentUser returns the authenticated user's profile
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := UserResponse{
		ID:                user.ID,
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		Tier:              user.Tier,
		CVLimit:           user.CVLimit,
		StorageLimitBytes: user.StorageLimitBytes,
		StorageUsedBytes:  user.StorageUsedBytes,
		Locale:            user.Locale,
	}

	httputil.JSON(w, http.StatusOK, response)
}

// UpdateCurrentUser updates the authenticated user's profile
func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.MustGetUserID(ctx)

	var req UpdateUserRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Update(ctx, userID, req.DisplayName, req.Locale)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	response := UserResponse{
		ID:                user.ID,
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		Tier:              user.Tier,
		CVLimit:           user.CVLimit,
		StorageLimitBytes: user.StorageLimitBytes,
		StorageUsedBytes:  user.StorageUsedBytes,
		Locale:            user.Locale,
	}

	httputil.JSON(w, http.StatusOK, response)
}
