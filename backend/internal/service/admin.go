package service

import (
	"context"
	cRand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"slices"
	"time"
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

func (a *Admin) GetTerritories(ctx context.Context) ([]domain.Territory, error) {
	return a.territoryStore.ListTerritories(ctx)
}

func (a *Admin) SetTerritory(ctx context.Context, territory domain.Territory) error {
	if territory.ID == "" {
		return AdminError{"id is required"}
	}
	for _, island := range territory.Islands {
		if island.ID == "" {
			return AdminError{"empty island id in island list"}
		}
		if island.Name == "" {
			return AdminError{"empty island name in island list"}
		}
	}
	if territory.StartIsland == "" {
		return AdminError{"invalid territory startIsland"}
	}
	isInIslands := func(id string) bool {
		return slices.ContainsFunc(territory.Islands, func(island domain.Island) bool {
			return island.ID == id
		})
	}
	if !isInIslands(territory.StartIsland) {
		return AdminError{fmt.Sprintf("startIsland %q not found in island list", territory.StartIsland)}
	}
	for _, e := range territory.Edges {
		if e.From == "" || e.To == "" {
			return AdminError{fmt.Sprintf("empty edge.from or edge.to: %v", e)}
		}
		if !isInIslands(e.From) {
			return AdminError{fmt.Sprintf("edge.from %q is not in island list", e.From)}
		}
		if !isInIslands(e.To) {
			return AdminError{fmt.Sprintf("edge.to %q is not in island list", e.To)}
		}
	}
	for _, r := range territory.RefuelIslands {
		if !isInIslands(r.ID) {
			return AdminError{fmt.Sprintf("refuelIsland %q not found in island list", r.ID)}
		}
	}
	for _, t := range territory.TerminalIslands {
		if !isInIslands(t.ID) {
			return AdminError{fmt.Sprintf("terminalIsland %q not found in island list", t.ID)}
		}
	}
	for islandID, prerequisites := range territory.IslandPrerequisites {
		if !isInIslands(islandID) {
			return AdminError{fmt.Sprintf("island %q in prerequisites not found in island list", islandID)}
		}
		for _, p := range prerequisites {
			if !isInIslands(p) {
				return AdminError{fmt.Sprintf("prerequisite %q not found in island list", p)}
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

func (a *Admin) GetIslandHeader(ctx context.Context, islandId string) (domain.IslandHeader, error) {
	return a.islandStore.GetIslandHeader(ctx, islandId)
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
	ID              string   `json:"id"`
	Text            string   `json:"text"`
	InputType       string   `json:"inputType"`
	InputAccept     []string `json:"inputAccept"`
	KnowledgeAmount int32    `json:"knowledgeAmount"`
	RewardSource    string   `json:"rewardSource,omitempty"`
	Context         string   `json:"correctionHintMessage,omitempty"`
}

func (a *Admin) SetBookAndBindToIsland(ctx context.Context, islandId string, input BookInput) (BookInput, error) {
	territoryId, err := a.islandStore.GetTerritory(ctx, islandId)
	if err != nil {
		return input, err
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
		return input, AdminError{fmt.Sprintf("invalid poolId %q", poolId)}
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
	if input.BookId == "" {
		input.BookId = domain.NewID(domain.ResourceTypeBook)
	} else if !domain.IdHasType(input.BookId, domain.ResourceTypeBook) {
		return input, AdminError{fmt.Sprintf("invalid bookId %q", input.BookId)}
	}
	book := domain.Book{ID: input.BookId, Components: make([]domain.BookComponent, 0)}
	var questions []domain.BookQuestion
	var treasures []domain.Treasure
	for i, c := range input.Components {
		if c.IFrame != nil {
			if c.IFrame.Url == "" {
				return input, AdminError{fmt.Sprintf("empty url for book %q iframe component at index %d", book.ID, i)}
			}
			book.Components = append(book.Components, domain.BookComponent{IFrame: c.IFrame})
			continue
		}
		if c.Question != nil {
			if c.Question.InputType == "" {
				return input, AdminError{fmt.Sprintf("empty inputType for book %q question at index %d", book.ID, i)}
			}
			if c.Question.InputType == "file" && len(c.Question.InputAccept) == 0 {
				return input, AdminError{fmt.Sprintf("empty inputAccept for book %q question at index %d", book.ID, i)}
			}
			if c.Question.KnowledgeAmount < 0 {
				return input, AdminError{fmt.Sprintf("negative knowledgeAmount for book %q question at index %d", book.ID, i)}
			}
			if !domain.IsValidRewardSource(c.Question.RewardSource) {
				return input, AdminError{fmt.Sprintf("invalid reward source %q", c.Question.RewardSource)}
			}
			if c.Question.Text == "" {
				return input, AdminError{fmt.Sprintf("empty text for book %q question at index %d", book.ID, i)}
			}
			if c.Question.ID == "" {
				c.Question.ID = domain.NewID(domain.ResourceTypeQuestion)
			} else if !domain.IdHasType(c.Question.ID, domain.ResourceTypeQuestion) {
				return input, AdminError{fmt.Sprintf("invalid question id %q", c.Question.ID)}
			}
			questions = append(questions, domain.BookQuestion{
				QuestionID:      c.Question.ID,
				BookID:          input.BookId,
				Text:            c.Question.Text,
				InputType:       c.Question.InputType,
				InputAccept:     c.Question.InputAccept,
				KnowledgeAmount: c.Question.KnowledgeAmount,
				RewardSource:    c.Question.RewardSource,
				Context:         c.Question.Context,
			})
			book.Components = append(book.Components, domain.BookComponent{Question: &domain.QuestionPlaceholder{ID: c.Question.ID}})
			continue
		}
		return input, AdminError{fmt.Sprintf("unknown component for book %q at index %d", book.ID, i)}
	}
	for _, t := range input.Treasures {
		if t.ID == "" {
			t.ID = domain.NewID(domain.ResourceTypeTreasure)
		} else if !domain.IdHasType(t.ID, domain.ResourceTypeTreasure) {
			return input, AdminError{fmt.Sprintf("invalid treasure id %q", t.ID)}
		}
		treasures = append(treasures, domain.Treasure{ID: t.ID, BookID: input.BookId})
	}
	err := a.islandStore.SetBook(ctx, book)
	if err != nil {
		return input, fmt.Errorf("failed to set book: %w", err)
	}
	err = a.questionStore.BindQuestionsToBook(ctx, book.ID, questions)
	if err != nil {
		return input, fmt.Errorf("failed to bind questions to book: %w", err)
	}
	err = a.treasureStore.BindTreasuresToBook(ctx, book.ID, treasures)
	if err != nil {
		return input, fmt.Errorf("failed to bind treasures to book: %w", err)
	}
	return input, nil
}

func (a *Admin) GetBook(ctx context.Context, bookId string) (BookInput, error) {
	book, err := a.islandStore.GetBook(ctx, bookId)
	if err != nil {
		return BookInput{}, err
	}
	bookQuestions, err := a.questionStore.GetQuestions(ctx, bookId)
	if err != nil {
		return BookInput{}, err
	}
	treasures, err := a.treasureStore.GetTreasures(ctx, bookId)
	if err != nil {
		return BookInput{}, err
	}

	islandInputQuestions := make(map[string]IslandInputQuestion)
	for _, q := range bookQuestions {
		islandInputQuestions[q.QuestionID] = IslandInputQuestion{
			ID:              q.QuestionID,
			Text:            q.Text,
			InputType:       q.InputType,
			InputAccept:     q.InputAccept,
			KnowledgeAmount: q.KnowledgeAmount,
			RewardSource:    q.RewardSource,
			Context:         q.Context,
		}
	}
	result := BookInput{
		BookId: book.ID,
	}
	for _, c := range book.Components {
		if c.IFrame != nil {
			result.Components = append(result.Components, &BookInputComponent{
				IFrame: c.IFrame,
			})
			continue
		}
		if c.Question != nil {
			q, ok := islandInputQuestions[c.Question.ID]
			if !ok {
				q.ID = c.Question.ID
			}
			result.Components = append(result.Components, &BookInputComponent{Question: &q})
			delete(islandInputQuestions, c.Question.ID)
			continue
		}
	}
	for _, q := range islandInputQuestions {
		result.Components = append(result.Components, &BookInputComponent{Question: &q})
	}
	for _, t := range treasures {
		result.Treasures = append(result.Treasures, &BookTreasureComponent{
			ID: t.ID,
		})
	}

	return result, nil
}

type PoolOutput struct {
	ID    string   `json:"id"`
	Books []string `json:"books"`
}

func (a *Admin) GetPools(ctx context.Context) ([]PoolOutput, error) {
	var result []PoolOutput
	for _, poolId := range domain.PoolIds() {
		books, err := a.islandStore.GetBooksInPool(ctx, poolId)
		if err != nil {
			return result, err
		}
		result = append(result, PoolOutput{
			ID:    poolId,
			Books: books,
		})
	}
	return result, nil
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
	if _, err := a.territoryStore.GetTerritoryByID(ctx, binding.TerritoryId); err != nil {
		return binding, err
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
			binding.EmptyIslands = append(binding.EmptyIslands, h.ID)
		}
	}
	settings, err := a.islandStore.GetTerritoryPoolSettings(ctx, territoryId)
	if errors.Is(err, domain.ErrPoolSettingsNotFound) {
		err = nil
		settings = domain.TerritoryPoolSettings{}
	}
	if err != nil {
		return binding, err
	}
	binding.PoolSettings = settings
	return binding, nil
}

func (a *Admin) SetTerritoryIslandBindings(ctx context.Context, bindings TerritoryIslandBindings) (TerritoryIslandBindings, error) {
	pooledCount := int32(len(bindings.PooledIslands))
	if pooledCount != bindings.PoolSettings.TotalCount() {
		return bindings, AdminError{fmt.Sprintf("number of pooled islands don't match pool settings: %d vs %d", pooledCount, bindings.PoolSettings.TotalCount())}
	}
	if _, err := a.territoryStore.GetTerritoryByID(ctx, bindings.TerritoryId); err != nil {
		return bindings, err
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

type User struct {
	Name              string `json:"name"`
	Username          string `json:"username"`
	Password          string `json:"password,omitempty"`
	StartingTerritory string `json:"startingTerritory,omitempty"`
	MeetLink          string `json:"meetLink"`
}

func (a *Admin) CreateUser(ctx context.Context, user User) (User, error) {
	if user.Username == "" {
		return User{}, AdminError{"username is required"}
	}
	if user.Password == "" {
		b := make([]byte, 8)
		_, _ = cRand.Read(b)
		user.Password = base64.RawURLEncoding.EncodeToString(b)
	}
	if user.StartingTerritory == "" {
		return User{}, AdminError{"startingTerritory is required"}
	}

	startingTerritory, err := a.territoryStore.GetTerritoryByID(ctx, user.StartingTerritory)
	if errors.Is(err, domain.ErrTerritoryNotFound) {
		return user, AdminError{fmt.Sprintf("failed to find starting territory %q", user.StartingTerritory)}
	}
	if err != nil {
		return user, err
	}

	hp, err := domain.HashPassword(user.Password)
	if err != nil {
		return user, err
	}

	u := &domain.User{
		ID:             rand.Int31(),
		Username:       user.Username,
		Name:           user.Name,
		MeetLink:       user.MeetLink,
		HashedPassword: hp,
	}
	if err := a.userStore.Create(ctx, u); err != nil {
		return user, err
	}
	return user, a.playerStore.Create(ctx, domain.NewPlayer(u.ID, startingTerritory))
}

func (a *Admin) GetUsers(ctx context.Context) ([]User, error) {
	users, err := a.userStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]User, 0, len(users))
	for _, u := range users {
		result = append(result, User{
			Name:     u.Name,
			Username: u.Username,
			MeetLink: u.MeetLink,
		})
	}
	return result, nil
}

func (a *Admin) Login(_ context.Context, username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", domain.ErrUserNotFound
	}
	if username != a.cfg.AdminUsername || password != a.cfg.AdminPassword {
		return "", domain.ErrUserNotFound
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"admin": true,
		"iat":   float64(time.Now().UTC().UnixNano()) / 1e9,
	})
	tokenString, err := token.SignedString(a.cfg.TokenSigningKeyBytes())
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func (a *Admin) ValidateToken(_ context.Context, tokenStr string) bool {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return a.cfg.TokenSigningKeyBytes(), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return false
	}
	if !token.Valid {
		return false
	}
	if iat, err := token.Claims.GetIssuedAt(); err != nil || time.Since(iat.Time) > 6*time.Hour {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	v, ok := claims["admin"]
	if !ok {
		return false
	}
	isAdmin, _ := v.(bool)
	return isAdmin
}

type AdminError struct {
	text string
}

func (e AdminError) Error() string {
	return e.text
}
