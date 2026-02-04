package postgres

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/repository/postgres/queries"
)

func TestAssetRepo_Create_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user first (asset has foreign key to user)
		user := createTestUserInTx(t, ctx, q)

		// Create a new asset
		width := 800
		height := 600
		asset := &domain.Asset{
			UserID:           user.ID,
			Filename:         "profile-image.jpg",
			OriginalFilename: "my-photo.jpg",
			MimeType:         "image/jpeg",
			SizeBytes:        2048,
			StoragePath:      user.ID + "/2024-01/profile-image.jpg",
			StorageBackend:   domain.StorageBackendLocal,
			Checksum:         "sha256-abc123def456",
			Width:            &width,
			Height:           &height,
			Metadata:         json.RawMessage(`{"source": "upload"}`),
		}

		err := repo.Create(ctx, asset)
		require.NoError(t, err)

		// Verify the asset was created with correct values
		assert.NotEqual(t, uuid.Nil, asset.ID, "ID should be generated")
		assert.NotEmpty(t, asset.CreatedAt, "CreatedAt should be set")
		assert.Equal(t, user.ID, asset.UserID)
		assert.Equal(t, "profile-image.jpg", asset.Filename)
		assert.Equal(t, "my-photo.jpg", asset.OriginalFilename)
		assert.Equal(t, "image/jpeg", asset.MimeType)
		assert.Equal(t, int64(2048), asset.SizeBytes)
		assert.Equal(t, &width, asset.Width)
		assert.Equal(t, &height, asset.Height)

		// Verify asset can be retrieved
		fetched, err := repo.GetByID(ctx, asset.ID)
		require.NoError(t, err)
		assert.Equal(t, asset.Filename, fetched.Filename)
		assert.Equal(t, asset.Checksum, fetched.Checksum)
	})
}

func TestAssetRepo_Create_WithoutDimensions(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Create asset without width/height (simulates failed dimension extraction)
		asset := &domain.Asset{
			UserID:           user.ID,
			Filename:         "corrupted.jpg",
			OriginalFilename: "photo.jpg",
			MimeType:         "image/jpeg",
			SizeBytes:        10240,
			StoragePath:      user.ID + "/images/corrupted.jpg",
			StorageBackend:   domain.StorageBackendLocal,
			Checksum:         "sha256-corrupted123",
			Width:            nil,
			Height:           nil,
		}

		err := repo.Create(ctx, asset)
		require.NoError(t, err)

		// Verify dimensions are nil
		fetched, err := repo.GetByID(ctx, asset.ID)
		require.NoError(t, err)
		assert.Nil(t, fetched.Width)
		assert.Nil(t, fetched.Height)
	})
}

func TestAssetRepo_GetByID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user and asset
		user := createTestUserInTx(t, ctx, q)
		created := createTestAssetInTx(t, ctx, q, user.ID)

		// Retrieve the asset by ID
		asset, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, asset.ID)
		assert.Equal(t, created.UserID, asset.UserID)
		assert.Equal(t, created.Filename, asset.Filename)
		assert.Equal(t, created.MimeType, asset.MimeType)
		assert.Equal(t, int64(created.SizeBytes), asset.SizeBytes)
	})
}

func TestAssetRepo_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &AssetRepository{q: queries.New(tx)}

		// Try to get a non-existent asset
		nonExistentID := uuid.New()
		asset, err := repo.GetByID(ctx, nonExistentID)

		assert.Nil(t, asset)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestAssetRepo_GetByChecksum_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user and asset with specific checksum
		user := createTestUserInTx(t, ctx, q)
		testChecksum := "sha256-unique-checksum-123"
		created := createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.Checksum = testChecksum
		})

		// Retrieve the asset by checksum
		asset, err := repo.GetByChecksum(ctx, user.ID, testChecksum)
		require.NoError(t, err)

		assert.Equal(t, created.ID, asset.ID)
		assert.Equal(t, testChecksum, asset.Checksum)
	})
}

