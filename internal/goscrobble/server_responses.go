package goscrobble

import (
	"encoding/json"
	"errors"
	"net/http"
)

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
