package models

import (
	"fmt"
	"log"
	"time"
)

type Chirp struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userid"`
	Content    string    `json:"content"`
	Location   string    `json:"location"`
	Created_At time.Time `json:"created_at"`
}

// GET All Chirps 
func (m DBModel) GetAllChirps() ([]Chirp, error) {
	rows, err := m.DB.Query("SELECT * FROM chirps")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var chirps []Chirp
	
	for rows.Next() {
		var chirp Chirp
		err := rows.Scan(&chirp.ID, &chirp.UserID, &chirp.Content, &chirp.Location, &chirp.Created_At)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		chirps = append(chirps, chirp)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return chirps, nil
}

// GET One User
func (m DBModel) GetChirpByID(id string) (Chirp, error) {
	row := m.DB.QueryRow("SELECT * FROM chirps WHERE id = ?", id)

	var chirp Chirp
	err := row.Scan(&chirp.ID, &chirp.UserID, &chirp.Content, &chirp.Location, &chirp.Created_At)
	if err != nil {
		log.Fatal(err)
	}

	return chirp, nil
}

// POST
func (m DBModel) CreateNewChirp(chirp Chirp) int64 {
	res, err := m.DB.Exec("INSERT INTO chirps (userid, content, location) VALUES(?, ?, ?)", &chirp.UserID, &chirp.Content, &chirp.Location)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

// PUT
func (m DBModel) UpdateChirp(chirp Chirp, id string) int64 {
	res, err := m.DB.Exec("UPDATE chirps SET userid = ?, content = ?, location = ? WHERE id = ?", &chirp.UserID, &chirp.Content, &chirp.Location, id)
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

// DELETE
func (m DBModel) DeleteChirp(id string) int64{
	res, err := m.DB.Exec("DELETE FROM chirps WHERE id = ?", id)
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