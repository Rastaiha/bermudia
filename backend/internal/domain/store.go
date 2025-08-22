package domain

import (
	"context"
	"errors"
)

var (
	ErrIslandNotFound    = errors.New("island not found")
	ErrTerritoryNotFound = errors.New("territory not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrPlayerConflict    = errors.New("player update conflict")
)

type TerritoryStore interface {
	CreateTerritory(ctx context.Context, territory *Territory) error
	GetTerritoryByID(ctx context.Context, territoryID string) (*Territory, error)
	ListTerritories(ctx context.Context) ([]Territory, error)
}

type IslandStore interface {
	SetContent(ctx context.Context, id string, content *IslandRawContent) error
	ReserveIDForTerritory(ctx context.Context, territoryId, islandId string) error
	GetByID(ctx context.Context, id string) (*IslandRawContent, error)
	GetTerritory(ctx context.Context, id string) (string, error)
	// GetOrCreateUserComponent gets the component if it exists, otherwise it creates a new component by generating a NewID for the resource type.
	// If the IdHasType returns false for the exiting ResourceID, it returns an error.
	GetOrCreateUserComponent(ctx context.Context, islandId string, userId int32, componentId string, resourceType ResourceType) (UserComponent, error)
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
	// SetQuestion creates the question if it does not exist, otherwise updates all fields based on id
	SetQuestion(ctx context.Context, question Question) error
	GetQuestion(ctx context.Context, id string) (Question, error)
	// GetOrCreateAnswer gets the Answer if it exists
	// otherwise creates an Answer with the given ID and zero value for other fields (except timestamps).
	GetOrCreateAnswer(ctx context.Context, userId int32, answerID string, questionID string) (Answer, error)
	// SubmitAnswer updates the existing Answer with the given args and sets the answer status to AnswerStatusPending.
	// If the answer is in AnswerStatusCorrect status, it returns ErrSubmitToCorrectAnswer error.
	// If the answer is in AnswerStatusPending status, it returns ErrSubmitToPendingAnswer error.
	SubmitAnswer(ctx context.Context, answerId string, userId int32, fileID, filename string) (Answer, error)
}
