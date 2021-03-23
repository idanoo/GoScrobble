package goscrobble

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleRequests - Boot HTTP server
func HandleRequests() {
	// creates a new instance of a mux router
	httpRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	httpRouter.HandleFunc("/", serveFrontend)
	httpRouter.HandleFunc("/api/v1", serveEndpoint)
	httpRouter.HandleFunc("/api/v1/scrobble/jellyfin", serveEndpoint)
	httpRouter.HandleFunc("/api/v1/jellyfin", serveEndpoint)

	// Serve HTTP Server
	log.Fatal(http.ListenAndServe(":42069", httpRouter))
}

// serveFrontend - Handle / queries
func serveFrontend(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func serveEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}")
}
