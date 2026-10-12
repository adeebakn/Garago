-- These indexes support workflow-integrity checks and common relationship lookups.
-- PostgreSQL does not create indexes automatically for foreign-key columns.
CREATE INDEX IF NOT EXISTS idx_vehicles_customer_id ON vehicles (customer_id);
CREATE INDEX IF NOT EXISTS idx_quotations_customer_id ON quotations (customer_id);
CREATE INDEX IF NOT EXISTS idx_quotations_vehicle_id ON quotations (vehicle_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_customer_id ON work_orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_vehicle_id ON work_orders (vehicle_id);
CREATE INDEX IF NOT EXISTS idx_invoices_customer_id ON invoices (customer_id);
CREATE INDEX IF NOT EXISTS idx_invoices_work_order_id ON invoices (work_order_id);
CREATE INDEX IF NOT EXISTS idx_payments_invoice_id ON payments (invoice_id);
