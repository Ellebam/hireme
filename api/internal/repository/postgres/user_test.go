package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/repository/postgres/queries"
)

func TestUserRepo_Create_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		// Create repository with transaction
		repo := &UserRepository{
			db: nil, // Not used when q is set
			q:  queries.New(tx),
		}

		// Create a new user
		user := &domain.User{
			ID:                "test-user-create-001",
			ExternalID:        "ext-create-001",
			Provider:          domain.ProviderGoogle,
			Email:             "create-test@example.com",
			EmailVerified:     true,
			DisplayName:       "Create Test User",
			Tier:              domain.TierFree,
			CVLimit:           1,
			StorageLimitBytes: 5 * 1024 * 1024,
			Locale:            domain.LocaleEN,
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Verify the user was created with correct values
		assert.NotEmpty(t, user.CreatedAt, "CreatedAt should be set")
		assert.NotEmpty(t, user.UpdatedAt, "UpdatedAt should be set")
		assert.Equal(t, "test-user-create-001", user.ID)
		assert.Equal(t, "create-test@example.com", user.Email)
		assert.Equal(t, domain.TierFree, user.Tier)
		assert.Equal(t, int64(0), user.StorageUsedBytes, "StorageUsedBytes should default to 0")

		// Verify user can be retrieved
		fetched, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, fetched.Email)
		assert.Equal(t, user.DisplayName, fetched.DisplayName)
	})
}

func TestUserRepo_GetByID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user first
		created := createTestUserInTx(t, ctx, q)

		// Retrieve the user by ID
		user, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, created.Email, user.Email)
		assert.Equal(t, created.ExternalID, user.ExternalID)
		assert.Equal(t, created.Provider, user.Provider)
	})
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &UserRepository{q: queries.New(tx)}

		// Try to get a non-existent user
		user, err := repo.GetByID(ctx, "non-existent-user-id")

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_GetByEmail_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user with specific email
		testEmail := "email-lookup-test@example.com"
		created := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = testEmail
		})

		// Retrieve the user by email
		user, err := repo.GetByEmail(ctx, testEmail)
		require.NoError(t, err)

		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, testEmail, user.Email)
	})
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &UserRepository{q: queries.New(tx)}

		// Try to get a user with non-existent email
		user, err := repo.GetByEmail(ctx, "nonexistent@example.com")

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_GetByExternalID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user with specific external ID
		testProvider := domain.ProviderGoogle
		testExternalID := "google-ext-12345"
		created := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Provider = testProvider
			p.ExternalID = testExternalID
		})

		// Retrieve the user by external ID
		user, err := repo.GetByExternalID(ctx, testProvider, testExternalID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, testProvider, user.Provider)
		assert.Equal(t, testExternalID, user.ExternalID)
	})
}

func TestUserRepo_GetByExternalID_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &UserRepository{q: queries.New(tx)}

		// Try to get a user with non-existent external ID
		user, err := repo.GetByExternalID(ctx, domain.ProviderGoogle, "non-existent-external-id")

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_Update_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user
		created := createTestUserInTx(t, ctx, q)

		// Prepare updated user
		user := userToDomain(created)
		user.DisplayName = "Updated Name"
		user.Email = "updated-email@example.com"
		user.Locale = domain.LocaleDE

		// Update the user
		err := repo.Update(ctx, user)
		require.NoError(t, err)

		// Verify the update
		assert.Equal(t, "Updated Name", user.DisplayName)
		assert.Equal(t, "updated-email@example.com", user.Email)
		assert.Equal(t, domain.LocaleDE, user.Locale)

		// Fetch again to verify persistence
		fetched, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", fetched.DisplayName)
		assert.Equal(t, "updated-email@example.com", fetched.Email)
		assert.Equal(t, domain.LocaleDE, fetched.Locale)
	})
}

func TestUserRepo_Update_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &UserRepository{q: queries.New(tx)}

		// Try to update a non-existent user
		user := &domain.User{
			ID:          "non-existent-user-id",
			Email:       "test@example.com",
			DisplayName: "Test",
			Locale:      domain.LocaleEN,
		}

		err := repo.Update(ctx, user)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_UpdateStorageUsed_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user
		created := createTestUserInTx(t, ctx, q)

		// Initial storage should be 0 (or whatever default is)
		user, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		initialStorage := user.StorageUsedBytes

		// Update storage used
		newStorageUsed := int64(1024 * 1024) // 1MB
		err = repo.UpdateStorageUsed(ctx, created.ID, newStorageUsed)
		require.NoError(t, err)

		// Verify the update
		updated, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.NotEqual(t, initialStorage, updated.StorageUsedBytes)
		assert.Equal(t, newStorageUsed, updated.StorageUsedBytes)
	})
}

func TestUserRepo_UpdateStorageUsed_Incremental(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user
		created := createTestUserInTx(t, ctx, q)

		// Set initial storage
		err := repo.UpdateStorageUsed(ctx, created.ID, 500)
		require.NoError(t, err)

		// Verify first update
		user, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(500), user.StorageUsedBytes)

		// Set to a new value (not incremental add, just replacement)
		err = repo.UpdateStorageUsed(ctx, created.ID, 1500)
		require.NoError(t, err)

		// Verify second update
		user, err = repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(1500), user.StorageUsedBytes)
	})
}

func TestUserRepo_Delete_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create test user
		created := createTestUserInTx(t, ctx, q)

		// Verify user exists
		_, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		// Delete the user
		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Verify user is deleted
		user, err := repo.GetByID(ctx, created.ID)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_Create_DuplicateProviderExternalID(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &UserRepository{q: q}

		// Create first user
		provider := domain.ProviderGoogle
		externalID := "google-123"
		createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Provider = provider
			p.ExternalID = externalID
		})

		// Try to create second user with same provider + external_id
		user := &domain.User{
			ID:                "different-user-id",
			ExternalID:        externalID, // Same external_id
			Provider:          provider,   // Same provider
			Email:             "different@example.com",
			EmailVerified:     true,
			DisplayName:       "Duplicate User",
			Tier:              domain.TierFree,
			CVLimit:           1,
			StorageLimitBytes: 5 * 1024 * 1024,
			Locale:            domain.LocaleEN,
		}

		err := repo.Create(ctx, user)
		assert.Error(t, err, "Should fail due to duplicate provider+external_id")
	})
}

func TestUserRepo_TierValues(t *testing.T) {
	pool := setupTestDB(t)

	testCases := []struct {
		name    string
		tier    string
		cvLimit int32
	}{
		{"free tier", domain.TierFree, 1},
		{"pro tier", domain.TierPro, 5},
		{"power tier", domain.TierPower, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
				q := queries.New(tx)
				repo := &UserRepository{q: q}

				// Create user with specific tier
				created := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
					p.Tier = &tc.tier
					p.CvLimit = &tc.cvLimit
				})

				// Retrieve and verify
				user, err := repo.GetByID(ctx, created.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.tier, user.Tier)
				assert.Equal(t, int(tc.cvLimit), user.CVLimit)
			})
		})
	}
}
