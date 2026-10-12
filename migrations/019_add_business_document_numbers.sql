CREATE SEQUENCE work_order_number_seq START 1;
CREATE SEQUENCE invoice_number_seq START 1;

ALTER TABLE work_orders ADD COLUMN work_order_number TEXT;
UPDATE work_orders SET work_order_number = 'WO-' || TO_CHAR(COALESCE(created_at, NOW()), 'YYYY') || '-' || LPAD(id::TEXT, 4, '0');
ALTER TABLE work_orders ALTER COLUMN work_order_number SET NOT NULL;
ALTER TABLE work_orders ADD CONSTRAINT work_orders_work_order_number_unique UNIQUE (work_order_number);
ALTER TABLE work_orders ALTER COLUMN work_order_number SET DEFAULT ('WO-' || TO_CHAR(CURRENT_DATE, 'YYYY') || '-' || LPAD(nextval('work_order_number_seq')::TEXT, 4, '0'));
SELECT setval('work_order_number_seq', COALESCE((SELECT MAX(id) FROM work_orders), 0) + 1, false);

ALTER TABLE invoices ADD COLUMN invoice_number TEXT;
UPDATE invoices SET invoice_number = 'INV-' || TO_CHAR(COALESCE(created_at, NOW()), 'YYYY') || '-' || LPAD(id::TEXT, 4, '0');
ALTER TABLE invoices ALTER COLUMN invoice_number SET NOT NULL;
ALTER TABLE invoices ADD CONSTRAINT invoices_invoice_number_unique UNIQUE (invoice_number);
ALTER TABLE invoices ALTER COLUMN invoice_number SET DEFAULT ('INV-' || TO_CHAR(CURRENT_DATE, 'YYYY') || '-' || LPAD(nextval('invoice_number_seq')::TEXT, 4, '0'));
SELECT setval('invoice_number_seq', COALESCE((SELECT MAX(id) FROM invoices), 0) + 1, false);
