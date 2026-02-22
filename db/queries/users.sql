-- name: CreateUser :one
INSERT INTO users (
    username,
    full_name,
    password_hash,
    is_active
) VALUES (
    sqlc.arg(username),
    sqlc.arg(full_name),
    sqlc.arg(password_hash),
    sqlc.arg(is_active)
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: ListActiveUsers :many
SELECT * FROM users WHERE is_active = TRUE ORDER BY created_at DESC;


-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role)
VALUES ($1, $2);

-- name: GetUserRoles :many
SELECT role FROM user_roles WHERE user_id = $1;
