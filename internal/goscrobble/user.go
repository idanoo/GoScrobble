package goscrobble

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bCryptCost = 16

type User struct {
	UUID       string    `json:"uuid"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedIp  net.IP    `json:"created_ip"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedIP net.IP    `jsos:"modified_ip"`
	Username   string    `json:"username"`
	Password   []byte    `json:"password"`
	Email      string    `json:"email"`
	Verified   bool      `json:"verified"`
	Active     bool      `json:"active"`
	Admin      bool      `json:"admin"`
}

type UserResponse struct {
	UUID       string    `json:"uuid"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedIp  net.IP    `json:"created_ip"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedIP net.IP    `jsos:"modified_ip"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Verified   bool      `json:"verified"`
}

// RegisterRequest - Incoming JSON
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest - Incoming JSON
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse - JWT issued
type LoginResponse struct {
	Token string `json:"token"`
}

// createUser - Called from API
func createUser(req *RegisterRequest, ip net.IP) error {
	// Check if user already exists..
	if len(req.Password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}

	// Check Username is set
	if req.Username == "" {
		return errors.New("A username is required")
	}

	// Check username is valid
	if !isUsernameValid(req.Username) {
		return errors.New("Username contains invalid characters")
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

	return insertUser(req.Username, req.Email, hash, ip)
}

func loginUser(logReq *LoginRequest, ip net.IP) ([]byte, error) {
	var resp []byte
	var user User

	if logReq.Username == "" {
		return resp, errors.New("A username is required")
	}

	if logReq.Password == "" {
		return resp, errors.New("A password is required")
	}

	if strings.Contains(logReq.Username, "@") {
		err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), `username`, `email`, `password`, `admin` FROM `users` WHERE `email` = ? AND `active` = 1",
			logReq.Username).Scan(&user.UUID, &user.Username, &user.Email, &user.Password, &user.Admin)
		if err != nil {
			if err == sql.ErrNoRows {
				return resp, errors.New("Invalid Username or Password")
			}
		}
	} else {
		err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), `username`, `email`, `password`, `admin` FROM `users` WHERE `username` = ? AND `active` = 1",
			logReq.Username).Scan(&user.UUID, &user.Username, &user.Email, &user.Password, &user.Admin)
		if err == sql.ErrNoRows {
			return resp, errors.New("Invalid Username or Password")
		}
	}

	if !isValidPassword(logReq.Password, user) {
		return resp, errors.New("Invalid Username or Password")
	}

	// Issue JWT + Response
	token, err := generateJWTToken(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return resp, errors.New("Error logging in")
	}

	loginResp := LoginResponse{
		Token: token,
	}

	resp, _ = json.Marshal(&loginResp)
	return resp, nil
}

// insertUser - Does the dirtywork!
func insertUser(username string, email string, password []byte, ip net.IP) error {
	token := generateToken(32)
	_, err := db.Exec("INSERT INTO users (uuid, created_at, created_ip, modified_at, modified_ip, username, email, password, token) "+
		"VALUES (UUID_TO_BIN(UUID(), true),NOW(),?,NOW(),?,?,?,?,?)", ip, ip, username, email, password, token)

	return err
}

func updateUser(uuid string, field string, value string, ip net.IP) error {
	_, err := db.Exec("UPDATE users SET ? = ?, modified_at = NOW(), modified_ip = ? WHERE uuid = ?", field, value, uuid, ip)

	return err
}

func updateUserDirect(uuid string, field string, value string) error {
	_, err := db.Exec("UPDATE users SET ? = ? WHERE uuid = ?", field, value, uuid)

	return err
}

// hashPassword - Returns bcrypt hash
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
}

// isValidPassword - Checks if password is valid
func isValidPassword(password string, user User) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return false
	}

	return true
}

// userAlreadyExists - Returns bool indicating if a record exists for either username or email
// Using two look ups to make use of DB indexes.
func userAlreadyExists(req *RegisterRequest) bool {
	count, err := getDbCount("SELECT COUNT(*) FROM users WHERE username = ?", req.Username)
	if err != nil {
		fmt.Printf("Error querying for duplicate users: %v", err)
		return true
	}

	if count > 0 {
		return true
	}

	if req.Email != "" {
		// Only run email check if there's an email...
		count, err = getDbCount("SELECT COUNT(*) FROM users WHERE email = ?", req.Email)
	}

	if err != nil {
		fmt.Printf("Error querying for duplicate users: %v", err)
		return true
	}

	return count > 0
}

func getUser(uuid string) (User, error) {
	var user User
	err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), `created_at`, `created_ip`, `modified_at`, `modified_ip`, `username`, `email`, `password`, `verified`, `admin` FROM `users` WHERE `uuid` = UUID_TO_BIN(?, true) AND `active` = 1",
		uuid).Scan(&user.UUID, &user.CreatedAt, &user.CreatedIp, &user.ModifiedAt, &user.ModifiedIP, &user.Username, &user.Email, &user.Password, &user.Verified, &user.Admin)

	if err == sql.ErrNoRows {
		return user, errors.New("Invalid JWT Token")
	}

	return user, nil
}

func getUserByUsername(username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT BIN_TO_UUID(`uuid`, true), `created_at`, `created_ip`, `modified_at`, `modified_ip`, `username`, `email`, `password`, `verified`, `admin` FROM `users` WHERE `username` = ? AND `active` = 1",
		username).Scan(&user.UUID, &user.CreatedAt, &user.CreatedIp, &user.ModifiedAt, &user.ModifiedIP, &user.Username, &user.Email, &user.Password, &user.Verified, &user.Admin)

	if err == sql.ErrNoRows {
		return user, errors.New("Invalid Username")
	}

	return user, nil
}
