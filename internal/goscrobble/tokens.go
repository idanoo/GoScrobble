package goscrobble

import (
	"errors"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type RefreshToken struct {
	UUID   string
	User   string
	Token  string
	Expiry time.Time
}

// generateToken - Generates a unique token for user input
func generateToken(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func getUserUuidForToken(token string) (string, error) {
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

func insertRefreshToken(userUuid string, token string) error {
	uuid := newUUID()
	_, err := db.Exec("INSERT INTO `refresh_tokens` (`uuid`, `user`, `token`) VALUES (UUID_TO_BIN(?,true),UUID_TO_BIN(?,true),?)",
		uuid, userUuid, token)

	return err
}

func deleteRefreshToken(token string) error {
	_, err := db.Exec("DELETE FROM `refresh_tokens` WHERE `token` = ?", token)

	return err
}

func isValidRefreshToken(refreshTokenStr string) (User, error) {
	var refresh RefreshToken
	err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), BIN_TO_UUID(`user`, true), `token`, `expiry` FROM `refresh_tokens` WHERE `token` = ?",
		refreshTokenStr).Scan(&refresh.UUID, &refresh.User, &refresh.Token, &refresh.Expiry)
	if err != nil {
		return User{}, errors.New("Invalid Refresh Token")
	}

	// Validate Expiry
	if refresh.Expiry.Unix() < time.Now().Unix() {
		return User{}, errors.New("Invalid Refresh Token")
	}

	user, err := getUserByUUID(refresh.User)
	if err != nil {
		return User{}, errors.New("Invalid Refresh Token")
	}

	return user, nil
}
