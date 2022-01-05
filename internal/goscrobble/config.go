package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Config struct {
	Setting map[string]string `json:"configs"`
}

type ServerInfo struct {
	Version             string `json:"version"`
	RegistrationEnabled string `json:"registration_enabled"`
}

func getAllConfigs() (Config, error) {
	config := Config{}
	configs := make(map[string]string)

	rows, err := db.Query(`SELECT key, value FROM config`)
	if err != nil {
		log.Printf("Failed to fetch config: %+v", err)
		return config, errors.New("Failed to fetch configs")
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value string
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Printf("Failed to fetch config: %+v", err)
			return config, errors.New("Failed to fetch configs")
		}

		// Append
		configs[key] = value
	}

	// Assign the data to the parent
	config.Setting = configs

	err = rows.Err()
	if err != nil {
		log.Printf("Failed to fetch config: %+v", err)
		return config, errors.New("Failed to fetch configs")
	}

	return config, nil
}

func updateConfigValue(key string, value string) error {
	_, err := db.Exec(`UPDATE config SET value = $1 WHERE key = $2`, value, key)
	if err != nil {
		fmt.Printf("Failed to update config: %+v", err)
		return errors.New("Failed to update config value.")
	}

	return nil
}

func getConfigValue(key string) (string, error) {
	var value string

	err := db.QueryRow(`SELECT value FROM config WHERE key = $1`,
		key).Scan(&value)

	if err == sql.ErrNoRows {
		return value, errors.New("Config key doesn't exist")
	}

	return value, nil
}
