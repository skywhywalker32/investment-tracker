package models

import "time"

type Stock struct {
	ID         int       `json:"id"`
	Ticker     string    `json:"ticker"`
	StockName  string    `json:"stock_name"`
	StockPrice float64   `json:"stock_price"`
	Currency   string    `json:"currency"`
	UpdatedAt  time.Time `json:"updated_at"`
	Source     string    `json:"source"`
}
