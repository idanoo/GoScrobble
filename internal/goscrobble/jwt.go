package goscrobble

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtToken - Store token from .env
var JwtToken []byte

// JwtExpiry - Expiry in seconds
var JwtExpiry time.Duration

type CustomClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func generateJWTToken(user User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["sub"] = user.UUID
	atClaims["username"] = user.Username
	atClaims["email"] = user.Email
	atClaims["admin"] = user.Admin
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(JwtExpiry).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	token, err := at.SignedString(JwtToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

// verifyToken - Verifies the JWT is valid
func verifyJWTToken(token string) (CustomClaims, error) {
	// Initialize a new instance of `Claims`
	claims := CustomClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtToken, nil
	})

	// Verify Signature
	if err != nil {
		return claims, err
	}

	// Verify expiry
	err = claims.Valid()
	if err != nil {
		return claims, err
	}

	return claims, err
}

func getClaims(token *jwt.Token) CustomClaims {
	claims, _ := token.Claims.(CustomClaims)
	return claims
}
