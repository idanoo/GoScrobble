package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/idanoo/goscrobble/internal/goscrobble"
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

	// Store Refresh expiry
	goscrobble.RefereshExpiry = (86400 * 7)
	refreshExpiryStr := os.Getenv("REFRESH_EXPIRY")
	if refreshExpiryStr != "" {
		i, err := strconv.ParseFloat(refreshExpiryStr, 64)
		if err != nil {
			panic("Invalid REFRESH_EXPIRY value")
		}

		goscrobble.RefereshExpiry = time.Duration(i) * time.Second
	}

	goscrobble.DataDirectory = os.Getenv("DATA_DIRECTORY")
	if goscrobble.DataDirectory == "" {
		panic("DATA_DIRECTORY required in .env")
	}

	goscrobble.FrontendDirectory = os.Getenv("FRONTEND_DIRECTORY")
	if goscrobble.FrontendDirectory == "" {
		panic("FRONTEND_DIRECTORY required in .env")
	}

	goscrobble.ApiDocsDirectory = os.Getenv("API_DOCS_DIRECTORY")
	if goscrobble.ApiDocsDirectory == "" {
		panic("API_DOCS_DIRECTORY required in .env")
	}

	// Ignore reverse proxies
	goscrobble.ReverseProxies = strings.Split(os.Getenv("REVERSE_PROXIES"), ",")

	goscrobble.DevMode = false
	devModeString := os.Getenv("DEV_MODE")
	if strings.ToLower(devModeString) == "true" {
		goscrobble.DevMode = true
	}

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

	// Start background workers if not DevMode
	if !goscrobble.DevMode {
		go goscrobble.StartBackgroundWorkers()
		defer goscrobble.EndBackgroundWorkers()
	} else {
		fmt.Printf("Running in DevMode. No background workers running")
		fmt.Println("")
	}

	// Boot up API webserver \o/
	goscrobble.HandleRequests(port)
}
