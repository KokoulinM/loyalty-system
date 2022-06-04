package models

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Order       string    `json:"order"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float64   `json:"sum"`
}

func (w Withdraw) MarshalJSON() ([]byte, error) {
	type WithdrawAlias Withdraw
	aliasValue := struct {
		WithdrawAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawAlias: WithdrawAlias(w),
		ProcessedAt:   w.ProcessedAt.Format(time.RFC3339),
	}
	return json.Marshal(aliasValue)
}
