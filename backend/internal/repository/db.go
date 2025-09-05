package repository

import (
	"database/sql"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func ConnectToPostgres(cfg config.Postgres) (*sql.DB, error) {
	url := fmt.Sprintf("user=%s password=%s host=%s port=%d database=%s sslmode=%s",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.SSLMode,
	)
	db, err := sql.Open("pgx", url)
	return db, err
}

type scannable interface {
	Scan(dest ...any) error
}

func n[T any](v T) sql.Null[T] {
	return sql.Null[T]{Valid: !reflect.ValueOf(v).IsZero(), V: v}
}
