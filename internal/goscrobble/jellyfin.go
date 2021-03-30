package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
)

// ParseJellyfinInput - Transform API data into a common struct
func ParseJellyfinInput(userUUID string, data map[string]interface{}, ip net.IP, tx *sql.Tx) error {
	if data["ItemType"] != "Audio" {
		return errors.New("Media type not audio")
	}

	// Safety Checks
	if data["Artist"] == nil {
		return errors.New("Missing artist data")
	}

	if data["Album"] == nil {
		return errors.New("Missing album data")
	}

	if data["Name"] == nil {
		return errors.New("Missing track data")
	}

	// Insert artist if not exist
	artist, err := insertArtist(fmt.Sprintf("%s", data["Artist"]), fmt.Sprintf("%s", data["Provider_musicbrainzartist"]), tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map artist")
	}

	// Insert album if not exist
	artists := []string{artist.Uuid}
	album, err := insertAlbum(fmt.Sprintf("%s", data["Album"]), fmt.Sprintf("%s", data["Provider_musicbrainzalbum"]), artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert album if not exist
	track, err := insertTrack(fmt.Sprintf("%s", data["Name"]), fmt.Sprintf("%s", data["Provider_musicbrainztrack"]), album.Uuid, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert album if not exist
	err = insertScrobble(userUUID, track.Uuid, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	_ = album
	_ = artist
	_ = track

	// Insert track if not exist
	return nil
}
