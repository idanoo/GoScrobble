package goscrobble

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// spaStaticHandler - Handles static imges
type spaStaticHandler struct {
	staticPath string
	indexPath  string
}

// spaHandler - Handles Single Page Applications (React)
type spaHandler struct {
	staticPath string
	indexPath  string
}

// apiDocHandler - Handles API Docs
type apiDocHandler struct {
	staticPath string
	indexPath  string
}

// apiDocRedirectHandler - Handles redirect to add trailing slash (Fixes relative URLs...)
type apiDocRedirectHandler struct {
}

// ServeHTTP - Frontend React server
func (h spaStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.ServeFile(w, r, filepath.Join(h.staticPath, "img/placeholder.jpg"))
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve placeholder
		http.ServeFile(w, r, filepath.Join(h.staticPath, "img/placeholder.jpg"))
		return
	} else if err != nil {
		http.ServeFile(w, r, filepath.Join(h.staticPath, "img/placeholder.jpg"))
		return
	}

	// otherwise, use http.FileServer to serve the static images
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// ServeHTTP - Frontend React server
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

// ServeHTTP - API Redirect to add trailing slash
func (h apiDocRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/docs/", 301)
	return
}

// ServeHTTP - API Docs
func (h apiDocHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// If we failed to get the absolute path respond with a 400 bad request and return
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, strings.Replace(path, "/docs", "", 1))

	info, err := os.Stat(path)
	if info.IsDir() {
		// Root directory, server index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Strip /docs from request and serve folder
	handler := http.FileServer(http.Dir(h.staticPath))
	http.StripPrefix("/docs", handler).ServeHTTP(w, r)
}
