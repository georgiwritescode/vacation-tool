package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {

	cfg := mysql.Config{
		User:                 "portal",
		Passwd:               "password123",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "vacation_tool",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected! Databse listening on port :3306")
	return db, nil
}
