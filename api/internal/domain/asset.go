package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Asset represents an uploaded file (typically an image)
type Asset struct {
	ID               uuid.UUID
	UserID           string
	Filename         string
	OriginalFilename string
	MimeType         string
	SizeBytes        int64
	StoragePath      string
	StorageBackend   string
	Checksum         string
	Width            *int
	Height           *int
	Metadata         json.RawMessage
	CreatedAt        time.Time
}

// Storage backend constants
const (
	StorageBackendLocal = "local"
	StorageBackendR2    = "r2"
)

// Allowed MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

// IsAllowedImageType checks if the MIME type is allowed
func IsAllowedImageType(mimeType string) bool {
	return AllowedImageTypes[mimeType]
}

// GetExtension returns the file extension for a MIME type
func GetExtension(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

// ExportJob represents a document export job
type ExportJob struct {
	ID           uuid.UUID
	UserID       string
	CVID         uuid.UUID
	Format       string
	Status       string
	ResultPath   *string
	ErrorMessage *string
	CreatedAt    time.Time
	StartedAt    *time.Time
	CompletedAt  *time.Time
}

// Export format constants
const (
	ExportFormatPDF  = "pdf"
	ExportFormatDOCX = "docx"
	ExportFormatJSON = "json"
	ExportFormatYAML = "yaml"
)

// Export status constants
const (
	ExportStatusPending    = "pending"
	ExportStatusProcessing = "processing"
	ExportStatusCompleted  = "completed"
	ExportStatusFailed     = "failed"
)

// IsValidExportFormat checks if the format is supported
func IsValidExportFormat(format string) bool {
	switch format {
	case ExportFormatPDF, ExportFormatDOCX, ExportFormatJSON, ExportFormatYAML:
		return true
	default:
		return false
	}
}
