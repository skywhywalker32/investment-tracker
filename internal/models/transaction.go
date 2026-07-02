package models

import "time"

type OperationType string

const (
	OpBuy  OperationType = "buy"
	OpSell OperationType = "sell"
)

type Transaction struct {
	ID            int64         `json:"id"`
	UserID        int           `json:"user_id"`
	TickerID      int           `json:"ticker_id"`
	OperationType OperationType `json:"operation_type"`
	Qty           int           `json:"qty"`
	Price         float64       `json:"price"`
	Currency      string        `json:"currency"`
	CreatedAt     time.Time     `json:"created_at"`
}
