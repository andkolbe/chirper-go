package driver

import (
	"database/sql"
	"log"
	"time"

	// We don't use anything in the mysql package directly, which means that the Go compiler will raise an error if we try to import it normally
	// But we need the mysql package's init() function to run so that our driver can register itself with database/sql
	// We get around this by aliasing the package name to the blank identifier
	// This means mysql.init() still gets executed, but the alias is harmlessly discarded (and our code runs error-free)
	// This approach is standard for most of Go's SQL drivers
	_ "github.com/go-sql-driver/mysql"
)

// creates a new database for the application
func DBConnect(dsn string) (*sql.DB, error) {
	// initialise a new sql.DB object by calling sql.Open()
	// the sql.DB object it returns is not a database connection – it's an abstraction representing a pool of underlying connections
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

	// Set the maximum number of concurrently open connections (in-use + idle) to 10. 
	// Setting this to less than or equal to 0 will mean there is no maximum limit (which is also the default setting)
	db.SetMaxOpenConns(10)

	// Set the maximum number of concurrently idle connections to 5. Setting this to less than or equal to 0 will mean that no idle connections are retained
	db.SetMaxIdleConns(5)


	// Set the maximum lifetime of a connection to 5 minutes
	// Setting it to 0 means that there is no maximum lifetime and the connection is reused forever (which is the default behavior)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}