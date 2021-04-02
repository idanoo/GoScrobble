package goscrobble

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type jsonResponse struct {
	Err   string `json:"error,omitempty"`
	Msg   string `json:"message,omitempty"`
	Valid bool   `json:"valid,omitempty"`
}

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
	v1.HandleFunc("/ingress/jellyfin", limitMiddleware(tokenMiddleware(handleIngress), lightLimiter)).Methods("POST")
	v1.HandleFunc("/ingress/multiscrobbler", limitMiddleware(tokenMiddleware(handleIngress), lightLimiter)).Methods("POST")

	// JWT Auth - Own profile only (Uses uuid in JWT)
	v1.HandleFunc("/user", limitMiddleware(jwtMiddleware(fetchUser), lightLimiter)).Methods("GET")
	// v1.HandleFunc("/user", jwtMiddleware(fetchScrobbleResponse)).Methods("PATCH")
	v1.HandleFunc("/user/spotify", limitMiddleware(jwtMiddleware(getSpotifyClientID), lightLimiter)).Methods("GET")
	v1.HandleFunc("/user/spotify", limitMiddleware(jwtMiddleware(deleteSpotifyLink), lightLimiter)).Methods("DELETE")
	v1.HandleFunc("/user/{uuid}/scrobbles", jwtMiddleware(fetchScrobbleResponse)).Methods("GET")

	// Config auth
	v1.HandleFunc("/config", limitMiddleware(adminMiddleware(fetchConfig), standardLimiter)).Methods("GET")
	v1.HandleFunc("/config", limitMiddleware(adminMiddleware(postConfig), standardLimiter)).Methods("POST")

	// No Auth
	v1.HandleFunc("/stats", limitMiddleware(handleStats, lightLimiter)).Methods("GET")
	v1.HandleFunc("/profile/{username}", limitMiddleware(fetchProfile, lightLimiter)).Methods("GET")

	v1.HandleFunc("/register", limitMiddleware(handleRegister, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/login", limitMiddleware(handleLogin, standardLimiter)).Methods("POST")
	v1.HandleFunc("/sendreset", limitMiddleware(handleSendReset, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/resetpassword", limitMiddleware(handleResetPassword, heavyLimiter)).Methods("POST")

	// Redirect from Spotify Oauth
	v1.HandleFunc("/link/spotify", limitMiddleware(postSpotifyReponse, lightLimiter))

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
	fmt.Println("")

	log.Fatal(http.ListenAndServe(":"+port, handler))
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

// handleSendReset - Does as it says!
func handleSendReset(w http.ResponseWriter, r *http.Request) {
	req := RegisterRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		throwBadReq(w, err.Error())
		return
	}

	if req.Email == "" {
		throwOkError(w, "Invalid Email")
		return
	}

	_ = getUserIp(r)
	user, err := getUserByEmail(req.Email)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	ip := getUserIp(r)
	err = user.sendResetEmail(ip)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	throwOkMessage(w, "Password reset email sent")
}

// handleSendReset - Does as it says!
func handleResetPassword(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := decodeJson(r.Body)
	if err != nil {
		throwInvalidJson(w)
		return
	}

	if bodyJson["password"] == nil {
		// validating
		valid, err := checkResetToken(fmt.Sprintf("%s", bodyJson["token"]))
		if err != nil {
			throwOkError(w, err.Error())
			return
		}
		jr := jsonResponse{
			Valid: valid,
		}
		msg, _ := json.Marshal(&jr)
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
		return
	} else {
		// resetting
		token := fmt.Sprintf("%s", bodyJson["token"])
		pw := fmt.Sprintf("%s", bodyJson["password"])
		if len(pw) < 8 {
			throwOkError(w, "Password must be at least 8 characters")
			return
		}

		ip := getUserIp(r)
		user, err := getUserByResetToken(token)
		if err != nil {
			throwOkError(w, err.Error())
			return
		}
		err = user.updatePassword(pw, ip)
		if err != nil {
			throwOkError(w, err.Error())
			return
		}

		throwOkMessage(w, "Password updated successfully!")
		return
	}
}

// serveEndpoint - API stuffs
func handleIngress(w http.ResponseWriter, r *http.Request, userUuid string) {
	bodyJson, err := decodeJson(r.Body)
	if err != nil {
		throwInvalidJson(w)
		return
	}

	ip := getUserIp(r)
	tx, _ := db.Begin()

	ingressType := strings.Replace(r.URL.Path, "/api/v1/ingress/", "", 1)

	switch ingressType {
	case "jellyfin":
		err := ParseJellyfinInput(userUuid, bodyJson, ip, tx)
		if err != nil {
			fmt.Printf("Err? %+v", err)
			tx.Rollback()
			throwOkError(w, err.Error())
			return
		}
	case "multiscrobbler":
		err := ParseMultiScrobblerInput(userUuid, bodyJson, ip, tx)
		if err != nil {
			tx.Rollback()
			throwOkError(w, err.Error())
			return
		}
	default:
		tx.Rollback()
		throwBadReq(w, "Unknown ingress type")
	}

	err = tx.Commit()
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	throwOkMessage(w, "success")
	return
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

	//
	oauth, err := getOauthToken(user.UUID, "spotify")
	if err == nil {
		user.SpotifyUsername = oauth.Username
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

// postSpotifyResponse - Oauth Response from Spotify
func postSpotifyReponse(w http.ResponseWriter, r *http.Request) {
	err := connectSpotifyResponse(r)

	if err != nil {
		throwOkError(w, "Failed to connect to spotify")
		return
	}

	http.Redirect(w, r, os.Getenv("GOSCROBBLE_DOMAIN")+"/user", 302)
}

// getSpotifyClientID - Returns public spotify APP ID
func getSpotifyClientID(w http.ResponseWriter, r *http.Request, u string, v string) {
	key, err := getConfigValue("SPOTIFY_APP_ID")
	if err != nil {
		throwOkError(w, "Failed to get Spotify ID")
		return

	}
	response := LoginResponse{
		Token: key,
	}

	resp, _ := json.Marshal(&response)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// deleteSpotifyLink - Unlinks spotify account
func deleteSpotifyLink(w http.ResponseWriter, r *http.Request, u string, v string) {
	err := removeOauthToken(u, "spotify")
	if err != nil {
		fmt.Println(err)
		throwOkError(w, "Failed to unlink spotify account")
		return
	}

	throwOkMessage(w, "Spotify account successfully unlinked")
}
