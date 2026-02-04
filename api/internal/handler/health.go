package handler

import (
	"net/http"

	"github.com/ellebam/hireme/api/pkg/httputil"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// Health returns the application health status
// This endpoint is used for liveness probes
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0", // TODO: inject from build
	}
	httputil.JSON(w, http.StatusOK, response)
}

// ReadyResponse represents the readiness check response
type ReadyResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// Ready returns the application readiness status
// This endpoint is used for readiness probes
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	// TODO: Add actual health checks for dependencies
	// - Database connection
	// - Storage backend
	// - Gotenberg (if configured)

	response := ReadyResponse{
		Status: "ready",
		Services: map[string]string{
			"database": "ok",
			"storage":  "ok",
		},
	}
	httputil.JSON(w, http.StatusOK, response)
}
