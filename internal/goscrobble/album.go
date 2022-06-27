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
func insertAlbum(name string, mbid string, spotifyId string, artists []string, img string, tx *sql.Tx) (Album, error) {
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
		err := insertNewAlbum(&album, name, mbid, spotifyId, img, tx)
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

	if album.Img == "" {
		if img != "" {
			err := importImage(album.UUID, img)
			if err != nil {
				log.Printf("Failed to import image: %+v. For Album: %s", err, album.Name)
				return album, nil
			}

			album.Img = "pending"
			_ = album.updateAlbum("img", "pending", tx)
		}
	}

	return album, nil
}

func getAlbumByCol(col string, val string, tx *sql.Tx) Album {
	var album Album
	err := tx.QueryRow(
		`SELECT uuid, name, COALESCE(desc, ''), COALESCE(img,''), mbid, spotify_id FROM albums WHERE "`+col+`" = $1`,
		val).Scan(&album.UUID, &album.Name, &album.Desc, &album.Img, &album.MusicBrainzID, &album.SpotifyID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching albums: %+v", err)
		}
	}

	return album
}

func insertNewAlbum(album *Album, name string, mbid string, spotifyId string, img string, tx *sql.Tx) error {
	album.UUID = newUUID()
	album.Name = name
	album.MusicBrainzID = mbid
	album.SpotifyID = spotifyId
	album.Img = img

	_, err := tx.Exec(`INSERT INTO albums (uuid, name, mbid, spotify_id, img) `+
		`VALUES ($1,$2,$3,$4,$5)`, album.UUID, album.Name, album.MusicBrainzID, album.SpotifyID, album.Img)

	return err
}

func (album *Album) linkAlbumToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec(`INSERT INTO album_artist (album, artist) `+
			`VALUES ($1,$2)`, album.UUID, artist)
		if err != nil {
			return err
		}
	}

	return err
}

func (album *Album) updateAlbum(col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE albums SET "`+col+`" = $1 WHERE uuid = $2`, val, album.UUID)

	return err
}

func getAlbumByUUID(uuid string) (Album, error) {
	var album Album
	err := db.QueryRow(`SELECT uuid, name, COALESCE(desc,''), COALESCE(img,''), mbid, spotify_id FROM albums WHERE uuid = $1`,
		uuid).Scan(&album.UUID, &album.Name, &album.Desc, &album.Img, &album.MusicBrainzID, &album.SpotifyID)

	if err != nil {
		return album, errors.New("Invalid UUID")
	}

	return album, nil
}

// getTopUsersForAlbumUUID  - Returns list of top users for a track
func getTopUsersForAlbumUUID(albumUUID string, limit int, page int) (TopUserResponse, error) {
	response := TopUserResponse{}
	var count int

	total, err := getDbCount(
		"SELECT COUNT(*) FROM `scrobbles` "+
			"JOIN `track_album` ON `track_album`.`track` = `scrobbles`.`track` "+
			"WHERE `track_album`.`album` = UUID_TO_BIN(?, true);", albumUUID)

	if err != nil {
		log.Printf("Failed to fetch scrobble count: %+v", err)
		return response, errors.New("Failed to fetch combined scrobbles")
	}

	rows, err := db.Query(
		"SELECT BIN_TO_UUID(`scrobbles`.`user`, true), `users`.`username`, COUNT(*) "+
			"FROM `scrobbles` "+
			"JOIN `users` ON `scrobbles`.`user` = `users`.`uuid` "+
			"JOIN `track_album` ON `track_album`.`track` = `scrobbles`.`track` "+
			"WHERE `track_album`.`album` = UUID_TO_BIN(?, true) "+
			"GROUP BY `scrobbles`.`user` "+
			"ORDER BY COUNT(*) DESC LIMIT ?",
		albumUUID, limit)

	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return response, errors.New("Failed to fetch combined scrobbles")
	}
	defer rows.Close()

	for rows.Next() {
		item := TopUserResponseItem{}
		err := rows.Scan(&item.UserUUID, &item.UserName, &item.Count)
		if err != nil {
			log.Printf("Failed to fetch scrobbles: %+v", err)
			return response, errors.New("Failed to fetch combined scrobbles")
		}
		count++
		response.Items = append(response.Items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return response, errors.New("Failed to fetch scrobbles")
	}

	response.Meta.Count = count
	response.Meta.Total = total
	response.Meta.Page = page

	return response, nil
}
