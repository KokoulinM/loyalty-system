package postgres

import (
	"context"
	"time"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
)

func (db *PostgresDatabase) CreateWithdraw(ctx context.Context, withdraw models.Withdraw, userID string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var enough bool

	row := tx.QueryRowContext(ctx, "SELECT balance >= $1 FROM users WHERE id = $2", withdraw.Sum, userID)

	err = row.Scan(&enough)
	if err != nil {
		return err
	}

	if !enough {
		return handlers.NewErrorWithDB(err, "NotEnoughBalanceForWithdraw")
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1, spend = spend + $1 WHERE id=$2", withdraw.Sum, userID)
	if err != nil {
		return err
	}

	query := `INSERT INTO withdrawals (user_id, order_number, status, sum, processed_at) VALUES ($1, $2, $3, $4, $5)`

	_, err = db.conn.ExecContext(ctx, query, userID, withdraw.Order, "PROCESSED", withdraw.Sum, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *PostgresDatabase) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawOrder, error) {
	var withdrawals []models.WithdrawOrder

	query := `SELECT order_number, sum, processed_at FROM withdrawals
					 WHERE user_id = $1 ORDER BY processed_at ASC`

	rows, err := db.conn.QueryContext(ctx, query, userID)
	if err != nil {
		return withdrawals, err
	}

	defer rows.Close()

	for rows.Next() {
		var withdraw models.WithdrawOrder

		err := rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return withdrawals, err
		}

		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
