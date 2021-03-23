package goscrobble

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type dbItem struct {
	conn    *sql.DB
	version *int
}

var db dbItem

// InitDb - Boots up a DB connection
func InitDb() {
	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbName := os.Getenv("MYSQL_DB")

	dbConn, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
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

	var vers int = 0
	db = dbItem{
		conn:    dbConn,
		version: &vers,
	}

	if !checkIfDbUpToDate() {
		fmt.Printf("Database not up to date.. triggering migrations!\n")
	} else {
		fmt.Printf("Database up to date!\n")

	}

}

// CloseDbConn - Closes DB connection
func CloseDbConn() {
	db.conn.Close()
}

// checkIfDbUpToDate - Checks if we need to run migrations
func checkIfDbUpToDate() bool {
	fmt.Printf("Code version: %v. ", DBVersion)

	dbVers := db.getDbVersion()
	fmt.Printf("DB version: %v.\n", dbVers)

	if dbVers == DBVersion {
		return true
	} else if dbVers > DBVersion {
		panic("!!Warning!! Your database is newer than the code. Please update!")
	}

	return false
}

// getDbVersion - Gets version of schema or generate basic schema
func (db dbItem) getDbVersion() int {
	stmtOut, err := db.conn.Prepare("SELECT version FROM goscrobble WHERE id = ? ")
	defer stmtOut.Close()

	if err != nil {
		// We can assume this is a fresh database - Lets config it!
		return runMigrations(DBVersion)
	}

	err = stmtOut.QueryRow(1).Scan(db.version)
	if err != nil {
		panic(err.Error())
	}

	return *db.version
}

func runMigrations(latestVersion int) int {
	driver, err := mysql.WithInstance(db.conn, &mysql.Config{})
	if err != nil {
		log.Fatalf("Unable to run migrations! %v", err)
	}

	m, _ := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"mysql",
		driver,
	)

	m.Steps(2)

	return latestVersion
}
