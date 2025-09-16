package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
    coin INT4 NOT NULL,
    red_key INT4 NOT NULL,
    blue_key INT4 NOT NULL,
    golden_key INT4 NOT NULL,
    master_key INT4 NOT NULL,
    visited_territories TEXT NOT NULL,
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

func (s sqlPlayerRepository) Create(ctx context.Context, player domain.Player) (err error) {
	initialUpdatedAt := time.Time{}

	current, err := s.Get(ctx, player.UserId)
	if err == nil && current.UpdatedAt.UTC().UnixMilli() == initialUpdatedAt.UTC().UnixMilli() {
		player.UpdatedAt = initialUpdatedAt
		err = s.update(ctx, nil, current, player)
		if errors.Is(err, domain.ErrPlayerConflict) {
			return nil
		}
		return err
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	visitedTerritories, err := json.Marshal(player.VisitedTerritories)
	if err != nil {
		return fmt.Errorf("failed to marshal visited territories: %w", err)
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO players (user_id, at_territory, at_island, anchored, fuel, fuel_cap, coin, red_key, blue_key, golden_key, master_key, visited_territories, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT DO NOTHING ;`,
		n(player.UserId), n(player.AtTerritory), n(player.AtIsland), player.Anchored, n(player.Fuel), n(player.FuelCap), player.Coin, player.RedKey, player.BlueKey, player.GoldenKey, player.MasterKey, visitedTerritories, initialUpdatedAt,
	)
	return err
}

func (s sqlPlayerRepository) Get(ctx context.Context, userId int32) (domain.Player, error) {
	var visitedTerritories []byte
	var p domain.Player
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id, at_territory, at_island, anchored, fuel, fuel_cap, coin, red_key, blue_key, golden_key, master_key, visited_territories, updated_at FROM players WHERE user_id = $1`,
		userId,
	).Scan(&p.UserId, &p.AtTerritory, &p.AtIsland, &p.Anchored, &p.Fuel, &p.FuelCap, &p.Coin, &p.RedKey, &p.BlueKey, &p.GoldenKey, &p.MasterKey, &visitedTerritories, &p.UpdatedAt)

	if err != nil {
		return domain.Player{}, fmt.Errorf("failed to get player from db: %w", err)
	}
	if err := json.Unmarshal(visitedTerritories, &p.VisitedTerritories); err != nil {
		return domain.Player{}, fmt.Errorf("failed to unmarshal visited territories: %w", err)
	}
	return p, nil
}

func (s sqlPlayerRepository) Update(ctx context.Context, tx domain.Tx, old, updated domain.Player) error {
	updated.UpdatedAt = time.Now().UTC()
	return s.update(ctx, tx, old, updated)
}

// Update updates a player row if and only if all fields match "old".
// UserId is never updated.
func (s sqlPlayerRepository) update(ctx context.Context, tx domain.Tx, old, updated domain.Player) error {
	if tx == nil {
		tx = s.db
	}
	visitedTerritories, err := json.Marshal(updated.VisitedTerritories)
	if err != nil {
		return fmt.Errorf("failed to marshal visited territories: %w", err)
	}
	cmd, err := tx.ExecContext(ctx,
		`UPDATE players
		 SET at_territory = $1, at_island = $2, anchored = $3, fuel = $4, fuel_cap = $5, coin = $6, red_key = $7, blue_key = $8, golden_key = $9, master_key = $10, visited_territories = $11, updated_at = $12
		 WHERE user_id = $13 AND updated_at = $14`,
		n(updated.AtTerritory), n(updated.AtIsland), updated.Anchored, updated.Fuel, n(updated.FuelCap), updated.Coin, updated.RedKey, updated.BlueKey, updated.GoldenKey, updated.MasterKey, visitedTerritories, updated.UpdatedAt,
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
