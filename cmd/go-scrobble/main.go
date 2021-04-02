package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/idanoo/go-scrobble/internal/goscrobble"
)

func init() {
	// Set UTC for errything
	os.Setenv("TZ", "Etc/UTC")
}

func main() {
	fmt.Println("Starting goscrobble")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Store JWT secret
	goscrobble.JwtToken = []byte(os.Getenv("JWT_SECRET"))

	// Store JWT expiry
	goscrobble.JwtExpiry = 86400
	jwtExpiryStr := os.Getenv("JWT_EXPIRY")
	if jwtExpiryStr != "" {
		i, err := strconv.ParseFloat(jwtExpiryStr, 64)
		if err != nil {
			panic("Invalid JWT_EXPIRY value")
		}

		goscrobble.JwtExpiry = time.Duration(i) * time.Second
	}

	// Ignore reverse proxies
	goscrobble.ReverseProxies = strings.Split(os.Getenv("REVERSE_PROXIES"), ",")

	// Store port
	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}

	// Boot up DB connection
	goscrobble.InitDb()
	defer goscrobble.CloseDbConn()

	// Boot up Redis connection
	goscrobble.InitRedis()
	defer goscrobble.CloseRedisConn()

	// Start background workers
	go goscrobble.StartBackgroundWorkers()
	defer goscrobble.EndBackgroundWorkers()

	// Boot up API webserver \o/
	goscrobble.HandleRequests(port)
}
