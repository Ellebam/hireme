-- name: GetCVByID :one
SELECT * FROM cvs WHERE id = $1;

-- name: GetCVByUserID :one
-- Gets the user's active CV (for single-CV users)
SELECT * FROM cvs 
WHERE user_id = $1 AND is_active = true
ORDER BY updated_at DESC
LIMIT 1;

-- name: ListCVsByUserID :many
SELECT * FROM cvs 
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: CountCVsByUserID :one
SELECT COUNT(*) FROM cvs WHERE user_id = $1;

-- name: CreateCV :one
INSERT INTO cvs (
    user_id,
    title,
    schema_version,
    content
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateCV :one
UPDATE cvs SET
    title = COALESCE($2, title),
    content = COALESCE($3, content),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateCVContent :one
UPDATE cvs SET
    content = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetCVActive :exec
UPDATE cvs SET
    is_active = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteCV :exec
DELETE FROM cvs WHERE id = $1;

-- name: DeleteCVsByUserID :exec
DELETE FROM cvs WHERE user_id = $1;
