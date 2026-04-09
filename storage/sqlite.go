package storage

import (
	"database/sql"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() (*SQLiteRepository, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(filepath.Dir(exe), "rates.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

func runMigrations(db *sql.DB) error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := fs.ReadFile(migrationsFS, "migrations/"+entry.Name())
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(data)); err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) InsertRate(timestamp int64, bestAsk, bestBid string) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO rates (timestamp, best_ask, best_bid) VALUES (?, ?, ?)`,
		timestamp, bestAsk, bestBid,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
