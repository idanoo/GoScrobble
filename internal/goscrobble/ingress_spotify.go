package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// updateSpotifyData - Pull data for all users
func updateSpotifyData() {
	// Lets ignore if not configured
	val, _ := getConfigValue("SPOTIFY_APP_SECRET")
	if val == "" {
		return
	}

	// Get all active users with a spotify token
	users, err := getAllSpotifyUsers()
	if err != nil {
		fmt.Printf("Failed to fetch spotify users")
		return
	}

	for _, user := range users {
		user.updateSpotifyPlaydata()
	}
}

func getSpotifyAuthHandler() spotify.Authenticator {
	appId, _ := getConfigValue("SPOTIFY_APP_ID")
	appSecret, _ := getConfigValue("SPOTIFY_APP_SECRET")

	redirectUrl := os.Getenv("GOSCROBBLE_DOMAIN") + "/api/v1/link/spotify"
	if redirectUrl == "http://localhost:3000/api/v1/link/spotify" {
		// Handle backend on a different port
		redirectUrl = "http://localhost:42069/api/v1/link/spotify"
	}

	auth := spotify.NewAuthenticator(redirectUrl,
		spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadCurrentlyPlaying)

	auth.SetAuthInfo(appId, appSecret)

	return auth
}

func connectSpotifyResponse(r *http.Request) error {
	urlParams := r.URL.Query()
	userUuid := urlParams["state"][0]

	// TODO: Add validation user exists here

	auth := getSpotifyAuthHandler()
	token, err := auth.Token(userUuid, r)
	if err != nil {
		fmt.Printf("%+v", err)
		return err
	}

	// Get displayName
	client := auth.NewClient(token)
	client.AutoRetry = true
	spotifyUser, err := client.CurrentUser()

	// Lets pull in last 30 minutes
	time := time.Now().UTC().Add(-(time.Duration(30) * time.Minute))
	err = insertOauthToken(userUuid, "spotify", token.AccessToken, token.RefreshToken, token.Expiry, spotifyUser.DisplayName, time)
	if err != nil {
		fmt.Printf("%+v", err)
		return err
	}

	return nil
}

func (user *User) updateSpotifyPlaydata() {
	dbToken, err := user.getSpotifyTokens()
	if err != nil {
		fmt.Printf("No spotify token for user: %+v %+v", user.Username, err)
		return
	}

	token := new(oauth2.Token)
	token.AccessToken = dbToken.AccessToken
	token.RefreshToken = dbToken.RefreshToken
	token.Expiry = dbToken.Expiry
	token.TokenType = "Bearer"

	auth := getSpotifyAuthHandler()
	client := auth.NewClient(token)
	client.AutoRetry = true

	// Only fetch tracks since last sync
	opts := spotify.RecentlyPlayedOptions{
		AfterEpochMs: dbToken.LastSynced.UnixNano() / int64(time.Millisecond),
	}

	// We want the next sync timestamp from before we call
	// so we don't end up with a few seconds gap
	currTime := time.Now()
	items, err := client.PlayerRecentlyPlayedOpt(&opts)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Unable to get recently played tracks for user: %+v", user.Username)
		return
	}

	for _, v := range items {
		if !checkIfSpotifyAlreadyScrobbled(user.UUID, v) {
			tx, _ := db.Begin()
			err := ParseSpotifyInput(user.UUID, v, client, tx)
			if err != nil {
				fmt.Printf("Failed to insert Spotify scrobble: %+v", err)
				tx.Rollback()
				break
			}
			tx.Commit()
			fmt.Printf("Updated spotify track: %+v", v.Track.Name)
		}
	}

	// Check if token has changed.. if so, save it to db
	currToken, err := client.Token()
	err = insertOauthToken(user.UUID, "spotify", currToken.AccessToken, currToken.RefreshToken, currToken.Expiry, dbToken.Username, currTime)
	if err != nil {
		fmt.Printf("Failed to update spotify token in database")
	}

}

func checkIfSpotifyAlreadyScrobbled(userUuid string, data spotify.RecentlyPlayedItem) bool {
	return checkIfScrobbleExists(userUuid, data.PlayedAt, "spotify")
}

// ParseSpotifyInput - Transform API data
func ParseSpotifyInput(userUUID string, data spotify.RecentlyPlayedItem, client spotify.Client, tx *sql.Tx) error {
	artists := make([]string, 0)
	albumartists := make([]string, 0)

	// Insert track artists
	for _, artist := range data.Track.Artists {
		artist, err := insertArtist(artist.Name, "", artist.ID.String(), tx)

		if err != nil {
			log.Printf("%+v", err)
			return errors.New("Failed to map artist: " + artist.Name)
		}
		artists = append(artists, artist.Uuid)
	}

	// Get full track data (album / track info)
	fulltrack, err := client.GetTrack(data.Track.ID)
	if err != nil {
		fmt.Printf("Failed to get full track info from spotify: %+v", data.Track.Name)
		return errors.New("Failed to get full track info from spotify: " + data.Track.Name)
	}

	// Insert album artists
	for _, artist := range fulltrack.Album.Artists {
		albumartist, err := insertArtist(artist.Name, "", artist.ID.String(), tx)

		if err != nil {
			log.Printf("%+v", err)
			return errors.New("Failed to map album: " + artist.Name)
		}
		albumartists = append(albumartists, albumartist.Uuid)
	}

	// Insert album if not exist
	album, err := insertAlbum(fulltrack.Album.Name, "", fulltrack.Album.ID.String(), albumartists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert track if not exist
	length := int(fulltrack.Duration / 60)
	track, err := insertTrack(fulltrack.Name, length, "", fulltrack.ID.String(), album.Uuid, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert scrobble if not exist
	ip := net.ParseIP("0.0.0.0")
	err = insertScrobble(userUUID, track.Uuid, "spotify", data.PlayedAt, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	return nil
}
