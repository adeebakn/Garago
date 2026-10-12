package app

import (
	"errors"
	"net/http"
	"strconv"

	db "garago/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

/*And we're using SQLC functions that already exist:

queries.GetWorkOrder(...)
queries.GetQuotation(...)
queries.CreateInvoice(...)*/

func createInvoice(c *gin.Context, queries *db.Queries) {
	var request CreateInvoiceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid invoice data",
		})
		return
	}
	if request.WorkOrderID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A work order is required"})
		return
	}

	// Get work order
	workOrder, err := queries.GetWorkOrder(
		c,
		int32(request.WorkOrderID),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Work order not found",
		})
		return
	}
	if workOrder.Status == WorkOrderCancelled {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cancelled work orders cannot be invoiced",
		})
		return
	}

	// Get quotation
	quotation, err := queries.GetQuotation(
		c,
		workOrder.QuotationID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Quotation not found",
		})
		return
	}

	// Create invoice
	invoice, err := queries.CreateInvoice(
		c,
		db.CreateInvoiceParams{
			WorkOrderID: int32(request.WorkOrderID),
			CustomerID:  quotation.CustomerID,
			TotalAmount: quotation.TotalAmount,
			Status:      InvoiceUnpaid,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "An invoice already exists for this work order"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create invoice",
		})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func getInvoices(c *gin.Context, queries *db.Queries) {
	invoices, err := queries.ListInvoices(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get invoices",
		})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

func getInvoiceByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid invoice ID",
		})
		return
	}

	invoice, err := queries.GetInvoice(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invoice not found",
		})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// deleteInvoiceByID archives an invoice instead of removing its financial history.
func deleteInvoiceByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	var hasPayments bool
	if err := databasePool.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM payments WHERE invoice_id = $1 AND deleted_at IS NULL)`, int32(id)).Scan(&hasPayments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check payment history"})
		return
	}
	if hasPayments {
		c.JSON(http.StatusConflict, gin.H{"error": "An invoice with payments cannot be archived. Keep it for financial history."})
		return
	}

	rows, err := queries.SoftDeleteInvoice(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive invoice"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Invoice archived successfully"})
}

func updateInvoiceStatus(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid invoice ID",
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

	// Validate invoice status
	if request.Status != InvoiceUnpaid &&
		request.Status != InvoicePartial &&
		request.Status != InvoicePaid {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid invoice status",
		})
		return
	}

	invoice, err := queries.UpdateInvoiceStatus(
		c,
		db.UpdateInvoiceStatusParams{
			Status: request.Status,
			ID:     int32(id),
		},
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invoice not found",
		})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func getInvoicesByStatus(c *gin.Context, queries *db.Queries) {
	status := c.Param("status")

	if status != InvoiceUnpaid &&
		status != InvoicePartial &&
		status != InvoicePaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice status"})
		return
	}

	invoices, err := queries.GetInvoicesByStatus(c, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get invoices"})
		return
	}

	if invoices == nil {
		invoices = []db.GetInvoicesByStatusRow{}
	}

	c.JSON(http.StatusOK, invoices)
}
