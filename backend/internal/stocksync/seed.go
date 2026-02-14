package stocksync

import (
	"log"
	"time"

	"stockapp/internal/stocks"

	"gorm.io/gorm"
)

var mockStocks = []stocks.Stock{
	{Symbol: "AAPL", Name: "Apple Inc.", Type: "ad", Currency: "USD", CurrentPrice: 185.92, HighPrice: 199.62, LowPrice: 164.08, OpenPrice: 184.22, PrevClose: 183.63, UpdatedAt: time.Now()},
	{Symbol: "MSFT", Name: "Microsoft Corp.", Type: "ad", Currency: "USD", CurrentPrice: 402.56, HighPrice: 410.22, LowPrice: 380.12, OpenPrice: 400.10, PrevClose: 398.45, UpdatedAt: time.Now()},
	{Symbol: "GOOGL", Name: "Alphabet Inc.", Type: "ad", Currency: "USD", CurrentPrice: 145.32, HighPrice: 155.00, LowPrice: 130.45, OpenPrice: 144.10, PrevClose: 146.20, UpdatedAt: time.Now()},
	{Symbol: "AMZN", Name: "Amazon.com Inc.", Type: "ad", Currency: "USD", CurrentPrice: 172.44, HighPrice: 180.12, LowPrice: 150.34, OpenPrice: 170.15, PrevClose: 171.22, UpdatedAt: time.Now()},
	{Symbol: "NVDA", Name: "NVIDIA Corp.", Type: "ad", Currency: "USD", CurrentPrice: 726.13, HighPrice: 750.00, LowPrice: 600.00, OpenPrice: 715.00, PrevClose: 710.15, UpdatedAt: time.Now()},
	{Symbol: "TSLA", Name: "Tesla Inc.", Type: "ad", Currency: "USD", CurrentPrice: 191.59, HighPrice: 215.00, LowPrice: 175.45, OpenPrice: 190.00, PrevClose: 188.12, UpdatedAt: time.Now()},
	{Symbol: "META", Name: "Meta Platforms Inc.", Type: "ad", Currency: "USD", CurrentPrice: 473.32, HighPrice: 485.96, LowPrice: 450.12, OpenPrice: 470.00, PrevClose: 468.10, UpdatedAt: time.Now()},
	{Symbol: "NFLX", Name: "Netflix Inc.", Type: "ad", Currency: "USD", CurrentPrice: 580.12, HighPrice: 600.00, LowPrice: 550.00, OpenPrice: 575.00, PrevClose: 578.00, UpdatedAt: time.Now()},
	{Symbol: "AMD", Name: "Advanced Micro Devices", Type: "ad", Currency: "USD", CurrentPrice: 175.44, HighPrice: 185.00, LowPrice: 160.00, OpenPrice: 170.00, PrevClose: 172.00, UpdatedAt: time.Now()},
	{Symbol: "PYPL", Name: "PayPal Holdings", Type: "ad", Currency: "USD", CurrentPrice: 62.15, HighPrice: 75.00, LowPrice: 55.00, OpenPrice: 60.00, PrevClose: 61.50, UpdatedAt: time.Now()},
}

func SeedMockData(db *gorm.DB) error {
	log.Println(">>> Seeding mock data...")

	for _, s := range mockStocks {
		var existing stocks.Stock
		result := db.Where("symbol = ?", s.Symbol).First(&existing)

		if result.Error == nil {
			db.Model(&existing).Updates(stocks.Stock{
				CurrentPrice: s.CurrentPrice,
				HighPrice:    s.HighPrice,
				LowPrice:     s.LowPrice,
				OpenPrice:    s.OpenPrice,
				PrevClose:    s.PrevClose,
				UpdatedAt:    time.Now(),
			})
		} else {
			if err := db.Create(&s).Error; err != nil {
				log.Printf("Failed to create %s: %v", s.Symbol, err)
				return err
			}
			existing = s
		}

		now := time.Now()
		for i := 10; i >= 0; i-- {
			offset := float64(i) * 1.5
			price := existing.CurrentPrice - offset
			if price < 0 {
				price = 10.0
			}
			db.Create(&stocks.StockPrice{
				StockID:   existing.ID,
				Price:     price,
				Timestamp: now.Add(time.Duration(-i) * time.Hour),
			})
		}
		log.Printf(">>> Seeded %s (with history)", existing.Symbol)
	}

	log.Println(">>> seedMockData completed successfully")
	return nil
}
