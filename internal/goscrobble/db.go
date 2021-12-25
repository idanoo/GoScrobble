package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDb - Boots up a DB connection
func InitDb() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASS")
	dbName := os.Getenv("POSTGRES_DB")

	dbConn, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, 5432, dbUser, dbPass, dbName),
	)

	if err != nil {
		panic(err)
	}

	dbConn.SetConnMaxLifetime(time.Minute * 3)
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(10)

	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}

	db = dbConn

	runMigrations()
}

// CloseDbConn - Closes DB connection
func CloseDbConn() {
	db.Close()
}

func runMigrations() {
	fmt.Println("Checking database migrations")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Unable to run migrations! %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(fmt.Errorf("Error fetching DB Migrations %v", err))
	}

	err = m.Up()
	if err != nil {
		// Skip 'no change'. This is fine. Everything is fine.
		if err.Error() == "no change" {
			fmt.Println("Database already up to date")
			return
		}

		panic(fmt.Errorf("Error running DB Migrations %v", err))
	}

	fmt.Println("Database migrations complete")
}

func getDbCount(query string, args ...interface{}) (int, error) {
	var result int
	err := db.QueryRow(query, args...).Scan(&result)
	if err != nil {
		return 0, errors.New("Error fetching data")
	}

	return result, nil
}
