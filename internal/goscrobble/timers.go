package goscrobble

import (
	"fmt"
	"time"
)

func ClearTokenTimer() {
	go func() {
		for now := range time.Tick(time.Second) {
			fmt.Println(now)
		}
	}()
}
