package goscrobble

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	// Cache key
	json, _ := json.Marshal(data)
	redisKey := getMd5(string(json) + userUUID)
	if getRedisKeyExists(redisKey) {
		fmt.Printf("Prevented duplicate entry!")
		return nil
	}

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

	// Add cache key for the duration of the song *2 since we're caching the start time too
	ttl := time.Duration(data.Duration*2) * time.Second
	setRedisValTtl(redisKey, "1", ttl)

	return nil
}
