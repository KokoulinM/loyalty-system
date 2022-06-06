package postgres

import "database/sql"

type PostgresDatabase struct {
	conn *sql.DB
}

func New(db *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		conn: db,
	}
}