func TestAssetRepo_GetByChecksum_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create a user
		user := createTestUserInTx(t, ctx, q)

		// Try to get an asset with non-existent checksum
		asset, err := repo.GetByChecksum(ctx, user.ID, "non-existent-checksum")

		assert.Nil(t, asset)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestAssetRepo_GetByChecksum_DifferentUser(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create two users
		user1 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "asset-user1@example.com"
		})
		user2 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "asset-user2@example.com"
		})

		// Create asset for user1 with specific checksum
		testChecksum := "sha256-user1-checksum"
		createTestAssetInTx(t, ctx, q, user1.ID, func(p *queries.CreateAssetParams) {
			p.Checksum = testChecksum
		})

		// Try to find the asset by user2 - should not find it
		asset, err := repo.GetByChecksum(ctx, user2.ID, testChecksum)

		assert.Nil(t, asset)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestAssetRepo_ListByUserID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Create multiple assets for the user
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.Filename = "image1.jpg"
		})
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.Filename = "image2.jpg"
		})
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.Filename = "image3.jpg"
		})

		// List assets for user
		assets, err := repo.ListByUserID(ctx, user.ID)
		require.NoError(t, err)

		assert.Len(t, assets, 3, "Should return all 3 assets")

		// Verify all assets belong to the user
		for _, asset := range assets {
			assert.Equal(t, user.ID, asset.UserID)
		}
	})
}

func TestAssetRepo_ListByUserID_Empty(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create user but no assets
		user := createTestUserInTx(t, ctx, q)

		// List assets for user with no assets
		assets, err := repo.ListByUserID(ctx, user.ID)
		require.NoError(t, err)

		assert.Empty(t, assets, "Should return empty slice for user with no assets")
	})
}

func TestAssetRepo_ListByUserID_OnlyUsersAssets(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create two users
		user1 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "list-user1@example.com"
		})
		user2 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "list-user2@example.com"
		})

		// Create assets for both users
		createTestAssetInTx(t, ctx, q, user1.ID, func(p *queries.CreateAssetParams) {
			p.Filename = "user1-asset.jpg"
		})
		createTestAssetInTx(t, ctx, q, user2.ID, func(p *queries.CreateAssetParams) {
			p.Filename = "user2-asset.jpg"
		})

		// List should only return user1's assets
		assets, err := repo.ListByUserID(ctx, user1.ID)
		require.NoError(t, err)

		assert.Len(t, assets, 1, "Should only return user1's assets")
		assert.Equal(t, "user1-asset.jpg", assets[0].Filename)
	})
}

func TestAssetRepo_GetTotalSizeByUserID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Initially should be 0
		totalSize, err := repo.GetTotalSizeByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), totalSize)

		// Create assets with known sizes
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 1000
		})
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 2000
		})
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 3000
		})

		// Total should be 6000
		totalSize, err = repo.GetTotalSizeByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(6000), totalSize)
	})
}

func TestAssetRepo_GetTotalSizeByUserID_OnlyUsersAssets(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create two users
		user1 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "size-user1@example.com"
		})
		user2 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "size-user2@example.com"
		})

		// Create assets for both users
		createTestAssetInTx(t, ctx, q, user1.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 1000
		})
		createTestAssetInTx(t, ctx, q, user1.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 2000
		})
		createTestAssetInTx(t, ctx, q, user2.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 5000
		})

		// User1's total should be 3000
		totalSize, err := repo.GetTotalSizeByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), totalSize, "Should only sum user1's assets")

		// User2's total should be 5000
		totalSize2, err := repo.GetTotalSizeByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(5000), totalSize2, "User2 should have 5000 bytes")
	})
}

func TestAssetRepo_Delete_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user and asset
		user := createTestUserInTx(t, ctx, q)
		created := createTestAssetInTx(t, ctx, q, user.ID)

		// Verify asset exists
		_, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		// Delete the asset
		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Verify asset is deleted
		asset, err := repo.GetByID(ctx, created.ID)
		assert.Nil(t, asset)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestAssetRepo_Delete_NonExistent(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &AssetRepository{q: queries.New(tx)}

		// Delete a non-existent asset (should not error)
		err := repo.Delete(ctx, uuid.New())
		// Note: pgx DELETE doesn't return error for non-existent rows
		assert.NoError(t, err)
	})
}

