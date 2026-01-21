package worker

import (
	"log"
	"time"

	"stockapp/internal/api"
)

func StartWorker() {
	go func() {
		log.Println("Starting Background Worker")
		
		// Initial fetch
		log.Println("Running initial fetch...")
		if err := api.FetchAndStoreStocks(); err != nil {
			log.Printf("Error in initial fetch: %v", err)
		} else {
			log.Println("Initial fetch completed")
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Running scheduled fetch...")
			if err := api.FetchAndStoreStocks(); err != nil {
				log.Printf("Error in scheduled fetch: %v", err)
			}
		}
	}()
}
