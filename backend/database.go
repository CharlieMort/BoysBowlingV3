package main

import (
	"database/sql"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Frame struct {
	Id        int    `json:"id"`
	PlayerId  int    `json:"playerId"`
	GameId    int    `json:"gameId"`
	Total     int    `json:"total"`
	Scorecard string `json:"scorecard"`
	ImgPath   string `json:"imgPath"`
}

func SetupDatabase() {
	var err error
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./bowling.db"
	}
	db, err = sql.Open("sqlite", dbPath)
	CheckNilError(err, "Error opening DB file")

	userQuery := `
		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE
		);
	`
	db.Exec(userQuery)

	gamesQuery := `
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY,
			name TEXT,
			date TEXT
		);
	`
	db.Exec(gamesQuery)

	framesQuery := `
		CREATE TABLE IF NOT EXISTS frames (
			id INTEGER PRIMARY KEY,
			playerId INTEGER,
			gameId INTEGER,
			total INTEGER,
			scorecard TEXT,
			imgPath TEXT
		);
	`
	db.Exec(framesQuery)
	// InsertUser("Charlie")
	// InsertUser("John")
	// InsertUser("Arthur")
	// InsertUser("Ani")
}

func InsertFrame(playerId int, gameId int, total int, scorecard string, imgPath string) (int64, error) {
	query := `
		INSERT INTO frames (
			playerId,
			gameId,
			total,
			scorecard,
			imgPath
		) VALUES ( ?, ?, ?, ?, ?);
	`
	res, err := db.Exec(query, playerId, gameId, total, scorecard, imgPath)
	CheckNilError(err, "failed inserting frame")
	LastModifiedFrames = time.Now()
	return res.LastInsertId()
}

func InsertUser(name string) (int64, error) {
	query := `
		INSERT INTO players (
			name
		) VALUES ( ? );
	`
	res, err := db.Exec(query, name)
	CheckNilError(err, "failed inserting user")
	return res.LastInsertId()
}

func InsertGame(name string, date string) (int64, error) {
	query := `
		INSERT INTO games (
			name,
			date
		) VALUES ( ?, ? );
	`
	res, err := db.Exec(query, name, date)
	CheckNilError(err, "failed inserting game")
	return res.LastInsertId()
}

func GetFrames() []Frame {
	rows, err := db.Query("SELECT * FROM frames;")
	CheckNilError(err, "failed getting frames")
	var frames []Frame
	for rows.Next() {
		var frame Frame
		rows.Scan(&frame.Id, &frame.PlayerId, &frame.GameId, &frame.Total, &frame.Scorecard, &frame.ImgPath)
		frames = append(frames, frame)
	}
	LastFetchedFrames = time.Now()
	return frames
}

func GetUserIdFromName(name string) int {
	query := `SELECT id FROM players WHERE name = ?;`
	row := db.QueryRow(query, name)
	var id int
	row.Scan(&id)
	return id
}

func GetNameFromId(id int) string {
	query := `SELECT name FROM players WHERE id = ?;`
	row := db.QueryRow(query, id)
	var name string
	row.Scan(&name)
	return name
}
