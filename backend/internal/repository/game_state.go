package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"strconv"
)

const gameStateSchema = `
CREATE TABLE IF NOT EXISTS game_state (
	key VARCHAR(255) PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default values if they don't exist
INSERT INTO game_state (key, value) VALUES ('is_paused', 'false')
ON CONFLICT (key) DO NOTHING;
`

const (
	gameStateKeyIsPaused = "is_paused"
)

type sqlGameStateRepository struct {
	db *sql.DB
}

func NewSqlGameStateRepository(db *sql.DB) (domain.GameStateStore, error) {
	_, err := db.Exec(gameStateSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create game_state table: %w", err)
	}
	return sqlGameStateRepository{db: db}, nil
}

func (s sqlGameStateRepository) GetIsPaused(ctx context.Context) (bool, error) {
	var value string
	err := s.db.QueryRowContext(ctx,
		`SELECT value FROM game_state WHERE key = $1`,
		gameStateKeyIsPaused,
	).Scan(&value)

	if err != nil {
		return false, fmt.Errorf("failed to get is_paused from db: %w", err)
	}

	isPaused, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("failed to parse is_paused value: %w", err)
	}

	return isPaused, nil
}

func (s sqlGameStateRepository) SetIsPaused(ctx context.Context, isPaused bool) error {
	value := strconv.FormatBool(isPaused)

	_, err := s.db.ExecContext(ctx,
		`UPDATE game_state SET value = $1, updated_at = CURRENT_TIMESTAMP WHERE key = $2`,
		value, gameStateKeyIsPaused,
	)

	if err != nil {
		return fmt.Errorf("failed to set is_paused in db: %w", err)
	}

	return nil
}
