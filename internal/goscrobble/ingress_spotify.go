package goscrobble

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// updateSpotifyData - Pull data for all users
func updateSpotifyData() {
	// Lets ignore if not configured
	val, _ := getConfigValue("SPOTIFY_API_SECRET")
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

// setSpotifyEnvVars - Temp fix for V2 Spotify SDK
func setSpotifyEnvVars() {
	appId, _ := getConfigValue("SPOTIFY_API_ID")
	appSecret, _ := getConfigValue("SPOTIFY_API_SECRET")
	os.Setenv("SPOTIFY_ID", appId)
	os.Setenv("SPOTIFY_SECRET", appSecret)
}

func getSpotifyAuthHandler() spotifyauth.Authenticator {
	setSpotifyEnvVars()
	redirectUrl := os.Getenv("GOSCROBBLE_DOMAIN") + "/api/v1/link/spotify"
	if redirectUrl == "http://localhost:3000/api/v1/link/spotify" {
		// Handle backend on a different port if running in dev-env
		redirectUrl = "http://localhost:42069/api/v1/link/spotify"
	}

	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadRecentlyPlayed, spotifyauth.ScopeUserReadCurrentlyPlaying),
	)

	return *auth
}

func connectSpotifyResponse(r *http.Request) error {
	urlParams := r.URL.Query()
	userUUID := urlParams["state"][0]

	_, err := getUserByUUID(userUUID)
	if err != nil {
		return err
	}

	auth := getSpotifyAuthHandler()
	token, err := auth.Token(r.Context(), userUUID, r)
	if err != nil {
		return err
	}

	// Get displayName
	client := spotify.New(auth.Client(r.Context(), token))
	spotifyUser, err := client.CurrentUser(r.Context())

	// Lets pull in last 30 minutes
	time := time.Now().UTC().Add(-(time.Duration(30) * time.Minute))
	err = insertOauthToken(userUUID, "spotify", token.AccessToken, token.RefreshToken, token.Expiry, spotifyUser.DisplayName, time, "")
	if err != nil {
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
	ctx := context.Background()
	client := spotify.New(auth.Client(ctx, token))

	// Only fetch tracks since last sync
	opts := spotify.RecentlyPlayedOptions{
		AfterEpochMs: dbToken.LastSynced.UnixNano() / int64(time.Millisecond),
	}

	// We want the next sync timestamp from before we call
	// so we don't end up with a few seconds gap
	currTime := time.Now()
	items, err := client.PlayerRecentlyPlayedOpt(ctx, &opts)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Unable to get recently played tracks for user: %+v", user.Username)
		return
	}

	for _, v := range items {
		if !checkIfSpotifyAlreadyScrobbled(user.UUID, v) {
			tx, _ := db.Begin()
			err := ParseSpotifyInput(ctx, user.UUID, v, client, tx)
			if err != nil {
				fmt.Printf("Failed to insert Spotify scrobble: %+v", err)
				tx.Rollback()
				break
			}
			tx.Commit()
		}
	}

	// Check if token has changed.. if so, save it to db
	currToken, err := client.Token()
	err = insertOauthToken(user.UUID, "spotify", currToken.AccessToken, currToken.RefreshToken, currToken.Expiry, dbToken.Username, currTime, "")
	if err != nil {
		fmt.Printf("Failed to update spotify token in database")
	}

}

func checkIfSpotifyAlreadyScrobbled(userUuid string, data spotify.RecentlyPlayedItem) bool {
	return checkIfScrobbleExists(userUuid, data.PlayedAt, "spotify")
}

