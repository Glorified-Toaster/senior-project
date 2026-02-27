-- name: CreateUser :one
INSERT INTO users (
    username,
    full_name,
    password_hash,
    role,
    is_active
) VALUES (
    sqlc.arg(username),
    sqlc.arg(full_name),
    sqlc.arg(password_hash),
    sqlc.arg(role),
    sqlc.arg(is_active)
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: ListActiveUsers :many
SELECT * FROM users WHERE is_active = TRUE ORDER BY created_at DESC;

-- name: ListUsersByRole :many
SELECT * FROM users WHERE role = $1 ORDER BY created_at DESC;

-- name: UpdateUserRole :exec
UPDATE users SET role = $2 WHERE id = $1;
