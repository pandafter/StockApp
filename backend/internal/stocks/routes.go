package stocks

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	stocks := r.Group("/stocks")
	{
		stocks.GET("", GetStocksHandler)
		stocks.GET("/:symbol", GetStockDetailHandler)
		stocks.POST("/:symbol/watchlist", ToggleWatchlistHandler)
	}
	r.GET("/recommendation", GetRecommendationHandler)
}
