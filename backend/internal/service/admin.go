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
	questionStore  domain.QuestionStore
}

func NewAdmin(territoryStore domain.TerritoryStore, islandStore domain.IslandStore, userStore domain.UserStore, playerStore domain.PlayerStore, questionStore domain.QuestionStore) *Admin {
	return &Admin{
		territoryStore: territoryStore,
		islandStore:    islandStore,
		userStore:      userStore,
		playerStore:    playerStore,
		questionStore:  questionStore,
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

func (a *Admin) SetIsland(ctx context.Context, id string, input domain.IslandInputContent) (domain.IslandInputContent, error) {
	raw := &domain.IslandRawContent{Components: make([]domain.IslandRawComponent, 0)}
	var questions []domain.Question
	for i, c := range input.Components {
		if c.ID == "" || !domain.IdHasType(c.ID, domain.ResourceTypeComponent) {
			c.ID = domain.NewID(domain.ResourceTypeComponent)
		}
		if c.IFrame != nil {
			if c.IFrame.Url == "" {
				return input, fmt.Errorf("empty url for island %q iframe component at index %d", id, i)
			}
			raw.Components = append(raw.Components, domain.IslandRawComponent{ID: c.ID, IFrame: c.IFrame})
			continue
		}
		if c.Question != nil {
			if c.Question.InputType == "" {
				return input, fmt.Errorf("empty inputType for island %q question at index %d", id, i)
			}
			if c.Question.InputType == "file" && len(c.Question.InputAccept) == 0 {
				return input, fmt.Errorf("empty inputAccept for island %q question at index %d", id, i)
			}
			if c.Question.KnowledgeAmount <= 0 {
				return input, fmt.Errorf("non-positive knowledgeAmount for island %q question at index %d", id, i)
			}
			if c.Question.Text == "" {
				return input, fmt.Errorf("empty text for island %q question at index %d", id, i)
			}
			if c.Question.ID == "" || !domain.IdHasType(c.Question.ID, domain.ResourceTypeQuestion) {
				c.Question.ID = domain.NewID(domain.ResourceTypeQuestion)
			}
			questions = append(questions, *c.Question)
			raw.Components = append(raw.Components, domain.IslandRawComponent{ID: c.ID, Question: &domain.QuestionComponent{QuestionID: c.Question.ID}})
			continue
		}
		return input, fmt.Errorf("unknown component for island %q at index %d", id, i)
	}
	for _, q := range questions {
		err := a.questionStore.SetQuestion(ctx, q)
		if err != nil {
			return input, err
		}
	}
	err := a.islandStore.SetContent(ctx, id, raw)
	if err != nil {
		return input, err
	}
	return input, nil
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
