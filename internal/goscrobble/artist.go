package goscrobble

import (
	"database/sql"
	"errors"
	"log"
)

type Artist struct {
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Img           string `json:"img"`
	MusicBrainzID string `json:"mbid"`
	SpotifyID     string `json:"spotify_id"`
}

// insertArtist - This will return if it exists or create it based on MBID > Name
func insertArtist(name string, mbid string, spotifyId string, tx *sql.Tx) (Artist, error) {
	artist := Artist{}

	// Try locate our artist
	if mbid != "" {
		artist = getArtistByCol("mbid", mbid, tx)
	} else if spotifyId != "" {
		artist = getArtistByCol("spotify_id", spotifyId, tx)
	}

	// If it didn't match above, lookup by name
	if artist.UUID == "" {
		artist = getArtistByCol("name", name, tx)
	}

	// If we can't find it. Lets add it!
	if artist.UUID == "" {
		err := insertNewArtist(&artist, name, mbid, spotifyId, tx)
		if err != nil {
			log.Printf("Error inserting artist: %+v", err)
			return artist, errors.New("Failed to insert artist")
		}

		if artist.UUID == "" {
			return artist, errors.New("Unable to fetch artist!")
		}
	}

	// Updates these values if we match earlier!M
	if artist.MusicBrainzID == "" {
		if artist.MusicBrainzID != mbid {
			artist.MusicBrainzID = mbid
			artist.updateArtist("mbid", mbid, tx)
		}
	}

	if artist.SpotifyID == "" {
		if artist.SpotifyID != spotifyId {
			artist.SpotifyID = spotifyId
			artist.updateArtist("spotify_id", spotifyId, tx)
		}
	}

	return artist, nil
}

func getArtistByCol(col string, val string, tx *sql.Tx) Artist {
	var artist Artist
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, IFNULL(`desc`,''), IFNULL(`img`,''), `mbid`, `spotify_id` FROM `artists` WHERE `"+col+"` = ?",
		val).Scan(&artist.UUID, &artist.Name, &artist.Desc, &artist.Img, &artist.MusicBrainzID, &artist.SpotifyID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching artists: %+v", err)
		}
	}

	return artist
}

func insertNewArtist(artist *Artist, name string, mbid string, spotifyId string, tx *sql.Tx) error {
	artist.UUID = newUUID()
	artist.Name = name
	artist.MusicBrainzID = mbid
	artist.SpotifyID = spotifyId

	_, err := tx.Exec("INSERT INTO `artists` (`uuid`, `name`, `mbid`, `spotify_id`) "+
		"VALUES (UUID_TO_BIN(?, true),?,?,?)", artist.UUID, artist.Name, artist.MusicBrainzID, artist.SpotifyID)

	return err
}

func (artist *Artist) updateArtist(col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE `artists` SET `"+col+"` = ? WHERE `uuid` = UUID_TO_BIN(?,true)", val, artist.UUID)

	return err
}

func getArtistByUUID(uuid string) (Artist, error) {
	var artist Artist
	err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), `name`, IFNULL(`desc`, ''), IFNULL(`img`,''), `mbid`, `spotify_id` FROM `artists` WHERE `uuid` = UUID_TO_BIN(?, true)",
		uuid).Scan(&artist.UUID, &artist.Name, &artist.Desc, &artist.Img, &artist.MusicBrainzID, &artist.SpotifyID)

	if err == sql.ErrNoRows {
		return artist, errors.New("Invalid UUID")
	}

	return artist, nil
}
