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
