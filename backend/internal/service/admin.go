package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"slices"
)

type Admin struct {
	territoryStore domain.TerritoryStore
	islandStore    domain.IslandStore
	userStore      domain.UserStore
	playerStore    domain.PlayerStore
}

func NewAdmin(territoryStore domain.TerritoryStore, islandStore domain.IslandStore, userStore domain.UserStore, playerStore domain.PlayerStore) *Admin {
	return &Admin{
		territoryStore: territoryStore,
		islandStore:    islandStore,
		userStore:      userStore,
		playerStore:    playerStore,
	}
}

func (a *Admin) SetTerritory(ctx context.Context, territory domain.Territory) error {
	for _, island := range territory.Islands {
		if island.ID == "" {
			return fmt.Errorf("empty island id in island list")
		}
	}
	if territory.StartIsland == "" {
		return errors.New("invalid territory startIsland")
	}
	if !slices.ContainsFunc(territory.Islands, func(island domain.Island) bool {
		return island.ID == territory.StartIsland
	}) {
		return fmt.Errorf("startIsland %q not found in island list", territory.StartIsland)
	}
	for _, e := range territory.Edges {
		if e.From == "" || e.To == "" {
			return fmt.Errorf("empty edge.from or edge.to: %v", e)
		}
		if !slices.ContainsFunc(territory.Islands, func(island domain.Island) bool {
			return island.ID == e.From
		}) {
			return fmt.Errorf("edge.from %q is not in island list", e.From)
		}
		if !slices.ContainsFunc(territory.Islands, func(island domain.Island) bool {
			return island.ID == e.To
		}) {
			return fmt.Errorf("edge.to %q is not in island list", e.To)
		}
	}

	for _, island := range territory.Islands {
		if err := a.islandStore.ReserveIDForTerritory(ctx, territory.ID, island.ID); err != nil {
			return err
		}
	}

	return a.territoryStore.CreateTerritory(ctx, &territory)
}

func (a *Admin) SetIsland(ctx context.Context, id string, islandContent domain.IslandContent) error {
	return a.islandStore.SetContent(ctx, id, &islandContent)
}

func (a *Admin) CreateUser(ctx context.Context, id int32, username, password string) error {
	territories, err := a.territoryStore.ListTerritories(ctx)
	if err != nil {
		return err
	}
	if len(territories) == 0 {
		return errors.New("no territory found")
	}
	startingTerritory := territories[int(id)%len(territories)]
	hp, err := domain.HashPassword(password)
	if err != nil {
		return err
	}
	user := domain.User{
		ID:             id,
		Username:       username,
		HashedPassword: hp,
	}
	if err := a.userStore.Create(ctx, &user); err != nil {
		return err
	}
	return a.playerStore.Create(ctx, domain.NewPlayer(user.ID, &startingTerritory))
}
