package domain

import (
	"context"
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
)

type TerritoryStore interface {
	SetTerritory(ctx context.Context, territory *Territory) error
	GetTerritoryByID(ctx context.Context, territoryID string) (*Territory, error)
	ListTerritories(ctx context.Context) ([]Territory, error)
}

type IslandStore interface {
	SetBook(ctx context.Context, book Book) error
	GetBook(ctx context.Context, bookId string) (*Book, error)
	SetIslandHeader(ctx context.Context, territoryId string, header IslandHeader) error
	ReserveIDForTerritory(ctx context.Context, territoryId, islandId string) error
	GetBookOfIsland(ctx context.Context, islandId string, userId int32) (string, error)
	GetTerritory(ctx context.Context, id string) (string, error)
	GetIslandHeadersByTerritory(ctx context.Context, territoryId string) ([]IslandHeader, error)
	SetTerritoryPoolSettings(ctx context.Context, territoryId string, settings TerritoryPoolSettings) error
	GetTerritoryPoolSettings(ctx context.Context, territoryId string) (TerritoryPoolSettings, error)
	AddBookToPool(ctx context.Context, poolId string, bookId string) error
	AssignBookToIslandFromPool(ctx context.Context, territoryId string, islandId string, userId int32) (bookId string, err error)
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id int32) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type PlayerStore interface {
	Create(ctx context.Context, player Player) error
	Get(ctx context.Context, userId int32) (Player, error)
	Update(ctx context.Context, old, updated Player) error
}

type QuestionStore interface {
	BindQuestionsToBook(ctx context.Context, bookId string, questions []BookQuestion) error
	// GetOrCreateAnswer gets the Answer if it exists
	// otherwise creates an Answer with the given ID and zero value for other fields (except timestamps).
	GetOrCreateAnswer(ctx context.Context, userId int32, questionID string) (Answer, error)
	// SubmitAnswer updates the existing Answer with the given args and sets the answer status to AnswerStatusPending.
	// If the answer is in AnswerStatusCorrect status, it returns ErrSubmitToCorrectAnswer error.
	// If the answer is in AnswerStatusPending status, it returns ErrSubmitToPendingAnswer error.
	SubmitAnswer(ctx context.Context, userId int32, questionId string, fileID, filename, textContent string) (Answer, error)
	GetKnowledgeBars(ctx context.Context, userId int32) ([]KnowledgeBar, error)
	HasAnsweredIsland(ctx context.Context, userId int32, islandId string) (bool, error)
	GetBookOfQuestion(ctx context.Context, questionId string) (string, error)
	CreateCorrection(ctx context.Context, Correction Correction) error
	ApplyCorrection(ctx context.Context, correction Correction, ifBefore time.Time) (int32, bool, error)
	GetUnappliedCorrections(ctx context.Context) ([]Correction, error)
}
