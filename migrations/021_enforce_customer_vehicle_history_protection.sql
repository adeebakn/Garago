-- Enforce audit-history protection at the database boundary so maintenance
-- scripts, older API processes, and future code cannot bypass the rule.
CREATE OR REPLACE FUNCTION prevent_customer_history_deletion()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.deleted_at IS NULL
       AND NEW.deleted_at IS NOT NULL
       AND (
           EXISTS (SELECT 1 FROM vehicles WHERE customer_id = OLD.id)
           OR EXISTS (SELECT 1 FROM quotations WHERE customer_id = OLD.id)
           OR EXISTS (SELECT 1 FROM work_orders WHERE customer_id = OLD.id)
           OR EXISTS (SELECT 1 FROM invoices WHERE customer_id = OLD.id)
       ) THEN
        RAISE EXCEPTION 'customer cannot be deleted because linked history exists'
            USING ERRCODE = '23503';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS protect_customer_history ON customers;
CREATE TRIGGER protect_customer_history
BEFORE UPDATE OF deleted_at ON customers
FOR EACH ROW
EXECUTE FUNCTION prevent_customer_history_deletion();

CREATE OR REPLACE FUNCTION prevent_vehicle_history_deletion()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.deleted_at IS NULL
       AND NEW.deleted_at IS NOT NULL
       AND (
           EXISTS (SELECT 1 FROM quotations WHERE vehicle_id = OLD.id)
           OR EXISTS (SELECT 1 FROM work_orders WHERE vehicle_id = OLD.id)
       ) THEN
        RAISE EXCEPTION 'vehicle cannot be deleted because linked history exists'
            USING ERRCODE = '23503';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS protect_vehicle_history ON vehicles;
CREATE TRIGGER protect_vehicle_history
BEFORE UPDATE OF deleted_at ON vehicles
FOR EACH ROW
EXECUTE FUNCTION prevent_vehicle_history_deletion();
