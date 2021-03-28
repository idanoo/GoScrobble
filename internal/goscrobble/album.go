package goscrobble

import (
	"database/sql"
	"errors"
	"log"
)

type Album struct {
	Uuid          string         `json:"uuid"`
	Name          string         `json:"name"`
	Desc          sql.NullString `json:"desc"`
	Img           sql.NullString `json:"img"`
	MusicBrainzID sql.NullString `json:"mbid"`
}

// insertAlbum - This will return if it exists or create it based on MBID > Name
func insertAlbum(name string, mbid string, artists []string, tx *sql.Tx) (Album, error) {
	album := Album{}

	if mbid != "" {
		album = fetchAlbum("mbid", mbid, tx)
		if album.Uuid == "" {
			err := insertNewAlbum(name, mbid, tx)
			if err != nil {
				log.Printf("Error inserting album via MBID %s  %+v", name, err)
				return album, errors.New("Failed to insert album")
			}

			album = fetchAlbum("mbid", mbid, tx)
			album.linkAlbumToArtists(artists, tx)
		}
	} else {
		album = fetchAlbum("name", name, tx)
		if album.Uuid == "" {
			err := insertNewAlbum(name, mbid, tx)
			if err != nil {
				log.Printf("Error inserting album via Name %s %+v", name, err)
				return album, errors.New("Failed to insert album")
			}

			album = fetchAlbum("name", name, tx)
			album.linkAlbumToArtists(artists, tx)
		}
	}

	if album.Uuid == "" {
		return album, errors.New("Unable to fetch album!")
	}

	return album, nil
}

func fetchAlbum(col string, val string, tx *sql.Tx) Album {
	var album Album
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, `desc`, `img`, `mbid` FROM `albums` WHERE `"+col+"` = ?",
		val).Scan(&album.Uuid, &album.Name, &album.Desc, &album.Img, &album.MusicBrainzID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching albums: %+v", err)
		}
	}

	return album
}

func insertNewAlbum(name string, mbid string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `albums` (`uuid`, `name`, `mbid`) "+
		"VALUES (UUID_TO_BIN(UUID(), true),?,?)", name, mbid)

	return err
}

func (album *Album) linkAlbumToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec("INSERT INTO `track_artist` (`track`, `artist`) "+
			"VALUES (UUID_TO_BIN(?, true), UUID_TO_BIN(?, true)", album.Uuid, artist)
		if err != nil {
			return err
		}
	}

	return err
}
