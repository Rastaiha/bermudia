package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
)

const playersSchema = `
CREATE TABLE IF NOT EXISTS players (
	user_id INT4 PRIMARY KEY,
	at_territory VARCHAR(255) NOT NULL,
	at_island VARCHAR(255) NOT NULL,
	fuel INT4 NOT NULL,
    fuel_cap INT4 NOT NULL
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
		`INSERT INTO players (user_id, at_territory, at_island, fuel, fuel_cap) VALUES ($1, $2, $3, $4, $5)`,
		n(player.UserId), n(player.AtTerritory), n(player.AtIsland), n(player.Fuel), n(player.FuelCap),
	)
	return err
}

func (s sqlPlayerRepository) Get(ctx context.Context, userId int32) (domain.Player, error) {
	var p domain.Player
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id, at_territory, at_island, fuel, fuel_cap FROM players WHERE user_id = $1`,
		userId,
	).Scan(&p.UserId, &p.AtTerritory, &p.AtIsland, &p.Fuel, &p.FuelCap)

	if err != nil {
		return domain.Player{}, fmt.Errorf("failed to get player from db: %w", err)
	}
	return p, nil
}

// Update updates a player row if and only if all fields match "old".
// UserId is never updated.
func (s sqlPlayerRepository) Update(ctx context.Context, old, updated domain.Player) error {
	cmd, err := s.db.ExecContext(ctx,
		`UPDATE players
		 SET at_territory = $1, at_island = $2, fuel = $3, fuel_cap = $4
		 WHERE user_id = $5
		   AND at_territory = $6
		   AND at_island = $7
		   AND fuel = $8
		   AND fuel_cap = $9`,
		n(updated.AtTerritory), n(updated.AtIsland), updated.Fuel, n(updated.FuelCap),
		old.UserId, old.AtTerritory, old.AtIsland, old.Fuel, old.FuelCap,
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
