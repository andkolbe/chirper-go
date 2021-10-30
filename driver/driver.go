package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// holds the database connection pool
// if we wanted to add other types of databases, we could add them here too
// type DB struct {
// 	SQL *sql.DB
// }

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifeTime = 5 * time.Minute

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

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetConnMaxIdleTime(maxIdleDBConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	return db, nil
}