package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OrganizationSettings is the single shop profile used throughout this installation.
// The application deliberately has one row (id = 1); multi-tenant support is out of scope.
type OrganizationSettings struct {
	ID      int16  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

func getOrganizationSettings(c *gin.Context) {
	var settings OrganizationSettings
	err := databasePool.QueryRow(c, `SELECT id, name, address, phone, email FROM organization_settings WHERE id = 1`).Scan(
		&settings.ID, &settings.Name, &settings.Address, &settings.Phone, &settings.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load organization settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func updateOrganizationSettings(c *gin.Context) {
	var input OrganizationSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization settings"})
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization name is required"})
		return
	}

	var settings OrganizationSettings
	err := databasePool.QueryRow(c, `
		INSERT INTO organization_settings (id, name, address, phone, email, updated_at)
		VALUES (1, $1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, address = EXCLUDED.address, phone = EXCLUDED.phone,
			email = EXCLUDED.email, updated_at = NOW()
		RETURNING id, name, address, phone, email`, input.Name, input.Address, input.Phone, input.Email,
	).Scan(&settings.ID, &settings.Name, &settings.Address, &settings.Phone, &settings.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save organization settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}
