package service

import (
	"context"
	"github.com/Rastaiha/bermudia/internal/domain"
)

type Player struct {
	playerStore              domain.PlayerStore
	territoryStore           domain.TerritoryStore
	playerUpdateEventHandler func(event *domain.PlayerUpdateEvent)
}

func NewPlayer(playerStore domain.PlayerStore, territoryStore domain.TerritoryStore) *Player {
	return &Player{playerStore: playerStore, territoryStore: territoryStore}
}

func (p *Player) GetPlayer(ctx context.Context, user *domain.User) (domain.Player, error) {
	return p.playerStore.Get(ctx, user.ID)
}

func (p *Player) TravelCheck(ctx context.Context, user *domain.User, fromIsland, toIsland string) (*domain.TravelCheckResult, error) {
	player, err := p.playerStore.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	territory, err := p.territoryStore.GetTerritoryByID(ctx, player.AtTerritory)
	if err != nil {
		return nil, err
	}

	checkResult := domain.TravelCheck(player, fromIsland, toIsland, territory)
	return &checkResult, nil
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

	event, err := domain.Travel(player, fromIsland, toIsland, territory)
	if err != nil {
		return err
	}
	if err := p.playerStore.Update(ctx, player, *event.Player); err != nil {
		return err
	}
	p.sendPlayerUpdateEvent(event)

	return nil
}

func (p *Player) RefuelCheck(ctx context.Context, userId int32) (*domain.RefuelCheckResult, error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	territory, err := p.territoryStore.GetTerritoryByID(ctx, player.AtTerritory)
	if err != nil {
		return nil, err
	}

	checkResult := domain.RefuelCheck(player, territory)
	return &checkResult, nil
}

func (p *Player) Refuel(ctx context.Context, userId int32, amount int32) error {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	territory, err := p.territoryStore.GetTerritoryByID(ctx, player.AtTerritory)
	if err != nil {
		return err
	}

	event, err := domain.Refuel(player, territory, amount)
	if err != nil {
		return err
	}
	if err := p.playerStore.Update(ctx, player, *event.Player); err != nil {
		return err
	}
	p.sendPlayerUpdateEvent(event)

	return nil
}

func (p *Player) OnPlayerUpdate(eventHandler func(event *domain.PlayerUpdateEvent)) {
	p.playerUpdateEventHandler = eventHandler
}

func (p *Player) sendPlayerUpdateEvent(event *domain.PlayerUpdateEvent) {
	if p.playerUpdateEventHandler != nil {
		p.playerUpdateEventHandler(event)
	}
}
