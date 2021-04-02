package goscrobble

import (
	"database/sql"
	"fmt"
	"net"
	"time"
)

type MultiScrobblerInput struct {
	Artists  []string  `json:"artists"`
	Album    string    `json:"album"`
	Track    string    `json:"track"`
	PlayedAt time.Time `json:"playDate"`
	Duration string    `json:"duration"`
}

// ParseMultiScrobblerInput - Transform API data
func ParseMultiScrobblerInput(userUUID string, data map[string]interface{}, ip net.IP, tx *sql.Tx) error {
	// Debugging
	fmt.Printf("%+v", data)

	// // Safety Checks
	// if data["artists"] == nil {
	// 	return errors.New("Missing artist data")
	// }

	// if data["album"] == nil {
	// 	return errors.New("Missing album data")
	// }

	// if data["track"] == nil {
	// 	return errors.New("Missing track data")
	// }

	// // Insert track artists
	// for _, artist := range data["artists"] {
	// 	artist, err := insertArtist(artist.Name, "", artist.ID.String(), tx)

	// 	if err != nil {
	// 		log.Printf("%+v", err)
	// 		return errors.New("Failed to map artist: " + artist.Name)
	// 	}
	// 	artists = append(artists, artist.Uuid)
	// }
	// // Insert album if not exist
	// artists := []string{artist.Uuid}
	// album, err := insertAlbum(fmt.Sprintf("%s", data["Album"]), fmt.Sprintf("%s", data["Provider_musicbrainzalbum"]), "", artists, tx)
	// if err != nil {
	// 	log.Printf("%+v", err)
	// 	return errors.New("Failed to map album")
	// }

	// // Insert track if not exist
	// length := timestampToSeconds(fmt.Sprintf("%s", data["RunTime"]))
	// track, err := insertTrack(fmt.Sprintf("%s", data["Name"]), length, fmt.Sprintf("%s", data["Provider_musicbrainztrack"]), "", album.Uuid, artists, tx)
	// if err != nil {
	// 	log.Printf("%+v", err)
	// 	return errors.New("Failed to map track")
	// }

	// // Insert scrobble if not exist
	// timestamp := time.Now()
	// fmt.Println(timestamp)
	// err = insertScrobble(userUUID, track.Uuid, "jellyfin", timestamp, ip, tx)
	// if err != nil {
	// 	log.Printf("%+v", err)
	// 	return errors.New("Failed to map track")
	// }

	// Insert track if not exist
	return nil
}
