package stocks

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStocksHandler(c *gin.Context) {
	search := c.Query("search")
	sortBy := c.Query("sort_by")
	order := c.Query("order")
	watchlistOnly := c.Query("watchlist") == "true"

	result, err := GetStocks(search, sortBy, order, watchlistOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetStockDetailHandler(c *gin.Context) {
	symbol := c.Param("symbol")

	result, err := GetStockDetail(symbol)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetRecommendationHandler(c *gin.Context) {
	result, err := GetRecommendation()
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

func ToggleWatchlistHandler(c *gin.Context) {
	symbol := c.Param("symbol")

	result, err := ToggleWatchlistStatus(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
