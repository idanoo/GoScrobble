package goscrobble

import (
	"errors"
	"fmt"
	"log"
)

type Config struct {
	Setting map[string]string `json:"configs"`
}

func getAllConfigs() (Config, error) {
	config := Config{}
	configs := make(map[string]string)

	rows, err := db.Query("SELECT `key`, `value` FROM `config`")
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
	_, err := db.Exec("UPDATE `config` SET `value` = ? WHERE `key` = ?", value, key)
	if err != nil {
		fmt.Printf("Failed to update config: %+v", err)
		return errors.New("Failed to update config value.")
	}

	return nil
}
