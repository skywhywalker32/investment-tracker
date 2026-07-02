package models

import "time"

type Portfolio struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TickerID  int       `json:"ticker_id"`
	Qty       int       `json:"qty"`
	UpdatedAt time.Time `json:"updated_at"`
}
