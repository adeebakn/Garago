package main

import (
	"fmt"
	"net/http"
	"strconv"

	db "garago/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func createCustomer(c *gin.Context, queries *db.Queries) {
	var customer Customer

	//Read JSON from the request
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Customer Data",
		})
		return
	}

	//Save customer to database

	newCustomer, err := queries.CreateCustomer(
		c,
		db.CreateCustomerParams{
			Name:  customer.Name,
			Phone: customer.Phone,
			Email: pgtype.Text{
				String: customer.Email,
				Valid:  customer.Email != "",
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create customer",
		})
		return
	}

	c.JSON(http.StatusCreated, newCustomer)

}

// get customers handler

func getCustomers(c *gin.Context, queries *db.Queries) {

	customers, err := queries.ListCustomers(c)
	//customers, err := queries.ListCustomers(c)
	if err != nil {
		fmt.Println("ListCustomers error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

//get customer by id

func getCustomerByID(c *gin.Context, queries *db.Queries) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	customer, err := queries.GetCustomer(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	c.JSON(http.StatusOK, customer)
}

//update customer by id

func updateCustomerByID(c *gin.Context, queries *db.Queries) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	var customer Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer data",
		})
		return
	}

	updatedCustomer, err := queries.UpdateCustomer(
		c,
		db.UpdateCustomerParams{
			Name:  customer.Name,
			Phone: customer.Phone,
			Email: pgtype.Text{
				String: customer.Email,
				Valid:  customer.Email != "",
			},
			ID: int32(id),
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

func deleteCustomerByID(c *gin.Context, queries *db.Queries) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
		})
		return
	}

	err = queries.DeleteCustomer(c, int32(id))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer deleted successfully",
	})
}
