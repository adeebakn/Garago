-- One work order produces one bill. Installments are stored in payments.
ALTER TABLE invoices
ADD CONSTRAINT uq_invoices_work_order UNIQUE (work_order_id);
