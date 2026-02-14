package stocksync

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"stockapp/internal/stocks"

	"gorm.io/gorm"
)

const (
	BaseURL   = "https://api.karenai.click/swechallenge/list"
	LoginURL  = "https://api.karenai.click/swechallenge/login"
	AuthToken = "swe-valid-token-1768932226032"
)

type APIResponse struct {
	Items    []APIStock `json:"items"`
	NextPage string     `json:"next_page"`
}

type APIStock struct {
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Currency     string  `json:"currency"`
	CurrentPrice float64 `json:"current_price"`
	HighPrice    float64 `json:"high_price"`
	LowPrice     float64 `json:"low_price"`
	OpenPrice    float64 `json:"open_price"`
	PrevClose    float64 `json:"prev_close"`
}

func Login() error {
	log.Println("Sending login signal...")

	req, err := http.NewRequest("POST", LoginURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send login signal: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Warning: Login signal returned status %s. Continuing anyway...", resp.Status)
	} else {
		log.Println("Login signal sent successfully")
	}

	return nil
}

func FetchAndStoreStocks(db *gorm.DB) error {
	log.Println(">>> Starting FetchAndStoreStocks process...")

	if err := Login(); err != nil {
		log.Printf("Login signal error (non-fatal): %v", err)
	}

	nextPage := ""
	maxPages := 20
	pageCount := 0

	for {
		pageCount++
		if pageCount > maxPages {
			log.Println("Max pages reached, stopping fetch")
			break
		}

		url := BaseURL
		if nextPage != "" {
			url = fmt.Sprintf("%s?next_page=%s", BaseURL, nextPage)
		}

		log.Printf("Fetching page %d: %s", pageCount, url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+AuthToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("API returned status: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return err
		}

		if len(apiResp.Items) == 0 {
			log.Println("No items found on this page")
			break
		}

		for i, item := range apiResp.Items {
			if err := db.Transaction(func(tx *gorm.DB) error {
				// We need to use stocks package with the DB from transaction
				// For simplicity, do upsert in transaction
				var existing stocks.Stock
				result := tx.Where("symbol = ?", item.Symbol).First(&existing)

				if result.Error == nil {
					existing.CurrentPrice = item.CurrentPrice
					existing.HighPrice = item.HighPrice
					existing.LowPrice = item.LowPrice
					existing.OpenPrice = item.OpenPrice
					existing.PrevClose = item.PrevClose
					if err := tx.Save(&existing).Error; err != nil {
						return err
					}
					return tx.Create(&stocks.StockPrice{
						StockID:   existing.ID,
						Price:     item.CurrentPrice,
						Timestamp: time.Now(),
					}).Error
				}

				newStock := stocks.Stock{
					Symbol:       item.Symbol,
					Name:         item.Name,
					Type:         item.Type,
					Currency:     item.Currency,
					CurrentPrice: item.CurrentPrice,
					HighPrice:    item.HighPrice,
					LowPrice:     item.LowPrice,
					OpenPrice:    item.OpenPrice,
					PrevClose:    item.PrevClose,
					UpdatedAt:    time.Now(),
				}
				if err := tx.Create(&newStock).Error; err != nil {
					return err
				}
				return tx.Create(&stocks.StockPrice{
					StockID:   newStock.ID,
					Price:     item.CurrentPrice,
					Timestamp: time.Now(),
				}).Error
			}); err != nil {
				log.Printf("Error storing %s: %v", item.Symbol, err)
				continue
			}

			if i%10 == 0 {
				log.Printf("Processed %d stocks on current page", i+1)
			}
		}

		nextPage = apiResp.NextPage
		if nextPage == "" {
			break
		}
	}

	log.Println(">>> FetchAndStoreStocks process completed")
	return nil
}
