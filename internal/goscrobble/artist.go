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

type TopArtist struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Img   string `json:"img"`
	Plays int    `json:"plays"`
}

type TopArtists struct {
	Artists map[int]TopArtist `json:"artists"`
}

// insertArtist - This will return if it exists or create it based on MBID > Name
func insertArtist(name string, mbid string, spotifyId string, img string, tx *sql.Tx) (Artist, error) {
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
		err := insertNewArtist(&artist, name, mbid, spotifyId, img, tx)
		if err != nil {
			log.Printf("Error inserting artist: %+v", err)
			return artist, errors.New("Failed to insert artist")
		}

		if artist.UUID == "" {
			return artist, errors.New("Unable to fetch artist!")
		}
	}

	// Updates these values if we have the data!
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

	if artist.Img == "" {
		if img != "" {
			artist.Img = img
			artist.updateArtist("img", img, tx)
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

func insertNewArtist(artist *Artist, name string, mbid string, spotifyId string, img string, tx *sql.Tx) error {
	artist.UUID = newUUID()
	artist.Name = name
	artist.MusicBrainzID = mbid
	artist.SpotifyID = spotifyId
	artist.Img = img

	_, err := tx.Exec("INSERT INTO `artists` (`uuid`, `name`, `mbid`, `spotify_id`, `img`) "+
		"VALUES (UUID_TO_BIN(?, true),?,?,?,?)", artist.UUID, artist.Name, artist.MusicBrainzID, artist.SpotifyID, artist.Img)

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

func getTopArtists(userUuid string) (TopArtists, error) {
	var topArtist TopArtists

	rows, err := db.Query("SELECT BIN_TO_UUID(`artists`.`uuid`, true), `artists`.`name`, IFNULL(BIN_TO_UUID(`artists`.`uuid`, true),''), count(*) "+
		"FROM `scrobbles` "+
		"JOIN `tracks` ON `tracks`.`uuid` = `scrobbles`.`track` "+
		"JOIN track_artist ON track_artist.track = tracks.uuid "+
		"JOIN artists ON track_artist.artist = artists.uuid "+
		"WHERE `scrobbles`.`user` = UUID_TO_BIN(?, true) "+
		"GROUP BY `artists`.`uuid` "+
		"ORDER BY count(*) DESC "+
		"LIMIT 14;",
		userUuid)
	if err != nil {
		log.Printf("Failed to fetch top artist: %+v", err)
		return topArtist, errors.New("Failed to fetch top artist")
	}
	defer rows.Close()

	i := 1
	tempArtists := make(map[int]TopArtist)

	for rows.Next() {
		var artist TopArtist
		err := rows.Scan(&artist.UUID, &artist.Name, &artist.Img, &artist.Plays)
		if err != nil {
			log.Printf("Failed to fetch artist: %+v", err)
			return topArtist, errors.New("Failed to fetch artist")
		}

		tempArtists[i] = artist
		i++
	}

	topArtist.Artists = tempArtists

	return topArtist, nil
}

// getTopUsersForArtistUUID  - Returns list of top users for a track
func getTopUsersForArtistUUID(trackUUID string, limit int, page int) (TopUserResponse, error) {
	response := TopUserResponse{}

	// TODO: Implement

	// var count int

	// total, err := getDbCount(
	// 	"SELECT COUNT(*) FROM `scrobbles` WHERE `track` = UUID_TO_BIN(?, true) GROUP BY `track`, `user`", trackUUID)

	// if err != nil {
	// 	log.Printf("Failed to fetch scrobble count: %+v", err)
	// 	return response, errors.New("Failed to fetch combined scrobbles")
	// }

	// rows, err := db.Query(
	// 	"SELECT BIN_TO_UUID(`scrobbles`.`user`, true), `users`.`username`, COUNT(*) "+
	// 		"FROM `scrobbles` "+
	// 		"JOIN `users` ON `scrobbles`.`user` = `users`.`uuid` "+
	// 		"WHERE `track` = UUID_TO_BIN(?, true) "+
	// 		"GROUP BY `scrobbles`.`user` "+
	// 		"ORDER BY COUNT(*) DESC LIMIT ?",
	// 	trackUUID, limit)

	// if err != nil {
	// 	log.Printf("Failed to fetch scrobbles: %+v", err)
	// 	return response, errors.New("Failed to fetch combined scrobbles")
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	item := TopUserResponseItem{}
	// 	err := rows.Scan(&item.UserUUID, &item.UserName, &item.Count)
	// 	if err != nil {
	// 		log.Printf("Failed to fetch scrobbles: %+v", err)
	// 		return response, errors.New("Failed to fetch combined scrobbles")
	// 	}
	// 	count++
	// 	response.Items = append(response.Items, item)
	// }

	// err = rows.Err()
	// if err != nil {
	// 	log.Printf("Failed to fetch scrobbles: %+v", err)
	// 	return response, errors.New("Failed to fetch scrobbles")
	// }

	// response.Meta.Count = count
	// response.Meta.Total = total
	// response.Meta.Page = page

	return response, nil
}
