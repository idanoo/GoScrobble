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

// insertScrobble - This will return if it exists or create it based on MBID > Name
func insertScrobble(user string, track string, ip net.IP, tx *sql.Tx) error {
	err := insertNewScrobble(user, track, ip, tx)
	if err != nil {
		log.Printf("Error inserting scrobble %s  %+v", user, err)
		return errors.New("Failed to insert scrobble!")
	}

	return nil
}

func fetchScrobble(col string, val string, tx *sql.Tx) Scrobble {
	var scrobble Scrobble
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `created_at`, `created_ip`, `user`, `track` FROM `scrobbles` WHERE `"+col+"` = ?",
		val).Scan(&scrobble.Uuid, &scrobble.CreatedAt, &scrobble.CreatedIp, &scrobble.User, &scrobble.Track)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching scrobbles: %+v", err)
		}
	}

	return scrobble
}

func insertNewScrobble(user string, track string, ip net.IP, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `scrobbles` (`uuid`, `created_at`, `created_ip`, `user`, `track`) "+
		"VALUES (UUID_TO_BIN(UUID(), true), NOW(), ?, UUID_TO_BIN(?, true),UUID_TO_BIN(?, true))", ip, user, track)

	return err
}
