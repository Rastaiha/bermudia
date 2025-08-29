package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"time"
)

const playersSchema = `
CREATE TABLE IF NOT EXISTS players (
	user_id INT4 PRIMARY KEY,
	at_territory VARCHAR(255) NOT NULL,
	at_island VARCHAR(255) NOT NULL,
    anchored BOOLEAN NOT NULL,
	fuel INT4 NOT NULL,
    fuel_cap INT4 NOT NULL,
    coins INT4 NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
`

type sqlPlayerRepository struct {
	db *sql.DB
}

func NewSqlPlayerRepository(db *sql.DB) (domain.PlayerStore, error) {
	_, err := db.Exec(playersSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create players table: %w", err)
	}
	return sqlPlayerRepository{db: db}, nil
}

func (s sqlPlayerRepository) Create(ctx context.Context, player domain.Player) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO players (user_id, at_territory, at_island, anchored, fuel, fuel_cap, coins, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		n(player.UserId), n(player.AtTerritory), n(player.AtIsland), player.Anchored, n(player.Fuel), n(player.FuelCap), player.Coins, time.Now().UTC(),
	)
	return err
}

func (s sqlPlayerRepository) Get(ctx context.Context, userId int32) (domain.Player, error) {
	var p domain.Player
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id, at_territory, at_island, anchored, fuel, fuel_cap, coins, updated_at FROM players WHERE user_id = $1`,
		userId,
	).Scan(&p.UserId, &p.AtTerritory, &p.AtIsland, &p.Anchored, &p.Fuel, &p.FuelCap, &p.Coins, &p.UpdatedAt)

	if err != nil {
		return domain.Player{}, fmt.Errorf("failed to get player from db: %w", err)
	}
	return p, nil
}

// Update updates a player row if and only if all fields match "old".
// UserId is never updated.
func (s sqlPlayerRepository) Update(ctx context.Context, old, updated domain.Player) error {
	updated.UpdatedAt = time.Now().UTC()
	cmd, err := s.db.ExecContext(ctx,
		`UPDATE players
		 SET at_territory = $1, at_island = $2, anchored = $3, fuel = $4, fuel_cap = $5, coins = $6, updated_at = $7
		 WHERE user_id = $8 AND updated_at = $9`,
		n(updated.AtTerritory), n(updated.AtIsland), updated.Anchored, updated.Fuel, n(updated.FuelCap), updated.Coins, updated.UpdatedAt,
		old.UserId, old.UpdatedAt,
	)
	if err != nil {
		return err
	}
	rows, err := cmd.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// nothing updated -> either player not found, or old didn’t match
		return domain.ErrPlayerConflict
	}
	return nil
}
