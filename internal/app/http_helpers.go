package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func pathID(c *gin.Context, resource string) (int32, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + resource + " ID"})
		return 0, false
	}
	return int32(id), true
}

func apiError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
