-- name: CreateInvoice :one
INSERT INTO invoices (work_order_id, customer_id, total_amount, status)
VALUES ($1, $2, $3, $4)
RETURNING id, invoice_number, work_order_id, customer_id, total_amount, status, created_at, deleted_at;

-- name: ListInvoices :many
SELECT id, invoice_number, work_order_id, customer_id, total_amount, status, created_at, deleted_at
FROM invoices
WHERE deleted_at IS NULL
ORDER BY id;

-- name: GetInvoice :one
SELECT id, invoice_number, work_order_id, customer_id, total_amount, status, created_at, deleted_at
FROM invoices
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetInvoiceByWorkOrder :one
SELECT id, work_order_id, customer_id, total_amount, status, created_at, deleted_at
FROM invoices
WHERE work_order_id = $1 AND deleted_at IS NULL;

-- name: GetInvoiceForUpdate :one
SELECT id, work_order_id, customer_id, total_amount, status, created_at, deleted_at
FROM invoices
WHERE id = $1 AND deleted_at IS NULL
FOR UPDATE;

-- name: UpdateInvoiceStatus :one
UPDATE invoices SET status = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING id, invoice_number, work_order_id, customer_id, total_amount, status, created_at, deleted_at;

-- name: SoftDeleteInvoice :execrows
UPDATE invoices SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetWorkOrderForInvoice :one
SELECT id, quotation_id, customer_id, vehicle_id, status, due_date, created_at, deleted_at
FROM work_orders
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetInvoicesByStatus :many
SELECT id, invoice_number, work_order_id, customer_id, total_amount, status, created_at, deleted_at
FROM invoices
WHERE status = $1 AND deleted_at IS NULL
ORDER BY id;

-- name: CountInvoicesByStatus :one
SELECT COUNT(*) FROM invoices
WHERE status = $1 AND deleted_at IS NULL;
