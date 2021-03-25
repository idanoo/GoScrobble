package goscrobble

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bCryptCost = 16

type User struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	password  []byte
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	Active    bool   `json:"active"`
	Admin     bool   `json:"admin"`
}

// RegisterRequest - Incoming JSON
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// createUser - Called from API
func createUser(req *RegisterRequest) error {
	// Check if user already exists..
	if len(req.Password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}

	// Check username is set
	if req.Username == "" {
		return errors.New("A username is required")
	}

	// If set an email.. validate it!
	if req.Email != "" {
		if !isEmailValid(req.Email) {
			return errors.New("Invalid email address")
		}
	}

	// Check if user or email exists!
	if userAlreadyExists(req) {
		return errors.New("Username or email already exists")
	}

	// Lets hashit!
	hash, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	return insertUser(req.Username, req.Email, hash)
}

// insertUser - Does the dirtywork!
func insertUser(username string, email string, password []byte) error {
	_, err := db.Exec("INSERT INTO users (uuid, created_at, username, email, password) VALUES (UUID_TO_BIN(UUID(), true),NOW(),?,?,?)", username, email, password)

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
func userAlreadyExists(req *RegisterRequest) bool {
	var userExists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&userExists)
	if userExists > 0 {
		return true
	}

	if req.Email != "" {
		// Only run email check if there's an email...
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&userExists)
	}

	if err != nil {
		fmt.Printf("Error querying for duplicate users: %v", err)
		return true
	}

	return userExists > 0
}
