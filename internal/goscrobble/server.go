package goscrobble

import (
	"encoding/json"
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

// HandleRequests - Boot HTTP!
func HandleRequests() {
	// Create a new router
	r := mux.NewRouter().StrictSlash(true)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	// STATIC TOKEN AUTH
	// httpRouter.HandleFunc("/api/v1/ingress/jellyfin", serveEndpoint)

	// JWT SESSION AUTH?
	// httpRouter.HandleFunc("/api/v1/profile/{id}", serveEndpoint)

	// NO AUTH
	v1.HandleFunc("/register", serveEndpoint).Methods("POST")
	v1.HandleFunc("/login", serveEndpoint).Methods("POST")
	v1.HandleFunc("/logout", serveEndpoint).Methods("POST")

	// This just prevents it serving frontend over /api
	r.PathPrefix("/api")

	// SERVE FRONTEND - NO AUTH
	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	// Serve it up!
	log.Fatal(http.ListenAndServe(":42069", r))
}

func serveEndpoint(w http.ResponseWriter, r *http.Request) {
	var jsonInput map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jsonInput)
	if err != nil {
		// If we can't decode. Lets tell them nicely.
		http.Error(w, "{\"error\":\"Invalid JSON\"}", http.StatusBadRequest)
		return
	}

	// Lets trick 'em for now ;) ;)
	fmt.Fprintf(w, "{}")
}

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
