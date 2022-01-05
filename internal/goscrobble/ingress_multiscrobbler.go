package goscrobble

import (
	"database/sql"
	"errors"
	"log"
	"net"
	"time"
)

type MultiScrobblerRequest struct {
	Artists  []string  `json:"artists"`
	Album    string    `json:"album"`
	Track    string    `json:"track"`
	PlayedAt time.Time `json:"playDate"`
	Duration int       `json:"duration"`
}

// ParseMultiScrobblerInput - Transform API data
func ParseMultiScrobblerInput(userUUID string, data MultiScrobblerRequest, ip net.IP, tx *sql.Tx) error {
	// Custom cache key - never log the same song twice in a row for now... (:
	lastPlayedTitle := getUserLastPlayed(userUUID)
	if lastPlayedTitle == data.Track+":"+data.Album {
		// If it matches last played song, skip it
		return nil
	}

	artists := make([]string, 0)
	albumartists := make([]string, 0)

	// Insert track artists
	for _, artist := range data.Artists {
		artist, err := insertArtist(artist, "", "", "", tx)

		if err != nil {
			log.Printf("%+v", err)
			return errors.New("Failed to map artist: " + artist.Name)
		}
		artists = append(artists, artist.UUID)
	}

	// Insert album if not exist
	album, err := insertAlbum(data.Album, "", "", albumartists, "", tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert track if not exist
	track, err := insertTrack(data.Track, data.Duration, "", "", album.UUID, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert scrobble if not exist
	err = insertScrobble(userUUID, track.UUID, "multiscrobbler", data.PlayedAt, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	setUserLastPlayed(userUUID, data.Track+":"+data.Album)

	return nil
}
