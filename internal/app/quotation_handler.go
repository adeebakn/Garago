package app

import (
	"net/http"
	"strconv"

	db "garago/db/sqlc"

	"github.com/gin-gonic/gin"
)

func createQuotation(c *gin.Context, queries *db.Queries) {
	var request CreateQuotationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quotation data",
		})
		return
	}

	if len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quotation must contain at least one item",
		})
		return
	}
	if request.CustomerID <= 0 || request.VehicleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer and vehicle are required"})
		return
	}

	vehicle, err := queries.GetVehicle(c, int32(request.VehicleID))
	if err != nil || vehicle.CustomerID != int32(request.CustomerID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle does not belong to the selected customer"})
		return
	}

	var total float64

	// First calculate total using service prices
	for _, item := range request.Items {
		if item.ServiceID <= 0 || item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Every quotation item needs a service and a positive quantity"})
			return
		}

		service, err := queries.GetService(c, int32(item.ServiceID))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Service not found",
			})
			return
		}

		amount := service.Price * float64(item.Quantity)
		total += amount
	}
	if request.AdvanceRequested < 0 || request.AdvanceRequested > total {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Advance requested must be between ₹0 and the quotation total"})
		return
	}

	// A transaction makes the quotation header and every item one unit of work.
	ctx := c.Request.Context()
	tx, err := databasePool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start quotation transaction"})
		return
	}
	defer tx.Rollback(ctx)
	txQueries := queries.WithTx(tx)

	quotation, err := txQueries.CreateQuotation(
		c,
		db.CreateQuotationParams{
			CustomerID:       int32(request.CustomerID),
			VehicleID:        int32(request.VehicleID),
			Status:           QuotationDraft,
			TotalAmount:      total,
			AdvanceRequested: request.AdvanceRequested,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create quotation",
		})
		return
	}

	// Add quotation items
	for _, item := range request.Items {

		service, err := txQueries.GetService(c, int32(item.ServiceID))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Service not found",
			})
			return
		}

		amount := service.Price * float64(item.Quantity)

		_, err = txQueries.AddQuotationItem(
			c,
			db.AddQuotationItemParams{
				QuotationID: int32(quotation.ID),
				ServiceID:   int32(item.ServiceID),
				Quantity:    int32(item.Quantity),
				Price:       service.Price,
				Amount:      amount,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add quotation item",
			})
			return
		}
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quotation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"quotation": quotation,
		"message":   "Quotation created successfully",
	})
}

func getQuotationByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quotation ID",
		})
		return
	}

	quotation, err := queries.GetQuotation(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Quotation not found",
		})
		return
	}

	items, err := queries.GetQuotationItems(c, int32(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get quotation items",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quotation": quotation,
		"items":     items,
	})
}

// updateQuotationByID replaces the line items only while the quotation is still a draft.
// Once approved, it becomes the agreed work scope and is intentionally immutable.
func updateQuotationByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quotation ID"})
		return
	}
	var request CreateQuotationRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.CustomerID <= 0 || request.VehicleID <= 0 || len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer, vehicle, and at least one service are required"})
		return
	}
	quotation, err := queries.GetQuotation(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quotation not found"})
		return
	}
	if quotation.Status != QuotationDraft {
		c.JSON(http.StatusConflict, gin.H{"error": "Only draft quotations can be edited"})
		return
	}
	vehicle, err := queries.GetVehicle(c, int32(request.VehicleID))
	if err != nil || vehicle.CustomerID != int32(request.CustomerID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle does not belong to the selected customer"})
		return
	}
	total := 0.0
	services := make([]db.GetServiceRow, len(request.Items))
	for index, item := range request.Items {
		if item.ServiceID <= 0 || item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Every quotation item needs a service and a positive quantity"})
			return
		}
		service, err := queries.GetService(c, int32(item.ServiceID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		services[index] = service
		total += service.Price * float64(item.Quantity)
	}
	if request.AdvanceRequested < 0 || request.AdvanceRequested > total {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Advance requested must be between ₹0 and the quotation total"})
		return
	}
	ctx := c.Request.Context()
	tx, err := databasePool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start quotation update"})
		return
	}
	defer tx.Rollback(ctx)
	txQueries := queries.WithTx(tx)
	updated, err := txQueries.UpdateQuotation(c, db.UpdateQuotationParams{CustomerID: int32(request.CustomerID), VehicleID: int32(request.VehicleID), TotalAmount: total, AdvanceRequested: request.AdvanceRequested, ID: int32(id)})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Only draft quotations can be edited"})
		return
	}
	if err := txQueries.DeleteQuotationItems(c, int32(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to replace quotation services"})
		return
	}
	for index, item := range request.Items {
		service := services[index]
		if _, err := txQueries.AddQuotationItem(c, db.AddQuotationItemParams{QuotationID: int32(id), ServiceID: int32(item.ServiceID), Quantity: int32(item.Quantity), Price: service.Price, Amount: service.Price * float64(item.Quantity)}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quotation service"})
			return
		}
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quotation"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func getQuotations(c *gin.Context, queries *db.Queries) {
	quotations, err := queries.ListQuotations(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get quotations",
		})
		return
	}

	c.JSON(http.StatusOK, quotations)
}

func updateQuotationStatus(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quotation ID",
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
	if request.Status != QuotationDraft && request.Status != QuotationApproved && request.Status != QuotationRejected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quotation status"})
		return
	}
	var linkedWorkOrder bool
	if err := databasePool.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM work_orders WHERE quotation_id = $1 AND deleted_at IS NULL)`, int32(id)).Scan(&linkedWorkOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check quotation workflow"})
		return
	}
	if linkedWorkOrder {
		c.JSON(http.StatusConflict, gin.H{"error": "A quotation linked to a work order cannot have its status changed."})
		return
	}

	quotation, err := queries.UpdateQuotationStatus(
		c,
		db.UpdateQuotationStatusParams{
			Status: request.Status,
			ID:     int32(id),
		},
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Quotation not found",
		})
		return
	}

	c.JSON(http.StatusOK, quotation)
}

func deleteQuotationByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quotation ID",
		})
		return
	}
	quotation, err := queries.GetQuotation(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quotation not found"})
		return
	}
	if quotation.Status != QuotationDraft {
		c.JSON(http.StatusConflict, gin.H{"error": "Only draft quotations can be archived."})
		return
	}
	var linkedWorkOrder bool
	if err := databasePool.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM work_orders WHERE quotation_id = $1 AND deleted_at IS NULL)`, int32(id)).Scan(&linkedWorkOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check quotation workflow"})
		return
	}
	if linkedWorkOrder {
		c.JSON(http.StatusConflict, gin.H{"error": "A quotation linked to a work order cannot be archived."})
		return
	}

	rows, err := queries.SoftDeleteQuotation(c, int32(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete quotation",
		})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quotation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quotation archived successfully",
	})
}

func getQuotationsByStatus(c *gin.Context, queries *db.Queries) {
	status := c.Param("status")

	if status != QuotationDraft &&
		status != QuotationApproved &&
		status != QuotationRejected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quotation status"})
		return
	}

	quotations, err := queries.GetQuotationsByStatus(c, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get quotations"})
		return
	}

	if quotations == nil {
		quotations = []db.GetQuotationsByStatusRow{}
	}

	c.JSON(http.StatusOK, quotations)
}
