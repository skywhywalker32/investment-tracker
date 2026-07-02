package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	Nickname     string    `json:"nickname"`
	CreatedAt    time.Time `json:"created_at"`
}
