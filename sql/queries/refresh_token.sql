-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, hashed_token,user_id, expires_at)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE hashed_token = $1;

-- name: SetRefreshTokenInvalid :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: RevokeAllExpiredRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE expires_at < NOW();

-- name: RevokeAllRefreshTokensForUser :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = $1;

-- name: RevokeRefreshTokenByToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE hashed_token = $1;
