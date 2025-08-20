package service

import (
	"context"
	"github.com/Rastaiha/bermudia/internal/domain"
)

type Player struct {
	playerStore    domain.PlayerStore
	territoryStore domain.TerritoryStore
}

func NewPlayer(playerStore domain.PlayerStore, territoryStore domain.TerritoryStore) *Player {
	return &Player{playerStore: playerStore, territoryStore: territoryStore}
}

func (p *Player) GetPlayer(ctx context.Context, user *domain.User) (domain.Player, error) {
	return p.playerStore.Get(ctx, user.ID)
}

func (p *Player) Travel(ctx context.Context, user *domain.User, fromIsland string, toIsland string) error {
	player, err := p.playerStore.Get(ctx, user.ID)
	if err != nil {
		return err
	}
	territory, err := p.territoryStore.GetTerritoryByID(ctx, player.AtTerritory)
	if err != nil {
		return err
	}

	updatedPlayer := player
	if err := domain.Travel(&updatedPlayer, fromIsland, toIsland, territory); err != nil {
		return err
	}
	if err := p.playerStore.Update(ctx, player, updatedPlayer); err != nil {
		return err
	}

	return nil
}
