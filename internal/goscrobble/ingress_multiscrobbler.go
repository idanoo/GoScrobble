package goscrobble

import (
	"database/sql"
	"errors"
	"log"
	"net"
	"time"
)

type MultiScrobblerInput struct {
	Artists  []string  `json:"artists"`
	Album    string    `json:"album"`
	Track    string    `json:"track"`
	PlayedAt time.Time `json:"playDate"`
	Duration int       `json:"duration"`
}

// ParseMultiScrobblerInput - Transform API data
func ParseMultiScrobblerInput(userUUID string, data MultiScrobblerInput, ip net.IP, tx *sql.Tx) error {
	// Debugging
	artists := make([]string, 0)
	albumartists := make([]string, 0)

	// Insert track artists
	for _, artist := range data.Artists {
		artist, err := insertArtist(artist, "", "", tx)

		if err != nil {
			log.Printf("%+v", err)
			return errors.New("Failed to map artist: " + artist.Name)
		}
		artists = append(artists, artist.Uuid)
	}

	// Insert album if not exist
	album, err := insertAlbum(data.Album, "", "", albumartists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert track if not exist
	track, err := insertTrack(data.Track, data.Duration, "", "", album.Uuid, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert scrobble if not exist
	err = insertScrobble(userUUID, track.Uuid, "multiscrobbler", data.PlayedAt, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	return nil
}
