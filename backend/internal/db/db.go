package db

import (
	"log"

	"stockapp/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "postgresql://root@localhost:26257/defaultdb?sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
	err = DB.AutoMigrate(&models.Stock{}, &models.StockPrice{})
	if err != nil {
		log.Printf("Warning: Migration reported an error, attempting to continue: %v", err)
	} else {
		log.Println("Database migrated")
	}
}

func GetDB() *gorm.DB {
	return DB
}
