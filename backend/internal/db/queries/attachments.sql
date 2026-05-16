-- name: ListAttachments :many
SELECT id, ticket_id, uploaded_by, file_name, file_key, file_size, content_type, created_at
FROM attachments
WHERE ticket_id = ?;

-- name: GetAttachment :one
SELECT id, ticket_id, uploaded_by, file_name, file_key, file_size, content_type, created_at
FROM attachments
WHERE id = ? AND ticket_id = ? LIMIT 1;

-- name: CreateAttachment :execresult
INSERT INTO attachments (ticket_id, uploaded_by, file_name, file_key, file_size, content_type)
VALUES (?, ?, ?, ?, ?, ?);

-- name: DeleteAttachment :exec
DELETE FROM attachments WHERE id = ? AND ticket_id = ?;
