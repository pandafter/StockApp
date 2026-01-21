package main

import (
	"log"
	"time"

	"stockapp/internal/db"
	"stockapp/internal/handlers"
	"stockapp/internal/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Stock App Backend...")

	// 1. Initialize Database
	db.InitDB()
	log.Println("Database initialization finished")

	// 2. Start Background Worker
	worker.StartWorker()

	// 3. Setup Router
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	api := r.Group("/api")
	{
		api.GET("/stocks", handlers.GetStocks)
		api.GET("/stocks/:symbol", handlers.GetStockDetail)
		api.POST("/stocks/:symbol/watchlist", handlers.ToggleWatchlist)
		api.GET("/recommendation", handlers.GetRecommendation)
		api.GET("/seed", handlers.SeedData)
	}

	// 4. Run Server
	log.Println("Server listening on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
