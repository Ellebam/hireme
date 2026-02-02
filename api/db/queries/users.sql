-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByExternalID :one
SELECT * FROM users WHERE provider = $1 AND external_id = $2;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
    id,
    external_id,
    provider,
    email,
    email_verified,
    display_name,
    tier,
    cv_limit,
    storage_limit_bytes,
    locale
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    email = COALESCE($2, email),
    email_verified = COALESCE($3, email_verified),
    display_name = COALESCE($4, display_name),
    locale = COALESCE($5, locale),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserTier :one
UPDATE users SET
    tier = $2,
    cv_limit = $3,
    storage_limit_bytes = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserStorageUsed :exec
UPDATE users SET
    storage_used_bytes = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
