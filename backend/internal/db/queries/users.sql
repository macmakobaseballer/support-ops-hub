-- name: GetUser :one
SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsersActive :many
SELECT id, name, email, role, is_active, created_at, updated_at
FROM users
WHERE is_active = TRUE
ORDER BY name;

-- name: ListUsersAll :many
SELECT id, name, email, role, is_active, created_at, updated_at
FROM users
ORDER BY name;
