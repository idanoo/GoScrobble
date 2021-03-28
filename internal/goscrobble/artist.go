package goscrobble

import (
	"database/sql"
	"errors"
	"log"
)

type Artist struct {
	Uuid          string         `json:"uuid"`
	Name          string         `json:"name"`
	Desc          sql.NullString `json:"desc"`
	Img           sql.NullString `json:"img"`
	MusicBrainzID sql.NullString `json:"mbid"`
}

// insertArtist - This will return if it exists or create it based on MBID > Name
func insertArtist(name string, mbid string, tx *sql.Tx) (Artist, error) {
	artist := Artist{}

	if mbid != "" {
		artist = fetchArtist("mbid", mbid, tx)
		if artist.Uuid == "" {
			err := insertNewArtist(name, mbid, tx)
			if err != nil {
				log.Printf("Error inserting artist via MBID %s  %+v", name, err)
				return artist, errors.New("Failed to insert artist")
			}

			artist = fetchArtist("mbid", mbid, tx)
		}
	} else {
		artist = fetchArtist("name", name, tx)
		if artist.Uuid == "" {
			err := insertNewArtist(name, mbid, tx)
			if err != nil {
				log.Printf("Error inserting artist via Name %s %+v", name, err)
				return artist, errors.New("Failed to insert artist")
			}

			artist = fetchArtist("name", name, tx)
		}
	}

	if artist.Uuid == "" {
		return artist, errors.New("Unable to fetch artist!")
	}

	return artist, nil
}

func fetchArtist(col string, val string, tx *sql.Tx) Artist {
	var artist Artist
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, `desc`, `img`, `mbid` FROM `artists` WHERE `"+col+"` = ?",
		val).Scan(&artist.Uuid, &artist.Name, &artist.Desc, &artist.Img, &artist.MusicBrainzID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching artists: %+v", err)
		}
	}

	return artist
}

func insertNewArtist(name string, mbid string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `artists` (`uuid`, `name`, `mbid`) "+
		"VALUES (UUID_TO_BIN(UUID(), true),?,?)", name, mbid)

	return err
}
