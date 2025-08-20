package domain

import (
	"context"
	"errors"
)

var (
	ErrIslandNotFound    = errors.New("island not found")
	ErrTerritoryNotFound = errors.New("territory not found")
	ErrUserNotFound      = errors.New("user not found")
)

type TerritoryStore interface {
	CreateTerritory(ctx context.Context, territory *Territory) error
	GetTerritoryByID(ctx context.Context, territoryID string) (*Territory, error)
	ListTerritories(ctx context.Context) ([]Territory, error)
}

type IslandStore interface {
	SetContent(ctx context.Context, id string, content *IslandContent) error
	ReserveIDForTerritory(ctx context.Context, territoryId, islandId string) error
	GetByID(ctx context.Context, id string) (*IslandContent, error)
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id int32) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
