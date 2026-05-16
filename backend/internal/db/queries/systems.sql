-- name: GetSystem :one
SELECT id, name, customer_id, description, is_active, created_at, updated_at
FROM systems
WHERE id = ? LIMIT 1;

-- name: ListSystemsByCustomer :many
SELECT id, name, customer_id, description, is_active, created_at, updated_at
FROM systems
WHERE customer_id = ?
ORDER BY name;

-- name: ListSystemsAll :many
SELECT id, name, customer_id, description, is_active, created_at, updated_at
FROM systems
ORDER BY name;

-- name: ListSystemsActiveByCustomer :many
SELECT id, name, customer_id, description, is_active, created_at, updated_at
FROM systems
WHERE customer_id = ? AND is_active = TRUE
ORDER BY name;

-- name: ListSystemsActive :many
SELECT id, name, customer_id, description, is_active, created_at, updated_at
FROM systems
WHERE is_active = TRUE
ORDER BY name;

-- name: CreateSystem :execresult
INSERT INTO systems (name, customer_id, description, is_active)
VALUES (?, ?, ?, ?);

-- name: UpdateSystem :exec
UPDATE systems
SET name = ?, description = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
