package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize initialises the database
func Initialize(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil || db == nil {
		panic("Error connecting to database")
	}

	return db
}

// Migrate migrates the database
func Migrate(db *sql.DB) {
	sql := `
        CREATE TABLE IF NOT EXISTS devices(
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				uuid VARCHAR NOT NULL,
				btc_min INTEGER,
				btc_max INTEGER,
				eth_min INTEGER,
				eth_max INTEGER
        );
   `

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
