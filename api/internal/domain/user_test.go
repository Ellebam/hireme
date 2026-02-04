package domain

import (
	"testing"
)

func TestUser_CanCreateCV_BelowLimit(t *testing.T) {
	user := &User{
		CVLimit: 1,
	}

	if !user.CanCreateCV(0) {
		t.Error("expected CanCreateCV to return true when count is below limit")
	}
}

func TestUser_CanCreateCV_AtLimit(t *testing.T) {
	user := &User{
		CVLimit: 1,
	}

	if user.CanCreateCV(1) {
		t.Error("expected CanCreateCV to return false when count equals limit")
	}
}

func TestUser_CanCreateCV(t *testing.T) {
	tests := []struct {
		name         string
		cvLimit      int
		currentCount int
		expected     bool
	}{
		{
			name:         "below limit",
			cvLimit:      3,
			currentCount: 0,
			expected:     true,
		},
		{
			name:         "one below limit",
			cvLimit:      3,
			currentCount: 2,
			expected:     true,
		},
		{
			name:         "at limit",
			cvLimit:      3,
			currentCount: 3,
			expected:     false,
		},
		{
			name:         "above limit",
			cvLimit:      3,
			currentCount: 5,
			expected:     false,
		},
		{
			name:         "zero limit",
			cvLimit:      0,
			currentCount: 0,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{CVLimit: tt.cvLimit}
			got := user.CanCreateCV(tt.currentCount)
			if got != tt.expected {
				t.Errorf("CanCreateCV(%d) = %v, want %v", tt.currentCount, got, tt.expected)
			}
		})
	}
}

func TestUser_CanUploadAsset_HasSpace(t *testing.T) {
	user := &User{
		StorageUsedBytes:  0,
		StorageLimitBytes: 5 * 1024 * 1024, // 5MB
	}

	assetSize := int64(1 * 1024 * 1024) // 1MB
	if !user.CanUploadAsset(assetSize) {
		t.Error("expected CanUploadAsset to return true when there is space")
	}
}

func TestUser_CanUploadAsset_NoSpace(t *testing.T) {
	user := &User{
		StorageUsedBytes:  5 * 1024 * 1024, // 5MB used
		StorageLimitBytes: 5 * 1024 * 1024, // 5MB limit
	}

	assetSize := int64(1 * 1024 * 1024) // 1MB
	if user.CanUploadAsset(assetSize) {
		t.Error("expected CanUploadAsset to return false when there is no space")
	}
}

func TestUser_CanUploadAsset(t *testing.T) {
	const MB = 1024 * 1024

	tests := []struct {
		name      string
		used      int64
		limit     int64
		assetSize int64
		expected  bool
	}{
		{
			name:      "plenty of space",
			used:      0,
			limit:     5 * MB,
			assetSize: 1 * MB,
			expected:  true,
		},
		{
			name:      "exact fit",
			used:      4 * MB,
			limit:     5 * MB,
			assetSize: 1 * MB,
			expected:  true,
		},
		{
			name:      "one byte over",
			used:      4*MB + 1,
			limit:     5 * MB,
			assetSize: 1 * MB,
			expected:  false,
		},
		{
			name:      "no space left",
			used:      5 * MB,
			limit:     5 * MB,
			assetSize: 1,
			expected:  false,
		},
		{
			name:      "zero size asset",
			used:      5 * MB,
			limit:     5 * MB,
			assetSize: 0,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				StorageUsedBytes:  tt.used,
				StorageLimitBytes: tt.limit,
			}
			got := user.CanUploadAsset(tt.assetSize)
			if got != tt.expected {
				t.Errorf("CanUploadAsset(%d) = %v, want %v", tt.assetSize, got, tt.expected)
			}
		})
	}
}

func TestUser_RemainingStorage(t *testing.T) {
	const MB = 1024 * 1024

	tests := []struct {
		name     string
		used     int64
		limit    int64
		expected int64
	}{
		{
			name:     "partial usage",
			used:     3 * MB,
			limit:    5 * MB,
			expected: 2 * MB,
		},
		{
			name:     "no usage",
			used:     0,
			limit:    5 * MB,
			expected: 5 * MB,
		},
		{
			name:     "full usage",
			used:     5 * MB,
			limit:    5 * MB,
			expected: 0,
		},
		{
			name:     "over limit returns zero",
			used:     6 * MB,
			limit:    5 * MB,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				StorageUsedBytes:  tt.used,
				StorageLimitBytes: tt.limit,
			}
			got := user.RemainingStorage()
			if got != tt.expected {
				t.Errorf("RemainingStorage() = %d, want %d", got, tt.expected)
			}
		})
	}
}
