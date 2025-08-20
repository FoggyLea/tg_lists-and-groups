package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(sourceName string) error {
	var err error
	DB, err = sql.Open("postgres", sourceName)
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("database is not responding: %v", err)
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
