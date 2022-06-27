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
	URL          string    `json:"url"`
}

func getOauthToken(userUuid string, service string) (OauthToken, error) {
	var oauth OauthToken

	err := db.QueryRow(`SELECT "user", service, access_token, refresh_token, expiry, username, last_synced, url FROM oauth_tokens `+
		`WHERE "user" = $1 AND service = $2`,
		userUuid, service).Scan(&oauth.UserUUID, &oauth.Service, &oauth.AccessToken, &oauth.RefreshToken, &oauth.Expiry, &oauth.Username, &oauth.LastSynced, &oauth.URL)

	if err == sql.ErrNoRows {
		return oauth, errors.New("No token for user")
	}

	return oauth, nil
}

func insertOauthToken(userUuid string, service string, token string, refresh string, expiry time.Time, username string, lastSynced time.Time, url string) error {
	_, err := db.Exec(`REPLACE INTO oauth_tokens ("user", service, access_token, refresh_token, expiry, username, last_synced, url) `+
		`VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, userUuid, service, token, refresh, expiry, username, lastSynced, url)

	return err
}

func removeOauthToken(userUuid string, service string) error {
	_, err := db.Exec(`DELETE FROM oauth_tokens WHERE "user" = $1 AND service = $2`, userUuid, service)

	return err
}
