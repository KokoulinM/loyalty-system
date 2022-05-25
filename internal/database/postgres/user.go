package postgres

import (
	"context"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
)

func (db *PostgresDatabase) CreateUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)`

	_, err := db.conn.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}
