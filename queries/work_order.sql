-- name: CreateWorkOrder :one
INSERT INTO work_orders (quotation_id, customer_id, vehicle_id, status, due_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, work_order_number, quotation_id, customer_id, vehicle_id, status, due_date, started_at, completed_at, created_at, deleted_at;

-- name: ListWorkOrders :many
SELECT id, work_order_number, quotation_id, customer_id, vehicle_id, status, due_date, started_at, completed_at, created_at, deleted_at
FROM work_orders WHERE deleted_at IS NULL ORDER BY id;

-- name: GetWorkOrder :one
SELECT id, work_order_number, quotation_id, customer_id, vehicle_id, status, due_date, started_at, completed_at, created_at, deleted_at
FROM work_orders WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateWorkOrderStatus :one
UPDATE work_orders SET status = sqlc.arg(status)::text,
    started_at = CASE WHEN sqlc.arg(status)::text = 'In Progress' AND started_at IS NULL THEN NOW() ELSE started_at END,
    completed_at = CASE WHEN sqlc.arg(status)::text = 'Completed' AND completed_at IS NULL THEN NOW() ELSE completed_at END
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id, work_order_number, quotation_id, customer_id, vehicle_id, status, due_date, started_at, completed_at, created_at, deleted_at;

-- name: UpdatePendingWorkOrderDueDate :one
UPDATE work_orders SET due_date = $1
WHERE id = $2 AND status = 'Pending' AND deleted_at IS NULL
RETURNING id, work_order_number, quotation_id, customer_id, vehicle_id, status, due_date, started_at, completed_at, created_at, deleted_at;

-- name: SoftDeleteWorkOrder :execrows
UPDATE work_orders SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetWorkOrdersByStatus :many
SELECT id, work_order_number, quotation_id, customer_id, vehicle_id, status, due_date, started_at, completed_at, created_at, deleted_at
FROM work_orders WHERE status = $1 AND deleted_at IS NULL ORDER BY id;

-- name: CountWorkOrdersByStatus :one
SELECT COUNT(*) FROM work_orders
WHERE status = $1 AND deleted_at IS NULL;
