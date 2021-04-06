package goscrobble

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtToken - Store token from .env
var JwtToken []byte

// JwtExpiry - Expiry in seconds
var JwtExpiry time.Duration

// RefereshExpiry - Expiry for refresh token
var RefereshExpiry time.Duration

type CustomClaims struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Admin        bool   `json:"admin"`
	RefreshToken string `json:"refresh_token"`
	RefreshExp   int    `json:"refresh_exp"`
	jwt.StandardClaims
}

func generateJWTToken(user User, existingRefresh string) (string, error) {
	refreshToken := generateToken(64)

	atClaims := jwt.MapClaims{}
	atClaims["sub"] = user.UUID
	atClaims["username"] = user.Username
	atClaims["email"] = user.Email
	atClaims["admin"] = user.Admin
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(JwtExpiry).Unix()
	atClaims["refresh_token"] = refreshToken
	atClaims["refresh_exp"] = time.Now().Add(RefereshExpiry).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	token, err := at.SignedString(JwtToken)
	if err != nil {
		return "", err
	}

	// Store refresh token
	err = insertRefreshToken(user.UUID, refreshToken)
	if err != nil {
		fmt.Println(err)
		return token, errors.New("Failed to generate token")
	}

	if existingRefresh != "" {
		deleteRefreshToken(existingRefresh)
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
