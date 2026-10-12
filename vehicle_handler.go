package main

import (
	db "garago/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func createVehicle(c *gin.Context, queries *db.Queries) {

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
		db.CreateVehicleParams{
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

func getVehicles(c *gin.Context, queries *db.Queries) {

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
func getVehicleByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid vehicle ID",
		})
		return
	}

	vehicle, err := queries.GetVehicle(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Vehicle not found",
		})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

//get the vehicle by customer id

func getVehiclesByCustomerID(c *gin.Context, queries *db.Queries) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
		})
		return
	}

	vehicles, err := queries.GetVehiclesByCustomerID(c, int32(customerID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get vehicles",
		})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// update vehicle

func updateVehicleByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid vehicle ID",
		})
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
		db.UpdateVehicleParams{
			CustomerID: int32(vehicle.CustomerID),
			RegNo:      vehicle.RegNo,
			Vehicle:    vehicle.Vehicle,
			ID:         int32(id),
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

func deleteVehicleByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid vehicle ID",
		})
		return
	}

	err = queries.DeleteVehicle(c, int32(id))

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
