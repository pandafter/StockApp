package stocksync

import (
	"log"
	"time"

	"stockapp/internal/stocks"
	"stockapp/pkg/database"

	"gorm.io/gorm"
)

func StartWorker() {
	go func() {
		log.Println("Starting Background Worker")

		db := database.GetDB()
		if db == nil {
			log.Println("Database not ready, skipping initial fetch")
			return
		}

		// Initial fetch
		log.Println("Running initial fetch...")
		if err := runSync(db); err != nil {
			log.Printf("Error in initial fetch: %v", err)
		} else {
			log.Println("Initial fetch completed")
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Running scheduled fetch...")
			if err := runSync(db); err != nil {
				log.Printf("Error in scheduled fetch: %v", err)
			}
		}
	}()
}

func runSync(db *gorm.DB) error {
	if err := FetchAndStoreStocks(db); err != nil {
		return err
	}

	// If no data, seed mock
	count, err := stocks.CountStocks()
	if err != nil {
		log.Printf("Error counting stocks: %v", err)
		return err
	}

	if count == 0 {
		log.Println(">>> Database is empty. Seeding mock data...")
		return SeedMockData(db)
	}

	return nil
}
