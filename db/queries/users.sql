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

-- name: SoftDeleteUser :exec
UPDATE users 
SET deleted_at = NOW(), is_active = false
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: RestoreUser :exec
UPDATE users 
SET deleted_at = NULL, is_active = true
WHERE id = $1;

-- name: ListDeletedUsers :many
SELECT * FROM users 
WHERE deleted_at IS NOT NULL 
ORDER BY deleted_at DESC, id ASC
LIMIT $1 OFFSET $2;

-- name: DisableUser :exec
UPDATE users 
SET is_active = false
WHERE id = $1 AND deleted_at IS NULL;

-- name: EnableUser :exec
UPDATE users 
SET is_active = true
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT * FROM users 
WHERE username = $1 AND deleted_at IS NULL;

-- name: ListActiveUsers :many
SELECT * FROM users 
WHERE is_active = TRUE AND deleted_at IS NULL 
ORDER BY created_at DESC;

-- name: ListUsersByRole :many
SELECT * FROM users 
WHERE role = $1 AND deleted_at IS NULL 
ORDER BY created_at DESC;

-- name: ListAllUsers :many
SELECT * FROM users 
WHERE deleted_at IS NULL 
ORDER BY last_login DESC NULLS LAST, id ASC
LIMIT $1 OFFSET $2;

-- name: UpdateUserRole :exec
UPDATE users 
SET role = $2 
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserLastLogin :exec
UPDATE users 
SET last_login = NOW()
WHERE id = $1 AND deleted_at IS NULL;


-- name: SearchUsers :many
SELECT * FROM users
WHERE 
    deleted_at IS NULL 
    AND (
        username ILIKE '%' || sqlc.arg(search)::text || '%'
        OR full_name ILIKE '%' || sqlc.arg(search)::text || '%'
    )
ORDER BY last_login DESC NULLS LAST, id ASC
LIMIT $1 OFFSET $2;

-- name: SearchDeletedUsers :many
SELECT * FROM users
WHERE 
    deleted_at IS NOT NULL 
    AND (
        username ILIKE '%' || sqlc.arg(search)::text || '%'
        OR full_name ILIKE '%' || sqlc.arg(search)::text || '%'
    )
ORDER BY deleted_at DESC, id ASC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT count(*) FROM users 
WHERE deleted_at IS NULL;

-- name: CountDeletedUsers :one
SELECT count(*) FROM users 
WHERE deleted_at IS NOT NULL;
-- name: CountSearchUsers :one
SELECT count(*) FROM users
WHERE 
    deleted_at IS NULL 
    AND (
        username ILIKE '%' || sqlc.arg(search)::text || '%'
        OR full_name ILIKE '%' || sqlc.arg(search)::text || '%'
    );

-- name: CountSearchDeletedUsers :one
SELECT count(*) FROM users
WHERE 
    deleted_at IS NOT NULL 
    AND (
        username ILIKE '%' || sqlc.arg(search)::text || '%'
        OR full_name ILIKE '%' || sqlc.arg(search)::text || '%'
    );

-- name: ListAllInstructors :many
SELECT * FROM users 
WHERE role = 'INSTRUCTOR' AND deleted_at IS NULL
ORDER BY last_login DESC NULLS LAST, id ASC
LIMIT $1 OFFSET $2;

-- name: CountInstructors :one
SELECT count(*) FROM users 
WHERE role = 'INSTRUCTOR' AND deleted_at IS NULL;