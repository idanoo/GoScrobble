package goscrobble

import (
	"database/sql"
	"errors"
	"log"
	"net"
	"time"
)

type JellyfinRequest struct {
	Album                           string `json:"Album"`
	Artist                          string `json:"Artist"`
	ClientName                      string `json:"ClientName"`
	DeviceID                        string `json:"DeviceId"`
	DeviceName                      string `json:"DeviceName"`
	IsAutomated                     bool   `json:"IsAutomated"`
	IsPaused                        bool   `json:"IsPaused"`
	ItemID                          string `json:"ItemId"`
	ItemType                        string `json:"ItemType"`
	MediaSourceID                   string `json:"MediaSourceId"`
	Name                            string `json:"Name"`
	NotificationType                string `json:"NotificationType"`
	Overview                        string `json:"Overview"`
	PlaybackPosition                string `json:"PlaybackPosition"`
	PlaybackPositionTicks           int64  `json:"PlaybackPositionTicks"`
	ProviderMusicbrainzalbum        string `json:"Provider_musicbrainzalbum"`
	ProviderMusicbrainzalbumartist  string `json:"Provider_musicbrainzalbumartist"`
	ProviderMusicbrainzartist       string `json:"Provider_musicbrainzartist"`
	ProviderMusicbrainzreleasegroup string `json:"Provider_musicbrainzreleasegroup"`
	ProviderMusicbrainztrack        string `json:"Provider_musicbrainztrack"`
	RunTime                         string `json:"RunTime"`
	RunTimeTicks                    int64  `json:"RunTimeTicks"`
	ServerID                        string `json:"ServerId"`
	ServerName                      string `json:"ServerName"`
	ServerURL                       string `json:"ServerUrl"`
	ServerVersion                   string `json:"ServerVersion"`
	Timestamp                       string `json:"Timestamp"`
	UserID                          string `json:"UserId"`
	Username                        string `json:"Username"`
	UtcTimestamp                    string `json:"UtcTimestamp"`
	Year                            int64  `json:"Year"`
}

// ParseJellyfinInput - Transform API data into a common struct. Uses MusicBrainzID primarily
func ParseJellyfinInput(userUUID string, jf JellyfinRequest, ip net.IP, tx *sql.Tx) error {
	// Debugging
	// fmt.Printf("%+v", jf)

	// Prevents scrobbling same song twice!
	cacheKey := jf.UserID + ":" + jf.Name + ":" + jf.Artist + ":" + jf.Album + ":" + jf.ServerID
	redisKey := getMd5(cacheKey + userUUID)
	if getRedisKeyExists(redisKey) {
		return nil
	}

	if jf.ItemType != "Audio" {
		return errors.New("Media type not audio")
	}

	// Safety Checks
	if jf.Artist == "" {
		return errors.New("Missing artist data")
	}

	if jf.Album == "" {
		return errors.New("Missing album data")
	}

	if jf.Name == "" {
		return errors.New("Missing track data")
	}

	// Insert artist if not exist
	artist, err := insertArtist(jf.Artist, jf.ProviderMusicbrainzartist, "", tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map artist")
	}

	// Insert album if not exist
	artists := []string{artist.UUID}
	album, err := insertAlbum(jf.Album, jf.ProviderMusicbrainzalbum, "", artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert track if not exist
	length := timestampToSeconds(jf.RunTime)
	track, err := insertTrack(jf.Name, length, jf.ProviderMusicbrainztrack, "", album.UUID, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert scrobble if not exist
	timestamp := time.Now()
	err = insertScrobble(userUUID, track.UUID, "jellyfin", timestamp, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Add cache key!
	ttl := time.Duration(timestampToSeconds(jf.RunTime)) * time.Second
	setRedisValTtl(redisKey, "1", ttl)

	return nil
}
