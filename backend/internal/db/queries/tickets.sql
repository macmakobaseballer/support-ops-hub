-- name: GetTicket :one
SELECT id, title, description, type, priority, status,
       customer_id, system_id, assignee_id, created_by,
       received_at, created_at, updated_at
FROM tickets
WHERE id = ? LIMIT 1;

-- name: ListTickets :many
SELECT id, title, description, type, priority, status,
       customer_id, system_id, assignee_id, created_by,
       received_at, created_at, updated_at
FROM tickets
ORDER BY received_at DESC;

-- name: ListTicketsByStatus :many
SELECT id, title, description, type, priority, status,
       customer_id, system_id, assignee_id, created_by,
       received_at, created_at, updated_at
FROM tickets
WHERE status = ?
ORDER BY received_at DESC;

-- name: CreateTicket :execresult
INSERT INTO tickets (title, description, type, priority, status,
                     customer_id, system_id, assignee_id, created_by, received_at)
VALUES (?, ?, ?, ?, 'new', ?, ?, ?, ?, ?);

-- name: UpdateTicket :exec
UPDATE tickets
SET title = ?, description = ?, type = ?, priority = ?,
    assignee_id = ?, received_at = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateTicketStatus :exec
UPDATE tickets SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: UpdateTicketAssignee :exec
UPDATE tickets SET assignee_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: CountTicketsByStatus :one
SELECT COUNT(*) FROM tickets WHERE status = ?;

-- name: CountOpenTicketsBySystem :one
SELECT COUNT(*) FROM tickets
WHERE system_id = ? AND status IN ('new', 'in_progress', 'waiting');

-- name: GetTicketDetail :one
SELECT t.id, t.title, t.description, t.type, t.priority, t.status,
       t.customer_id, c.name AS customer_name,
       t.system_id,   s.name AS system_name,
       t.assignee_id, a.name AS assignee_name,
       t.created_by,  cb.name AS created_by_name,
       t.received_at, t.created_at, t.updated_at
FROM tickets t
JOIN customers c  ON c.id = t.customer_id
JOIN systems   s  ON s.id = t.system_id
LEFT JOIN users a ON a.id = t.assignee_id
JOIN users     cb ON cb.id = t.created_by
WHERE t.id = ? LIMIT 1;
