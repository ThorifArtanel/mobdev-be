package common

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DbConn() (*sql.DB, error) {
	conn, err := sql.Open("pgx", GetDBURL())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return conn, nil
}
