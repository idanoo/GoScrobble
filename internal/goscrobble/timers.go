package goscrobble

import (
	"fmt"
	"time"
)

var endTicker chan bool

func StartBackgroundWorkers() {
	updateSpotifyData()

	endTicker := make(chan bool)

	hourTicker := time.NewTicker(time.Hour)
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
