package goscrobble

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

// Limits to 1 req / 4 sec
var heavyLimiter = NewIPRateLimiter(0.25, 2)

// Limits to 5 req / sec
var standardLimiter = NewIPRateLimiter(5, 5)

// Limits to 10 req / sec
var lightLimiter = NewIPRateLimiter(10, 10)

// List of Reverse proxies
var ReverseProxies []string

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// HandleRequests - Boot HTTP!
func HandleRequests(port string) {
	// Create a new router
	r := mux.NewRouter().StrictSlash(true)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Static Token for /ingress
	v1.HandleFunc("/ingress/jellyfin", tokenMiddleware(handleIngress)).Methods("POST")

	// JWT Auth - PWN PROFILE ONLY.
	v1.HandleFunc("/user", jwtMiddleware(fetchUser)).Methods("GET")
	// v1.HandleFunc("/user", jwtMiddleware(fetchScrobbleResponse)).Methods("PATCH")
	v1.HandleFunc("/user/{uuid}/scrobbles", jwtMiddleware(fetchScrobbleResponse)).Methods("GET")

	// Config auth
	v1.HandleFunc("/config", adminMiddleware(fetchConfig)).Methods("GET")
	v1.HandleFunc("/config", adminMiddleware(postConfig)).Methods("POST")

	// No Auth
	v1.HandleFunc("/stats", handleStats).Methods("GET")
	v1.HandleFunc("/profile/{username}", limitMiddleware(fetchProfile, lightLimiter)).Methods("GET")

	v1.HandleFunc("/register", limitMiddleware(handleRegister, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/login", limitMiddleware(handleLogin, standardLimiter)).Methods("POST")

	// This just prevents it serving frontend stuff over /api
	r.PathPrefix("/api")

	// SERVE FRONTEND - NO AUTH
	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	c := cors.New(cors.Options{
		// Grrrr CORS. To clean up at a later date
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	handler := c.Handler(r)

	// Serve it up!
	fmt.Printf("Goscrobble listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// MIDDLEWARE RESPONSES
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

// throwOkError - Throws a 403
func throwOkError(w http.ResponseWriter, m string) {
	jr := jsonResponse{
		Err: m,
	}
	js, _ := json.Marshal(&jr)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// throwOkMessage - Throws a happy 200
func throwOkMessage(w http.ResponseWriter, m string) {
	jr := jsonResponse{
		Msg: m,
	}
	js, _ := json.Marshal(&jr)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// throwOkMessage - Throws a happy 200
func throwInvalidJson(w http.ResponseWriter) {
	jr := jsonResponse{
		Err: "Invalid JSON",
	}
	js, _ := json.Marshal(&jr)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(js)
}

// generateJsonMessage - Generates a message:str response
func generateJsonMessage(m string) []byte {
	jr := jsonResponse{
		Msg: m,
	}
	js, _ := json.Marshal(&jr)
	return js
}

// generateJsonError - Generates a err:str response
func generateJsonError(m string) []byte {
	jr := jsonResponse{
		Err: m,
	}
	js, _ := json.Marshal(&jr)
	return js
}

// MIDDLEWARE ACTIONS
// tokenMiddleware - Validates token to a user
func tokenMiddleware(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullToken := r.Header.Get("Authorization")
		authToken := strings.Replace(fullToken, "Bearer ", "", 1)
		if authToken == "" {
			throwUnauthorized(w, "A token is required")
			return
		}

		userUuid, err := getUserUuidForToken(authToken)
		if err != nil {
			throwUnauthorized(w, err.Error())
			return
		}

		next(w, r, userUuid)
	}
}

// jwtMiddleware - Validates middleware to a user
func jwtMiddleware(next func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
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

		next(w, r, claims.Subject, reqUuid)
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

		user, err := getUser(claims.Subject)
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
		throwOkError(w, err.Error())
		return
	}

	throwOkMessage(w, "User created succesfully. You may now login")
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
		throwOkError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleStats - Returns stats for homepage
func handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := getStats()
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	js, _ := json.Marshal(&stats)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// serveEndpoint - API stuffs
func handleIngress(w http.ResponseWriter, r *http.Request, userUuid string) {
	bodyJson, err := decodeJson(r.Body)
	if err != nil {
		throwInvalidJson(w)
		return
	}

	ingressType := strings.Replace(r.URL.Path, "/api/v1/ingress/", "", 1)

	switch ingressType {
	case "jellyfin":
		tx, _ := db.Begin()

		ip := getUserIp(r)
		err := ParseJellyfinInput(userUuid, bodyJson, ip, tx)
		if err != nil {
			// log.Printf("Error inserting track: %+v", err)
			tx.Rollback()
			throwOkError(w, err.Error())
			return
		}

		err = tx.Commit()
		if err != nil {
			throwOkError(w, err.Error())
			return
		}

		throwOkMessage(w, "success")
		return
	}

	throwBadReq(w, "Unknown ingress type")
}

// fetchUser - Return personal userprofile
func fetchUser(w http.ResponseWriter, r *http.Request, jwtUser string, reqUser string) {
	// We don't this var most of the time
	userFull, err := getUser(jwtUser)
	if err != nil {
		throwOkError(w, "Failed to fetch user information")
		return
	}

	jsonFull, _ := json.Marshal(&userFull)

	// Lets strip out vars we don't want to send.
	user := UserResponse{}
	err = json.Unmarshal(jsonFull, &user)
	if err != nil {
		throwOkError(w, "Failed to fetch user information")
		return
	}
	json, _ := json.Marshal(&user)

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// fetchScrobbles - Return an array of scrobbles
func fetchScrobbleResponse(w http.ResponseWriter, r *http.Request, jwtUser string, reqUser string) {
	resp, err := fetchScrobblesForUser(reqUser, 100, 1)
	if err != nil {
		throwOkError(w, "Failed to fetch scrobbles")
		return
	}

	// Fetch last 500 scrobbles
	json, _ := json.Marshal(&resp)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// fetchScrobbles - Return an array of scrobbles
func fetchConfig(w http.ResponseWriter, r *http.Request, jwtUser string) {
	config, err := getAllConfigs()
	if err != nil {
		throwOkError(w, "Failed to fetch scrobbles")
		return
	}

	json, _ := json.Marshal(&config)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// fetchScrobbles - Return an array of scrobbles
func postConfig(w http.ResponseWriter, r *http.Request, jwtUser string) {
	bodyJson, err := decodeJson(r.Body)
	if err != nil {
		throwInvalidJson(w)
		return
	}

	for k, v := range bodyJson {
		err = updateConfigValue(k, fmt.Sprintf("%s", v))
		if err != nil {
			throwOkError(w, err.Error())
			return
		}
	}

	throwOkMessage(w, "Config updated successfully")
}

// fetchProfile - Returns public user profile data
func fetchProfile(w http.ResponseWriter, r *http.Request) {
	var username string
	for k, v := range mux.Vars(r) {
		if k == "username" {
			username = v
		}
	}

	if username == "" {
		throwOkError(w, "Invalid Username")
		return
	}

	user, err := getUserByUsername(username)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	resp, err := getProfile(user)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&resp)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
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
