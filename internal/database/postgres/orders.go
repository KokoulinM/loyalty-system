package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/handlers"
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

func (db *PostgresDatabase) GetOrders(ctx context.Context, userID string) ([]models.ResponseOrderWithAccrual, error) {
	var result []models.ResponseOrderWithAccrual

	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id=$1 ORDER BY uploaded_at`

	rows, err := db.conn.QueryContext(ctx, query, userID)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var order models.ResponseOrderWithAccrual

		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return result, err
		}

		result = append(result, order)
	}

	return result, nil
}

func (db *PostgresDatabase) ChangeOrderStatus(ctx context.Context, order string, status string, accrual float64) error {
	userID, err := db.getUserIDByOrder(ctx, order)
	if err != nil {
		return err
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE orders SET accrual = $1, status = $2 WHERE number = $3`, accrual, status, order)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, accrual, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *PostgresDatabase) getUserIDByOrder(ctx context.Context, order string) (string, error) {
	var userID string

	query := db.conn.QueryRowContext(ctx, `SELECT user_id FROM orders WHERE number = $1`, order)

	err := query.Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
