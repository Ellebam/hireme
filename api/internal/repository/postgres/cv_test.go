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

func TestCVRepo_Create_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user first (CV has foreign key to user)
		user := createTestUserInTx(t, ctx, q)

		// Create a new CV
		content := json.RawMessage(`{"schemaVersion":"1.0.0","sections":[{"id":"personal","type":"personal","content":{}}]}`)
		cv := &domain.CV{
			UserID:        user.ID,
			Title:         "My First CV",
			SchemaVersion: "1.0.0",
			Content:       content,
		}

		err := repo.Create(ctx, cv)
		require.NoError(t, err)

		// Verify the CV was created with correct values
		assert.NotEqual(t, uuid.Nil, cv.ID, "ID should be generated")
		assert.NotEmpty(t, cv.CreatedAt, "CreatedAt should be set")
		assert.NotEmpty(t, cv.UpdatedAt, "UpdatedAt should be set")
		assert.Equal(t, user.ID, cv.UserID)
		assert.Equal(t, "My First CV", cv.Title)
		assert.Equal(t, "1.0.0", cv.SchemaVersion)

		// Verify CV can be retrieved
		fetched, err := repo.GetByID(ctx, cv.ID)
		require.NoError(t, err)
		assert.Equal(t, cv.Title, fetched.Title)
		assert.Equal(t, cv.SchemaVersion, fetched.SchemaVersion)
	})
}

func TestCVRepo_GetByID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user and CV
		user := createTestUserInTx(t, ctx, q)
		created := createTestCVInTx(t, ctx, q, user.ID)

		// Retrieve the CV by ID
		cv, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, cv.ID)
		assert.Equal(t, created.UserID, cv.UserID)
		assert.Equal(t, created.Title, cv.Title)
		assert.Equal(t, created.SchemaVersion, cv.SchemaVersion)
	})
}

func TestCVRepo_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &CVRepository{q: queries.New(tx)}

		// Try to get a non-existent CV
		nonExistentID := uuid.New()
		cv, err := repo.GetByID(ctx, nonExistentID)

		assert.Nil(t, cv)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestCVRepo_GetByUserID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user and CV
		user := createTestUserInTx(t, ctx, q)
		created := createTestCVInTx(t, ctx, q, user.ID, func(p *queries.CreateCVParams) {
			p.Title = "User's Active CV"
		})

		// Retrieve the CV by user ID
		cv, err := repo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, cv.ID)
		assert.Equal(t, user.ID, cv.UserID)
		assert.Equal(t, "User's Active CV", cv.Title)
	})
}

func TestCVRepo_GetByUserID_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create a user but no CV
		user := createTestUserInTx(t, ctx, q)

		// Try to get CV for user with no CVs
		cv, err := repo.GetByUserID(ctx, user.ID)

		assert.Nil(t, cv)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestCVRepo_ListByUserID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Create multiple CVs for the user
		createTestCVInTx(t, ctx, q, user.ID, func(p *queries.CreateCVParams) {
			p.Title = "CV 1"
		})
		createTestCVInTx(t, ctx, q, user.ID, func(p *queries.CreateCVParams) {
			p.Title = "CV 2"
		})
		createTestCVInTx(t, ctx, q, user.ID, func(p *queries.CreateCVParams) {
			p.Title = "CV 3"
		})

		// List CVs for user
		cvs, err := repo.ListByUserID(ctx, user.ID)
		require.NoError(t, err)

		assert.Len(t, cvs, 3, "Should return all 3 CVs")

		// Verify all CVs belong to the user
		for _, cv := range cvs {
			assert.Equal(t, user.ID, cv.UserID)
		}
	})
}

func TestCVRepo_ListByUserID_Empty(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create user but no CVs
		user := createTestUserInTx(t, ctx, q)

		// List CVs for user with no CVs
		cvs, err := repo.ListByUserID(ctx, user.ID)
		require.NoError(t, err)

		assert.Empty(t, cvs, "Should return empty slice for user with no CVs")
	})
}

func TestCVRepo_ListByUserID_OnlyUsersCVs(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create two users
		user1 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "user1@example.com"
		})
		user2 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "user2@example.com"
		})

		// Create CVs for both users
		createTestCVInTx(t, ctx, q, user1.ID, func(p *queries.CreateCVParams) {
			p.Title = "User1 CV"
		})
		createTestCVInTx(t, ctx, q, user2.ID, func(p *queries.CreateCVParams) {
			p.Title = "User2 CV"
		})

		// List should only return user1's CVs
		cvs, err := repo.ListByUserID(ctx, user1.ID)
		require.NoError(t, err)

		assert.Len(t, cvs, 1, "Should only return user1's CVs")
		assert.Equal(t, "User1 CV", cvs[0].Title)
	})
}

