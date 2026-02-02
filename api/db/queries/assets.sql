-- name: GetAssetByID :one
SELECT * FROM assets WHERE id = $1;

-- name: GetAssetByChecksum :one
SELECT * FROM assets 
WHERE user_id = $1 AND checksum = $2;

-- name: ListAssetsByUserID :many
SELECT * FROM assets 
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTotalAssetSizeByUserID :one
SELECT COALESCE(SUM(size_bytes), 0)::bigint as total_size
FROM assets 
WHERE user_id = $1;

-- name: CreateAsset :one
INSERT INTO assets (
    user_id,
    filename,
    original_filename,
    mime_type,
    size_bytes,
    storage_path,
    storage_backend,
    checksum,
    width,
    height,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: DeleteAsset :exec
DELETE FROM assets WHERE id = $1;

-- name: DeleteAssetsByUserID :exec
DELETE FROM assets WHERE user_id = $1;
