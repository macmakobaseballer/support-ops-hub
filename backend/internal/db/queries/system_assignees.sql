-- name: ListSystemAssignees :many
SELECT system_id, user_id
FROM system_assignees
WHERE system_id = ?;

-- name: CreateSystemAssignee :exec
INSERT INTO system_assignees (system_id, user_id) VALUES (?, ?);

-- name: DeleteSystemAssignees :exec
DELETE FROM system_assignees WHERE system_id = ?;

-- name: ListSystemAssigneesWithUsers :many
SELECT u.id, u.name
FROM system_assignees sa
JOIN users u ON u.id = sa.user_id
WHERE sa.system_id = ?
ORDER BY u.name;
