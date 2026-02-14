package stocks

import (
	"strings"

	"stockapp/pkg/database"
)

func FindStocks(search string, sortBy, order string, watchlistOnly bool) ([]Stock, error) {
	db := database.GetDB()
	var stocks []Stock

	query := db.Model(&Stock{})

	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(symbol) LIKE ? OR LOWER(name) LIKE ?", term, term)
	}

	if watchlistOnly {
		query = query.Where("in_watchlist = ?", true)
	}

	if sortBy != "" {
		if order != "desc" {
			order = "asc"
		}
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("symbol asc")
	}

	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

func FindBySymbol(symbol string) (*Stock, error) {
	db := database.GetDB()
	var stock Stock
	if err := db.Where("symbol = ?", symbol).First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func FindPricesByStockID(stockID string, limit int) ([]StockPrice, error) {
	db := database.GetDB()
	var prices []StockPrice
	if limit <= 0 {
		limit = 50
	}
	if err := db.Where("stock_id = ?", stockID).Order("timestamp desc").Limit(limit).Find(&prices).Error; err != nil {
		return nil, err
	}
	return prices, nil
}

func GetBestRecommendation() (*Stock, float64, error) {
	db := database.GetDB()
	var stocks []Stock
	if err := db.Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	if len(stocks) == 0 {
		return nil, 0, nil
	}

	var bestStock *Stock
	bestScore := -1.0

	for i := range stocks {
		s := &stocks[i]
		if s.CurrentPrice <= 0 || s.HighPrice == s.LowPrice {
			continue
		}
		potential := (s.HighPrice - s.CurrentPrice) / s.CurrentPrice
		if potential > bestScore {
			bestScore = potential
			bestStock = s
		}
	}

	if bestStock == nil {
		return &stocks[0], 0, nil
	}
	return bestStock, bestScore * 100, nil
}

func ToggleWatchlist(symbol string) (*Stock, error) {
	db := database.GetDB()
	var stock Stock
	if err := db.Where("symbol = ?", symbol).First(&stock).Error; err != nil {
		return nil, err
	}
	stock.InWatchlist = !stock.InWatchlist
	if err := db.Save(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func CountStocks() (int64, error) {
	db := database.GetDB()
	var count int64
	if err := db.Model(&Stock{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func SaveStock(stock *Stock) error {
	db := database.GetDB()
	return db.Save(stock).Error
}

func CreateStock(stock *Stock) error {
	db := database.GetDB()
	return db.Create(stock).Error
}

func CreateStockPrice(price *StockPrice) error {
	db := database.GetDB()
	return db.Create(price).Error
}

func UpsertStock(stock *Stock) (*Stock, error) {
	db := database.GetDB()
	var existing Stock
	result := db.Where("symbol = ?", stock.Symbol).First(&existing)

	if result.Error == nil {
		existing.CurrentPrice = stock.CurrentPrice
		existing.HighPrice = stock.HighPrice
		existing.LowPrice = stock.LowPrice
		existing.OpenPrice = stock.OpenPrice
		existing.PrevClose = stock.PrevClose
		existing.UpdatedAt = stock.UpdatedAt
		if err := db.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	if err := db.Create(stock).Error; err != nil {
		return nil, err
	}
	return stock, nil
}
