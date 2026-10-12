package apis

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var (
	errCustomerHasHistory = errors.New("customer has linked vehicle or business history")
	errVehicleHasHistory  = errors.New("vehicle has linked business history")
)

type customerHistory struct {
	vehicles   bool
	quotations bool
	workOrders bool
	invoices   bool
}

func (history customerHistory) blocksDeletion() bool {
	return history.vehicles || history.quotations || history.workOrders || history.invoices
}

// archiveCustomer locks the customer before checking every downstream relation.
// Archived rows are intentionally included: audit history must remain attributable.
func archiveCustomer(ctx context.Context, id int32) error {
	tx, err := databasePool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin customer archive: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := lockActiveRecord(ctx, tx, "customers", id); err != nil {
		return err
	}

	var history customerHistory
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM vehicles WHERE customer_id = $1),
		       EXISTS (SELECT 1 FROM quotations WHERE customer_id = $1),
		       EXISTS (SELECT 1 FROM work_orders WHERE customer_id = $1),
		       EXISTS (SELECT 1 FROM invoices WHERE customer_id = $1)`, id).Scan(
		&history.vehicles, &history.quotations, &history.workOrders, &history.invoices,
	)
	if err != nil {
		return fmt.Errorf("check customer history: %w", err)
	}
	if history.blocksDeletion() {
		return errCustomerHasHistory
	}

	if _, err := tx.Exec(ctx, `UPDATE customers SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`, id); err != nil {
		return fmt.Errorf("archive customer: %w", err)
	}
	return tx.Commit(ctx)
}

func archiveVehicle(ctx context.Context, id int32) error {
	tx, err := databasePool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin vehicle archive: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := lockActiveRecord(ctx, tx, "vehicles", id); err != nil {
		return err
	}

	var blocked bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM quotations WHERE vehicle_id = $1)
		    OR EXISTS (SELECT 1 FROM work_orders WHERE vehicle_id = $1)`, id).Scan(&blocked)
	if err != nil {
		return fmt.Errorf("check vehicle history: %w", err)
	}
	if blocked {
		return errVehicleHasHistory
	}

	if _, err := tx.Exec(ctx, `UPDATE vehicles SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`, id); err != nil {
		return fmt.Errorf("archive vehicle: %w", err)
	}
	return tx.Commit(ctx)
}

func lockActiveRecord(ctx context.Context, tx pgx.Tx, table string, id int32) error {
	// table is selected only by internal callers; it is never derived from request input.
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1 AND deleted_at IS NULL FOR UPDATE", table)
	var lockedID int32
	return tx.QueryRow(ctx, query, id).Scan(&lockedID)
}
