package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

func ConnectToSqlite() (*sql.DB, error) {
	p := filepath.Join(os.TempDir(), "sqlite3.db")
	db, err := sql.Open("sqlite3", p)
	return db, err
}
