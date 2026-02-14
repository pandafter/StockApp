package main

import (
	"log"

	"stockapp/internal/stocks"
	"stockapp/internal/stocksync"
	"stockapp/pkg/database"
	"stockapp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Stock App Backend...")

	// 1. Initialize Database
	if _, err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	database.InitDB(&stocks.Stock{}, &stocks.StockPrice{})
	log.Println("Database initialization finished")

	// 2. Start Background Worker
	stocksync.StartWorker()

	// 3. Setup Router
	r := gin.Default()
	r.Use(middleware.CORS())

	// API Routes
	api := r.Group("/api")
	{
		stocks.SetupRoutes(api)
		api.GET("/seed", stocksync.SeedDataHandler)
	}

	// 4. Run Server
	log.Println("Server listening on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
