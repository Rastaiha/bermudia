package domain

import (
	"context"
	"errors"
)

var (
	ErrIslandNotFound    = errors.New("island not found")
	ErrTerritoryNotFound = errors.New("territory not found")
)

type IslandStore interface {
	GetByID(ctx context.Context, id string) (*IslandContent, error)
}

type TerritoryStore interface {
	GetTerritoryByID(ctx context.Context, territoryID string) (*Territory, error)
	ListTerritories(ctx context.Context) ([]Territory, error)
}
