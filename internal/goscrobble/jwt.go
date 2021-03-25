package goscrobble

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// JwtToken - Store token from .env
var JwtToken []byte

// Store custom claims here
type Claims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

// verifyToken - Verifies the JWT is valid
func verifyToken(token string, w http.ResponseWriter) bool {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(JwtToken *jwt.Token) (interface{}, error) {
		return JwtToken, nil
	})

	if err != nil {
		log.Printf("%v", err)
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}

		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}
