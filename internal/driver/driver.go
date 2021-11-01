package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// creates a new database for the application
func DBConnect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err	
	}

	log.Println("Connected to db!")

	// because sql.Open doesn't actually check a connection, we also call DB.Ping() to make sure that everything works OK on startup
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Pinged db!")

	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}