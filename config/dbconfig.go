package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "hello.db")
	if err != nil {
		log.Fatal(err)
	}
	if errPing := db.Ping(); errPing != nil {
		log.Panicln(errPing)
	}

	return db
}
