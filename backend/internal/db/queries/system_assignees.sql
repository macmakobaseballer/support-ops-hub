-- name: ListSystemAssignees :many
SELECT system_id, user_id
FROM system_assignees
WHERE system_id = ?;

-- name: CreateSystemAssignee :exec
INSERT INTO system_assignees (system_id, user_id) VALUES (?, ?);

-- name: DeleteSystemAssignees :exec
DELETE FROM system_assignees WHERE system_id = ?;
