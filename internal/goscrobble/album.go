package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Album struct {
	Uuid          string         `json:"uuid"`
	Name          string         `json:"name"`
	Desc          sql.NullString `json:"desc"`
	Img           sql.NullString `json:"img"`
	MusicBrainzID sql.NullString `json:"mbid"`
	SpotifyID     sql.NullString `json:"spotify_id"`
}

// insertAlbum - This will return if it exists or create it based on MBID > Name
func insertAlbum(name string, mbid string, spotifyId string, artists []string, tx *sql.Tx) (Album, error) {
	album := Album{}

	// Try locate our album
	if mbid != "" {
		album = fetchAlbum("mbid", mbid, tx)
	} else if spotifyId != "" {
		album = fetchAlbum("spotify_id", spotifyId, tx)
	}

	// If it didn't match above, lookup by name
	if album.Uuid == "" {
		album = fetchAlbum("name", name, tx)
	}

	// If we can't find it. Lets add it!
	if album.Uuid == "" {
		err := insertNewAlbum(name, mbid, spotifyId, tx)
		if err != nil {
			return album, errors.New("Failed to insert album")
		}

		// Fetch the recently inserted album to get the UUID
		if mbid != "" {
			album = fetchAlbum("mbid", mbid, tx)
		} else if spotifyId != "" {
			album = fetchAlbum("spotify_id", spotifyId, tx)
		}

		if album.Uuid == "" {
			album = fetchAlbum("name", name, tx)
		}

		// Try linkem up
		err = album.linkAlbumToArtists(artists, tx)
		if err != nil {
			return album, errors.New("Unable to link albums!")
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

func insertNewAlbum(name string, mbid string, spotifyId string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `albums` (`uuid`, `name`, `mbid`, `spotify_id`) "+
		"VALUES (UUID_TO_BIN(UUID(), true),?,?,?)", name, mbid, spotifyId)

	return err
}

func (album *Album) linkAlbumToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec("INSERT INTO `album_artist` (`album`, `artist`) "+
			"VALUES (UUID_TO_BIN(?, true), UUID_TO_BIN(?, true))", album.Uuid, artist)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return err
}

func updateAlbum(uuid string, col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE `albums` SET `"+col+"` = ? WHERE `uuid` = UUID_TO_BIN(?,true)", val, uuid)

	return err
}
