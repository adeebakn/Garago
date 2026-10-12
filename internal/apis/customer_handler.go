package apis

import (
	"errors"
	"net/http"
	"strings"

	"garago/internal/database/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func createCustomer(c *gin.Context, queries *postgres.Queries) {
	var customer Customer

	//Read JSON from the request
	if err := c.ShouldBindJSON(&customer); err != nil {
		apiError(c, http.StatusBadRequest, "Invalid customer data")
		return
	}
	customer.Name = strings.TrimSpace(customer.Name)
	customer.Phone = strings.TrimSpace(customer.Phone)
	customer.Email = strings.ToLower(strings.TrimSpace(customer.Email))
	if customer.Name == "" || customer.Phone == "" {
		apiError(c, http.StatusBadRequest, "Customer name and phone are required")
		return
	}

	//Save customer to database

	newCustomer, err := queries.CreateCustomer(
		c,
		postgres.CreateCustomerParams{
			Name:  customer.Name,
			Phone: customer.Phone,
			Email: pgtype.Text{
				String: customer.Email,
				Valid:  customer.Email != "",
			},
		},
	)

	if err != nil {
		apiError(c, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	c.JSON(http.StatusCreated, newCustomer)

}

// get customers handler

func getCustomers(c *gin.Context, queries *postgres.Queries) {

	customers, err := queries.ListCustomers(c)
	//customers, err := queries.ListCustomers(c)
	if err != nil {
		apiError(c, http.StatusInternalServerError, "Failed to get customers")
		return
	}

	c.JSON(http.StatusOK, customers)
}

//get customer by id

func getCustomerByID(c *gin.Context, queries *postgres.Queries) {

	id, ok := pathID(c, "customer")
	if !ok {
		return
	}

	customer, err := queries.GetCustomer(c, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	c.JSON(http.StatusOK, customer)
}

//update customer by id

func updateCustomerByID(c *gin.Context, queries *postgres.Queries) {

	id, ok := pathID(c, "customer")
	if !ok {
		return
	}

	var customer Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer data",
		})
		return
	}
	customer.Name = strings.TrimSpace(customer.Name)
	customer.Phone = strings.TrimSpace(customer.Phone)
	customer.Email = strings.ToLower(strings.TrimSpace(customer.Email))
	if customer.Name == "" || customer.Phone == "" {
		apiError(c, http.StatusBadRequest, "Customer name and phone are required")
		return
	}

	updatedCustomer, err := queries.UpdateCustomer(
		c,
		postgres.UpdateCustomerParams{
			Name:  customer.Name,
			Phone: customer.Phone,
			Email: pgtype.Text{
				String: customer.Email,
				Valid:  customer.Email != "",
			},
			ID: id,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update customer",
		})
		return
	}

	c.JSON(http.StatusOK, updatedCustomer)

}

//delete customer by id

func deleteCustomerByID(c *gin.Context, queries *postgres.Queries) {

	id, ok := pathID(c, "customer")
	if !ok {
		return
	}

	err := archiveCustomer(c.Request.Context(), id)
	if errors.Is(err, errCustomerHasHistory) {
		apiError(c, http.StatusConflict, "Customer cannot be deleted because vehicles or business history are linked to this customer")
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		apiError(c, http.StatusNotFound, "Customer not found")
		return
	}
	if err != nil {
		apiError(c, http.StatusInternalServerError, "Failed to delete customer")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer deleted successfully",
	})
}
