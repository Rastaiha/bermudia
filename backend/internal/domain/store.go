package domain

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrIslandNotFound    = errors.New("island not found")
	ErrTerritoryNotFound = errors.New("territory not found")
	ErrUserNotFound      = errors.New("user not found")
)

type IslandStore interface {
	GetByID(ctx context.Context, id string) (*IslandContent, error)
}

type TerritoryStore interface {
	GetTerritoryByID(ctx context.Context, territoryID string) (*Territory, error)
	ListTerritories(ctx context.Context) ([]Territory, error)
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id int32) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

func CreateMockData(userStore UserStore, mockUsersPassword string) error {
	if mockUsersPassword == "" {
		return errors.New("mock users password is empty")
	}
	err := errors.Join(
		createMockUser(userStore, 100, "alice", mockUsersPassword),
	)
	return err
}

func createMockUser(store UserStore, id int32, username string, password string) error {
	hp, err := HashPassword(password)
	if err != nil {
		return err
	}
	err = store.Create(context.Background(), &User{
		ID:             id,
		Username:       username,
		HashedPassword: hp,
	})
	if err != nil {
		return fmt.Errorf("failed to create mock user: %w", err)
	}
	return nil
}
