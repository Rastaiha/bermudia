package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-co-op/gocron/v2"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Player struct {
	cfg                        config.Config
	db                         *sql.DB
	userStore                  domain.UserStore
	playerStore                domain.PlayerStore
	territoryStore             domain.TerritoryStore
	questionStore              domain.QuestionStore
	islandStore                domain.IslandStore
	treasureStore              domain.TreasureStore
	marketStore                domain.MarketStore
	inboxStore                 domain.InboxStore
	investStore                domain.InvestStore
	playerUpdateEventHandler   func(event *domain.FullPlayerUpdateEvent)
	tradeEventBroadcastHandler TradeEventBroadcastHandler
	inboxEventHandler          func(e *domain.InboxEvent)
	broadcastMessageHandler    MessageBroadcastHandler
	playerLocationsCache       *cache.Cache
	cron                       gocron.Scheduler
}

type TradeEventBroadcastHandler func(func(userId int32) *domain.TradeEvent)

type MessageBroadcastHandler func(func(userId int32) *domain.InboxMessageView)

func NewPlayer(cfg config.Config, db *sql.DB, userStore domain.UserStore, playerStore domain.PlayerStore, territoryStore domain.TerritoryStore, questionStore domain.QuestionStore, islandStore domain.IslandStore, treasureStore domain.TreasureStore, marketStore domain.MarketStore, inboxStore domain.InboxStore, investStore domain.InvestStore) *Player {
	return &Player{
		cfg:                  cfg,
		db:                   db,
		userStore:            userStore,
		playerStore:          playerStore,
		territoryStore:       territoryStore,
		questionStore:        questionStore,
		islandStore:          islandStore,
		treasureStore:        treasureStore,
		marketStore:          marketStore,
		inboxStore:           inboxStore,
		investStore:          investStore,
		playerLocationsCache: cache.New(20*time.Second, time.Minute),
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
	if p.cfg.DevMode {
		_, _ = p.cron.NewJob(gocron.DurationJob(1*time.Minute), gocron.NewTask(func() {
			go func() {
				now := time.Now().UTC()
				_, err := p.BroadcastMessage(
					context.Background(),
					fmt.Sprintf("ساعت %s به وقت لندن است.\nامروز %s", now.Format(time.TimeOnly), now.Weekday().String()),
				)
				if err != nil {
					slog.Error("failed to send mock broadcast message", "error", err)
				}
			}()
		}))
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
		hasAnsweredIsland, err := p.questionStore.HasAnsweredIsland(ctx, userId, pre)
		if err != nil {
			return false, fmt.Errorf("failed to check if user answered all island: %w", err)
		}
		if !hasAnsweredIsland {
			return false, nil
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
	if err := p.playerStore.Update(ctx, nil, oldPlayer, *event.Player); err != nil {
		return err
	}
	if err := p.sendPlayerUpdateEventErr(ctx, event); err != nil {
		return err
	}
	return nil
}

func (p *Player) sendPlayerUpdateEventErr(ctx context.Context, event *domain.PlayerUpdateEvent) error {
	fullPlayer, err := p.getFullPlayer(ctx, *event.Player)
	if err != nil {
		slog.Error("failed to send player update event",
			slog.String("error", err.Error()),
			slog.Int("userId", int(event.Player.UserId)),
		)
		return fmt.Errorf("failed to send player update event: %w", err)
	}

	if event.Reason != domain.PlayerUpdateEventInitial {
		err = p.playerStore.CreatePlayerEvent(ctx, event.Player.UserId, time.Now().UTC(), event.Reason, fullPlayer)
		if err != nil {
			slog.Error("failed to create player event",
				slog.String("error", err.Error()),
				slog.Int("userId", int(event.Player.UserId)),
			)
			return fmt.Errorf("failed to create player event: %w", err)
		}
	}

	p.playerUpdateEventHandler(&domain.FullPlayerUpdateEvent{
		Reason: event.Reason,
		Player: &fullPlayer,
	})
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
	portableIslands, err := p.islandStore.GetPortableIslands(ctx, player.UserId)
	if err != nil {
		return domain.FullPlayer{}, fmt.Errorf("failed to get portable islands: %w", err)
	}
	territories, err := p.territoryStore.ListTerritories(ctx)
	if err != nil {
		return domain.FullPlayer{}, fmt.Errorf("failed to get territories: %w", err)
	}
	books := make([]domain.FullPortableIsland, 0, len(portableIslands))
	for _, pi := range portableIslands {
		territoryName := ""
		for _, t := range territories {
			if t.ID == pi.TerritoryID {
				territoryName = t.Name
				break
			}
		}
		books = append(books, domain.FullPortableIsland{
			PortableIsland: pi,
			TerritoryName:  territoryName,
		})
	}
	knowledgeBars, err := p.questionStore.GetKnowledgeBars(ctx, player.UserId)
	if err != nil {
		return domain.FullPlayer{}, fmt.Errorf("failed to get knowledge bars: %w", err)
	}
	return domain.FullPlayer{
		Player:        player,
		KnowledgeBars: knowledgeBars,
		Books:         books,
	}, nil
}

func (p *Player) applyCorrections(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	corrections, err := p.questionStore.GetUnappliedCorrections(ctx, time.Now().Add(-p.cfg.CorrectionRevertWindow).UTC())
	if err != nil {
		slog.Error("failed to GetUnappliedCorrections from db", err)
		return
	}

	var wg sync.WaitGroup
	workerLimit := make(chan struct{}, 2)
	appliedCorrections := &atomic.Int64{}
	for _, c := range corrections {
		workerLimit <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-workerLimit
				wg.Done()
			}()
			ok, err := p.applyCorrection(ctx, c)
			if err != nil {
				slog.Error("failed to ApplyCorrection from db", err)
				return
			}
			if ok {
				appliedCorrections.Add(1)
			}
		}()
	}
	wg.Wait()

	if appliedCorrections.Load() > 0 {
		slog.Info("successfully applied corrections", slog.Int64("count", appliedCorrections.Load()))
	}
}

func (p *Player) applyCorrection(ctx context.Context, c domain.Correction) (ok bool, err error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var event *domain.PlayerUpdateEvent
	defer func() {
		if err != nil || event == nil {
			return
		}
		if err := p.sendPlayerUpdateEventErr(ctx, event); err != nil {
			err = fmt.Errorf("failed to send player update event: %w", err)
		}
	}()
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	var answer domain.Answer
	answer, ok, err = p.questionStore.ApplyCorrection(ctx, tx, time.Now().Add(-p.cfg.MinCorrectionDelay).UTC(), c)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	question, err := p.questionStore.GetQuestion(ctx, c.QuestionId)
	if err != nil {
		return false, err
	}

	currentPlayer, err := p.playerStore.Get(ctx, c.UserId)
	if err != nil {
		return false, err
	}

	pool, hasPool, err := p.islandStore.GetPoolOfBook(ctx, question.BookID)
	if err != nil {
		return false, err
	}

	event, reward, rewarded := domain.GetRewardOfCorrection(currentPlayer, question, c, pool, hasPool)
	if rewarded {
		if err := p.playerStore.Update(ctx, tx, currentPlayer, *event.Player); err != nil {
			return false, err
		}
	}

	islandHeader, err := p.islandStore.GetIslandHeaderByBookIdAndUserId(ctx, question.BookID, answer.UserID)
	if err != nil {
		return false, err
	}
	territory, err := p.territoryStore.GetTerritoryByID(ctx, islandHeader.TerritoryID)
	if err != nil {
		return false, err
	}

	err = p.createAndSendInboxMessage(ctx, tx, domain.InboxMessage{
		ID:        domain.NewID(domain.ResourceTypeInboxMessage),
		UserID:    c.UserId,
		CreatedAt: time.Now().UTC(),
		Content: domain.InboxMessageContent{
			NewCorrection: &domain.InboxMessageNewCorrection{
				TerritoryID:   territory.ID,
				TerritoryName: territory.Name,
				IslandID:      islandHeader.ID,
				IslandName:    islandHeader.Name,
				InputID:       answer.QuestionID,
				NewState:      domain.GetSubmissionState(question, answer),
				Reward:        reward,
			},
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
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

func (p *Player) UnlockTreasureCheck(ctx context.Context, userId int32, treasureId string) (*domain.UnlockTreasureCheckResult, error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	treasure, err := p.treasureStore.GetTreasure(ctx, treasureId)
	if err != nil {
		return nil, err
	}
	bookId, err := p.islandStore.GetBookOfIsland(ctx, player.AtIsland, userId)
	if errors.Is(err, domain.ErrNoBookAssignedFromPool) {
		err = nil
		bookId = ""
	} else if err != nil {
		return nil, err
	}
	userTreasure, err := p.treasureStore.GetUserTreasure(ctx, userId, treasureId)
	if err != nil {
		return nil, err
	}
	check := domain.UnlockTreasureCheck(player, treasure, userTreasure, bookId)
	return &check, nil
}

func (p *Player) UnlockTreasure(ctx context.Context, userId int32, treasureId string, chosenCost string) (*domain.IslandTreasure, error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	treasure, err := p.treasureStore.GetTreasure(ctx, treasureId)
	if err != nil {
		return nil, err
	}
	bookId, err := p.islandStore.GetBookOfIsland(ctx, player.AtIsland, userId)
	if errors.Is(err, domain.ErrNoBookAssignedFromPool) {
		err = nil
		bookId = ""
	} else if err != nil {
		return nil, err
	}
	userTreasure, err := p.treasureStore.GetUserTreasure(ctx, userId, treasureId)
	if err != nil {
		return nil, err
	}
	event, updatedUserTreasure, err := domain.UnlockTreasure(player, treasure, userTreasure, bookId, chosenCost)
	if err != nil {
		return nil, err
	}
	err = p.treasureStore.UpdateUserTreasure(ctx, userTreasure, updatedUserTreasure)
	if err != nil {
		return nil, err
	}
	islandTreasure := domain.GetIslandTreasureOfUserTreasure(updatedUserTreasure, true)
	return &islandTreasure, p.applyAndSendPlayerUpdateEvent(ctx, player, event)
}

func (p *Player) HandleNewPortableIsland(userId int32) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		player, err := p.playerStore.Get(ctx, userId)
		if err != nil {
			slog.Error("failed to get player for sending new portable island update", "error", err.Error())
			return
		}
		err = p.sendPlayerUpdateEventErr(ctx, &domain.PlayerUpdateEvent{
			Reason: domain.PlayerUpdateEventNewBook,
			Player: &player,
		})
		if err != nil {
			slog.Error("failed to send new portable island update", "error", err.Error())
		}
	}()
}

func (p *Player) MakeOfferCheck(ctx context.Context, userId int32) (*domain.MakeOfferCheckResult, error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	count, err := p.marketStore.GetOffersCountOfUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	check := domain.MakeOfferCheck(player, count)
	return &check, nil
}

func (p *Player) MakeOffer(ctx context.Context, offerer *domain.User, offered, requested domain.Cost) (result *domain.TradeOfferView, err error) {
	player, err := p.playerStore.Get(ctx, offerer.ID)
	if err != nil {
		return nil, err
	}
	count, err := p.marketStore.GetOffersCountOfUser(ctx, offerer.ID)
	if err != nil {
		return nil, err
	}
	event, tradeOffer, err := domain.MakeOffer(player, count, offered, requested)
	if err != nil {
		return nil, err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()
	err = p.marketStore.CreateOffer(ctx, tx, tradeOffer)
	if err != nil {
		return nil, err
	}
	err = p.playerStore.Update(ctx, tx, player, *event.Player)
	if err != nil {
		return nil, err
	}
	err = p.sendPlayerUpdateEventErr(ctx, event)
	if err != nil {
		return nil, err
	}
	p.tradeEventBroadcastHandler(func(userId int32) *domain.TradeEvent {
		return &domain.TradeEvent{
			NewOffer: &domain.NewOfferTradeEvent{
				Offer: domain.TradeOfferViewForPlayer(userId, offerer.Name, tradeOffer),
			},
		}
	})
	view := domain.TradeOfferViewForPlayer(offerer.ID, offerer.Name, tradeOffer)
	return &view, nil
}

func (p *Player) AcceptOffer(ctx context.Context, userId int32, tradeOfferId string) (err error) {
	acceptor, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return err
	}

	offer, err := p.marketStore.GetOffer(ctx, tradeOfferId)
	if err != nil {
		return err
	}

	offerer, err := p.playerStore.Get(ctx, offer.By)
	if err != nil {
		return err
	}

	acceptorEvent, offererEvent, err := domain.AcceptOffer(acceptor, offerer, offer)
	if err != nil {
		return err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	err = p.marketStore.DeleteOffer(ctx, tx, tradeOfferId)
	if err != nil {
		return err
	}

	err = p.playerStore.Update(ctx, tx, acceptor, *acceptorEvent.Player)
	if err != nil {
		return err
	}

	err = p.playerStore.Update(ctx, tx, offerer, *offererEvent.Player)
	if err != nil {
		return err
	}

	err = p.createAndSendInboxMessage(ctx, tx, domain.InboxMessage{
		ID:        domain.NewID(domain.ResourceTypeInboxMessage),
		UserID:    offer.By,
		CreatedAt: time.Now().UTC(),
		Content: domain.InboxMessageContent{
			OwnOfferAccepted: &domain.InboxMessageOwnOfferAccepted{
				Offer: domain.TradeOfferViewForPlayer(offer.By, "", offer),
			},
		},
	})
	if err != nil {
		return err
	}

	err = p.sendPlayerUpdateEventErr(ctx, acceptorEvent)
	if err != nil {
		return err
	}

	err = p.sendPlayerUpdateEventErr(ctx, offererEvent)
	if err != nil {
		return err
	}

	p.tradeEventBroadcastHandler(func(userId int32) *domain.TradeEvent {
		return &domain.TradeEvent{
			DeletedOffer: &domain.DeletedOfferTradeEvent{
				OfferID: offer.ID,
				ByMe:    userId == offer.By,
			},
		}
	})

	return nil
}

func (p *Player) DeleteOffer(ctx context.Context, userId int32, tradeOfferId string) (err error) {
	player, err := p.playerStore.Get(ctx, userId)
	if err != nil {
		return err
	}

	offer, err := p.marketStore.GetOffer(ctx, tradeOfferId)
	if err != nil {
		return err
	}

	event, err := domain.DeleteOffer(player, offer)
	if err != nil {
		return err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	err = p.marketStore.DeleteOffer(ctx, tx, tradeOfferId)
	if err != nil {
		return err
	}

	err = p.playerStore.Update(ctx, tx, player, *event.Player)
	if err != nil {
		return err
	}

	err = p.sendPlayerUpdateEventErr(ctx, event)
	if err != nil {
		return err
	}

	p.tradeEventBroadcastHandler(func(userId int32) *domain.TradeEvent {
		return &domain.TradeEvent{
			DeletedOffer: &domain.DeletedOfferTradeEvent{
				OfferID: offer.ID,
				ByMe:    userId == offer.By,
			},
		}
	})

	return nil
}

func (p *Player) GetTradeOffers(ctx context.Context, userId int32, filter domain.GetOffersByFilterType, offset int64, limit int) ([]domain.TradeOfferView, error) {
	limit = min(max(1, limit), 100)
	before := time.Now().UTC()
	if offset > 0 {
		before = time.UnixMilli(offset).UTC()
	}
	offers, err := p.marketStore.GetOffers(ctx, filter, userId, before, limit)
	if err != nil {
		return nil, err
	}

	result := make([]domain.TradeOfferView, 0)
	for _, offer := range offers {
		// Get the offerer's user to get their username
		offererUser, err := p.userStore.Get(ctx, offer.By)
		if err != nil {
			return nil, fmt.Errorf("failed to get offerer user: %w", err)
		}

		tradeOfferView := domain.TradeOfferViewForPlayer(userId, offererUser.Name, offer)
		result = append(result, tradeOfferView)
	}

	return result, nil
}

func (p *Player) OnTradeEventBroadcast(handler TradeEventBroadcastHandler) {
	p.tradeEventBroadcastHandler = handler
}

func (p *Player) GetInitialTradeEvent(_ context.Context) (*domain.TradeEvent, error) {
	return &domain.TradeEvent{
		Sync: &domain.SyncTradeEvent{
			Offset: fmt.Sprint(time.Now().UTC().UnixMilli()),
		},
	}, nil
}

func (p *Player) OnInboxEvent(handler func(e *domain.InboxEvent)) {
	p.inboxEventHandler = handler
}

func (p *Player) createAndSendInboxMessage(ctx context.Context, tx domain.Tx, msg domain.InboxMessage) error {
	err := p.inboxStore.CreateMessage(ctx, tx, msg)
	if err != nil {
		return fmt.Errorf("failed to create inbox message: %w", err)
	}
	view := domain.InboxMessageToView(msg)
	p.inboxEventHandler(&domain.InboxEvent{
		UserId:     msg.UserID,
		NewMessage: &view,
	})
	return nil
}

func (p *Player) GetInitialInboxEvent(_ context.Context, userId int32) (*domain.InboxEvent, error) {
	return &domain.InboxEvent{
		UserId: userId,
		Sync: &domain.SyncInboxEvent{
			Offset: fmt.Sprint(time.Now().UTC().UnixMilli()),
		},
	}, nil
}

func (p *Player) GetInboxMessages(ctx context.Context, userId int32, offset int64, limit int) ([]domain.InboxMessageView, error) {
	limit = min(max(1, limit), 100)
	before := time.Now().UTC()
	if offset > 0 {
		before = time.UnixMilli(offset).UTC()
	}
	messages, err := p.inboxStore.GetMessages(ctx, userId, before, limit)
	if err != nil {
		return nil, err
	}
	result := make([]domain.InboxMessageView, 0, len(messages))
	for _, m := range messages {
		result = append(result, domain.InboxMessageToView(m))
	}
	return result, nil
}

func (p *Player) OnBroadcastMessage(f MessageBroadcastHandler) {
	p.broadcastMessageHandler = f
}

func (p *Player) BroadcastMessage(ctx context.Context, text string) (int, error) {
	players, err := p.playerStore.GetAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get players: %w", err)
	}
	var errs []error
	count := 0
	messages := make(map[int32]domain.InboxMessage)

	for _, player := range players {
		msg := domain.InboxMessage{
			ID:        domain.NewID(domain.ResourceTypeInboxMessage),
			UserID:    player,
			CreatedAt: time.Now().UTC(),
			Content: domain.InboxMessageContent{
				Announcement: &domain.InboxMessageAnnouncement{Text: text},
			},
		}
		messages[player] = msg
		err := p.inboxStore.CreateMessage(ctx, nil, msg)
		errs = append(errs, err)
		if err == nil {
			count++
		}
	}
	p.broadcastMessageHandler(func(userId int32) *domain.InboxMessageView {
		msg, ok := messages[userId]
		if !ok {
			return nil
		}
		r := domain.InboxMessageToView(msg)
		return &r
	})

	return count, errors.Join(errs...)
}

func (p *Player) InvestCheck(ctx context.Context, user *domain.User) (*domain.InvestmentCheckResult, error) {
	player, err := p.playerStore.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	activeSession, err := p.investStore.GetActiveSession(ctx)
	if err != nil {
		return nil, err
	}

	var investments []domain.UserInvestment
	if activeSession != nil {
		investments, err = p.investStore.GetUserInvestments(ctx, activeSession.ID, user.ID)
		if err != nil {
			return nil, err
		}
	}

	result := domain.InvestCheck(activeSession, investments, player)
	return &result, nil
}

// Invest allows a player to make an investment
func (p *Player) Invest(ctx context.Context, user *domain.User, sessionId string, coinAmount int32) (*domain.UserInvestment, error) {
	player, err := p.playerStore.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	session, err := p.investStore.GetSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	investments, err := p.investStore.GetUserInvestments(ctx, session.ID, user.ID)
	if err != nil {
		return nil, err
	}

	event, userInvestment, err := domain.Invest(*session, investments, player, coinAmount)
	if err != nil {
		return nil, err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err == nil {
			err = p.sendPlayerUpdateEventErr(ctx, event)
		}
	}()
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	err = p.investStore.AddUserInvestment(ctx, tx, userInvestment)
	if err != nil {
		return nil, err
	}

	err = p.playerStore.Update(ctx, tx, player, *event.Player)
	if err != nil {
		return nil, err
	}

	return &userInvestment, nil
}

type PlayersLocation struct {
	IslandID string `json:"islandId"`
	Players  []struct {
		UserID int32  `json:"-"`
		Name   string `json:"name"`
	} `json:"players"`
}

func (p *Player) GetPlayerLocations(ctx context.Context, userId int32, territoryID string) ([]PlayersLocation, error) {
	filterForUser := func(result []PlayersLocation) []PlayersLocation {
		c := make([]PlayersLocation, 0, len(result))
		for _, location := range result {
			l := PlayersLocation{IslandID: location.IslandID}
			for _, player := range location.Players {
				if player.UserID != userId {
					l.Players = append(l.Players, player)
				}
			}
			if len(l.Players) > 0 {
				c = append(c, l)
			}
		}
		return c
	}

	v, ok := p.playerLocationsCache.Get(territoryID)
	if ok {
		result, _ := v.([]PlayersLocation)
		return filterForUser(result), nil
	}

	locations, err := p.playerStore.GetLocations(ctx, territoryID)
	if err != nil {
		return nil, err
	}
	if len(locations) == 0 {
		return nil, nil
	}

	result := make([]PlayersLocation, 0, len(locations))
	for island, players := range locations {
		location := PlayersLocation{
			IslandID: island,
		}
		for _, userId := range players {
			user, err := p.userStore.Get(ctx, userId)
			if err != nil {
				return nil, err
			}
			location.Players = append(location.Players, struct {
				UserID int32  `json:"-"`
				Name   string `json:"name"`
			}{UserID: userId, Name: user.Name})
		}
		result = append(result, location)
	}

	p.playerLocationsCache.SetDefault(territoryID, result)
	return filterForUser(result), nil
}

func (p *Player) ResolveInvestmentSession(ctx context.Context, sessionID string, coefficient float64) (affectedPlayers int, sumOfRewards int, err error) {
	session, err := p.investStore.GetSession(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}
	investments, err := p.investStore.GetAllUserInvestments(ctx, session.ID)
	if err != nil {
		return 0, 0, err
	}

	rewards, err := domain.ResolveInvestments(*session, investments, coefficient)
	if err != nil {
		return 0, 0, err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	var events []domain.PlayerUpdateEvent
	defer func() {
		if err == nil {
			for _, e := range events {
				_ = p.sendPlayerUpdateEventErr(ctx, &e)
			}
		}
	}()
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	err = p.investStore.MarkResolved(ctx, tx, session.ID)
	if err != nil {
		return 0, 0, err
	}

	for userId, coinCount := range rewards {
		player, err := p.playerStore.Get(ctx, userId)
		if err != nil {
			return 0, 0, err
		}
		event, ok := domain.GiveInvestmentReward(player, coinCount)
		if ok {
			if err := p.playerStore.Update(ctx, tx, player, *event.Player); err != nil {
				return 0, 0, err
			}
			sumOfRewards += int(coinCount)
			affectedPlayers++
			events = append(events, event)
		}
	}

	return
}
