package storage

import (
	"context"
	"fmt"
	"io"
)

// R2Storage implements Storage using Cloudflare R2 (S3-compatible)
//
// TODO: Implement when ready for production deployment
// Required dependency: go get github.com/aws/aws-sdk-go-v2/service/s3
// R2 endpoint format: https://<account_id>.r2.cloudflarestorage.com
type R2Storage struct {
	accountID string
	accessKey string
	secretKey string
	bucket    string
	publicURL string
	// TODO: Add s3.Client field when implementing
}

// R2Config holds Cloudflare R2 configuration
type R2Config struct {
	AccountID string
	AccessKey string
	SecretKey string
	Bucket    string
	PublicURL string
}

// NewR2Storage creates a new R2Storage instance
//
// TODO: Implementation steps:
// 1. Create AWS credentials from accessKey/secretKey
// 2. Configure S3 client with custom endpoint: https://{accountID}.r2.cloudflarestorage.com
// 3. Verify bucket exists or create it
func NewR2Storage(cfg R2Config) (*R2Storage, error) {
	if cfg.AccountID == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("R2 credentials are required")
	}

	return &R2Storage{
		accountID: cfg.AccountID,
		accessKey: cfg.AccessKey,
		secretKey: cfg.SecretKey,
		bucket:    cfg.Bucket,
		publicURL: cfg.PublicURL,
	}, nil
}

// Put stores a file in R2 and returns the storage key
//
// TODO: Implementation:
// 1. Use s3.PutObject with bucket and key
// 2. Set appropriate Content-Type header
// 3. Consider using multipart upload for files > 5MB
// 4. Return the key on success
func (s *R2Storage) Put(ctx context.Context, key string, reader io.Reader) (string, error) {
	return "", fmt.Errorf("R2Storage.Put not implemented: configure STORAGE_BACKEND=local for development")
}

// Get retrieves a file from R2
//
// TODO: Implementation:
// 1. Use s3.GetObject with bucket and path
// 2. Return the response body as io.ReadCloser
// 3. Handle NotFound errors appropriately
func (s *R2Storage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("R2Storage.Get not implemented: configure STORAGE_BACKEND=local for development")
}

// Delete removes a file from R2
//
// TODO: Implementation:
// 1. Use s3.DeleteObject with bucket and path
// 2. Return nil even if file doesn't exist (idempotent)
func (s *R2Storage) Delete(ctx context.Context, path string) error {
	return fmt.Errorf("R2Storage.Delete not implemented: configure STORAGE_BACKEND=local for development")
}

// Exists checks if a file exists in R2
//
// TODO: Implementation:
// 1. Use s3.HeadObject with bucket and path
// 2. Return true if no error, false if NotFound
func (s *R2Storage) Exists(ctx context.Context, path string) (bool, error) {
	return false, fmt.Errorf("R2Storage.Exists not implemented: configure STORAGE_BACKEND=local for development")
}

// GetPublicURL returns the public URL for an asset
// Returns empty string if R2_PUBLIC_URL is not configured
func (s *R2Storage) GetPublicURL(path string) string {
	if s.publicURL == "" {
		return ""
	}
	return s.publicURL + "/" + path
}
