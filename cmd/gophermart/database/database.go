package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/logger"
)

func Conn(driverName, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn can not be missing")
	}

	if driverName == "" {
		return nil, fmt.Errorf("driver name can not be missing")
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return db, err
	}

	log.Println("Connect to postgres")

	return db, nil
}

func Migrations(db *sql.DB, logger logger.Logger) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	migrationsPath := basePath + "/migrations"

	err := goose.Up(db, migrationsPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
