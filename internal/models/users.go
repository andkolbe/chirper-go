package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
)

// we can only use string and int safely because we set NOT NULL constraints on all of the columns on the table
// If the table contained nullable fields we would need to use the sql.NullString and sql.NullInt types instead
type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_At time.Time `json:"-"`
}

// Create a custom DBModel type which wraps the sql.DB connection pool
type DBModel struct {
	DB *sql.DB
}

// GET All Users
// Use a method on the custom UserModel type to run the SQL query
func (m DBModel) GetAllUsers() ([]User, error) {
	// fetch a result set from the userss table using the DB.Query() method and assign it to a rows variable
	rows, err := m.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Closing a result set properly is important
	// As long as a resultset is open it will keep the underlying database connection open – which in turn means the connection is not available to the pool
	// So if something goes wrong and the result set isn't closed it can rapidly lead to all the connections in your pool being used up
	// the defer statement should come after you check for an error from DB.Query
	// Otherwise, if DB.Query() returns an error, you'll get a panic trying to close a nil resultset
	defer rows.Close()

	var users []User
	// use rows.Next() to iterate through the rows in the result set
	//  This preps the first (and then each subsequent) row to be acted on by the rows.Scan() method
	// if iteration over all of the rows completes then the result set automatically closes itself and frees-up the connection
	for rows.Next() {
		var user User
		// We use the rows.Scan() method to copy the values from each field in the row to a new User object that we created
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// add the new User to the users slice we created earlier
		users = append(users, user)
	}
	// When our rows.Next() loop has finished we call rows.Err(). This returns any error that was encountered during the interation
	// It's important to call this – don't just assume that we completed a successful iteration over the whole resultset
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// GET One User
func (m DBModel) GetUserByID(id string) (User, error) {

	// Because we need to include untrusted input (the id variable) in our SQL query, we take advantage of placeholder parameters, passing in the value of our 
	// placeholder as the second argument to DB.QueryRow()
	// Behind the scenes, db.QueryRow (and also db.Query() and db.Exec()) work by creating a new prepared statement on the database, 
	// and subsequently execute that prepared statement using the placeholder parameters provided
	// This means that all three methods are safe from SQL injection when used correctly
	// Prepared statements are resilient against SQL injection, because parameter values, which are transmitted later using a different protocol, 
	// need not be correctly escaped
	// If the original statement template is not derived from external input, injection cannot occur
	row := m.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At)
	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}

// POST
func (m DBModel) RegisterNewUser(user User) int64 {

	// CreateHash returns a Argon2id hash of a plain-text password using the provided algorithm parameters
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	// DB.Exec(), like DB.Query() and DB.QueryRow(), is a variadic function, which means you can pass in as many parameters as you need
	// The db.Exec() method returns an object satisfying the sql.Result interface, which you can either use or discard with the blank identifier
	res, err := m.DB.Exec("INSERT INTO users (name, email, password) VALUES(?, ?, ?)", &user.Name, &user.Email, hash)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// The sql.Result() interface guarantees two methods: LastInsertId() – which is often used to return the value of an new auto increment id
	id, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

// PUT
func (m DBModel) UpdateUser(user User, id string) int64 {
	res, err := m.DB.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", &user.Name, &user.Email, &user.Password, id)
	if err != nil {
		log.Fatal(err)
	}

	// The sql.Result() interface guarantees two methods: RowsAffected() – which contains the number of rows that the statement affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// DELETE
func (m DBModel) DeleteUser(id string) int64 {
	res, err := m.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }
    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}

// Authenticate User
func (m DBModel) AuthenticateUser(email, password string) (int, error) {
	var id int 
	var hashedPassword string 

	row := m.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return 0, err
	}
	if !match {
		log.Panicln("incorrect password")
	} 

	return id, nil
}