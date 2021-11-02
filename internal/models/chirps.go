package models

import (
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