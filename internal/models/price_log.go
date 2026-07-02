package models

import "time"

type StockPriceLog struct {
	ID        int64     `json:"id"`
	TickerID  int       `json:"ticker_id"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
}
