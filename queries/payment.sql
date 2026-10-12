-- name: CreatePayment :one
INSERT INTO payments (invoice_id, amount, payment_method)
VALUES ($1, $2, $3)
RETURNING id, invoice_id, amount, payment_method, payment_date, deleted_at;

-- name: ListPayments :many
SELECT p.id, p.invoice_id, p.amount, p.payment_method, p.payment_date, p.deleted_at
FROM payments p JOIN invoices i ON i.id = p.invoice_id
WHERE p.deleted_at IS NULL AND i.deleted_at IS NULL
ORDER BY p.id;

-- name: GetPayment :one
SELECT id, invoice_id, amount, payment_method, payment_date, deleted_at
FROM payments WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPaymentsByInvoice :many
SELECT id, invoice_id, amount, payment_method, payment_date, deleted_at
FROM payments WHERE invoice_id = $1 AND deleted_at IS NULL
ORDER BY id;

-- name: SoftDeletePayment :execrows
UPDATE payments SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTotalRevenue :one
SELECT COALESCE(SUM(p.amount), 0)
FROM payments p JOIN invoices i ON i.id = p.invoice_id
WHERE p.deleted_at IS NULL AND i.deleted_at IS NULL;
