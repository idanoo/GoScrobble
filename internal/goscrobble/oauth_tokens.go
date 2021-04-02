package goscrobble

import (
	"database/sql"
	"errors"
	"time"
)

type OauthToken struct {
	UserUUID     string    `json:"user"`
	Service      string    `json:"service"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	Username     string    `json:"username"`
	LastSynced   time.Time `json:"last_synced"`
}

func getOauthToken(userUuid string, service string) (OauthToken, error) {
	var oauth OauthToken

	err := db.QueryRow("SELECT BIN_TO_UUID(`user`, true), `service`, `access_token`, `refresh_token`, `expiry`, `username`, `last_synced` FROM `oauth_tokens` "+
		"WHERE `user` = UUID_TO_BIN(?, true) AND `service` = ?",
		userUuid, service).Scan(&oauth.UserUUID, &oauth.Service, &oauth.AccessToken, &oauth.RefreshToken, &oauth.Expiry, &oauth.Username, &oauth.LastSynced)

	if err == sql.ErrNoRows {
		return oauth, errors.New("No token for user")
	}

	return oauth, nil
}

func insertOauthToken(userUuid string, service string, token string, refresh string, expiry time.Time, username string, lastSynced time.Time) error {
	_, err := db.Exec("REPLACE INTO `oauth_tokens` (`user`, `service`, `access_token`, `refresh_token`, `expiry`, `username`, `last_synced`) "+
		"VALUES (UUID_TO_BIN(?, true),?,?,?,?,?,?)", userUuid, service, token, refresh, expiry, username, lastSynced)

	return err
}

func removeOauthToken(userUuid string, service string) error {
	_, err := db.Exec("DELETE FROM `oauth_tokens` WHERE `user` = UUID_TO_BIN(?, true) AND `service` = ?", userUuid, service)

	return err
}