package service

import (
	"context"
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
)

type Player struct {
	playerStore              domain.PlayerStore
	territoryStore           domain.TerritoryStore
	questionStore            domain.QuestionStore
	playerUpdateEventHandler func(event *domain.FullPlayerUpdateEvent)
}

func NewPlayer(playerStore domain.PlayerStore, territoryStore domain.TerritoryStore, questionStore domain.QuestionStore) *Player {
	return &Player{
		playerStore:    playerStore,
		territoryStore: territoryStore,
		questionStore:  questionStore,
	}
}

func (p *Player) GetPlayer(ctx context.Context, user *domain.User) (domain.FullPlayer, error) {
	player, err := p.playerStore.Get(ctx, user.ID)
	if err != nil {
		return domain.FullPlayer{}, err
	}
	return p.getFullPlayer(ctx, player)
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
	p.sendPlayerUpdateEvent(ctx, event)

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
	p.sendPlayerUpdateEvent(ctx, event)

	return nil
}

func (p *Player) OnPlayerUpdate(eventHandler func(event *domain.FullPlayerUpdateEvent)) {
	p.playerUpdateEventHandler = eventHandler
}

func (p *Player) sendPlayerUpdateEvent(ctx context.Context, event *domain.PlayerUpdateEvent) {
	if p.playerUpdateEventHandler != nil {
		fullPlayer, err := p.getFullPlayer(ctx, *event.Player)
		if err != nil {
			slog.Error("failed to get the knowledge bars, missing event", err, "userId", event.Player.UserId)
			return
		}
		p.playerUpdateEventHandler(&domain.FullPlayerUpdateEvent{
			Reason: event.Reason,
			Player: &fullPlayer,
		})
	}
}

func (p *Player) getFullPlayer(ctx context.Context, player domain.Player) (domain.FullPlayer, error) {
	knowledgeBars, err := p.questionStore.GetKnowledgeBars(ctx, player.UserId)
	if err != nil {
		return domain.FullPlayer{}, err
	}
	return domain.FullPlayer{
		Player:        player,
		KnowledgeBars: knowledgeBars,
	}, nil
}
