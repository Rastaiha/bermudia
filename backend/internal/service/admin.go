package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"slices"
)

type Admin struct {
	cfg            config.Config
	territoryStore domain.TerritoryStore
	islandStore    domain.IslandStore
	userStore      domain.UserStore
	playerStore    domain.PlayerStore
	questionStore  domain.QuestionStore
	treasureStore  domain.TreasureStore
}

func NewAdmin(cfg config.Config, territoryStore domain.TerritoryStore, islandStore domain.IslandStore, userStore domain.UserStore, playerStore domain.PlayerStore, questionStore domain.QuestionStore, treasureStore domain.TreasureStore) *Admin {
	return &Admin{
		cfg:            cfg,
		territoryStore: territoryStore,
		islandStore:    islandStore,
		userStore:      userStore,
		playerStore:    playerStore,
		questionStore:  questionStore,
		treasureStore:  treasureStore,
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
	isInIslands := func(id string) bool {
		return slices.ContainsFunc(territory.Islands, func(island domain.Island) bool {
			return island.ID == id
		})
	}
	if !isInIslands(territory.StartIsland) {
		return fmt.Errorf("startIsland %q not found in island list", territory.StartIsland)
	}
	for _, e := range territory.Edges {
		if e.From == "" || e.To == "" {
			return fmt.Errorf("empty edge.from or edge.to: %v", e)
		}
		if !isInIslands(e.From) {
			return fmt.Errorf("edge.from %q is not in island list", e.From)
		}
		if !isInIslands(e.To) {
			return fmt.Errorf("edge.to %q is not in island list", e.To)
		}
	}
	for _, r := range territory.RefuelIslands {
		if !isInIslands(r.ID) {
			return fmt.Errorf("refuelIsland %q not found in island list", r.ID)
		}
	}
	for _, t := range territory.TerminalIslands {
		if !isInIslands(t.ID) {
			return fmt.Errorf("terminalIsland %q not found in island list", t.ID)
		}
	}
	for islandID, prerequisites := range territory.IslandPrerequisites {
		if !isInIslands(islandID) {
			return fmt.Errorf("island %q in prerequisites not found in island list", islandID)
		}
		for _, p := range prerequisites {
			if !isInIslands(p) {
				return fmt.Errorf("prerequisite %q not found in island list", p)
			}
		}
	}

	for _, island := range territory.Islands {
		if err := a.islandStore.ReserveIDForTerritory(ctx, territory.ID, island.ID, island.Name); err != nil {
			return err
		}
		if err := a.islandStore.ReserveIDForTerritory(ctx, territory.ID, island.ID, island.Name); err != nil {
			return err
		}
	}

	return a.territoryStore.SetTerritory(ctx, &territory)
}

type BookInput struct {
	BookId     string                   `json:"bookId"`
	Components []*BookInputComponent    `json:"components"`
	Treasures  []*BookTreasureComponent `json:"treasures"`
}

type BookTreasureComponent struct {
	ID string `json:"id"`
}

type BookInputComponent struct {
	IFrame   *domain.IslandIFrame `json:"iframe,omitempty"`
	Question *IslandInputQuestion `json:"question,omitempty"`
}

type IslandInputQuestion struct {
	domain.Question
	KnowledgeAmount int32  `json:"knowledgeAmount"`
	RewardSource    string `json:"rewardSource,omitempty"`
}

func (a *Admin) SetBookAndBindToIsland(ctx context.Context, islandId string, input BookInput) (BookInput, error) {
	territoryId, err := a.islandStore.GetTerritory(ctx, islandId)
	if err != nil {
		return input, fmt.Errorf("island %q does not have territory", islandId)
	}
	input, err = a.setBook(ctx, input)
	if err != nil {
		return input, err
	}
	err = a.islandStore.SetIslandHeader(ctx, domain.IslandHeader{
		ID:          islandId,
		TerritoryID: territoryId,
		BookID:      input.BookId,
		FromPool:    false,
	})
	if err != nil {
		return input, fmt.Errorf("failed to set island header: %w", err)
	}
	return input, nil
}

func (a *Admin) SetBookAndBindToPool(ctx context.Context, poolId string, input BookInput) (BookInput, error) {
	if !domain.IsPoolIdValid(poolId) {
		return input, fmt.Errorf("invalid poolId %q", poolId)
	}
	input, err := a.setBook(ctx, input)
	if err != nil {
		return input, err
	}
	err = a.islandStore.AddBookToPool(ctx, poolId, input.BookId)
	if err != nil {
		return input, fmt.Errorf("failed to add book to pool: %w", err)
	}
	return input, nil
}

func (a *Admin) setBook(ctx context.Context, input BookInput) (BookInput, error) {
	if input.BookId == "" || domain.IdHasType(input.BookId, domain.ResourceTypeBook) {
		input.BookId = domain.NewID(domain.ResourceTypeBook)
	}
	book := domain.Book{ID: input.BookId, Components: make([]domain.BookComponent, 0)}
	var questions []domain.BookQuestion
	for i, c := range input.Components {
		if c.IFrame != nil {
			if c.IFrame.Url == "" {
				return input, fmt.Errorf("empty url for book %q iframe component at index %d", book.ID, i)
			}
			book.Components = append(book.Components, domain.BookComponent{IFrame: c.IFrame})
			continue
		}
		if c.Question != nil {
			if c.Question.InputType == "" {
				return input, fmt.Errorf("empty inputType for book %q question at index %d", book.ID, i)
			}
			if c.Question.InputType == "file" && len(c.Question.InputAccept) == 0 {
				return input, fmt.Errorf("empty inputAccept for book %q question at index %d", book.ID, i)
			}
			if c.Question.KnowledgeAmount < 0 {
				return input, fmt.Errorf("negative knowledgeAmount for book %q question at index %d", book.ID, i)
			}
			if !domain.IsValidRewardSource(c.Question.RewardSource) {
				return input, fmt.Errorf("invalid reward source %q", c.Question.RewardSource)
			}
			if c.Question.Text == "" {
				return input, fmt.Errorf("empty text for book %q question at index %d", book.ID, i)
			}
			if c.Question.ID == "" || !domain.IdHasType(c.Question.ID, domain.ResourceTypeQuestion) {
				c.Question.ID = domain.NewID(domain.ResourceTypeQuestion)
			}
			questions = append(questions, domain.BookQuestion{
				QuestionID:      c.Question.ID,
				BookID:          input.BookId,
				Text:            c.Question.Text,
				KnowledgeAmount: c.Question.KnowledgeAmount,
				RewardSource:    c.Question.RewardSource,
			})
			book.Components = append(book.Components, domain.BookComponent{Question: &c.Question.Question})
			continue
		}
		return input, fmt.Errorf("unknown component for book %q at index %d", book.ID, i)
	}
	for _, t := range input.Treasures {
		if t.ID == "" || !domain.IdHasType(t.ID, domain.ResourceTypeTreasure) {
			t.ID = domain.NewID(domain.ResourceTypeTreasure)
		}
		book.Treasures = append(book.Treasures, domain.Treasure{ID: t.ID, BookID: input.BookId})
	}
	err := a.islandStore.SetBook(ctx, book)
	if err != nil {
		return input, fmt.Errorf("failed to set book: %w", err)
	}
	err = a.questionStore.BindQuestionsToBook(ctx, book.ID, questions)
	if err != nil {
		return input, fmt.Errorf("failed to bind questions to book: %w", err)
	}
	err = a.treasureStore.BindTreasuresToBook(ctx, book.ID, book.Treasures)
	if err != nil {
		return input, fmt.Errorf("failed to bind treasures to book: %w", err)
	}
	return input, nil
}

type TerritoryIslandBindings struct {
	TerritoryId   string                       `json:"territoryId"`
	EmptyIslands  []string                     `json:"emptyIslands"`
	PooledIslands []string                     `json:"pooledIslands"`
	PoolSettings  domain.TerritoryPoolSettings `json:"poolSettings"`
}

func (a *Admin) GetTerritoryIslandBindings(ctx context.Context, territoryId string) (TerritoryIslandBindings, error) {
	binding := TerritoryIslandBindings{
		TerritoryId: territoryId,
	}
	islands, err := a.islandStore.GetIslandHeadersByTerritory(ctx, territoryId)
	if err != nil {
		return binding, fmt.Errorf("failed to get island headers by territory %q: %w", territoryId, err)
	}
	for _, h := range islands {
		if h.FromPool {
			binding.PooledIslands = append(binding.PooledIslands, h.ID)
		}
		if !h.FromPool && h.BookID == "" {
			binding.PooledIslands = append(binding.EmptyIslands, h.ID)
		}
	}
	settings, err := a.islandStore.GetTerritoryPoolSettings(ctx, territoryId)
	if err != nil {
		return binding, err
	}
	binding.PoolSettings = settings
	return binding, nil
}

func (a *Admin) SetTerritoryIslandBindings(ctx context.Context, bindings TerritoryIslandBindings) (TerritoryIslandBindings, error) {
	pooledCount := int32(len(bindings.PooledIslands))
	if pooledCount != bindings.PoolSettings.TotalCount() {
		return bindings, fmt.Errorf("number of pooled islands don't match pool settings: %d vs %d", pooledCount, bindings.PoolSettings.TotalCount())
	}
	err := a.islandStore.SetTerritoryPoolSettings(ctx, bindings.TerritoryId, bindings.PoolSettings)
	if err != nil {
		return bindings, err
	}
	for _, id := range bindings.EmptyIslands {
		err := a.islandStore.SetIslandHeader(ctx, domain.IslandHeader{
			ID:          id,
			TerritoryID: bindings.TerritoryId,
			FromPool:    false,
			BookID:      "",
		})
		if err != nil {
			return bindings, fmt.Errorf("failed to set header for island %q: %w", id, err)
		}
	}
	for _, id := range bindings.PooledIslands {
		err := a.islandStore.SetIslandHeader(ctx, domain.IslandHeader{
			ID:          id,
			TerritoryID: bindings.TerritoryId,
			FromPool:    true,
			BookID:      "",
		})
		if err != nil {
			return bindings, fmt.Errorf("failed to set header for island %q: %w", id, err)
		}
	}
	return bindings, nil
}

func (a *Admin) CreateUser(ctx context.Context, id int32, username, password, startingTerritoryID string) error {
	territories, err := a.territoryStore.ListTerritories(ctx)
	if err != nil {
		return err
	}
	if len(territories) == 0 {
		return errors.New("no territory found")
	}
	startingTerritory := territories[int(id)%len(territories)]
	if startingTerritory.ID != "" {
		for _, t := range territories {
			if t.ID == startingTerritoryID {
				startingTerritory = t
				break
			}
		}
	}
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
