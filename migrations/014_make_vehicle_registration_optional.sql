ALTER TABLE vehicles
    ALTER COLUMN reg_no SET DEFAULT '';

CREATE UNIQUE INDEX vehicles_customer_registration_unique
    ON vehicles (customer_id, LOWER(reg_no))
    WHERE reg_no <> '' AND deleted_at IS NULL;
