ALTER TABLE quotations
    ADD COLUMN advance_requested DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE quotations
    ADD CONSTRAINT quotations_advance_requested_valid
    CHECK (advance_requested >= 0 AND advance_requested <= total_amount);
