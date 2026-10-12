package apis

import (
	"errors"
	"garago/internal/database/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strings"
)

func createVehicle(c *gin.Context, queries *postgres.Queries) {

	var vehicle Vehicle

	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle data"})
		return
	}
	vehicle.RegNo = strings.ToUpper(strings.TrimSpace(vehicle.RegNo))
	vehicle.Vehicle = strings.TrimSpace(vehicle.Vehicle)
	if vehicle.CustomerID <= 0 || vehicle.Vehicle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer and vehicle are required"})
		return
	}

	createdVehicle, err := queries.CreateVehicle(
		c,
		postgres.CreateVehicleParams{
			CustomerID: int32(vehicle.CustomerID),
			RegNo:      vehicle.RegNo,
			Vehicle:    vehicle.Vehicle,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "vehicles_customer_registration_unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "This registration number already exists for the selected customer"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create vehicle"})
		return
	}

	c.JSON(http.StatusCreated, createdVehicle)
}

//lsit all vehicles

func getVehicles(c *gin.Context, queries *postgres.Queries) {

	vehicles, err := queries.ListVehicles(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get vehicles",
		})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// get  vehicle by id
func getVehicleByID(c *gin.Context, queries *postgres.Queries) {
	id, ok := pathID(c, "vehicle")
	if !ok {
		return
	}

	vehicle, err := queries.GetVehicle(c, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Vehicle not found",
		})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

//get the vehicle by customer id

func getVehiclesByCustomerID(c *gin.Context, queries *postgres.Queries) {
	customerID, ok := pathID(c, "customer")
	if !ok {
		return
	}

	vehicles, err := queries.GetVehiclesByCustomerID(c, customerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get vehicles",
		})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// update vehicle

func updateVehicleByID(c *gin.Context, queries *postgres.Queries) {
	id, ok := pathID(c, "vehicle")
	if !ok {
		return
	}

	var vehicle Vehicle

	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid vehicle data",
		})
		return
	}
	vehicle.RegNo = strings.ToUpper(strings.TrimSpace(vehicle.RegNo))
	vehicle.Vehicle = strings.TrimSpace(vehicle.Vehicle)
	if vehicle.CustomerID <= 0 || vehicle.Vehicle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer and vehicle are required"})
		return
	}

	updatedVehicle, err := queries.UpdateVehicle(
		c,
		postgres.UpdateVehicleParams{
			CustomerID: int32(vehicle.CustomerID),
			RegNo:      vehicle.RegNo,
			Vehicle:    vehicle.Vehicle,
			ID:         id,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "vehicles_customer_registration_unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "This registration number already exists for the selected customer"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle"})
		return
	}

	c.JSON(http.StatusOK, updatedVehicle)
}

//delete vehicle

func deleteVehicleByID(c *gin.Context, queries *postgres.Queries) {
	id, ok := pathID(c, "vehicle")
	if !ok {
		return
	}

	err := archiveVehicle(c.Request.Context(), id)
	if errors.Is(err, errVehicleHasHistory) {
		apiError(c, http.StatusConflict, "Vehicle cannot be deleted because quotation or work-order history is linked to this vehicle")
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		apiError(c, http.StatusNotFound, "Vehicle not found")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vehicle deleted successfully",
	})
}
