package handlers

import (
	"net/http"

	"stockapp/internal/api"
	"stockapp/internal/commands"
	"stockapp/internal/db"

	"github.com/gin-gonic/gin"
)

func GetStocks(c *gin.Context) {
	cmd := &commands.GetStocksCommand{
		Search:        c.Query("search"),
		SortBy:        c.Query("sort_by"),
		Order:         c.Query("order"),
		WatchlistOnly: c.Query("watchlist") == "true",
	}

	result, err := cmd.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetStockDetail(c *gin.Context) {
	cmd := &commands.GetStockDetailCommand{
		Symbol: c.Param("symbol"),
	}

	result, err := cmd.Execute()
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetRecommendation(c *gin.Context) {
	cmd := &commands.GetRecommendationCommand{}

	result, err := cmd.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No stocks available"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func ToggleWatchlist(c *gin.Context) {
	symbol := c.Param("symbol")
	cmd := &commands.ToggleWatchlistCommand{Symbol: symbol}

	result, err := cmd.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SeedData(c *gin.Context) {
	if err := api.SeedMockData(db.GetDB()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database seeded successfully with dummy data"})
}
