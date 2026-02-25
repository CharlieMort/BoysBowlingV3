package main

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func SetupDatabase() {
	var err error
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./bowling.db"
	}
	db, err = sql.Open("sqlite", dbPath)
	CheckNilError(err, "Error opening DB file")

	userQuery := `
		CREATE TABLE IF NOT EXISTS players {
			id INTEGER PRIMARY KEY,
			name TEXT
		};
	`
	db.Exec(userQuery)

	gamesQuery := `
		CREATE TABLE IF NOT EXISTS games {
			id INTEGER PRIMARY KEY,
			name TEXT,
			date TEXT
		}
	`
	db.Exec(gamesQuery)

	framesQuery := `
		CREATE TABLE IF NOT EXISTS frames {
			id INTEGER PRIMARY KEY,
			playerId INTEGER,
			gameId INTEGER,
			total INTEGER,
			scorecard TEXT,
			imgPath TEXT
		};
	`
	db.Exec(framesQuery)
}
