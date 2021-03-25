package goscrobble

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// spaHandler - Handles Single Page Applications (React)
type spaHandler struct {
	staticPath string
	indexPath  string
}

type jsonResponse struct {
	Err string `json:"error,omitempty"`
	Msg string `json:"message,omitempty"`
}

// Limits to 1 req / 10 sec
var heavyLimiter = NewIPRateLimiter(0.1, 1)

// Limits to 5 req / sec
var standardLimiter = NewIPRateLimiter(5, 5)

// List of Reverse proxies
var ReverseProxies []string

// HandleRequests - Boot HTTP!
func HandleRequests() {
	// Create a new router
	r := mux.NewRouter().StrictSlash(true)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Static Token for /ingress
	v1.HandleFunc("/ingress/jellyfin", tokenMiddleware(serveEndpoint))

	// JWT Auth
	v1.HandleFunc("/profile/{id}", jwtMiddleware(serveEndpoint))

	// No Auth
	v1.HandleFunc("/register", limitMiddleware(handleRegister, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/login", limitMiddleware(handleLogin, standardLimiter)).Methods("POST")
	v1.HandleFunc("/logout", serveEndpoint).Methods("POST")

	// This just prevents it serving frontend stuff over /api
	r.PathPrefix("/api")

	// SERVE FRONTEND - NO AUTH
	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	// Serve it up!
	log.Fatal(http.ListenAndServe(":42069", r))
}

// MIDDLEWARE
// throwUnauthorized - Throws a 403
func throwUnauthorized(w http.ResponseWriter, m string) {
	jr := jsonResponse{
		Err: m,
	}
	js, _ := json.Marshal(&jr)
	err := errors.New(string(js))
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

// throwUnauthorized - Throws a 403
func throwBadReq(w http.ResponseWriter, m string) {
	jr := jsonResponse{
		Err: m,
	}
	js, _ := json.Marshal(&jr)
	err := errors.New(string(js))
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// generateJsonMessage - Generates a message:str response
func generateJsonMessage(m string) []byte {
	jr := jsonResponse{
		Msg: m,
	}
	js, _ := json.Marshal(&jr)
	return js
}

// tokenMiddleware - Validates token to a user
func tokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		throwUnauthorized(w, "Invalid API Token")
		return
		// next(res, req)
	}
}

// jwtMiddleware - Validates middleware to a user
func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		throwUnauthorized(w, "Invalid JWT Token")
		return
		// next(res, req)
	}
}

// limitMiddleware - Rate limits important stuff
func limitMiddleware(next http.HandlerFunc, limiter *IPRateLimiter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			msg := generateJsonMessage("Too many requests")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(msg)
			return
		}

		next(w, r)
	})
}

// API ENDPOINT HANDLING

// handleRegister - Does as it says!
func handleRegister(w http.ResponseWriter, r *http.Request) {
	regReq := RegisterRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&regReq)
	if err != nil {
		throwBadReq(w, err.Error())
		return
	}

	ip := getUserIp(r)
	err = createUser(&regReq, ip)
	if err != nil {
		throwBadReq(w, err.Error())
		return
	}

	msg := generateJsonMessage("User created succesfully")
	w.WriteHeader(http.StatusCreated)
	w.Write(msg)
}

// handleLogin - Does as it says!
func handleLogin(w http.ResponseWriter, r *http.Request) {
	logReq := LoginRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&logReq)
	if err != nil {
		throwBadReq(w, err.Error())
		return
	}

	ip := getUserIp(r)
	data, err := loginUser(&logReq, ip)
	if err != nil {
		throwBadReq(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// serveEndpoint - API stuffs
func serveEndpoint(w http.ResponseWriter, r *http.Request) {
	_, err := decodeJson(r.Body)
	if err != nil {
		// If we can't decode. Lets tell them nicely.
		http.Error(w, "{\"error\":\"Invalid JSON\"}", http.StatusBadRequest)
		return
	}

	// Lets trick 'em for now ;) ;)
	fmt.Fprintf(w, "{}")
}

// FRONTEND HANDLING

// ServerHTTP - Frontend server
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// If we failed to get the absolute path respond with a 400 bad request and return
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
