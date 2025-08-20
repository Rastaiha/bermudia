package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

func ConnectToSqlite() (*sql.DB, error) {
	p := filepath.Join(os.TempDir(), "bermudia_sqlite.db")
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	db, err := sql.Open("sqlite3", p)
	return db, err
}

type scannable interface {
	Scan(dest ...any) error
}
