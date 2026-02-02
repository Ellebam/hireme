package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Environment string

	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Storage  StorageConfig
	Export   ExportConfig
	Limits   LimitsConfig
	Features FeaturesConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	URL            string
	MaxConnections int
	EncryptionKey  string
}

type AuthConfig struct {
	BypassEnabled bool
	BypassUserID  string
	JWTSecret     string
	JWTExpiration int // hours

	Google GoogleOIDCConfig
}

type GoogleOIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type StorageConfig struct {
	Backend   string // "local" or "r2"
	LocalPath string

	// Cloudflare R2
	R2AccountID string
	R2AccessKey string
	R2SecretKey string
	R2Bucket    string
	R2PublicURL string
}

type ExportConfig struct {
	GotenbergURL string
}

type LimitsConfig struct {
	MaxAssetSizeBytes  int64
	MaxImageDimension  int
	AllowedImageTypes  []string
	MaxCVsFree         int
	MaxCVsPro          int
	MaxCVsPower        int
}

type FeaturesConfig struct {
	EnableAds        bool
	EnableExportPDF  bool
	EnableExportDOCX bool
}

type LogConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		Environment: getEnv("APP_ENV", "development"),

		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 30),
		},

		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", ""),
			MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 25),
			EncryptionKey:  getEnv("DB_ENCRYPTION_KEY", ""),
		},

		Auth: AuthConfig{
			BypassEnabled: getEnvBool("AUTH_BYPASS_ENABLED", false),
			BypassUserID:  getEnv("AUTH_BYPASS_USER_ID", "dev-user-001"),
			JWTSecret:     getEnv("JWT_SECRET", ""),
			JWTExpiration: getEnvInt("JWT_EXPIRATION_HOURS", 24),
			Google: GoogleOIDCConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/api/auth/callback/google"),
			},
		},

		Storage: StorageConfig{
			Backend:     getEnv("STORAGE_BACKEND", "local"),
			LocalPath:   getEnv("STORAGE_LOCAL_PATH", "./data/uploads"),
			R2AccountID: getEnv("R2_ACCOUNT_ID", ""),
			R2AccessKey: getEnv("R2_ACCESS_KEY", ""),
			R2SecretKey: getEnv("R2_SECRET_KEY", ""),
			R2Bucket:    getEnv("R2_BUCKET", "hireme-assets"),
			R2PublicURL: getEnv("R2_PUBLIC_URL", ""),
		},

		Export: ExportConfig{
			GotenbergURL: getEnv("GOTENBERG_URL", "http://localhost:3001"),
		},

		Limits: LimitsConfig{
			MaxAssetSizeBytes: getEnvInt64("MAX_ASSET_SIZE_BYTES", 2097152), // 2MB
			MaxImageDimension: getEnvInt("MAX_IMAGE_DIMENSION", 2000),
			AllowedImageTypes: []string{"image/jpeg", "image/png", "image/webp"},
			MaxCVsFree:        getEnvInt("MAX_CVS_FREE_TIER", 1),
			MaxCVsPro:         getEnvInt("MAX_CVS_PRO_TIER", 5),
			MaxCVsPower:       getEnvInt("MAX_CVS_POWER_TIER", 50),
		},

		Features: FeaturesConfig{
			EnableAds:        getEnvBool("FEATURE_ADS", false),
			EnableExportPDF:  getEnvBool("FEATURE_EXPORT_PDF", true),
			EnableExportDOCX: getEnvBool("FEATURE_EXPORT_DOCX", true),
		},

		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}

	// Validate required configuration
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if !c.Auth.BypassEnabled && c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required when auth bypass is disabled")
	}

	if c.Storage.Backend == "r2" {
		if c.Storage.R2AccountID == "" || c.Storage.R2AccessKey == "" || c.Storage.R2SecretKey == "" {
			return fmt.Errorf("R2 credentials are required when STORAGE_BACKEND=r2")
		}
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
