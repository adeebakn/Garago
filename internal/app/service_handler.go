package app

import (
	"net/http"
	"strconv"

	db "garago/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func createService(c *gin.Context, queries *db.Queries) {
	var service Service

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service data",
		})
		return
	}
	if service.Name == "" || service.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name and a positive price are required"})
		return
	}

	createdService, err := queries.CreateService(
		c,
		db.CreateServiceParams{
			Name: service.Name,
			Description: pgtype.Text{
				String: service.Description,
				Valid:  true,
			},
			Price: service.Price,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create service",
		})
		return
	}

	c.JSON(http.StatusCreated, createdService)
}

func getServices(c *gin.Context, queries *db.Queries) {
	services, err := queries.ListServices(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get services",
		})
		return
	}

	c.JSON(http.StatusOK, services)
}

func getServiceByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	service, err := queries.GetService(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Service not found",
		})
		return
	}

	c.JSON(http.StatusOK, service)
}

func updateServiceByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	var service Service

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service data",
		})
		return
	}
	if service.Name == "" || service.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name and a positive price are required"})
		return
	}

	updatedService, err := queries.UpdateService(
		c,
		db.UpdateServiceParams{
			Name: service.Name,
			Description: pgtype.Text{
				String: service.Description,
				Valid:  true,
			},
			Price: service.Price,
			ID:    int32(id),
		},
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Service not found",
		})
		return
	}

	c.JSON(http.StatusOK, updatedService)
}

func deleteServiceByID(c *gin.Context, queries *db.Queries) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	err = queries.DeleteService(c, int32(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete service",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service deleted successfully",
	})
}
