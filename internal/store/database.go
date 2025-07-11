package store

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres "+
		"dbname=postgres port=5439 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db open method %w\n", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(1 * time.Minute)

	fmt.Printf("Connected to database")
	return db, nil
}
