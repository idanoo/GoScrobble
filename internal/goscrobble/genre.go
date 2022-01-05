package goscrobble

import (
	"database/sql"
	"log"
)

type Genre struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func getGenreByUUID(uuid string) Genre {
	var genre Genre
	err := db.QueryRow(
		`SELECT uuid, name FROM artists WHERE uuid = $1`,
		uuid).Scan(&genre.UUID, &genre.Name)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching artists: %+v", err)
		}
	}

	return genre
}

func getGenreByName(name string) Genre {
	var genre Genre
	err := db.QueryRow(
		`SELECT uuid, name FROM artists WHERE name = $1`,
		name).Scan(&genre.UUID, &genre.Name)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching artists: %+v", err)
		}
	}

	return genre
}

func (genre *Genre) updateGenreName(name string, value string) error {
	_, err := db.Exec(`UPDATE genres SET name = $1 WHERE uuid = $2`, name, genre.UUID)

	return err
}