// ParseSpotifyInput - Transform API data
func ParseSpotifyInput(ctx context.Context, userUUID string, data spotify.RecentlyPlayedItem, client *spotify.Client, tx *sql.Tx) error {
	artists := make([]string, 0)
	albumartists := make([]string, 0)

	// Insert track artists
	for _, artist := range data.Track.Artists {
		fullArtist, err := client.GetArtist(ctx, artist.ID)
		img := ""
		if len(fullArtist.Images) > 0 {
			img = fullArtist.Images[0].URL
		}

		artist, err := insertArtist(artist.Name, "", artist.ID.String(), img, tx)

		if err != nil {
			log.Printf("%+v", err)
			return errors.New("Failed to map artist: " + artist.Name)
		}
		artists = append(artists, artist.UUID)
	}

	// Get full track data (album / track info)
	fulltrack, err := client.GetTrack(ctx, data.Track.ID)
	if err != nil {
		fmt.Printf("Failed to get full track info from spotify: %+v", data.Track.Name)
		return errors.New("Failed to get full track info from spotify: " + data.Track.Name)
	}

	// Insert album artists
	for _, artist := range fulltrack.Album.Artists {
		fullArtist, err := client.GetArtist(ctx, artist.ID)
		img := ""
		if len(fullArtist.Images) > 0 {
			img = fullArtist.Images[0].URL
		}

		albumartist, err := insertArtist(artist.Name, "", artist.ID.String(), img, tx)

		if err != nil {
			log.Printf("%+v", err)
			return errors.New("Failed to map album: " + artist.Name)
		}
		albumartists = append(albumartists, albumartist.UUID)
	}

	// Insert album if not exist
	albumImage := ""
	if len(fulltrack.Album.Images) > 0 {
		albumImage = fulltrack.Album.Images[0].URL
	}
	album, err := insertAlbum(fulltrack.Album.Name, "", fulltrack.Album.ID.String(), albumartists, albumImage, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map album")
	}

	// Insert track if not exist
	length := int(fulltrack.Duration / 1000)
	track, err := insertTrack(fulltrack.Name, length, "", fulltrack.ID.String(), album.UUID, artists, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	// Insert scrobble if not exist
	ip := net.ParseIP("0.0.0.0")
	err = insertScrobble(userUUID, track.UUID, "spotify", data.PlayedAt, ip, tx)
	if err != nil {
		log.Printf("%+v", err)
		return errors.New("Failed to map track")
	}

	return nil
}

// updateImageDataFromSpotify update artist/album images from spotify ;D
func (user *User) updateImageDataFromSpotify() {
	// Check that data is set before we attempt to pull
	val, _ := getConfigValue("SPOTIFY_API_SECRET")
	if val == "" {
		return
	}

	// TO BE REWORKED TO NOT USE A DAMN USER ARGHHH
	dbToken, err := user.getSpotifyTokens()
	if err != nil {
		return
	}

	token := new(oauth2.Token)
	token.AccessToken = dbToken.AccessToken
	token.RefreshToken = dbToken.RefreshToken
	token.Expiry = dbToken.Expiry
	token.TokenType = "Bearer"

	auth := getSpotifyAuthHandler()
	client := spotify.New(auth.Client(ctx, token))

	rows, err := db.Query(`SELECT uuid, name FROM artists WHERE COALESCE(img,'') NOT IN ('pending', 'complete') LIMIT 100`)
	if err != nil {
		log.Printf("Failed to fetch config: %+v", err)
		return
	}

	toUpdate := make(map[string]string)
	for rows.Next() {
		var uuid string
		var name string
		err := rows.Scan(&uuid, &name)
		if err != nil {
			log.Printf("Failed to fetch artists: %+v", err)
			rows.Close()
			return
		}
		res, err := client.Search(ctx, name, spotify.SearchTypeArtist)
		if len(res.Artists.Artists) > 0 {
			if len(res.Artists.Artists[0].Images) > 0 {
				toUpdate[uuid] = res.Artists.Artists[0].Images[0].URL
			}
		} else {
			_, _ = db.Exec("UPDATE `artists` SET `img` = 'none' WHERE `uuid` = UUID_TO_BIN(?,true)", uuid)
		}
	}
	rows.Close()

	var artist Artist
	tx, _ := db.Begin()
	for uuid, img := range toUpdate {
		artist, err = getArtistByUUID(uuid)
		if err != nil {
			continue
		}

		err = importImage(uuid, img)
		if err != nil {
			fmt.Printf("Failed to import image: %+v", err)
			continue
		}

		_ = artist.updateArtist("img", "pending", tx)
	}
	tx.Commit()

	rows, err = db.Query("SELECT uuid, name FROM albums WHERE COALESCE(img,'') NOT IN ('pending', 'complete') LIMIT 100")
	if err != nil {
		log.Printf("Failed to fetch config: %+v", err)
		return
	}

	toUpdate = make(map[string]string)
	for rows.Next() {
		var uuid string
		var name string
		err := rows.Scan(&uuid, &name)
		if err != nil {
			log.Printf("Failed to fetch albums: %+v", err)
			rows.Close()
			return
		}
		res, err := client.Search(ctx, name, spotify.SearchTypeAlbum)
		if len(res.Albums.Albums) > 0 {
			if len(res.Albums.Albums[0].Images) > 0 {
				toUpdate[uuid] = res.Albums.Albums[0].Images[0].URL
			}
		} else {
			_, _ = db.Exec("UPDATE `albums` SET `img` = 'none' WHERE `uuid` = UUID_TO_BIN(?,true)", uuid)
		}
	}
	rows.Close()

	var album Album
	tx, _ = db.Begin()
	for uuid, img := range toUpdate {
		album, err = getAlbumByUUID(uuid)
		err = importImage(uuid, img)
		if err != nil {
			continue
		}

		// Update
		_ = album.updateAlbum("img", "pending", tx)
	}

	tx.Commit()

	return
}
