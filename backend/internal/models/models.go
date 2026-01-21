package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stock struct {
	ID           string    `gorm:"primaryKey;type:uuid" json:"id"`
	Symbol       string    `gorm:"unique;not null" json:"symbol"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`          // "ad" (Common Stock), "re" (REIT), etc.
	Currency     string    `json:"currency"`      // "USD"
	CurrentPrice float64   `json:"current_price"` // Cached latest price
	HighPrice    float64   `json:"high_price"`
	LowPrice     float64   `json:"low_price"`
	OpenPrice    float64   `json:"open_price"`
	PrevClose    float64   `json:"prev_close"`
	InWatchlist  bool      `gorm:"default:false" json:"in_watchlist"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StockPrice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StockID   string    `gorm:"index;not null;type:uuid" json:"stock_id"` // FK
	Price     float64   `gorm:"not null" json:"price"`
	Timestamp time.Time `gorm:"index;not null" json:"timestamp"`
}

func (s *Stock) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return
}

type Recommendation struct {
	ID             string    `gorm:"primaryKey;type:uuid" json:"id"`
	StockID        string    `gorm:"index;not null;type:uuid" json:"stock_id"` // FK
	Recommendation string    `json:"recommendation"`
	Reason         string    `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
}
