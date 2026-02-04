package postgres

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/ellebam/hireme/api/internal/domain"
	"github.com/ellebam/hireme/api/internal/repository/postgres/queries"
)

// testDB holds the shared database connection pool for integration tests
var testDB *pgxpool.Pool

// setupTestDB returns a database connection pool for integration tests.
// It skips the test if no database URL is configured.
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// Try TEST_DATABASE_URL first, then fall back to DATABASE_URL
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		t.Skip("Skipping integration test: no database URL configured (set TEST_DATABASE_URL or DATABASE_URL)")
	}

	// Reuse existing connection if available
	if testDB != nil {
		// Verify connection is still alive
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := testDB.Ping(ctx); err == nil {
			return testDB
		}
		// Connection is dead, close it and create new one
		testDB.Close()
		testDB = nil
	}

	// Create new connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		t.Skipf("Skipping integration test: invalid database URL: %v", err)
	}

	// Configure pool for testing
	config.MaxConns = 5
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Skipf("Skipping integration test: could not connect to database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("Skipping integration test: could not ping database: %v", err)
	}

	testDB = pool
	return pool
}

// withTestTx runs a test function within a transaction that is rolled back after completion.
// This provides test isolation without affecting other tests.
func withTestTx(t *testing.T, pool *pgxpool.Pool, fn func(ctx context.Context, tx pgx.Tx)) {
	t.Helper()

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	require.NoError(t, err, "failed to begin transaction")

	defer func() {
		// Always rollback to ensure test isolation
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			t.Logf("Warning: failed to rollback transaction: %v", err)
		}
	}()

	fn(ctx, tx)
}

// testQueries creates a queries.Queries instance backed by a transaction
func testQueries(tx pgx.Tx) *queries.Queries {
	return queries.New(tx)
}

// createTestUser creates a test user in the database within the given transaction
func createTestUserInTx(t *testing.T, ctx context.Context, q *queries.Queries, overrides ...func(*queries.CreateUserParams)) queries.User {
	t.Helper()

	// Generate unique ID for this test
	id := uuid.New().String()
	emailVerified := true
	displayName := "Test User"
	tier := domain.TierFree
	cvLimit := int32(1)
	storageLimitBytes := int64(5 * 1024 * 1024) // 5MB
	locale := domain.LocaleEN

	params := queries.CreateUserParams{
		ID:                id,
		ExternalID:        "ext-" + id[:8],
		Provider:          domain.ProviderDevelopment,
		Email:             "test-" + id[:8] + "@example.com",
		EmailVerified:     &emailVerified,
		DisplayName:       &displayName,
		Tier:              &tier,
		CvLimit:           &cvLimit,
		StorageLimitBytes: &storageLimitBytes,
		Locale:            &locale,
	}

	// Apply any overrides
	for _, override := range overrides {
		override(&params)
	}

	user, err := q.CreateUser(ctx, params)
	require.NoError(t, err, "failed to create test user")

	return user
}

// createTestCVInTx creates a test CV in the database within the given transaction
func createTestCVInTx(t *testing.T, ctx context.Context, q *queries.Queries, userID string, overrides ...func(*queries.CreateCVParams)) queries.Cv {
	t.Helper()

	content := json.RawMessage(`{"schemaVersion":"1.0.0","sections":[]}`)

	params := queries.CreateCVParams{
		UserID:        userID,
		Title:         "Test CV",
		SchemaVersion: "1.0.0",
		Content:       content,
	}

	// Apply any overrides
	for _, override := range overrides {
		override(&params)
	}

	cv, err := q.CreateCV(ctx, params)
	require.NoError(t, err, "failed to create test CV")

	return cv
}

// createTestAssetInTx creates a test asset in the database within the given transaction
func createTestAssetInTx(t *testing.T, ctx context.Context, q *queries.Queries, userID string, overrides ...func(*queries.CreateAssetParams)) queries.Asset {
	t.Helper()

	// Generate unique filename
	filename := uuid.New().String() + ".jpg"
	metadata := []byte(`{}`)

	params := queries.CreateAssetParams{
		UserID:           userID,
		Filename:         filename,
		OriginalFilename: "photo.jpg",
		MimeType:         "image/jpeg",
		SizeBytes:        1024,
		StoragePath:      userID + "/2024-01/" + filename,
		StorageBackend:   domain.StorageBackendLocal,
		Checksum:         "sha256-" + uuid.New().String()[:16],
		Metadata:         metadata,
	}

	// Apply any overrides
	for _, override := range overrides {
		override(&params)
	}

	asset, err := q.CreateAsset(ctx, params)
	require.NoError(t, err, "failed to create test asset")

	return asset
}

// ptr is a helper to create pointers to values
func ptr[T any](v T) *T {
	return &v
}
