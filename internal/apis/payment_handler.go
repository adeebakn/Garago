package apis

import (
	"net/http"
	"strconv"

	"garago/internal/database/postgres"

	"github.com/gin-gonic/gin"
)

func validPaymentMethod(method string) bool {
	return method == PaymentCash || method == PaymentUPI || method == PaymentCard || method == PaymentBankTransfer
}

func paymentSummary(totalAmount float64, payments []postgres.Payment) (float64, float64, string) {
	var paid float64
	for _, payment := range payments {
		paid += payment.Amount
	}
	remaining := totalAmount - paid
	status := InvoiceUnpaid
	if paid > 0 && remaining > 0 {
		status = InvoicePartial
	} else if remaining == 0 {
		status = InvoicePaid
	}
	return paid, remaining, status
}

// createPayment records an installment and updates its invoice in one transaction.
func createPayment(c *gin.Context, queries *postgres.Queries) {
	var request CreatePaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.InvoiceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid invoice is required"})
		return
	}
	if request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount must be greater than 0"})
		return
	}
	if !validPaymentMethod(request.PaymentMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	ctx := c.Request.Context()
	tx, err := databasePool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start payment transaction"})
		return
	}
	defer tx.Rollback(ctx)
	txQueries := queries.WithTx(tx)

	// FOR UPDATE locks this invoice while its remaining balance is checked.
	invoice, err := txQueries.GetInvoiceForUpdate(c, int32(request.InvoiceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	payments, err := txQueries.GetPaymentsByInvoice(c, invoice.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment history"})
		return
	}
	paid, remaining, _ := paymentSummary(invoice.TotalAmount, payments)
	if request.Amount > remaining {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount exceeds remaining balance", "remaining_balance": remaining})
		return
	}

	payment, err := txQueries.CreatePayment(c, postgres.CreatePaymentParams{InvoiceID: invoice.ID, Amount: request.Amount, PaymentMethod: request.PaymentMethod})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}
	newRemaining := remaining - request.Amount
	newStatus := InvoicePartial
	if newRemaining == 0 {
		newStatus = InvoicePaid
	}
	updatedInvoice, err := txQueries.UpdateInvoiceStatus(c, postgres.UpdateInvoiceStatusParams{ID: invoice.ID, Status: newStatus})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice status"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"payment": payment, "invoice": updatedInvoice, "invoice_total": invoice.TotalAmount, "total_paid": paid + request.Amount, "remaining_balance": newRemaining})
}

func getPaymentsByInvoice(c *gin.Context, queries *postgres.Queries) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	payments, err := queries.GetPaymentsByInvoice(c, int32(invoiceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment history"})
		return
	}
	if payments == nil {
		payments = []postgres.Payment{}
	}
	c.JSON(http.StatusOK, payments)
}

func getPayments(c *gin.Context, queries *postgres.Queries) {
	payments, err := queries.ListPayments(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
		return
	}
	if payments == nil {
		payments = []postgres.Payment{}
	}
	c.JSON(http.StatusOK, payments)
}

func getPaymentByID(c *gin.Context, queries *postgres.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}
	payment, err := queries.GetPayment(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func deletePaymentByID(c *gin.Context, queries *postgres.Queries) {
	c.JSON(http.StatusConflict, gin.H{"error": "Recorded payments cannot be deleted. Use a future reversal flow to preserve audit history."})
	return
	/*
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
			return
		}
		ctx := c.Request.Context()
		tx, err := databasePool.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start payment transaction"})
			return
		}
		defer tx.Rollback(ctx)
		txQueries := queries.WithTx(tx)
		payment, err := txQueries.GetPayment(c, int32(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		invoice, err := txQueries.GetInvoiceForUpdate(c, payment.InvoiceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		if rows, err := txQueries.SoftDeletePayment(c, payment.ID); err != nil || rows == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
			return
		}
		payments, err := txQueries.GetPaymentsByInvoice(c, payment.InvoiceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment history"})
			return
		}
		paid, remaining, status := paymentSummary(invoice.TotalAmount, payments)
		if _, err := txQueries.UpdateInvoiceStatus(c, postgres.UpdateInvoiceStatusParams{ID: invoice.ID, Status: status}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice status"})
			return
		}
		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully", "total_paid": paid, "remaining_balance": remaining, "invoice_status": status})
	*/
}
