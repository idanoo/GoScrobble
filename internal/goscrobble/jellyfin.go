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
	log.Printf("%+v : %+v", userUUID, data)

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
