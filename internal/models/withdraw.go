package models

import "time"

type Withdraw struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Order       string    `json:"order_number"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float64   `json:"sum"`
}
