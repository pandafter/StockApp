package stocksync

import (
	"net/http"

	"stockapp/pkg/database"

	"github.com/gin-gonic/gin"
)

func SeedDataHandler(c *gin.Context) {
	if err := SeedMockData(database.GetDB()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database seeded successfully with dummy data"})
}
