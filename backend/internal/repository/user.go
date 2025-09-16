package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"strings"
)

const (
	userSchema = `
CREATE TABLE IF NOT EXISTS users (
    id INT4 PRIMARY KEY,
    username_display VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    meet_link VARCHAR(255) NOT NULL,
    hashed_password BYTEA NOT NULL
);`
)

type sqlUser struct {
	db *sql.DB
}

func NewSqlUser(db *sql.DB) (domain.UserStore, error) {
	_, err := db.Exec(userSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}
	return sqlUser{
		db: db,
	}, nil
}

func (s sqlUser) columns() string {
	return "SELECT id, username_display, meet_link, hashed_password FROM users"
}

func (s sqlUser) scan(row *sql.Row, user *domain.User) error {
	err := row.Scan(&user.ID, &user.Username, &user.MeetLink, &user.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	return err
}

func (s sqlUser) Create(ctx context.Context, user *domain.User) error {
	err := s.db.QueryRowContext(ctx, `INSERT INTO users (id, username_display, username, meet_link, hashed_password) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (username) DO UPDATE SET username_display = $2, username = $3, meet_link = $4, hashed_password = $5 RETURNING id`,
		n(user.ID),
		n(user.Username),
		n(strings.ToLower(user.Username)),
		user.MeetLink, // TODO: wrap in n
		user.HashedPassword,
	).Scan(&user.ID)
	return err
}

func (s sqlUser) Get(ctx context.Context, id int32) (*domain.User, error) {
	var result domain.User
	err := s.scan(s.db.QueryRowContext(ctx, s.columns()+" WHERE id = $1", id), &result)
	return &result, err
}

func (s sqlUser) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var result domain.User
	err := s.scan(s.db.QueryRowContext(ctx, s.columns()+" WHERE username = $1", strings.ToLower(username)), &result)
	return &result, err
}
