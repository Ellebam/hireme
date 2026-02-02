package domain

import "time"

// User represents a user account
type User struct {
	ID                string
	ExternalID        string
	Provider          string
	Email             string
	EmailVerified     bool
	DisplayName       string
	Tier              string
	CVLimit           int
	StorageLimitBytes int64
	StorageUsedBytes  int64
	Locale            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Tier constants
const (
	TierFree  = "free"
	TierPro   = "pro"
	TierPower = "power"
)

// Provider constants
const (
	ProviderGoogle      = "google"
	ProviderGitHub      = "github"
	ProviderDevelopment = "development"
)

// Locale constants
const (
	LocaleEN = "en"
	LocaleDE = "de"
)

// CanCreateCV checks if the user can create another CV
func (u *User) CanCreateCV(currentCount int) bool {
	return currentCount < u.CVLimit
}

// CanUploadAsset checks if the user can upload an asset of the given size
func (u *User) CanUploadAsset(assetSize int64) bool {
	return (u.StorageUsedBytes + assetSize) <= u.StorageLimitBytes
}

// RemainingStorage returns bytes remaining in user's storage quota
func (u *User) RemainingStorage() int64 {
	remaining := u.StorageLimitBytes - u.StorageUsedBytes
	if remaining < 0 {
		return 0
	}
	return remaining
}
