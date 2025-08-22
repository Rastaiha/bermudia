package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"reflect"
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

func n[T any](v T) sql.Null[T] {
	return sql.Null[T]{Valid: !reflect.ValueOf(v).IsZero(), V: v}
}
