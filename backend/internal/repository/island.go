package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
)

const (
	islandsSchema = `
CREATE TABLE IF NOT EXISTS islands (
    id VARCHAR(255) PRIMARY KEY,
    territory_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL
);
`
)

type sqlIslandRepository struct {
	db *sql.DB
}

func NewSqlIslandRepository(db *sql.DB) (domain.IslandStore, error) {
	_, err := db.Exec(islandsSchema)
	if err != nil {
		return nil, fmt.Errorf("create islands table: %w", err)
	}
	return sqlIslandRepository{
		db: db,
	}, nil
}

func (s sqlIslandRepository) SetContent(ctx context.Context, id string, content *domain.IslandContent) error {
	c, err := json.Marshal(content)
	if err != nil {
		return err
	}
	cmd, err := s.db.ExecContext(ctx, `UPDATE islands SET content = $1 WHERE id = $2`, c, id)
	if err != nil {
		return err
	}
	affected, err := cmd.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrIslandNotFound
	}
	return nil
}

func (s sqlIslandRepository) ReserveIDForTerritory(ctx context.Context, territoryId, islandId string) error {
	var actualTerritoryId string
	err := s.db.QueryRowContext(ctx, `INSERT INTO islands (id, territory_id, content) VALUES ($1, $2, $3) ON CONFLICT DO UPDATE SET id = EXCLUDED.id RETURNING territory_id ;`, islandId, territoryId, []byte("{}")).Scan(&actualTerritoryId)
	if err != nil {
		return err
	}
	if actualTerritoryId != territoryId {
		return fmt.Errorf("island_id %q is already taken by territory %q", islandId, actualTerritoryId)
	}
	return nil
}

func (s sqlIslandRepository) GetByID(ctx context.Context, id string) (*domain.IslandContent, error) {
	var content []byte
	err := s.db.QueryRowContext(ctx, `SELECT content FROM islands WHERE id = $1`, id).Scan(&content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrIslandNotFound
	}
	if err != nil {
		return nil, err
	}
	var result domain.IslandContent
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s sqlIslandRepository) GetTerritory(ctx context.Context, id string) (string, error) {
	var territoryId string
	err := s.db.QueryRowContext(ctx, `SELECT territory_id FROM islands WHERE id = $1`, id).Scan(&territoryId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrIslandNotFound
	}
	return territoryId, err
}
