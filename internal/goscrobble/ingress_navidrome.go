package goscrobble

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type NavidromeResponse struct {
	Response struct {
		Status        string `json:"status"`
		Version       string `json:"version"`
		Type          string `json:"type"`
		Serverversion string `json:"serverVersion"`
	} `json:"subsonic-response"`
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
		user.updateNavidromePlaydata()
	}
}

func (user *User) updateNavidromePlaydata() {
	_, err := user.getNavidromeTokens()
	if err != nil {
		fmt.Printf("No Navidrome token for user: %+v %+v", user.Username, err)
		return
	}
}

func validateNavidromeConnection(url string, username string, hash string, salt string) error {
	fmt.Printf("url:%s, username:%s, hash:%s, salt:%s", url, username, hash, salt)
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
