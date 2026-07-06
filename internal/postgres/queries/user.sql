-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, google_id, avatar_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, password_hash, google_id, avatar_url, email_verified, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, name, email, password_hash, google_id, avatar_url, email_verified, created_at, updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, google_id, avatar_url, email_verified, created_at, updated_at
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByGoogleID :one
SELECT id, name, email, password_hash, google_id, avatar_url, email_verified, created_at, updated_at
FROM users
WHERE google_id = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, password_hash = $4, google_id = $5, avatar_url = $6, email_verified = $7, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, name, email, password_hash, google_id, avatar_url, email_verified, created_at, updated_at;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, name, email, password_hash, google_id, avatar_url, email_verified, created_at, updated_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
