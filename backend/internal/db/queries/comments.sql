-- name: ListComments :many
SELECT id, ticket_id, author_id, body, created_at
FROM comments
WHERE ticket_id = ?
ORDER BY created_at ASC;

-- name: CreateComment :execresult
INSERT INTO comments (ticket_id, author_id, body) VALUES (?, ?, ?);
