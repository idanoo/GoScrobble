package goscrobble

import (
	"encoding/json"
	"io"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// decodeJson - Returns a map[string]interface{}
func decodeJson(body io.ReadCloser) (map[string]interface{}, error) {
	var jsonInput map[string]interface{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&jsonInput)

	return jsonInput, err
}

// isEmailValid - checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
