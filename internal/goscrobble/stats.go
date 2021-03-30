package goscrobble

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type StatsRequest struct {
	Users       int       `json:"users"`
	Scrobbles   int       `json:"scrobbles"`
	Tracks      int       `json:"tracks"`
	Artists     int       `json:"artists"`
	LastUpdated time.Time `json:"last_updated"`
}

func getStats() (StatsRequest, error) {
	js := getRedisVal("stats")
	statsReq := StatsRequest{}
	var err error
	if js != "" {
		// If cached, deserialize and return
		err = json.Unmarshal([]byte(js), &statsReq)
		if err != nil {
			log.Printf("Error unmarshalling stats json: %+v", err)
			return statsReq, errors.New("Error fetching stats")
		}

		// Check if older than 5 min - we want to update async for the next caller
		now := time.Now()
		if now.Sub(statsReq.LastUpdated) > time.Duration(5)*time.Minute {
			go goFetchStats()
		}
	} else {
		// If not cached, pull data then return
		statsReq, err = fetchStats()
		if err != nil {
			log.Printf("Error fetching stats: %+v", err)
			return statsReq, errors.New("Error fetching stats")
		}
	}

	return statsReq, nil
}

// goFetchStats - Async call
func goFetchStats() {
	_, _ = fetchStats()
}

func fetchStats() (StatsRequest, error) {
	statsReq := StatsRequest{}
	var err error

	statsReq.Users, err = getDbCount("SELECT COUNT(*) FROM `users` WHERE `active` = 1")
	if err != nil {
		log.Printf("Failed to fetch user count: %+v", err)
		return statsReq, errors.New("Failed to fetch stats")
	}

	statsReq.Scrobbles, err = getDbCount("SELECT COUNT(*) FROM `scrobbles`")
	if err != nil {
		log.Printf("Failed to fetch scrobble count: %+v", err)
		return statsReq, errors.New("Failed to fetch stats")
	}

	statsReq.Tracks, err = getDbCount("SELECT COUNT(*) FROM `tracks`")
	if err != nil {
		log.Printf("Failed to fetch track count: %+v", err)
		return statsReq, errors.New("Failed to fetch stats")
	}

	statsReq.Artists, err = getDbCount("SELECT COUNT(*) FROM `artists`")
	if err != nil {
		log.Printf("Failed to fetch artist count: %+v", err)
		return statsReq, errors.New("Failed to fetch stats")
	}

	// Mark the time this was last updated
	statsReq.LastUpdated = time.Now()

	b, err := json.Marshal(statsReq)
	if err != nil {
		return statsReq, errors.New("Failed to fetch stats")
	}

	err = setRedisVal("stats", string(b))
	if err != nil {
		return statsReq, errors.New("Failed to fetch stats")
	}

	return statsReq, nil
}
