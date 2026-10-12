package app

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	db "garago/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func createWorkOrder(c *gin.Context, queries *db.Queries) {
	var request CreateWorkOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work order data",
		})
		return
	}
	if request.QuotationID <= 0 || request.CustomerID <= 0 || request.VehicleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quotation, customer, and vehicle are required"})
		return
	}

	quotation, err := queries.GetQuotation(c, int32(request.QuotationID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quotation not found"})
		return
	}
	if quotation.Status != QuotationApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only approved quotations can create work orders"})
		return
	}
	if quotation.CustomerID != int32(request.CustomerID) || quotation.VehicleID != int32(request.VehicleID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Work order customer and vehicle must match the quotation"})
		return
	}

	dueDate, err := time.Parse("2006-01-02", request.DueDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid due date. Use YYYY-MM-DD",
		})
		return
	}

	workOrder, err := queries.CreateWorkOrder(
		c,
		db.CreateWorkOrderParams{
			QuotationID: int32(request.QuotationID),
			CustomerID:  int32(request.CustomerID),
			VehicleID:   int32(request.VehicleID),
			Status:      WorkOrderPending,
			DueDate: pgtype.Date{
				Time:  dueDate,
				Valid: true,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create work order",
		})
		return
	}

	c.JSON(http.StatusCreated, workOrder)
}

func getWorkOrders(c *gin.Context, queries *db.Queries) {
	workOrders, err := queries.ListWorkOrders(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get work orders",
		})
		return
	}

	c.JSON(http.StatusOK, workOrders)
}

func getWorkOrderByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work order ID",
		})
		return
	}

	workOrder, err := queries.GetWorkOrder(c, int32(id))

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Work order not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get work order",
		})
		return
	}

	c.JSON(http.StatusOK, workOrder)
}

func updateWorkOrderStatus(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work order ID",
		})
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status",
		})
		return
	}

	// Validate status
	if request.Status != WorkOrderPending &&
		request.Status != WorkOrderInProgress &&
		request.Status != WorkOrderCompleted &&
		request.Status != WorkOrderCancelled {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work order status",
		})
		return
	}
	current, err := queries.GetWorkOrder(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}
	if current.Status == WorkOrderCompleted && request.Status != WorkOrderCompleted {
		c.JSON(http.StatusConflict, gin.H{"error": "A completed work order cannot be moved backward or cancelled"})
		return
	}
	if current.Status == WorkOrderCancelled && request.Status != WorkOrderCancelled {
		c.JSON(http.StatusConflict, gin.H{"error": "A cancelled work order cannot be reopened"})
		return
	}
	if current.Status == WorkOrderInProgress && request.Status == WorkOrderPending {
		c.JSON(http.StatusConflict, gin.H{"error": "An in-progress work order cannot return to pending"})
		return
	}
	if request.Status == WorkOrderCancelled {
		tx, err := databasePool.Begin(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start cancellation"})
			return
		}
		defer tx.Rollback(c)
		txQueries := queries.WithTx(tx)
		if err := archiveUnpaidInvoiceForWorkOrder(c, tx, txQueries, int32(id)); err != nil {
			if err.Error() == "payment history exists" {
				c.JSON(http.StatusConflict, gin.H{"error": "This work order cannot be cancelled because payment has already been recorded."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not prepare work order cancellation"})
			return
		}
		workOrder, err := txQueries.UpdateWorkOrderStatus(c, db.UpdateWorkOrderStatusParams{Status: request.Status, ID: int32(id)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cancel work order"})
			return
		}
		if err := tx.Commit(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not complete work order cancellation"})
			return
		}
		c.JSON(http.StatusOK, workOrder)
		return
	}

	// Update database
	workOrder, err := queries.UpdateWorkOrderStatus(
		c,
		db.UpdateWorkOrderStatusParams{
			Status: request.Status,
			ID:     int32(id),
		},
	)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Work order not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update work order status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workOrder)
}

// updatePendingWorkOrder allows only the schedule to change while a work order is pending.
// Customer, vehicle, and quotation remain fixed because they define the approved scope.
func updatePendingWorkOrder(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order ID"})
		return
	}
	var request struct {
		DueDate string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order data"})
		return
	}
	dueDate, err := time.Parse("2006-01-02", request.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date. Use YYYY-MM-DD"})
		return
	}
	workOrder, err := queries.UpdatePendingWorkOrderDueDate(c, db.UpdatePendingWorkOrderDueDateParams{ID: int32(id), DueDate: pgtype.Date{Time: dueDate, Valid: true}})
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusConflict, gin.H{"error": "Only pending work orders can be edited"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update work order"})
		return
	}
	c.JSON(http.StatusOK, workOrder)
}

func deleteWorkOrderByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work order ID",
		})
		return
	}
	tx, err := databasePool.Begin(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start work order archive"})
		return
	}
	defer tx.Rollback(c)
	txQueries := queries.WithTx(tx)
	if err := archiveUnpaidInvoiceForWorkOrder(c, tx, txQueries, int32(id)); err != nil {
		if err.Error() == "payment history exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "This work order cannot be archived because payment has already been recorded."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check work order financial history"})
		return
	}

	rows, err := txQueries.SoftDeleteWorkOrder(c, int32(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete work order",
		})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}
	if err := tx.Commit(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not archive work order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Work order archived successfully",
	})
}

// archiveUnpaidInvoiceForWorkOrder keeps the work-order/invoice relationship consistent.
// It locks the invoice row, blocks operations after any payment, and archives a zero-payment invoice.
func archiveUnpaidInvoiceForWorkOrder(c *gin.Context, tx pgx.Tx, queries *db.Queries, workOrderID int32) error {
	var invoiceID int32
	err := tx.QueryRow(c, `SELECT id FROM invoices WHERE work_order_id = $1 AND deleted_at IS NULL FOR UPDATE`, workOrderID).Scan(&invoiceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	var hasPayments bool
	if err := tx.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM payments WHERE invoice_id = $1 AND deleted_at IS NULL)`, invoiceID).Scan(&hasPayments); err != nil {
		return err
	}
	if hasPayments {
		return errors.New("payment history exists")
	}
	_, err = queries.SoftDeleteInvoice(c, invoiceID)
	return err
}

func getWorkOrdersByStatus(c *gin.Context, queries *db.Queries) {
	status := c.Param("status")

	if status != WorkOrderPending &&
		status != WorkOrderInProgress &&
		status != WorkOrderCompleted &&
		status != WorkOrderCancelled {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work order status",
		})
		return
	}

	workOrders, err := queries.GetWorkOrdersByStatus(c, status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get work orders",
		})
		return
	}

	if workOrders == nil {
		workOrders = []db.GetWorkOrdersByStatusRow{}
	}

	c.JSON(http.StatusOK, workOrders)
}