func TestAssetRepo_Delete_UpdatesTotalSize(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Create two assets
		asset1 := createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 1000
		})
		createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
			p.SizeBytes = 2000
		})

		// Initial total should be 3000
		totalSize, err := repo.GetTotalSizeByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), totalSize)

		// Delete one asset
		err = repo.Delete(ctx, asset1.ID)
		require.NoError(t, err)

		// Total should now be 2000
		totalSize, err = repo.GetTotalSizeByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(2000), totalSize)
	})
}

func TestAssetRepo_MetadataPersistence(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &AssetRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Create asset with complex metadata
		complexMetadata := json.RawMessage(`{
			"source": "camera",
			"camera": {
				"make": "Canon",
				"model": "EOS 5D"
			},
			"exif": {
				"iso": 400,
				"aperture": "f/2.8",
				"exposureTime": "1/500"
			},
			"tags": ["portrait", "professional"]
		}`)

		asset := &domain.Asset{
			UserID:           user.ID,
			Filename:         "photo-with-metadata.jpg",
			OriginalFilename: "IMG_0001.jpg",
			MimeType:         "image/jpeg",
			SizeBytes:        4096,
			StoragePath:      user.ID + "/photos/photo-with-metadata.jpg",
			StorageBackend:   domain.StorageBackendLocal,
			Checksum:         "sha256-metadata-test",
			Metadata:         complexMetadata,
		}

		err := repo.Create(ctx, asset)
		require.NoError(t, err)

		// Fetch and verify metadata is preserved
		fetched, err := repo.GetByID(ctx, asset.ID)
		require.NoError(t, err)

		// Parse the metadata
		var metadata map[string]interface{}
		err = json.Unmarshal(fetched.Metadata, &metadata)
		require.NoError(t, err)

		assert.Equal(t, "camera", metadata["source"])

		camera := metadata["camera"].(map[string]interface{})
		assert.Equal(t, "Canon", camera["make"])
		assert.Equal(t, "EOS 5D", camera["model"])

		exif := metadata["exif"].(map[string]interface{})
		assert.Equal(t, float64(400), exif["iso"])
		assert.Equal(t, "f/2.8", exif["aperture"])

		tags := metadata["tags"].([]interface{})
		assert.Len(t, tags, 2)
		assert.Equal(t, "portrait", tags[0])
		assert.Equal(t, "professional", tags[1])
	})
}

func TestAssetRepo_DifferentStorageBackends(t *testing.T) {
	pool := setupTestDB(t)

	testCases := []struct {
		name    string
		backend string
	}{
		{"local storage", domain.StorageBackendLocal},
		{"r2 storage", domain.StorageBackendR2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
				q := queries.New(tx)
				repo := &AssetRepository{q: q}

				// Create test user
				user := createTestUserInTx(t, ctx, q)

				// Create asset with specific storage backend
				created := createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
					p.StorageBackend = tc.backend
				})

				// Retrieve and verify
				asset, err := repo.GetByID(ctx, created.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.backend, asset.StorageBackend)
			})
		})
	}
}

func TestAssetRepo_MimeTypes(t *testing.T) {
	pool := setupTestDB(t)

	testCases := []struct {
		name     string
		mimeType string
		filename string
	}{
		{"JPEG image", "image/jpeg", "photo.jpg"},
		{"PNG image", "image/png", "logo.png"},
		{"WebP image", "image/webp", "optimized.webp"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
				q := queries.New(tx)
				repo := &AssetRepository{q: q}

				// Create test user
				user := createTestUserInTx(t, ctx, q)

				// Create asset with specific mime type
				created := createTestAssetInTx(t, ctx, q, user.ID, func(p *queries.CreateAssetParams) {
					p.MimeType = tc.mimeType
					p.Filename = tc.filename
				})

				// Retrieve and verify
				asset, err := repo.GetByID(ctx, created.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.mimeType, asset.MimeType)
				assert.Equal(t, tc.filename, asset.Filename)
			})
		})
	}
}
