package models

type Order struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	UploadedAt string  `json:"uploaded_at"`
	Accrual    float64 `json:"accrual"`
}

type ResponseOrder struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}
