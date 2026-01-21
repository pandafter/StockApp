package commands

import (
	"strings"

	"stockapp/internal/db"
	"stockapp/internal/models"
)

type GetStocksCommand struct {
	Search        string
	SortBy        string
	Order         string
	WatchlistOnly bool
}

func (c *GetStocksCommand) Execute() (interface{}, error) {
	database := db.GetDB()
	var stocks []models.Stock

	query := database.Model(&models.Stock{})

	if c.Search != "" {
		searchTerm := "%" + strings.ToLower(c.Search) + "%"
		query = query.Where("LOWER(symbol) LIKE ? OR LOWER(name) LIKE ?", searchTerm, searchTerm)
	}

	if c.WatchlistOnly {
		query = query.Where("in_watchlist = ?", true)
	}

	if c.SortBy != "" {
		order := c.Order
		if order != "desc" {
			order = "asc"
		}
		query = query.Order(c.SortBy + " " + order)
	} else {
		query = query.Order("symbol asc")
	}

	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

// GetStockDetailCommand retrieves details and history for a specific stock
type GetStockDetailCommand struct {
	Symbol string
}

func (c *GetStockDetailCommand) Execute() (interface{}, error) {
	database := db.GetDB()
	var stock models.Stock

	if err := database.Where("symbol = ?", c.Symbol).First(&stock).Error; err != nil {
		return nil, err
	}

	var prices []models.StockPrice
	database.Where("stock_id = ?", stock.ID).Order("timestamp desc").Limit(50).Find(&prices)

	return map[string]interface{}{
		"stock":   stock,
		"history": prices,
	}, nil
}

// GetRecommendationCommand identifies the stock with the highest potential upside
type GetRecommendationCommand struct{}

func (c *GetRecommendationCommand) Execute() (interface{}, error) {
	database := db.GetDB()
	var stocks []models.Stock

	if err := database.Find(&stocks).Error; err != nil {
		return nil, err
	}

	if len(stocks) == 0 {
		return nil, nil // Or explicit error if preferred
	}

	bestStock := stocks[0]
	bestScore := -1.0

	for _, s := range stocks {
		if s.CurrentPrice <= 0 || s.HighPrice == s.LowPrice {
			continue
		}

		potential := (s.HighPrice - s.CurrentPrice) / s.CurrentPrice
		if potential > bestScore {
			bestScore = potential
			bestStock = s
		}
	}

	return map[string]interface{}{
		"recommendation":         bestStock,
		"potential_gain_percent": bestScore * 100,
		"reason":                 "Highest potential upside to 52-week high",
	}, nil
}

// ToggleWatchlistCommand updates the watchlist status of a stock
type ToggleWatchlistCommand struct {
	Symbol string
}

func (c *ToggleWatchlistCommand) Execute() (interface{}, error) {
	database := db.GetDB()
	var stock models.Stock

	if err := database.Where("symbol = ?", c.Symbol).First(&stock).Error; err != nil {
		return nil, err
	}

	stock.InWatchlist = !stock.InWatchlist
	if err := database.Save(&stock).Error; err != nil {
		return nil, err
	}

	return stock, nil
}
