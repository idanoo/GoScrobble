package goscrobble

import (
	"errors"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func getUserForToken(token string) (string, error) {
	var uuid string
	cachedKey := getRedisVal("user_token:" + token)
	if cachedKey == "" {
		err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true) FROM `users` WHERE `token` = ? AND `active` = 1", token).Scan(&uuid)
		if err != nil {
			return "", errors.New("Invalid Token")
		}
		setRedisVal("user_token:"+token, uuid)
	} else {
		uuid = cachedKey
	}

	return uuid, nil
}
