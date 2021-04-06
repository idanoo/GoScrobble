package goscrobble

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Limits to 1 req / 4 sec
var heavyLimiter = NewIPRateLimiter(0.25, 2)

// Limits to 5 req / sec
var standardLimiter = NewIPRateLimiter(5, 5)

// Limits to 10 req / sec
var lightLimiter = NewIPRateLimiter(10, 10)

// tokenMiddleware - Validates token to a user
func tokenMiddleware(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := ""
		urlParams := r.URL.Query()
		if val, ok := urlParams["key"]; ok {
			key = val[0]
		} else {
			throwUnauthorized(w, "No key parameter provided")
			return
		}

		if key == "" {
			throwUnauthorized(w, "A token is required")
			return
		}

		userUuid, err := getUserUuidForToken(key)
		if err != nil {
			throwUnauthorized(w, err.Error())
			return
		}

		next(w, r, userUuid)
	}
}

// jwtMiddleware - Validates middleware to a user
func jwtMiddleware(next func(http.ResponseWriter, *http.Request, CustomClaims, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullToken := r.Header.Get("Authorization")
		authToken := strings.Replace(fullToken, "Bearer ", "", 1)
		claims, err := verifyJWTToken(authToken)
		if err != nil {
			throwUnauthorized(w, "Invalid JWT Token")
			return
		}

		var reqUuid string
		for k, v := range mux.Vars(r) {
			if k == "uuid" {
				reqUuid = v
			}
		}

		next(w, r, claims, reqUuid)
	}
}

// adminMiddleware - Validates user is admin
func adminMiddleware(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullToken := r.Header.Get("Authorization")
		authToken := strings.Replace(fullToken, "Bearer ", "", 1)
		claims, err := verifyJWTToken(authToken)
		if err != nil {
			throwUnauthorized(w, "Invalid JWT Token")
			return
		}

		user, err := getUserByUUID(claims.Subject)
		if err != nil {
			throwUnauthorized(w, err.Error())
			return
		}

		if !user.Admin {
			throwUnauthorized(w, "User is not admin")
			return
		}

		next(w, r, claims.Subject)
	}
}

// limitMiddleware - Rate limits important stuff
func limitMiddleware(next http.HandlerFunc, limiter *IPRateLimiter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			jr := jsonResponse{
				Msg: "Too many requests",
			}
			msg, _ := json.Marshal(&jr)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(msg)
			return
		}

		next(w, r)
	})
}