func TestCVRepo_CountByUserID_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Initially should be 0
		count, err := repo.CountByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)

		// Create some CVs
		createTestCVInTx(t, ctx, q, user.ID)
		createTestCVInTx(t, ctx, q, user.ID)

		// Count should be 2
		count, err = repo.CountByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestCVRepo_CountByUserID_OnlyUsersCVs(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create two users
		user1 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "count-user1@example.com"
		})
		user2 := createTestUserInTx(t, ctx, q, func(p *queries.CreateUserParams) {
			p.Email = "count-user2@example.com"
		})

		// Create CVs for both users
		createTestCVInTx(t, ctx, q, user1.ID)
		createTestCVInTx(t, ctx, q, user1.ID)
		createTestCVInTx(t, ctx, q, user2.ID)

		// Count should only return user1's CV count
		count, err := repo.CountByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, count, "Should only count user1's CVs")

		// User2 should have 1
		count2, err := repo.CountByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, count2, "User2 should have 1 CV")
	})
}

func TestCVRepo_Update_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user and CV
		user := createTestUserInTx(t, ctx, q)
		created := createTestCVInTx(t, ctx, q, user.ID)

		// Prepare updated CV
		cv := cvToDomain(created)
		cv.Title = "Updated CV Title"
		newContent := json.RawMessage(`{"schemaVersion":"1.0.0","sections":[{"id":"updated","type":"personal","content":{}}]}`)
		cv.Content = newContent

		// Update the CV
		err := repo.Update(ctx, cv)
		require.NoError(t, err)

		// Verify the update
		assert.Equal(t, "Updated CV Title", cv.Title)

		// Fetch again to verify persistence
		fetched, err := repo.GetByID(ctx, cv.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated CV Title", fetched.Title)

		// Verify content was updated
		var content map[string]interface{}
		err = json.Unmarshal(fetched.Content, &content)
		require.NoError(t, err)
		sections := content["sections"].([]interface{})
		assert.Len(t, sections, 1)
	})
}

func TestCVRepo_Update_NotFound(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &CVRepository{q: queries.New(tx)}

		// Try to update a non-existent CV
		cv := &domain.CV{
			ID:      uuid.New(),
			Title:   "Non-existent CV",
			Content: json.RawMessage(`{}`),
		}

		err := repo.Update(ctx, cv)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestCVRepo_Delete_Success(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user and CV
		user := createTestUserInTx(t, ctx, q)
		created := createTestCVInTx(t, ctx, q, user.ID)

		// Verify CV exists
		_, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		// Delete the CV
		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Verify CV is deleted
		cv, err := repo.GetByID(ctx, created.ID)
		assert.Nil(t, cv)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestCVRepo_Delete_NonExistent(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		repo := &CVRepository{q: queries.New(tx)}

		// Delete a non-existent CV (should not error)
		err := repo.Delete(ctx, uuid.New())
		// Note: pgx DELETE doesn't return error for non-existent rows
		assert.NoError(t, err)
	})
}

func TestCVRepo_ContentPersistence(t *testing.T) {
	pool := setupTestDB(t)

	withTestTx(t, pool, func(ctx context.Context, tx pgx.Tx) {
		q := queries.New(tx)
		repo := &CVRepository{q: q}

		// Create test user
		user := createTestUserInTx(t, ctx, q)

		// Create CV with complex content
		complexContent := json.RawMessage(`{
			"schemaVersion": "1.0.0",
			"templateId": "modern",
			"sections": [
				{
					"id": "personal",
					"type": "personal",
					"order": 0,
					"visible": true,
					"content": {
						"firstName": "John",
						"lastName": "Doe",
						"email": "john@example.com"
					}
				},
				{
					"id": "experience",
					"type": "experience",
					"order": 1,
					"visible": true,
					"content": {
						"entries": [
							{
								"company": "Tech Corp",
								"position": "Engineer"
							}
						]
					}
				}
			]
		}`)

		cv := &domain.CV{
			UserID:        user.ID,
			Title:         "Complex CV",
			SchemaVersion: "1.0.0",
			Content:       complexContent,
		}

		err := repo.Create(ctx, cv)
		require.NoError(t, err)

		// Fetch and verify content is preserved
		fetched, err := repo.GetByID(ctx, cv.ID)
		require.NoError(t, err)

		// Parse the content
		var content map[string]interface{}
		err = json.Unmarshal(fetched.Content, &content)
		require.NoError(t, err)

		assert.Equal(t, "1.0.0", content["schemaVersion"])
		assert.Equal(t, "modern", content["templateId"])

		sections := content["sections"].([]interface{})
		assert.Len(t, sections, 2)

		personalSection := sections[0].(map[string]interface{})
		assert.Equal(t, "personal", personalSection["type"])

		personalContent := personalSection["content"].(map[string]interface{})
		assert.Equal(t, "John", personalContent["firstName"])
		assert.Equal(t, "Doe", personalContent["lastName"])
	})
}
