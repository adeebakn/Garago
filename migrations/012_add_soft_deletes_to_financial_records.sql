-- Retain financial and operational history when a manager removes a record.
ALTER TABLE quotations ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE work_orders ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- A soft-deleted invoice no longer counts as the active invoice for its work order.
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS uq_invoices_work_order;
CREATE UNIQUE INDEX IF NOT EXISTS uq_invoices_work_order_active
ON invoices (work_order_id)
WHERE deleted_at IS NULL;
