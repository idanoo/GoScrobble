package main

import (
	"git.m2.nz/go-scrobble/internal/goscrobble"
)

func main() {
	// // Boot up DB connection for life of application
	// goscrobble.InitDb()
	// defer goscrobble.CloseDbConn()

	// Boot up API webserver \o/
	goscrobble.HandleRequests()
}
