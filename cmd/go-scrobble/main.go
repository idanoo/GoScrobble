package main

import (
	"log"
	"os"

	"git.m2.nz/go-scrobble/internal/goscrobble"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Store JWT secret
	goscrobble.JwtToken = []byte(os.Getenv("JWT_SECRET"))

	// // Boot up DB connection for life of application
	goscrobble.InitDb()
	defer goscrobble.CloseDbConn()

	// Boot up API webserver \o/
	goscrobble.HandleRequests()
}
