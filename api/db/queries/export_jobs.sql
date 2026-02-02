-- name: GetExportJobByID :one
SELECT * FROM export_jobs WHERE id = $1;

-- name: ListPendingExportJobs :many
SELECT * FROM export_jobs 
WHERE status IN ('pending', 'processing')
ORDER BY created_at ASC
LIMIT $1;

-- name: ListExportJobsByUserID :many
SELECT * FROM export_jobs 
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: CreateExportJob :one
INSERT INTO export_jobs (
    user_id,
    cv_id,
    format,
    status
) VALUES (
    $1, $2, $3, 'pending'
) RETURNING *;

-- name: UpdateExportJobStatus :exec
UPDATE export_jobs SET
    status = $2,
    result_path = $3,
    error_message = $4,
    started_at = CASE WHEN $2 = 'processing' THEN NOW() ELSE started_at END,
    completed_at = CASE WHEN $2 IN ('completed', 'failed') THEN NOW() ELSE completed_at END
WHERE id = $1;

-- name: DeleteOldExportJobs :exec
DELETE FROM export_jobs
WHERE completed_at < NOW() - INTERVAL '7 days';
