package postgres

import (
	"context"
	"errors"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
)

func (db *PostgresDatabase) CreateOrder(ctx context.Context, order models.Order) error {
	query := `INSERT INTO orders (user_id, number, status, accrual) VALUES($1, $2, $3, $4)`

	_, err := db.conn.ExecContext(ctx, query, order.UserID, order.Number, order.Status, order.Accrual)

	var pgErr *pq.Error

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			existingOrder, err := db.getOrder(ctx, order.Number)
			if err != nil {
				return err
			}
			if existingOrder.UserID == order.UserID {
				return handlers.NewErrorWithDB(err, "OrderAlreadyRegisterByYou")
			}
			return handlers.NewErrorWithDB(err, "OrderAlreadyRegister")
		}

		if pgErr.Code == pgerrcode.UndefinedTable {
			return handlers.NewErrorWithDB(err, "UndefinedTable")
		}
	}

	return err
}

func (db *PostgresDatabase) getOrder(ctx context.Context, number string) (*models.Order, error) {
	order := &models.Order{}

	query := `SELECT id, user_id, number, status, uploaded_at, accrual FROM orders WHERE number=$1`

	row := db.conn.QueryRowContext(ctx, query, number)

	err := row.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &order.Accrual)
	if err != nil {
		return nil, err
	}

	return order, nil
}
