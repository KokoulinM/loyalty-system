package migratons

import (
	"database/sql"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/pressly/goose/v3"
)

func Migrations(db *sql.DB, logger logger.Logger) {
	err := goose.Up(db, "/")
	if err != nil {
		logger.Fatal(err.Error())
	}
}
