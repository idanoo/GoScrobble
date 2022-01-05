package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Scrobble struct {
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	CreatedIp net.IP    `json:"created_ip"`
	User      string    `json:"user"`
	Track     string    `json:"track"`
}

type ScrobbleResponse struct {
	Meta  ScrobbleResponseMeta   `json:"meta"`
	Items []ScrobbleResponseItem `json:"items"`
}

type ScrobbleResponseMeta struct {
	Count int `json:"count"`
	Total int `json:"total"`
	Page  int `json:"page"`
}

type ScrobbleResponseItem struct {
	UUID      string            `json:"uuid"`
	Timestamp time.Time         `json:"time"`
	Artist    ScrobbleTrackItem `json:"artist"`
	Album     string            `json:"album"`
	Track     ScrobbleTrackItem `json:"track"`
	Source    string            `json:"source"`
	User      ScrobbleTrackItem `json:"user"`
}

type ScrobbleTrackItem struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// insertScrobble - This will return if it exists or create it based on MBID > Name
func insertScrobble(user string, track string, source string, timestamp time.Time, ip net.IP, tx *sql.Tx) error {
	err := insertNewScrobble(user, track, source, timestamp, ip, tx)
	if err != nil {
		log.Printf("Error inserting scrobble %s  %+v", user, err)
		return errors.New("Failed to insert scrobble!")
	}

	return nil
}

func getScrobblesForUser(userUuid string, limit int, page int) (ScrobbleResponse, error) {
	scrobbleReq := ScrobbleResponse{}
	var count int

	// Yeah this isn't great. But for now.. it works! Cache later
	total, err := getDbCount(
		"SELECT COUNT(*) FROM `scrobbles` WHERE `user` = UUID_TO_BIN(?, true) ", userUuid)

	if err != nil {
		log.Printf("Failed to fetch scrobble count: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}

	rows, err := db.Query(
		"SELECT BIN_TO_UUID(`scrobbles`.`uuid`, true), `scrobbles`.`created_at`, BIN_TO_UUID(`artists`.`uuid`, true), `artists`.`name`, `albums`.`name`, BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name`, `scrobbles`.`source` FROM `scrobbles` "+
			"JOIN tracks ON scrobbles.track = tracks.uuid "+
			"JOIN track_artist ON track_artist.track = tracks.uuid "+
			"JOIN track_album ON track_album.track = tracks.uuid "+
			"JOIN artists ON track_artist.artist = artists.uuid "+
			"JOIN albums ON track_album.album = albums.uuid "+
			"JOIN users ON scrobbles.user = users.uuid "+
			"WHERE user = UUID_TO_BIN(?, true) "+
			"GROUP BY scrobbles.uuid, albums.uuid "+
			"ORDER BY scrobbles.created_at DESC LIMIT ?",
		userUuid, limit)

	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}
	defer rows.Close()

	for rows.Next() {
		item := ScrobbleResponseItem{}
		err := rows.Scan(&item.UUID, &item.Timestamp, &item.Artist.UUID, &item.Artist.Name, &item.Album, &item.Track.UUID, &item.Track.Name, &item.Source)
		if err != nil {
			log.Printf("Failed to fetch scrobbles: %+v", err)
			return scrobbleReq, errors.New("Failed to fetch scrobbles")
		}
		count++
		scrobbleReq.Items = append(scrobbleReq.Items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}

	scrobbleReq.Meta.Count = count
	scrobbleReq.Meta.Total = total
	scrobbleReq.Meta.Page = page

	return scrobbleReq, nil
}

func insertNewScrobble(user string, track string, source string, timestamp time.Time, ip net.IP, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `scrobbles` (`uuid`, `created_at`, `created_ip`, `user`, `track`, `source`) "+
		"VALUES (UUID_TO_BIN(?, true), ?, ?, UUID_TO_BIN(?, true), UUID_TO_BIN(?, true), ?)", newUUID(), timestamp, ip, user, track, source)

	return err
}

func checkIfScrobbleExists(userUuid string, timestamp time.Time, source string) bool {
	count, err := getDbCount("SELECT COUNT(*) FROM `scrobbles` WHERE `user` = UUID_TO_BIN(?, true) AND `created_at` = ? AND `source` = ?",
		userUuid, timestamp, source)

	if err != nil {
		fmt.Printf("Error fetching scrobble exists count: %+v", err)
		return true
	}

	return count != 0
}

func getRecentScrobbles() (ScrobbleResponse, error) {
	scrobbleReq := ScrobbleResponse{}
	var count int
	limit := 50

	rows, err := db.Query(
		"SELECT BIN_TO_UUID(`scrobbles`.`uuid`, true), `scrobbles`.`created_at`, BIN_TO_UUID(`artists`.`uuid`, true), `artists`.`name`, `albums`.`name`, BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name`, `scrobbles`.`source`, BIN_TO_UUID(`scrobbles`.`user`, true), `users`.`username` FROM `scrobbles` "+
			"JOIN tracks ON scrobbles.track = tracks.uuid "+
			"JOIN track_artist ON track_artist.track = tracks.uuid "+
			"JOIN track_album ON track_album.track = tracks.uuid "+
			"JOIN artists ON track_artist.artist = artists.uuid "+
			"JOIN albums ON track_album.album = albums.uuid "+
			"JOIN users ON scrobbles.user = users.uuid "+
			"GROUP BY scrobbles.uuid, albums.uuid "+
			"ORDER BY scrobbles.created_at DESC LIMIT ?", limit)

	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}
	defer rows.Close()

	for rows.Next() {
		item := ScrobbleResponseItem{}
		err := rows.Scan(&item.UUID, &item.Timestamp, &item.Artist.UUID, &item.Artist.Name, &item.Album, &item.Track.UUID, &item.Track.Name, &item.Source, &item.User.UUID, &item.User.Name)
		if err != nil {
			log.Printf("Failed to fetch scrobbles: %+v", err)
			return scrobbleReq, errors.New("Failed to fetch scrobbles")
		}
		count++
		scrobbleReq.Items = append(scrobbleReq.Items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}

	scrobbleReq.Meta.Count = count
	scrobbleReq.Meta.Total = 50
	scrobbleReq.Meta.Page = 1

	return scrobbleReq, nil
}
