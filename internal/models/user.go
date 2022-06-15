package models

type UserBalance struct {
	Balance float64 `json:"current"`
	Spent   float64 `json:"withdrawn"`
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	UserBalance
}
