package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrIslandNotFound             = errors.New("island not found")
	ErrTerritoryNotFound          = errors.New("territory not found")
	ErrUserNotFound               = errors.New("user not found")
	ErrPlayerConflict             = errors.New("player update conflict")
	ErrAnswerNotPending           = errors.New("answer not in pending state")
	ErrQuestionNotRelatedToIsland = errors.New("question not related to island")
	ErrInvalidIslandHeader        = errors.New("invalid island header")
	ErrPoolSettingExhausted       = errors.New("pool setting exhausted")
	ErrBookPoolExhausted          = errors.New("book pool exhausted")
	ErrNoBookAssignedFromPool     = errors.New("no book assigned from pool")
	ErrEmptyIsland                = errors.New("empty island")
	ErrUserTreasureConflict       = errors.New("user treasure update conflict")
	ErrAlreadyApplied             = errors.New("already applied")
	ErrOfferAlreadyDeleted        = errors.New("offer already deleted")
	ErrInvalidFilter              = errors.New("invalid filter")
)

type Tx interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type TerritoryStore interface {
	SetTerritory(ctx context.Context, territory *Territory) error
	GetTerritoryByID(ctx context.Context, territoryID string) (*Territory, error)
	ListTerritories(ctx context.Context) ([]Territory, error)
}

type IslandStore interface {
	SetBook(ctx context.Context, book Book) error
	GetBook(ctx context.Context, bookId string) (*Book, error)
	SetIslandHeader(ctx context.Context, header IslandHeader) error
	ReserveIDForTerritory(ctx context.Context, territoryId, islandId, islandName string) error
	GetBookOfIsland(ctx context.Context, islandId string, userId int32) (string, error)
	GetTerritory(ctx context.Context, id string) (string, error)
	GetIslandHeader(ctx context.Context, islandId string) (IslandHeader, error)
	GetIslandHeadersByTerritory(ctx context.Context, territoryId string) ([]IslandHeader, error)
	GetIslandHeaderByBookIdAndUserId(ctx context.Context, bookId string, userId int32) (IslandHeader, error)
	SetTerritoryPoolSettings(ctx context.Context, territoryId string, settings TerritoryPoolSettings) error
	GetTerritoryPoolSettings(ctx context.Context, territoryId string) (TerritoryPoolSettings, error)
	AddBookToPool(ctx context.Context, poolId string, bookId string) error
	GetPoolOfBook(ctx context.Context, bookId string) (poolId string, found bool, err error)
	AssignBookToIslandFromPool(ctx context.Context, territoryId string, islandId string, userId int32) (bookId string, err error)
	IsIslandPortable(ctx context.Context, userId int32, islandId string) (bool, error)
	AddPortableIsland(ctx context.Context, userId int32, islandId string) (bool, error)
	GetPortableIslands(ctx context.Context, userId int32) (result []PortableIsland, err error)
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id int32) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type PlayerStore interface {
	Create(ctx context.Context, player Player) error
	Get(ctx context.Context, userId int32) (Player, error)
	Update(ctx context.Context, tx Tx, old, updated Player) error
	GetAll(ctx context.Context) ([]int32, error)
	CreatePlayerEvent(ctx context.Context, userId int32, createdAt time.Time, reason string, player FullPlayer) error
	GetLocations(ctx context.Context, territoryID string) (map[string][]int32, error)
}

type QuestionStore interface {
	BindQuestionsToBook(ctx context.Context, bookId string, questions []BookQuestion) error
	// GetOrCreateAnswer gets the Answer if it exists
	// otherwise creates an Answer with the given ID and zero value for other fields (except timestamps).
	GetOrCreateAnswer(ctx context.Context, userId int32, questionID string) (Answer, error)
	GetAnswer(ctx context.Context, userId int32, questionId string) (Answer, error)
	GetPendingAnswers(ctx context.Context, ifBefore time.Time) ([]Answer, error)
	MarkHelpRequest(ctx context.Context, userId int32, questionId string) error
	SetHelpState(ctx context.Context, userId int32, questionId string, state HelpState) error
	// SubmitAnswer updates the existing Answer with the given args and sets the answer status to AnswerStatusPending.
	// If the answer is in AnswerStatusCorrect status, it returns ErrSubmitToCorrectAnswer error.
	// If the answer is in AnswerStatusPending status, it returns ErrSubmitToPendingAnswer error.
	SubmitAnswer(ctx context.Context, userId int32, questionId, fileID, filename, textContent string, lastUpdatedAt time.Time) (Answer, error)
	GetKnowledgeBars(ctx context.Context, userId int32) ([]KnowledgeBar, error)
	HasAnsweredIsland(ctx context.Context, userId int32, islandId string) (bool, error)
	GetQuestion(ctx context.Context, questionId string) (BookQuestion, error)
	CreateCorrection(ctx context.Context, Correction Correction) error
	ApplyCorrection(ctx context.Context, tx Tx, ifBefore time.Time, correction Correction) (Answer, bool, error)
	GetUnappliedCorrections(ctx context.Context, before time.Time) ([]Correction, error)
	UpdateCorrectionNewStatus(ctx context.Context, id string, newStatus AnswerStatus) error
	UpdateCorrectionFeedback(ctx context.Context, id string, feedback string) (AnswerStatus, error)
	FinalizeCorrection(ctx context.Context, id string) error
}

type TreasureStore interface {
	BindTreasuresToBook(ctx context.Context, bookId string, treasures []Treasure) error
	GetOrCreateUserTreasure(ctx context.Context, userId int32, treasureId string) (UserTreasure, error)
	GetTreasure(ctx context.Context, treasureId string) (Treasure, error)
	GetUserTreasure(ctx context.Context, userId int32, treasureId string) (UserTreasure, error)
	UpdateUserTreasure(ctx context.Context, old UserTreasure, updated UserTreasure) error
}

type GetOffersByFilterType string

const (
	GetOffersByAll    GetOffersByFilterType = ""
	GetOffersByMe     GetOffersByFilterType = "me"
	GetOffersByOthers GetOffersByFilterType = "others"
)

type MarketStore interface {
	CreateOffer(ctx context.Context, tx Tx, offer TradeOffer) error
	// DeleteOffer soft-deletes the offer
	DeleteOffer(ctx context.Context, tx Tx, offerId string) error
	GetOffer(ctx context.Context, offerId string) (TradeOffer, error)
	GetOffers(ctx context.Context, byFilter GetOffersByFilterType, userId int32, before time.Time, limit int) ([]TradeOffer, error)
	GetOffersCountOfUser(ctx context.Context, userId int32) (int, error)
}

type InboxStore interface {
	CreateMessage(ctx context.Context, tx Tx, msg InboxMessage) error
	GetMessages(ctx context.Context, userId int32, before time.Time, limit int) ([]InboxMessage, error)
}

type GameStateStore interface {
	GetIsPaused(ctx context.Context) (bool, error)
	SetIsPaused(ctx context.Context, isPaused bool) error
}

type InvestStore interface {
	// GetActiveSession returns the closest investment session where EndAt is in the future
	GetActiveSession(ctx context.Context) (*InvestmentSession, error)
	GetSession(ctx context.Context, id string) (*InvestmentSession, error)

	// CreateInvestmentSession creates a new investment session
	CreateInvestmentSession(ctx context.Context, tx Tx, session InvestmentSession) error

	// AddUserInvestment adds a user's investment to a session
	AddUserInvestment(ctx context.Context, tx Tx, investment UserInvestment) error

	// GetUserInvestments returns all investments for a specific user in a specific session
	GetUserInvestments(ctx context.Context, sessionID string, userID int32) ([]UserInvestment, error)
}
