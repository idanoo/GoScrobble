package goscrobble

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bCryptCost = 16

type User struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	password []byte
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Active   bool   `json:"active"`
	Admin    bool   `json:"admin"`
}

// createUser - Called from API
func createUser(username string, email string, password string) error {
	// Check if user already exists..
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}

	// Check username is set
	if username == "" {
		return errors.New("A username is required")
	}

	// Check if user or email exists!
	if userAlreadyExists(username, email) {
		return errors.New("Username or email already exists")
	}

	// Lets hashit!
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	return insertUser(username, email, hash)
}

// insertUser - Does the dirtywork!
func insertUser(username string, email string, password []byte) error {
	_, err := db.Exec("INSERT INTO users (uuid, username, email, password) VALUES (UUID_TO_BIN(UUID(), true),'?','?','?')", username, email, password)

	return err
}

// hashPassword - Returns bcrypt hash
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
}

// isValidPassword - Checks if password is valid
func isValidPassword(password string, user User) bool {
	err := bcrypt.CompareHashAndPassword(user.password, []byte(password))
	if err != nil {
		return false
	}

	return true
}

// userAlreadyExists - Returns bool indicating if a record exists for either username or email
// Using two look ups to make use of DB indexes.
func userAlreadyExists(username string, email string) bool {
	var usernameCount, emailCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = '?'", username).Scan(&usernameCount)
	// Only run email check if there's an email...
	if email != "" {
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = '?'", email).Scan(&emailCount)
	} else {
		emailCount = 0
	}

	if err != nil {
		fmt.Printf("Error querying for duplicate users: %v", err)
		return true
	}

	count := usernameCount + emailCount

	// If there is more than one.. Return true. User exists.
	return count != 0
}
