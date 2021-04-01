package goscrobble

import (
	"database/sql"
	"errors"
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

type ScrobbleRequest struct {
	Meta  ScrobbleRequestMeta   `json:"meta"`
	Items []ScrobbleRequestItem `json:"items"`
}

type ScrobbleRequestMeta struct {
	Count int `json:"count"`
	Total int `json:"total"`
	Page  int `json:"page"`
}

type ScrobbleRequestItem struct {
	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"time"`
	Artist    string    `json:"artist"`
	Album     string    `json:"album"`
	Track     string    `json:"track"`
}

// insertScrobble - This will return if it exists or create it based on MBID > Name
func insertScrobble(user string, track string, source string, ip net.IP, tx *sql.Tx) error {
	err := insertNewScrobble(user, track, source, ip, tx)
	if err != nil {
		log.Printf("Error inserting scrobble %s  %+v", user, err)
		return errors.New("Failed to insert scrobble!")
	}

	return nil
}

func fetchScrobblesForUser(userUuid string, limit int, page int) (ScrobbleRequest, error) {
	scrobbleReq := ScrobbleRequest{}
	var count int

	// Yeah this isn't great. But for now.. it works! Cache later
	total, err := getDbCount(
		"SELECT COUNT(*) FROM `scrobbles` "+
			"JOIN tracks ON scrobbles.track = tracks.uuid "+
			"JOIN track_artist ON track_artist.track = tracks.uuid "+
			"JOIN track_album ON track_album.track = tracks.uuid "+
			"JOIN artists ON track_artist.artist = artists.uuid "+
			"JOIN albums ON track_album.album = albums.uuid "+
			"JOIN users ON scrobbles.user = users.uuid "+
			"WHERE user = UUID_TO_BIN(?, true)",
		userUuid)

	if err != nil {
		log.Printf("Failed to fetch scrobble count: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}

	rows, err := db.Query(
		"SELECT BIN_TO_UUID(`scrobbles`.`uuid`, true), `scrobbles`.`created_at`, `artists`.`name`, `albums`.`name`,`tracks`.`name` FROM `scrobbles` "+
			"JOIN tracks ON scrobbles.track = tracks.uuid "+
			"JOIN track_artist ON track_artist.track = tracks.uuid "+
			"JOIN track_album ON track_album.track = tracks.uuid "+
			"JOIN artists ON track_artist.artist = artists.uuid "+
			"JOIN albums ON track_album.album = albums.uuid "+
			"JOIN users ON scrobbles.user = users.uuid "+
			"WHERE user = UUID_TO_BIN(?, true) "+
			"ORDER BY scrobbles.created_at DESC LIMIT ?",
		userUuid, limit)
	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return scrobbleReq, errors.New("Failed to fetch scrobbles")
	}
	defer rows.Close()

	for rows.Next() {
		item := ScrobbleRequestItem{}
		err := rows.Scan(&item.UUID, &item.Timestamp, &item.Artist, &item.Album, &item.Track)
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

func insertNewScrobble(user string, track string, source string, ip net.IP, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `scrobbles` (`uuid`, `created_at`, `created_ip`, `user`, `track`, `source`) "+
		"VALUES (UUID_TO_BIN(UUID(), true), NOW(), ?, UUID_TO_BIN(?, true),UUID_TO_BIN(?, true), ?)", ip, user, track, source)

	return err
}
