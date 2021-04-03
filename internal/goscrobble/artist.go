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
	MusicBrainzID string         `json:"mbid"`
	SpotifyID     string         `json:"spotify_id"`
}

// insertArtist - This will return if it exists or create it based on MBID > Name
func insertArtist(name string, mbid string, spotifyId string, tx *sql.Tx) (Artist, error) {
	artist := Artist{}

	// Try locate our artist
	if mbid != "" {
		artist = fetchArtist("mbid", mbid, tx)
	} else if spotifyId != "" {
		artist = fetchArtist("spotify_id", spotifyId, tx)
	}

	// If it didn't match above, lookup by name
	if artist.Uuid == "" {
		artist = fetchArtist("name", name, tx)
	}

	// If we can't find it. Lets add it!
	if artist.Uuid == "" {
		err := insertNewArtist(name, mbid, spotifyId, tx)
		if err != nil {
			log.Printf("Error inserting artist: %+v", err)
			return artist, errors.New("Failed to insert artist")
		}
	}

	// Fetch the recently inserted artist to get the UUID
	if mbid != "" {
		artist = fetchArtist("mbid", mbid, tx)
	} else if spotifyId != "" {
		artist = fetchArtist("spotify_id", spotifyId, tx)
	}

	if artist.Uuid == "" {
		artist = fetchArtist("name", name, tx)
	}

	if artist.Uuid == "" {
		return artist, errors.New("Unable to fetch artist!")
	}

	if artist.MusicBrainzID != mbid {
		artist.MusicBrainzID = mbid
		artist.updateArtist("mbid", mbid, tx)
	}

	if artist.SpotifyID != spotifyId {
		artist.SpotifyID = spotifyId
		artist.updateArtist("spotify_id", spotifyId, tx)
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

func insertNewArtist(name string, mbid string, spotifyId string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `artists` (`uuid`, `name`, `mbid`, `spotify_id`) "+
		"VALUES (UUID_TO_BIN(UUID(), true),?,?,?)", name, mbid, spotifyId)

	return err
}

func (artist *Artist) updateArtist(col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE `artists` SET `"+col+"` = ? WHERE `uuid` = UUID_TO_BIN(?,true)", val, artist.Uuid)

	return err
}
