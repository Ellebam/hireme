package handler

import (
	"github.com/yourusername/hireme/api/internal/config"
	"github.com/yourusername/hireme/api/internal/service"
)

// Handler holds all HTTP handlers and their dependencies
type Handler struct {
	config       *config.Config
	userService  *service.UserService
	cvService    *service.CVService
	assetService *service.AssetService
}

// Dependencies contains all services needed by handlers
type Dependencies struct {
	Config       *config.Config
	UserService  *service.UserService
	CVService    *service.CVService
	AssetService *service.AssetService
}

// New creates a new Handler with the given dependencies
func New(deps Dependencies) *Handler {
	return &Handler{
		config:       deps.Config,
		userService:  deps.UserService,
		cvService:    deps.CVService,
		assetService: deps.AssetService,
	}
}
