package goscrobble

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type NavidromeResponse struct {
	Response struct {
		Status        string `json:"status"`
		Version       string `json:"version"`
		Type          string `json:"type"`
		ServerVersion string `json:"serverVersion"`
		NowPlaying    struct {
			Entry []NavidromeNowPlaying `json:"entry"`
		} `json:"nowPlaying"`
	} `json:"subsonic-response"`
}

type NavidromeNowPlaying struct {
	ID          string    `json:"id"`
	Parent      string    `json:"parent"`
	IsDir       bool      `json:"isDir"`
	Title       string    `json:"title"`
	Album       string    `json:"album"`
	Artist      string    `json:"artist"`
	Track       int       `json:"track"`
	Year        int       `json:"year"`
	Genre       string    `json:"genre"`
	CoverArt    string    `json:"coverArt"`
	Size        int       `json:"size"`
	ContentType string    `json:"contentType"`
	Suffix      string    `json:"suffix"`
	Duration    int       `json:"duration"`
	BitRate     int       `json:"bitRate"`
	Path        string    `json:"path"`
	DiscNumber  int       `json:"discNumber"`
	Created     time.Time `json:"created"`
	AlbumID     string    `json:"albumId"`
	ArtistID    string    `json:"artistId"`
	Type        string    `json:"type"`
	IsVideo     bool      `json:"isVideo"`
	Username    string    `json:"username"`
	PlayerID    int       `json:"playerId"`
	PlayerName  string    `json:"playerName"`
}

// updateSpotifyData - Pull data for all users
func updateNavidromeData() {
	// Get all active users with a spotify token
	users, err := getAllNavidromeUsers()
	if err != nil {
		fmt.Printf("Failed to fetch navidrome users")
		return
	}

	for _, user := range users {
		tx, err := db.Begin()
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = user.updateNavidromePlaydata(tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			continue
		}
		tx.Commit()
	}
}

func (user *User) updateNavidromePlaydata(tx *sql.Tx) error {
	tokens, err := user.getNavidromeTokens()
	if err != nil {
		fmt.Printf("No Navidrome token for user: %+v %+v", user.Username, err)
		return errors.New("Failed to fetch Navidrome Tokens")
	}

	response, err := getNavidromeNowPlaying(&tokens)
	fmt.Println(response)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to fetch Navidrome Tokens %+v", err))
	}

	ip := net.ParseIP("0.0.0.0")
	for _, v := range response.Response.NowPlaying.Entry {
		err = ParseNavidromeInput(user.UUID, v, ip, tx)
	}

	if err != nil {
		return err
	}

	return nil
}

func getNavidromeNowPlaying(token *OauthToken) (NavidromeResponse, error) {
	response := NavidromeResponse{}
	resp, err := http.Get(token.URL + "/rest/getNowPlaying?u=" + token.Username + "&t=" + token.AccessToken + "&s=" + token.RefreshToken + "&c=GoScrobble&v=1.16.1&f=json")
	if err != nil {
		return response, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func validateNavidromeConnection(url string, username string, hash string, salt string) error {
	resp, err := http.Get(url + "/rest/ping.view?u=" + username + "&t=" + hash + "&s=" + salt + "&c=GoScrobble&v=1.16.1&f=json")
	if err != nil {
		return err
	}

	response := NavidromeResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	if response.Response.Status == "ok" {
		return nil
	}

	return errors.New("Failed to validate")
}

// ParseNavidromeInput - Transform API data
func ParseNavidromeInput(userUUID string, data NavidromeNowPlaying, ip net.IP, tx *sql.Tx) error {
	// Cache key
	json := fmt.Sprintf("%s:%s:%s", data.ID, data.Parent, userUUID)
	redisKey := getMd5(json)
	if getRedisKeyExists(redisKey) {
		return nil
	}

	artists := make([]string, 0)

	// Insert track artists
	artist, err := insertArtist(data.Artist, "", "", "", tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map artist: " + artist.Name)
	}
	artists = append(artists, artist.UUID)

	// Insert album if not exist
	album, err := insertAlbum(data.Album, "", "", artists, "", tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert track if not exist
	track, err := insertTrack(data.Title, data.Duration, "", "", album.UUID, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert scrobble if not exist
	timeNow := time.Now()
	err = insertScrobble(userUUID, track.UUID, "navidrome", timeNow, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Todo: Find a better way to check dupes
	ttl := time.Duration(data.Duration*2) * time.Second
	setRedisValTtl(redisKey, "1", ttl)

	return nil
}
