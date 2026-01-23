-- name: StoreSessionID :one
INSERT INTO session_ids(id, created_at, updated_at, user_id, expires_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetSessionIDByID :one
SELECT * FROM session_ids
WHERE id = $1;

-- name: SetSessionIDInvalid :exec
UPDATE session_ids
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: RevokeAllExpiredSessionIDs :exec
UPDATE session_ids
SET revoked_at = NOW(), updated_at = NOW()
WHERE expires_at < NOW();

-- name: RevokeAllSessionsForUser :exec
DELETE FROM session_ids
WHERE user_id = $1;