package apis

import (
	"net/http"

	"garago/internal/database/postgres"

	"github.com/gin-gonic/gin"
)

func getDashboard(c *gin.Context, queries *postgres.Queries) {

	pending, err := queries.CountWorkOrdersByStatus(c, WorkOrderPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending work orders"})
		return
	}

	inProgress, err := queries.CountWorkOrdersByStatus(c, WorkOrderInProgress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get in-progress work orders"})
		return
	}

	completed, err := queries.CountWorkOrdersByStatus(c, WorkOrderCompleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get completed work orders"})
		return
	}

	unpaid, err := queries.CountInvoicesByStatus(c, InvoiceUnpaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unpaid invoices"})
		return
	}

	partial, err := queries.CountInvoicesByStatus(c, InvoicePartial)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get partial invoices"})
		return
	}

	paid, err := queries.CountInvoicesByStatus(c, InvoicePaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get paid invoices"})
		return
	}

	totalRevenue, err := queries.GetTotalRevenue(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total revenue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_work_orders":     pending,
		"in_progress_work_orders": inProgress,
		"completed_work_orders":   completed,
		"unpaid_invoices":         unpaid,
		"partial_invoices":        partial,
		"paid_invoices":           paid,
		"total_revenue":           totalRevenue,
	})
}
