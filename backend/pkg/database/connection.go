package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://root@localhost:26257/defaultdb?sslmode=disable"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return DB, nil
}

func InitDB(models ...interface{}) {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	if len(models) > 0 {
		if err := DB.AutoMigrate(models...); err != nil {
			log.Printf("Warning: Migration reported an error, attempting to continue: %v", err)
		} else {
			log.Println("Database migrated")
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}
