package database

import (
	"database/sql"
	"fmt"
	"log"
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
