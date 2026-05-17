-- name: GetAllStatusCounts :one
SELECT
  COUNT(CASE WHEN status = 'new'         THEN 1 END) AS new_count,
  COUNT(CASE WHEN status = 'in_progress' THEN 1 END) AS in_progress_count,
  COUNT(CASE WHEN status = 'waiting'     THEN 1 END) AS waiting_count,
  COUNT(CASE WHEN status = 'done'        THEN 1 END) AS done_count
FROM tickets;

-- name: GetMonthlyTrend :many
SELECT DATE_FORMAT(received_at, '%Y-%m') AS month, COUNT(*) AS count
FROM tickets
WHERE received_at >= DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL 5 MONTH)
GROUP BY DATE_FORMAT(received_at, '%Y-%m')
ORDER BY month ASC;

-- name: GetOpenTypeBreakdown :many
SELECT type, COUNT(*) AS count
FROM tickets
WHERE status != 'done'
GROUP BY type;

-- name: GetOpenPriorityBreakdown :many
SELECT priority, COUNT(*) AS count
FROM tickets
WHERE status != 'done'
GROUP BY priority;

-- name: GetTopSystemsByOpenCount :many
SELECT t.system_id, s.name AS system_name, COUNT(*) AS count
FROM tickets t
JOIN systems s ON s.id = t.system_id
WHERE t.status IN ('new', 'in_progress', 'waiting')
GROUP BY t.system_id, s.name
ORDER BY count DESC
LIMIT 6;

-- name: GetCustomerAllStatusCounts :one
SELECT
  COUNT(CASE WHEN status = 'new'         THEN 1 END) AS new_count,
  COUNT(CASE WHEN status = 'in_progress' THEN 1 END) AS in_progress_count,
  COUNT(CASE WHEN status = 'waiting'     THEN 1 END) AS waiting_count,
  COUNT(CASE WHEN status = 'done'        THEN 1 END) AS done_count
FROM tickets
WHERE customer_id = ?;

-- name: GetCustomerMonthlyTrend :many
SELECT DATE_FORMAT(received_at, '%Y-%m') AS month, COUNT(*) AS count
FROM tickets
WHERE customer_id = ?
  AND received_at >= DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL 5 MONTH)
GROUP BY DATE_FORMAT(received_at, '%Y-%m')
ORDER BY month ASC;

-- name: GetCustomerOpenTypeBreakdown :many
SELECT type, COUNT(*) AS count
FROM tickets
WHERE customer_id = ? AND status != 'done'
GROUP BY type;

-- name: GetCustomerSystemsBreakdown :many
SELECT t.system_id, s.name AS system_name,
       COUNT(CASE WHEN t.status = 'new'         THEN 1 END) AS new_count,
       COUNT(CASE WHEN t.status = 'in_progress' THEN 1 END) AS in_progress_count,
       COUNT(CASE WHEN t.status = 'waiting'     THEN 1 END) AS waiting_count,
       COUNT(CASE WHEN t.status = 'done'        THEN 1 END) AS done_count,
       COUNT(*) AS total
FROM tickets t
JOIN systems s ON s.id = t.system_id
WHERE t.customer_id = ?
GROUP BY t.system_id, s.name
ORDER BY total DESC;
