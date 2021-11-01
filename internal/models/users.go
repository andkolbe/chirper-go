package models

import (
	"database/sql"
	"log"
	"time"
)

// we can only use string and int safely because we set NOT NULL constraints on all of the columns on the table
type User struct {
	ID         string       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_At time.Time `json:"created_at"`
}

// Create a custom UserModel type which wraps the sql.DB connection pool
type UserModel struct {
	DB *sql.DB
}

// Use a method on the custom UserModel type to run the SQL query.
func (m UserModel) GetAllUsers() ([]User, error) {
	rows, err := m.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (m UserModel) GetUserByID(id string) (User, error) {
	
	row := m.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At)
	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}

func (m UserModel) CreateNewUser(user User) {

	// add password hashing 
	
	_, err := m.DB.Exec("INSERT INTO users (name, email, password) VALUES(?, ?, ?)", &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Fatal(err)
	}
}

func (m UserModel) UpdateUser(user User) {
	_, err := m.DB.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", &user.Name, &user.Email, &user.Password, &user.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func (m UserModel) DeleteUser(id string) {
	_, err := m.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

}