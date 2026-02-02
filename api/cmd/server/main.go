package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourusername/hireme/api/internal/config"
	"github.com/yourusername/hireme/api/internal/handler"
	appMiddleware "github.com/yourusername/hireme/api/internal/middleware"
	"github.com/yourusername/hireme/api/internal/repository/postgres"
	"github.com/yourusername/hireme/api/internal/service"
	"github.com/yourusername/hireme/api/internal/validator"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logging
	setupLogging(cfg)

	slog.Info("starting HireMe API",
		"environment", cfg.Environment,
		"port", cfg.Server.Port,
	)

	// Connect to database
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	cvRepo := postgres.NewCVRepository(db)
	assetRepo := postgres.NewAssetRepository(db)

	// Initialize validator
	cvValidator, err := validator.NewCVValidator()
	if err != nil {
		slog.Error("failed to create CV validator", "error", err)
		os.Exit(1)
	}

	// Initialize services
	userSvc := service.NewUserService(userRepo)
	cvSvc := service.NewCVService(cvRepo, userRepo, cvValidator)
	assetSvc := service.NewAssetService(assetRepo, userRepo, cfg)

	// Initialize handlers
	h := handler.New(handler.Dependencies{
		Config:       cfg,
		UserService:  userSvc,
		CVService:    cvSvc,
		AssetService: assetSvc,
	})

	// Setup router
	r := setupRouter(cfg, h)

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func setupLogging(cfg *config.Config) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.Log.Level),
	}

	if cfg.Log.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func setupRouter(cfg *config.Config, h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(appMiddleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://*.hireme.io"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health endpoints (no auth)
	r.Get("/health", h.Health)
	r.Get("/ready", h.Ready)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth middleware
		r.Use(appMiddleware.Auth(cfg))

		// User routes
		r.Get("/users/me", h.GetCurrentUser)
		r.Patch("/users/me", h.UpdateCurrentUser)

		// CV routes
		r.Get("/cv", h.GetCV)
		r.Post("/cv", h.CreateCV)
		r.Put("/cv/{id}", h.UpdateCV)
		r.Delete("/cv/{id}", h.DeleteCV)

		// Asset routes
		r.Post("/assets", h.UploadAsset)
		r.Get("/assets/{id}", h.GetAsset)
		r.Delete("/assets/{id}", h.DeleteAsset)

		// Export routes
		r.Post("/export/{format}", h.CreateExport)
		r.Get("/export/{id}", h.GetExport)
	})

	return r
}
