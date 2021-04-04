package goscrobble

import (
	"database/sql"
	"errors"
	"log"
)

type Album struct {
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Img           string `json:"img"`
	MusicBrainzID string `json:"mbid"`
	SpotifyID     string `json:"spotify_id"`
}

// insertAlbum - This will return if it exists or create it based on MBID > Name
func insertAlbum(name string, mbid string, spotifyId string, artists []string, tx *sql.Tx) (Album, error) {
	album := Album{}

	// Try locate our album
	if mbid != "" {
		album = getAlbumByCol("mbid", mbid, tx)
	} else if spotifyId != "" {
		album = getAlbumByCol("spotify_id", spotifyId, tx)
	}

	// If it didn't match above, lookup by name
	if album.UUID == "" {
		album = getAlbumByCol("name", name, tx)
	}

	// If we can't find it. Lets add it!
	if album.UUID == "" {
		err := insertNewAlbum(&album, name, mbid, spotifyId, tx)
		if err != nil {
			return album, errors.New("Failed to insert album")
		}

		if album.UUID == "" {
			return album, errors.New("Failed to fetch album")
		}

		// Try linkem up
		err = album.linkAlbumToArtists(artists, tx)
		if err != nil {
			return album, errors.New("Unable to link albums!")
		}
	}

	// Updates these values if we match earlier!
	if album.MusicBrainzID != mbid {
		album.MusicBrainzID = mbid
		album.updateAlbum("mbid", mbid, tx)
	}

	if album.SpotifyID != spotifyId {
		album.SpotifyID = spotifyId
		album.updateAlbum("spotify_id", spotifyId, tx)
	}

	return album, nil
}

func getAlbumByCol(col string, val string, tx *sql.Tx) Album {
	var album Album
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, IFNULL(`desc`, ''), IFNULL(`img`,''), `mbid`, `spotify_id` FROM `albums` WHERE `"+col+"` = ?",
		val).Scan(&album.UUID, &album.Name, &album.Desc, &album.Img, &album.MusicBrainzID, &album.SpotifyID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching albums: %+v", err)
		}
	}

	return album
}

func insertNewAlbum(album *Album, name string, mbid string, spotifyId string, tx *sql.Tx) error {
	album.UUID = newUUID()
	album.Name = name
	album.MusicBrainzID = mbid
	album.SpotifyID = spotifyId

	_, err := tx.Exec("INSERT INTO `albums` (`uuid`, `name`, `mbid`, `spotify_id`) "+
		"VALUES (UUID_TO_BIN(?, true),?,?,?)", album.UUID, album.Name, album.MusicBrainzID, album.SpotifyID)

	return err
}

func (album *Album) linkAlbumToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec("INSERT INTO `album_artist` (`album`, `artist`) "+
			"VALUES (UUID_TO_BIN(?, true), UUID_TO_BIN(?, true))", album.UUID, artist)
		if err != nil {
			return err
		}
	}

	return err
}

func (album *Album) updateAlbum(col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE `albums` SET `"+col+"` = ? WHERE `uuid` = UUID_TO_BIN(?,true)", val, album.UUID)

	return err
}

func getAlbumByUUID(uuid string) (Album, error) {
	var album Album
	err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), `name`, IFNULL(`desc`,''), IFNULL(`img`,''), `mbid`, `spotify_id` FROM `albums` WHERE `uuid` = UUID_TO_BIN(?, true)",
		uuid).Scan(&album.UUID, &album.Name, &album.Desc, &album.Img, &album.MusicBrainzID, &album.SpotifyID)

	if err != nil {
		return album, errors.New("Invalid UUID")
	}

	return album, nil
}
