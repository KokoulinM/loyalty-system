package database

import (
	"database/sql"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
)

type Database struct {
	conn *sql.DB
	cfg  config.Config
}

func DatabaseRepository(db *sql.DB, cfg config.Config) *Database {
	return &Database{
		conn: db,
		cfg:  cfg,
	}
}

func (db *Database) AddUser(user models.User) error {
	return nil
}
