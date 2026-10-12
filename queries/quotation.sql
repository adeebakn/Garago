-- name: CreateQuotation :one
INSERT INTO quotations (customer_id, vehicle_id, status, total_amount, advance_requested)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, customer_id, vehicle_id, status, total_amount, advance_requested, created_at, deleted_at;

-- name: ListQuotations :many
SELECT id, customer_id, vehicle_id, status, total_amount, advance_requested, created_at, deleted_at
FROM quotations WHERE deleted_at IS NULL ORDER BY id;

-- name: GetQuotation :one
SELECT id, customer_id, vehicle_id, status, total_amount, advance_requested, created_at, deleted_at
FROM quotations WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateQuotationStatus :one
UPDATE quotations SET status = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING id, customer_id, vehicle_id, status, total_amount, advance_requested, created_at, deleted_at;

-- name: UpdateQuotation :one
UPDATE quotations
SET customer_id = $1, vehicle_id = $2, total_amount = $3, advance_requested = $4
WHERE id = $5 AND status = 'Draft' AND deleted_at IS NULL
RETURNING id, customer_id, vehicle_id, status, total_amount, advance_requested, created_at, deleted_at;

-- name: DeleteQuotationItems :exec
DELETE FROM quotation_items WHERE quotation_id = $1;

-- name: AddQuotationItem :one
INSERT INTO quotation_items (quotation_id, service_id, quantity, price, amount)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, quotation_id, service_id, quantity, price, amount;

-- name: GetQuotationItems :many
SELECT id, quotation_id, service_id, quantity, price, amount
FROM quotation_items WHERE quotation_id = $1 ORDER BY id;

-- name: SoftDeleteQuotation :execrows
UPDATE quotations SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetQuotationsByStatus :many
SELECT id, customer_id, vehicle_id, status, total_amount, advance_requested, created_at, deleted_at
FROM quotations WHERE status = $1 AND deleted_at IS NULL ORDER BY id;
