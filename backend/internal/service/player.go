package service

import (
	"context"
	"fmt"
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
	islandStore              domain.IslandStore
	playerUpdateEventHandler func(event *domain.FullPlayerUpdateEvent)
	cron                     gocron.Scheduler
}

func NewPlayer(cfg config.Config, playerStore domain.PlayerStore, territoryStore domain.TerritoryStore, questionStore domain.QuestionStore, islandStore domain.IslandStore) *Player {
	return &Player{
		cfg:            cfg,
		playerStore:    playerStore,
		territoryStore: territoryStore,
		questionStore:  questionStore,
		islandStore:    islandStore,
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
	isDestinationIslandUnlocked, err := p.isIslandUnlocked(ctx, user.ID, *territory, toIsland)
	if err != nil {
		return nil, err
	}

	checkResult := domain.TravelCheck(player, fromIsland, toIsland, territory, isDestinationIslandUnlocked)
	return &checkResult, nil
}

func (p *Player) isIslandUnlocked(ctx context.Context, userId int32, territory domain.Territory, islandId string) (bool, error) {
	prerequisites := territory.IslandPrerequisites[islandId]
	for _, pre := range prerequisites {
		answers, questionsCount, err := p.islandStore.GetUserAnswerComponents(ctx, userId, pre)
		if err != nil {
			return false, err
		}
		if len(answers) < questionsCount {
			return false, nil
		}
		for _, answerId := range answers {
			answer, err := p.questionStore.GetAnswer(ctx, answerId)
			if err != nil {
				return false, fmt.Errorf("failed to get answer while checking isIslandUnlocked: %w", err)
			}
			if answer.Status != domain.AnswerStatusCorrect {
				return false, nil
			}
		}
	}
	return true, nil
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
	isDestinationIslandUnlocked, err := p.isIslandUnlocked(ctx, user.ID, *territory, toIsland)
	if err != nil {
		return err
	}

	event, err := domain.Travel(player, fromIsland, toIsland, territory, isDestinationIslandUnlocked)
	if err != nil {
		return err
	}
	return p.applyAndSendPlayerUpdateEvent(ctx, player, event)
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
	return p.applyAndSendPlayerUpdateEvent(ctx, player, event)
}

func (p *Player) AnchorCheck(ctx context.Context, userId int32, islandID string) (*domain.AnchorCheckResult, error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	checkResult := domain.AnchorCheck(player, islandID)
	return &checkResult, nil
}

func (p *Player) Anchor(ctx context.Context, userId int32, islandID string) error {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	event, err := domain.Anchor(player, islandID)
	if err != nil {
		return err
	}
	return p.applyAndSendPlayerUpdateEvent(ctx, player, event)
}

func (p *Player) OnPlayerUpdate(eventHandler func(event *domain.FullPlayerUpdateEvent)) {
	p.playerUpdateEventHandler = eventHandler
}

func (p *Player) applyAndSendPlayerUpdateEvent(ctx context.Context, oldPlayer domain.Player, event *domain.PlayerUpdateEvent) error {
	if err := p.playerStore.Update(ctx, oldPlayer, *event.Player); err != nil {
		return err
	}
	if err := p.sendPlayerUpdateEventErr(ctx, event); err != nil {
		slog.Error("failed to send event",
			slog.String("error", err.Error()),
			slog.Int("userId", int(event.Player.UserId)),
		)
	}
	return nil
}

func (p *Player) sendPlayerUpdateEventErr(ctx context.Context, event *domain.PlayerUpdateEvent) error {
	if p.playerUpdateEventHandler != nil {
		fullPlayer, err := p.getFullPlayer(ctx, *event.Player)
		if err != nil {
			return err
		}
		p.playerUpdateEventHandler(&domain.FullPlayerUpdateEvent{
			Reason: event.Reason,
			Player: &fullPlayer,
		})
	}
	return nil
}

func (p *Player) SendInitialEvents(ctx context.Context, userId int32) error {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	return p.sendPlayerUpdateEventErr(ctx, &domain.PlayerUpdateEvent{
		Reason: domain.PlayerUpdateEventInitial,
		Player: &player,
	})
}

func (p *Player) getFullPlayer(ctx context.Context, player domain.Player) (domain.FullPlayer, error) {
	knowledgeBars, err := p.questionStore.GetKnowledgeBars(ctx, player.UserId)
	if err != nil {
		return domain.FullPlayer{}, fmt.Errorf("failed to get knowledge bars: %w", err)
	}
	return domain.FullPlayer{
		Player:        player,
		KnowledgeBars: knowledgeBars,
	}, nil
}

func (p *Player) applyCorrections(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

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
			err = p.sendPlayerUpdateEventErr(ctx, &domain.PlayerUpdateEvent{
				Reason: domain.PlayerUpdateEventCorrection,
				Player: &player,
			})
			if err != nil {
				slog.Error("failed to sendPlayerUpdateEventErr after applying correction", slog.String("error", err.Error()))
			}
		}()
	}
	wg.Wait()

	if appliedCorrections.Load() > 0 {
		slog.Info("successfully applied corrections", slog.Int64("count", appliedCorrections.Load()))
	}
}

func (p *Player) MigrateCheck(ctx context.Context, userId int32) (*domain.MigrateCheckResult, error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	territories, err := p.territoryStore.ListTerritories(ctx)
	if err != nil {
		return nil, err
	}
	currentTerritory, err := p.territoryStore.GetTerritoryByID(ctx, player.AtTerritory)
	if err != nil {
		return nil, err
	}
	knowledgeBars, err := p.questionStore.GetKnowledgeBars(ctx, userId)
	if err != nil {
		return nil, err
	}

	check := domain.MigrateCheck(player, knowledgeBars, *currentTerritory, territories)
	return &check, nil
}

func (p *Player) Migrate(ctx context.Context, userId int32, toTerritory string) error {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	territories, err := p.territoryStore.ListTerritories(ctx)
	if err != nil {
		return err
	}
	currentTerritory, err := p.territoryStore.GetTerritoryByID(ctx, player.AtTerritory)
	if err != nil {
		return err
	}
	knowledgeBars, err := p.questionStore.GetKnowledgeBars(ctx, userId)
	if err != nil {
		return err
	}

	event, err := domain.Migrate(player, knowledgeBars, *currentTerritory, territories, toTerritory)
	if err != nil {
		return err
	}
	return p.applyAndSendPlayerUpdateEvent(ctx, player, event)
}
