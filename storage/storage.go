package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(filepath.Dir(exe), "rates.db")

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rates (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp INTEGER NOT NULL,
		best_ask  TEXT    NOT NULL,
		best_bid  TEXT    NOT NULL
	)`)

	return err
}

func InsertRate(timestamp int64, bestAsk, bestBid string) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO rates (timestamp, best_ask, best_bid) VALUES (?, ?, ?)`,
		timestamp, bestAsk, bestBid,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
