package goscrobble

import (
	"database/sql"
	"log"
)

type Genre struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func getGenre(uuid string) Genre {
	var genre Genre
	err := db.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name` FROM `artists` WHERE `uuid` = UUID_TO_BIN(?,true)",
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
		"SELECT BIN_TO_UUID(`uuid`, true), `name` FROM `artists` WHERE `name` = ?",
		name).Scan(&genre.UUID, &genre.Name)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching artists: %+v", err)
		}
	}

	return genre
}

func (genre *Genre) updateGenreName(name string, value string) error {
	_, err := db.Exec("UPDATE `genres` SET `name` = ? WHERE uuid = UUID_TO_BIN(?, true)", name, genre.UUID)

	return err
}
