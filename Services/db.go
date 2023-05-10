package services

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type RowScanner interface {
    Scan(dest ...interface{}) error
}

var db_conn *sql.DB

func GetDB() *sql.DB {
	var err error
	if db_conn != nil {
		return db_conn
	}
    db_conn, err = sql.Open("mysql", os.Getenv("DSN"))
    if err != nil {
        log.Printf("failed to connect: %v", err)
    }
    if err := db_conn.Ping(); err != nil {
        log.Printf("failed to ping: %v", err)
    }
	return db_conn
}