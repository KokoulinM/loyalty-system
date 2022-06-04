package models

import (
	"encoding/json"
	"time"
)

type WithdrawOrder struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Withdraw struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
	WithdrawOrder
}

func (w WithdrawOrder) MarshalJSON() ([]byte, error) {
	type WithdrawAlias WithdrawOrder
	aliasValue := struct {
		WithdrawAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawAlias: WithdrawAlias(w),
		ProcessedAt:   w.ProcessedAt.Format(time.RFC3339),
	}
	return json.Marshal(aliasValue)
}
