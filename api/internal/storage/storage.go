package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ellebam/hireme/api/internal/config"
)

// Storage defines the interface for file storage backends
type Storage interface {
	// Put stores a file and returns the storage path
	Put(ctx context.Context, key string, reader io.Reader) (string, error)

	// Get retrieves a file by its storage path
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes a file by its storage path
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)
}

// LocalStorage implements Storage using the local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("creating storage directory: %w", err)
	}

	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) Put(ctx context.Context, key string, reader io.Reader) (string, error) {
	path := filepath.Join(s.basePath, key)

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating directory: %w", err)
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Copy content
	if _, err := io.Copy(file, reader); err != nil {
		_ = os.Remove(path) // Clean up on error
		return "", fmt.Errorf("writing file: %w", err)
	}

	return key, nil
}

func (s *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return file, nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("deleting file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("checking file: %w", err)
	}

	return true, nil
}

// NewStorage creates the appropriate storage backend based on configuration
func NewStorage(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Backend {
	case "local":
		return NewLocalStorage(cfg.LocalPath)
	case "r2":
		return NewR2Storage(R2Config{
			AccountID: cfg.R2AccountID,
			AccessKey: cfg.R2AccessKey,
			SecretKey: cfg.R2SecretKey,
			Bucket:    cfg.R2Bucket,
			PublicURL: cfg.R2PublicURL,
		})
	default:
		return NewLocalStorage(cfg.LocalPath)
	}
}
