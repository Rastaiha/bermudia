package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rastaiha/bermudia/internal/domain"
)

const (
	territorySchema = `
CREATE TABLE IF NOT EXISTS territories (
	id VARCHAR(255) PRIMARY KEY,
    start_island VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
`
)

type sqlTerritoryRepository struct {
	db *sql.DB
}

func NewSqlTerritoryRepository(db *sql.DB) (domain.TerritoryStore, error) {
	_, err := db.Exec(territorySchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create territories table: %w", err)
	}
	return sqlTerritoryRepository{
		db: db,
	}, nil
}

type territoryContent struct {
	Name            string                `json:"name"`
	BackgroundAsset string                `json:"backgroundAsset"`
	Islands         []domain.Island       `json:"islands"`
	Edges           []domain.Edge         `json:"edges"`
	RefuelIslands   []domain.RefuelIsland `json:"refuelIslands"`
}

func (s sqlTerritoryRepository) columns() string {
	return "SELECT id, start_island, content FROM territories"
}

func (s sqlTerritoryRepository) scan(row scannable, territory *domain.Territory) error {
	var content []byte
	err := row.Scan(&territory.ID, &territory.StartIsland, &content)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrTerritoryNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to get territory from db: %w", err)
	}
	var tc territoryContent
	err = json.Unmarshal(content, &tc)
	if err != nil {
		return err
	}
	territory.Name = tc.Name
	territory.BackgroundAsset = tc.BackgroundAsset
	territory.Islands = tc.Islands
	territory.Edges = tc.Edges
	territory.RefuelIslands = tc.RefuelIslands
	return nil
}

func (s sqlTerritoryRepository) CreateTerritory(ctx context.Context, territory *domain.Territory) error {
	tc := territoryContent{
		Name:            territory.Name,
		BackgroundAsset: territory.BackgroundAsset,
		Islands:         territory.Islands,
		Edges:           territory.Edges,
		RefuelIslands:   territory.RefuelIslands,
	}
	content, err := json.Marshal(tc)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO territories (id, start_island, content, updated_at) VALUES ($1, $2, $3, $4)`, n(territory.ID), n(territory.StartIsland), content, time.Now().UTC())
	return err
}

func (s sqlTerritoryRepository) GetTerritoryByID(ctx context.Context, territoryID string) (*domain.Territory, error) {
	var t domain.Territory
	err := s.scan(s.db.QueryRowContext(ctx, s.columns()+" WHERE id = $1", territoryID), &t)
	return &t, err
}

// ListTerritories returns all available territories from embedded files
func (s sqlTerritoryRepository) ListTerritories(ctx context.Context) (result []domain.Territory, err error) {
	rows, err := s.db.QueryContext(ctx, s.columns())
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		err = errors.Join(err, closeErr)
	}()
	for rows.Next() {
		var t domain.Territory
		if err := s.scan(rows, &t); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}
