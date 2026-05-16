-- name: GetCustomer :one
SELECT id, name, notes, is_active, created_at, updated_at
FROM customers
WHERE id = ? LIMIT 1;

-- name: ListCustomersActive :many
SELECT id, name, notes, is_active, created_at, updated_at
FROM customers
WHERE is_active = TRUE
ORDER BY name;

-- name: ListCustomersAll :many
SELECT id, name, notes, is_active, created_at, updated_at
FROM customers
ORDER BY name;

-- name: CreateCustomer :execresult
INSERT INTO customers (name, notes, is_active)
VALUES (?, ?, ?);

-- name: UpdateCustomer :exec
UPDATE customers
SET name = ?, notes = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
