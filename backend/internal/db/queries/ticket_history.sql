-- name: ListTicketHistory :many
SELECT id, ticket_id, changed_by, field_name, old_value, new_value, changed_at
FROM ticket_history
WHERE ticket_id = ?
ORDER BY changed_at ASC;

-- name: CreateTicketHistory :exec
INSERT INTO ticket_history (ticket_id, changed_by, field_name, old_value, new_value)
VALUES (?, ?, ?, ?, ?);
