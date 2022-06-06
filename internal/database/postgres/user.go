package postgres

import (
	"context"
	"errors"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (db *PostgresDatabase) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	query := `INSERT INTO users (first_name, last_name, login, password) VALUES ($1, $2, $3, $4)`

	_, err := db.conn.ExecContext(ctx, query, user.FirstName, user.LastName, user.Login, user.Password)

	var pgErr *pq.Error

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return nil, handlers.NewErrorWithDB(err, "UniqConstraint")
		}

		if pgErr.Code == pgerrcode.UndefinedTable {
			return nil, handlers.NewErrorWithDB(err, "UndefinedTable")
		}
	}

	resultUser, err := db.getUserByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}

	return resultUser, err
}

func (db *PostgresDatabase) getUserByLogin(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, login, first_name, last_name FROM users WHERE login=$1`

	row := db.conn.QueryRowContext(ctx, query, login)

	err := row.Scan(&user.ID, &user.Login, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *PostgresDatabase) CheckPassword(ctx context.Context, user models.User) (*models.User, error) {
	u := &models.User{}

	query := `SELECT id, login, first_name, last_name FROM users WHERE login=$1 AND password=$2`

	row := db.conn.QueryRowContext(ctx, query, user.Login, user.Password)

	err := row.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, errors.New("wrong login or password")
	}

	if u.ID == "" {
		return nil, errors.New("wrong login or password")
	}

	return u, nil
}
