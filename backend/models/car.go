package models

import "time"

type Car struct {
	ID          int       `json:"id"`
	Model       string    `json:"model"`
	Title       string    `json:"title"`
	Price       string    `json:"price"`
	Category    string    `json:"category"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
