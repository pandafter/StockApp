package stocks

import "gorm.io/gorm"

func GetStocks(search, sortBy, order string, watchlistOnly bool) ([]Stock, error) {
	return FindStocks(search, sortBy, order, watchlistOnly)
}

func GetStockDetail(symbol string) (map[string]interface{}, error) {
	stock, err := FindBySymbol(symbol)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	prices, err := FindPricesByStockID(stock.ID, 50)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stock":   stock,
		"history": prices,
	}, nil
}

func GetRecommendation() (map[string]interface{}, error) {
	stock, score, err := GetBestRecommendation()
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, nil
	}
	return map[string]interface{}{
		"recommendation":         stock,
		"potential_gain_percent": score,
		"reason":                 "Highest potential upside to 52-week high",
	}, nil
}

func ToggleWatchlistStatus(symbol string) (*Stock, error) {
	return ToggleWatchlist(symbol)
}

