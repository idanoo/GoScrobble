package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var db *sql.DB

// InitDb - Boots up a DB connection
func InitDb() {
	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbName := os.Getenv("MYSQL_DB")

	dbConn, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName+"?multiStatements=true")
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
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Unable to run migrations! %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
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
