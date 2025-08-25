package service

import (
	"context"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Player struct {
	cfg                      config.Config
	playerStore              domain.PlayerStore
	territoryStore           domain.TerritoryStore
	questionStore            domain.QuestionStore
	playerUpdateEventHandler func(event *domain.FullPlayerUpdateEvent)
	cron                     gocron.Scheduler
}

func NewPlayer(cfg config.Config, playerStore domain.PlayerStore, territoryStore domain.TerritoryStore, questionStore domain.QuestionStore) *Player {
	return &Player{
		cfg:            cfg,
		playerStore:    playerStore,
		territoryStore: territoryStore,
		questionStore:  questionStore,
	}
}

func (p *Player) Start() {
	var err error
	p.cron, err = gocron.NewScheduler(gocron.WithLimitConcurrentJobs(1, gocron.LimitModeReschedule))
	if err != nil {
		panic(err)
	}
	_, err = p.cron.NewJob(gocron.DurationJob(p.cfg.CorrectionJobInterval), gocron.NewTask(p.applyCorrections))
	if err != nil {
		panic(err)
	}
	p.cron.Start()
}

func (p *Player) Stop() {
	if err := p.cron.Shutdown(); err != nil {
		slog.Error("failed to stop cron", err)
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

func (p *Player) applyCorrections(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	slog.Info("applying corrections...")

	corrections, err := p.questionStore.GetUnappliedCorrections(ctx)
	if err != nil {
		slog.Error("failed to GetUnappliedCorrections from db", err)
		return
	}

	var wg sync.WaitGroup
	workerLimit := make(chan struct{}, 2)

	appliedCorrections := &atomic.Int64{}
	updatedUserIDs := make(map[int32]struct{})
	lock := sync.Mutex{}
	for _, c := range corrections {
		workerLimit <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-workerLimit
				wg.Done()
			}()
			userId, ok, err := p.questionStore.ApplyCorrection(ctx, c, time.Now().Add(-p.cfg.MinCorrectionDelay).UTC())
			if err != nil {
				slog.Error("failed to ApplyCorrection from db", err)
				return
			}
			if ok {
				appliedCorrections.Add(1)
			}
			if ok && c.IsCorrect {
				lock.Lock()
				defer lock.Unlock()
				updatedUserIDs[userId] = struct{}{}
			}
		}()
	}
	wg.Wait()

	for userId := range updatedUserIDs {
		workerLimit <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-workerLimit
				wg.Done()
			}()
			player, err := p.playerStore.Get(ctx, userId)
			if err != nil {
				slog.Error("failed to Get player from db", err)
				return
			}
			p.sendPlayerUpdateEvent(ctx, &domain.PlayerUpdateEvent{
				Reason: domain.PlayerUpdateEventCorrection,
				Player: &player,
			})
		}()
	}
	wg.Wait()

	slog.Info("successfully applied corrections", slog.Int64("count", appliedCorrections.Load()))
}
