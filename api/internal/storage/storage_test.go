package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage_Put_Success(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("test file content")
	key := "test-file.txt"

	// Put the file
	returnedKey, err := storage.Put(ctx, key, bytes.NewReader(content))
	require.NoError(t, err)
	assert.Equal(t, key, returnedKey)

	// Verify file exists on disk
	fullPath := filepath.Join(tmpDir, key)
	_, err = os.Stat(fullPath)
	assert.NoError(t, err, "file should exist on disk")

	// Verify content matches
	diskContent, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Equal(t, content, diskContent)
}

func TestLocalStorage_Put_CreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("nested file content")
	key := "nested/deep/path/test-file.txt"

	// Put the file in a nested path
	returnedKey, err := storage.Put(ctx, key, bytes.NewReader(content))
	require.NoError(t, err)
	assert.Equal(t, key, returnedKey)

	// Verify parent directory was created
	parentDir := filepath.Join(tmpDir, "nested", "deep", "path")
	info, err := os.Stat(parentDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "parent directory should be created")

	// Verify file exists
	fullPath := filepath.Join(tmpDir, key)
	_, err = os.Stat(fullPath)
	assert.NoError(t, err, "file should exist on disk")
}

func TestLocalStorage_Get_Success(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("content to retrieve")
	key := "retrieve-me.txt"

	// Put the file first
	_, err = storage.Put(ctx, key, bytes.NewReader(content))
	require.NoError(t, err)

	// Get the file
	reader, err := storage.Get(ctx, key)
	require.NoError(t, err)
	defer func() { _ = reader.Close() }()

	// Read and verify content
	retrievedContent, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, content, retrievedContent)
}

func TestLocalStorage_Get_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Try to get a non-existent file
	reader, err := storage.Get(ctx, "non-existent-file.txt")
	assert.Error(t, err)
	assert.Nil(t, reader)
	assert.Contains(t, err.Error(), "file not found")
}

func TestLocalStorage_Delete_Success(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("delete me")
	key := "to-delete.txt"

	// Put the file first
	_, err = storage.Put(ctx, key, bytes.NewReader(content))
	require.NoError(t, err)

	// Verify file exists
	fullPath := filepath.Join(tmpDir, key)
	_, err = os.Stat(fullPath)
	require.NoError(t, err, "file should exist before deletion")

	// Delete the file
	err = storage.Delete(ctx, key)
	require.NoError(t, err)

	// Verify file is removed
	_, err = os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err), "file should be removed after deletion")
}

func TestLocalStorage_Delete_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Delete a non-existent file (should be idempotent - no error)
	err = storage.Delete(ctx, "non-existent-file.txt")
	assert.NoError(t, err, "deleting non-existent file should not return error (idempotent)")
}

func TestLocalStorage_Exists_True(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("i exist")
	key := "existing-file.txt"

	// Put the file first
	_, err = storage.Put(ctx, key, bytes.NewReader(content))
	require.NoError(t, err)

	// Check existence
	exists, err := storage.Exists(ctx, key)
	require.NoError(t, err)
	assert.True(t, exists, "file should exist")
}

func TestLocalStorage_Exists_False(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorage(tmpDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Check existence of non-existent file
	exists, err := storage.Exists(ctx, "non-existent-file.txt")
	require.NoError(t, err)
	assert.False(t, exists, "file should not exist")
}

func TestNewLocalStorage_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	newPath := filepath.Join(tmpDir, "new-storage-dir")

	// Verify directory doesn't exist yet
	_, err := os.Stat(newPath)
	require.True(t, os.IsNotExist(err))

	// Create storage - should create the directory
	storage, err := NewLocalStorage(newPath)
	require.NoError(t, err)
	require.NotNil(t, storage)

	// Verify directory was created
	info, err := os.Stat(newPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
