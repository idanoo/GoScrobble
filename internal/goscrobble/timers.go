package goscrobble

import (
	"fmt"
	"time"
)

var endTicker chan bool

func StartBackgroundWorkers() {
	endTicker := make(chan bool)

	minuteTicker := time.NewTicker(time.Duration(60) * time.Second)
	hourTicker := time.NewTicker(time.Duration(1) * time.Hour)

	go func() {
		for {
			select {
			case <-endTicker:
				fmt.Println("Stopping background workers")
				return
			case <-minuteTicker.C:
				// Update playdata from Spotify
				go updateSpotifyData()

				// Update playdate from Navidrome
				go updateNavidromeData()

				// There should only be whatever new images are in pending
				go resizeImages()
			case <-hourTicker.C:
				// Clear old password reset tokens
				go clearOldResetTokens()

				// Attempt to pull missing images from spotify - hackerino version!
				user, _ := getUserByUsername("idanoo")
				go user.updateImageDataFromSpotify()
			}
		}
	}()
}

func EndBackgroundWorkers() {
	endTicker <- true
}
