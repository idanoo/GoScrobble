package goscrobble

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

// Directories
var FrontendDirectory string
var DataDirectory string
var ApiDocsDirectory string

// RequestRequest - Incoming JSON!
type RequestRequest struct {
	URL      string `json:"url"`
	Token    string `json:"token"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

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
	v1.HandleFunc("/user", limitMiddleware(jwtMiddleware(getUser), lightLimiter)).Methods("GET")
	v1.HandleFunc("/user", limitMiddleware(jwtMiddleware(patchUser), lightLimiter)).Methods("PATCH")
	v1.HandleFunc("/user/navidrome", limitMiddleware(jwtMiddleware(postNavidrome), lightLimiter)).Methods("POST")
	v1.HandleFunc("/user/navidrome", limitMiddleware(jwtMiddleware(deleteNavidrome), lightLimiter)).Methods("DELETE")
	v1.HandleFunc("/user/spotify", limitMiddleware(jwtMiddleware(getSpotifyClientID), lightLimiter)).Methods("GET")
	v1.HandleFunc("/user/spotify", limitMiddleware(jwtMiddleware(deleteSpotify), lightLimiter)).Methods("DELETE")
	v1.HandleFunc("/user/{uuid}/scrobbles", jwtMiddleware(getScrobbles)).Methods("GET")

	// Config auth
	v1.HandleFunc("/config", limitMiddleware(adminMiddleware(getConfig), standardLimiter)).Methods("GET")
	v1.HandleFunc("/config", limitMiddleware(adminMiddleware(postConfig), standardLimiter)).Methods("POST")

	// No Auth
	v1.HandleFunc("/stats", limitMiddleware(handleStats, lightLimiter)).Methods("GET")
	v1.HandleFunc("/profile/{username}", limitMiddleware(getProfile, lightLimiter)).Methods("GET")

	v1.HandleFunc("/artists/top/{uuid}", limitMiddleware(getArtists, lightLimiter)).Methods("GET")
	v1.HandleFunc("/artists/{uuid}", limitMiddleware(getArtist, lightLimiter)).Methods("GET")

	v1.HandleFunc("/albums/top/{uuid}", limitMiddleware(getArtists, lightLimiter)).Methods("GET")
	v1.HandleFunc("/albums/{uuid}", limitMiddleware(getAlbum, lightLimiter)).Methods("GET")

	v1.HandleFunc("/tracks/top/{uuid}", limitMiddleware(getTracks, lightLimiter)).Methods("GET")           // User UUID - Top Tracks
	v1.HandleFunc("/tracks/{uuid}", limitMiddleware(getTrack, lightLimiter)).Methods("GET")                // Track UUID
	v1.HandleFunc("/tracks/{uuid}/top", limitMiddleware(getTopUsersForTrack, lightLimiter)).Methods("GET") // TrackUUID - Top Listeners

	v1.HandleFunc("/register", limitMiddleware(handleRegister, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/login", limitMiddleware(handleLogin, standardLimiter)).Methods("POST")
	v1.HandleFunc("/sendreset", limitMiddleware(handleSendReset, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/resetpassword", limitMiddleware(handleResetPassword, heavyLimiter)).Methods("POST")
	v1.HandleFunc("/serverinfo", getServerInfo).Methods("GET")
	v1.HandleFunc("/refresh", limitMiddleware(handleTokenRefresh, standardLimiter)).Methods("POST")

	// Redirect from Spotify Oauth
	v1.HandleFunc("/link/spotify", limitMiddleware(postSpotifyReponse, lightLimiter))

	// This just prevents it serving frontend stuff over /api
	r.PathPrefix("/api")

	// Serve Images
	spaStatic := spaStaticHandler{staticPath: DataDirectory}
	r.PathPrefix("/img").Handler(spaStatic)

	// Serve API Docs
	apiDocs := apiDocHandler{staticPath: ApiDocsDirectory, indexPath: "index.html"}
	r.PathPrefix("/docs/").Handler(apiDocs)

	// This is a really terrible work around to Slate using relative paths and not
	// picking up the css/img files when you don't have a trailing slash\
	apiDocRedirect := apiDocRedirectHandler{}
	r.PathPrefix("/docs").Handler(apiDocRedirect)

	// Serve Frontend
	spa := spaHandler{staticPath: FrontendDirectory + string(os.PathSeparator) + "build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	// Setup CORS rules
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})
	handler := c.Handler(r)

	// Serve it!
	fmt.Printf("Goscrobble listening on port %s", port)
	fmt.Println("")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// API ENDPOINT HANDLING
// handleRegister - Does as it says!
func handleRegister(w http.ResponseWriter, r *http.Request) {
	registrationEnabled, err := getConfigValue("REGISTRATION_ENABLED")
	if err != nil {
		log.Printf("%+v", err)
		throwOkError(w, "Registration is currently disabled")
		return
	}

	if registrationEnabled == "0" {
		throwOkError(w, "Registration is currently disabled")
		return
	}

	regReq := RequestRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&regReq)
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
	logReq := RequestRequest{}
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

// handleTokenRefresh - Refresh access token based on refresh token
func handleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	logReq := RequestRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&logReq)
	user, err := isValidRefreshToken(logReq.Token)
	if err != nil {
		throwOkError(w, "Invalid refresh token")
		return
	}

	// Issue JWT + Response
	token, err := generateJWTToken(user, logReq.Token)
	if err != nil {
		throwOkError(w, "Failed to refresh Token")
		return
	}

	loginResp := RequestResponse{
		Token: token,
	}

	resp, _ := json.Marshal(&loginResp)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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
	req := RequestRequest{}
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
	var err error
	ip := getUserIp(r)
	tx, _ := db.Begin()

	ingressType := strings.Replace(r.URL.Path, "/api/v1/ingress/", "", 1)

	switch ingressType {
	case "jellyfin":
		jfInput := JellyfinRequest{}
		err := json.NewDecoder(r.Body).Decode(&jfInput)
		if err != nil {
			throwInvalidJson(w)
			return
		}

		err = ParseJellyfinInput(userUuid, jfInput, ip, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			throwOkError(w, err.Error())
			return
		}
	case "multiscrobbler":
		msInput := MultiScrobblerRequest{}
		err := json.NewDecoder(r.Body).Decode(&msInput)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			throwOkError(w, err.Error())
			return
		}

		err = ParseMultiScrobblerInput(userUuid, msInput, ip, tx)
		if err != nil {
			fmt.Println(err)
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

// getUser - Return personal userprofile
func getUser(w http.ResponseWriter, r *http.Request, claims CustomClaims, reqUser string) {
	jwtUser := claims.Subject
	userFull, err := getUserByUUID(jwtUser)
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

	oauthNavi, err := getOauthToken(user.UUID, "navidrome")
	if err == nil {
		user.NavidromeURL = oauthNavi.URL
	}

	oauthSpotify, err := getOauthToken(user.UUID, "spotify")
	if err == nil {
		user.SpotifyUsername = oauthSpotify.Username
	}

	json, _ := json.Marshal(&user)

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// patchUser - Update specific values
func patchUser(w http.ResponseWriter, r *http.Request, claims CustomClaims, reqUser string) {
	jwtUser := claims.Subject

	userFull, err := getUserByUUID(jwtUser)
	if err != nil {
		throwOkError(w, "Failed to fetch user information")
		return
	}

	bodyJson, _ := decodeJson(r.Body)

	ip := getUserIp(r)
	for k, v := range bodyJson {
		val := fmt.Sprintf("%s", v)
		if k == "timezone" {
			if isValidTimezone(val) {
				userFull.updateUser("timezone", val, ip)
			}
		} else if k == "token" {
			token := generateToken(32)
			userFull.updateUser("token", token, ip)
		}
	}

	throwOkMessage(w, "User updated successfully")
}

// getScrobbles - Return an array of scrobbles
func getScrobbles(w http.ResponseWriter, r *http.Request, claims CustomClaims, reqUser string) {
	resp, err := getScrobblesForUser(reqUser, 100, 1)
	if err != nil {
		throwOkError(w, "Failed to fetch scrobbles")
		return
	}

	// Fetch last 500 scrobbles
	json, _ := json.Marshal(&resp)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getConfig - Return an array of scrobbles
func getConfig(w http.ResponseWriter, r *http.Request, jwtUser string) {
	config, err := getAllConfigs()
	if err != nil {
		throwOkError(w, "Failed to fetch scrobbles")
		return
	}

	json, _ := json.Marshal(&config)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// postConfig - Return an array of scrobbles
func postConfig(w http.ResponseWriter, r *http.Request, jwtUser string) {
	bodyJson, err := decodeJson(r.Body)
	if err != nil {
		throwInvalidJson(w)
		return
	}

	for k, v := range bodyJson {
		val := fmt.Sprintf("%s", v)
		err = updateConfigValue(k, val)
		if err != nil {
			throwOkError(w, err.Error())
			return
		}
		setRedisVal(k, val)
	}

	throwOkMessage(w, "Config updated successfully")
}

// getProfile - Returns public user profile data
func getProfile(w http.ResponseWriter, r *http.Request) {
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

	resp, err := getProfileForUser(user)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&resp)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getArtist - Returns artist data
func getArtist(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	artist, err := getArtistByUUID(uuid)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&artist)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getAlbum - Returns album data
func getAlbum(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	album, err := getAlbumByUUID(uuid)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&album)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getTrack - Returns track data
func getTrack(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	// Load track obj
	track, err := getTrackByUUID(uuid)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	// Load in Album/Artist info
	err = track.loadExtraTrackInfo()
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&track)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getArtists - Returns artist data for a user
func getArtists(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	track, err := getTopArtists(uuid)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&track)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getAlbums - Returns album data for a user
func getAlbums(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	album, err := getAlbumByUUID(uuid)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&album)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getTracks - Returns track data for a user
func getTracks(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	track, err := getTopTracks(uuid)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&track)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// getTopUsersForTrack - I suck at naming. Returns top users that have scrobbled this track.
func getTopUsersForTrack(w http.ResponseWriter, r *http.Request) {
	var uuid string
	for k, v := range mux.Vars(r) {
		if k == "uuid" {
			uuid = v
		}
	}

	if uuid == "" {
		throwOkError(w, "Invalid UUID")
		return
	}

	userList, err := getTopUsersForTrackUUID(uuid, 10, 1)
	if err != nil {
		throwOkError(w, err.Error())
		return
	}

	json, _ := json.Marshal(&userList)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// postSpotifyResponse - Oauth Response from Spotify
func postSpotifyReponse(w http.ResponseWriter, r *http.Request) {
	err := connectSpotifyResponse(r)

	if err != nil {
		log.Printf("Post Spotify Response: %s", err)
		throwOkError(w, "Failed to connect to spotify")
		return
	}

	http.Redirect(w, r, os.Getenv("GOSCROBBLE_DOMAIN")+"/user", 302)
}

// getSpotifyClientID - Returns public spotify APP ID
func getSpotifyClientID(w http.ResponseWriter, r *http.Request, claims CustomClaims, v string) {
	key, err := getConfigValue("SPOTIFY_API_ID")
	if err != nil {
		throwOkError(w, "Failed to get Spotify ID")
		return

	}
	response := RequestResponse{
		Token: key,
	}

	resp, _ := json.Marshal(&response)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// deleteSpotify - Unlinks spotify account
func deleteSpotify(w http.ResponseWriter, r *http.Request, claims CustomClaims, v string) {
	jwtUser := claims.Subject
	err := removeOauthToken(jwtUser, "spotify")
	if err != nil {
		fmt.Println(err)
		throwOkError(w, "Failed to unlink spotify account")
		return
	}

	throwOkMessage(w, "Spotify account successfully unlinked")
}

// postNavidrome - Submits data for navidrome URL/User/Password
func postNavidrome(w http.ResponseWriter, r *http.Request, claims CustomClaims, v string) {
	jwtUser := claims.Subject

	request := RequestRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		throwOkError(w, "Invalid JSON")
		return
	}

	// hash password with salt
	salt := generateToken(32)
	hash := getMd5(request.Password + salt)

	err = validateNavidromeConnection(request.URL, request.Username, hash, salt)
	if err != nil {
		throwOkError(w, "Failed to validate credentials")
		return
	}

	// Lets set this back 30min
	time := time.Now().UTC().Add(-(time.Duration(30) * time.Minute))
	err = insertOauthToken(jwtUser, "navidrome", hash, salt, time, request.Username, time, request.URL)
	if err != nil {
		throwOkError(w, "Failed to save Navidome token")
		return
	}

	throwOkMessage(w, "Successfully saved!")
}

// deleteNavidrome - Unlinks Navidrome account
func deleteNavidrome(w http.ResponseWriter, r *http.Request, claims CustomClaims, v string) {
	jwtUser := claims.Subject
	err := removeOauthToken(jwtUser, "navidrome")
	if err != nil {
		fmt.Println(err)
		throwOkError(w, "Failed to unlink navidrome account")
		return
	}

	throwOkMessage(w, "Navidrome account successfully unlinked")
}

func getServerInfo(w http.ResponseWriter, r *http.Request) {
	registrationEnabled, err := getConfigValue("REGISTRATION_ENABLED")
	if err != nil {
		log.Printf("%+v", err)
		registrationEnabled = "0"
	}

	info := ServerInfo{
		Version:             "0.1.1",
		RegistrationEnabled: registrationEnabled,
	}

	js, _ := json.Marshal(&info)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
