package goscrobble

import (
	"fmt"
	"time"
)

var endTicker chan bool

func StartBackgroundWorkers() {

	endTicker := make(chan bool)

	hourTicker := time.NewTicker(time.Duration(1) * time.Hour)
	minuteTicker := time.NewTicker(time.Duration(60) * time.Second)

	go func() {
		for {
			select {
			case <-endTicker:
				fmt.Println("Stopping background workers")
				return
			case <-hourTicker.C:
				// Clear old password reset tokens
				clearOldResetTokens()

				// Attempt to pull missing images from spotify - hackerino version!
				user, _ := getUserByUsername("idanoo")
				user.updateImageDataFromSpotify()
			case <-minuteTicker.C:
				// Update playdata from spotify
				updateSpotifyData()
			}
		}
	}()
}

func EndBackgroundWorkers() {
	endTicker <- true
}
