package mock

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

//go:embed territories
var territoryFiles embed.FS

//go:embed islands
var islandFiles embed.FS

func CreateMockData(userStore domain.UserStore, territoryStore domain.TerritoryStore, islandStore domain.IslandStore, mockUsersPassword string) error {
	if mockUsersPassword == "" {
		return errors.New("mock users password is empty")
	}
	if err := errors.Join(
		createMockUser(userStore, 100, "alice", mockUsersPassword),
	); err != nil {
		return fmt.Errorf("failed to create mock users: %w", err)
	}
	if err := createMockTerritories(territoryStore, islandStore); err != nil {
		return fmt.Errorf("failed to create mock territories: %w", err)
	}
	if err := createMockIslands(islandStore); err != nil {
		return fmt.Errorf("failed to create mock islands: %w", err)
	}
	return nil
}

func createMockUser(store domain.UserStore, id int32, username string, password string) error {
	hp, err := domain.HashPassword(password)
	if err != nil {
		return err
	}
	err = store.Create(context.Background(), &domain.User{
		ID:             id,
		Username:       username,
		HashedPassword: hp,
	})
	if err != nil {
		return fmt.Errorf("failed to create mock user: %w", err)
	}
	return nil
}

func createMockTerritories(territoryStore domain.TerritoryStore, islandStore domain.IslandStore) error {
	ctx := context.Background()
	return fs.WalkDir(territoryFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if id == "" {
			return errors.New("invalid territory id")
		}

		content, err := territoryFiles.ReadFile(path)
		if err != nil {
			return err
		}

		var territory domain.Territory
		if err := json.Unmarshal(content, &territory); err != nil {
			return err
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
			if island.ID == "" {
				return fmt.Errorf("empty island id in island list")
			}
			if err := islandStore.ReserveIDForTerritory(ctx, id, island.ID); err != nil {
				return err
			}
		}

		return territoryStore.CreateTerritory(ctx, &territory)
	})
}

func createMockIslands(islandStore domain.IslandStore) error {
	ctx := context.Background()
	return fs.WalkDir(islandFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if id == "" {
			return errors.New("invalid island id")
		}
		content, err := islandFiles.ReadFile(path)
		if err != nil {
			return err
		}
		var islandContent domain.IslandContent
		if err := json.Unmarshal(content, &islandContent); err != nil {
			return err
		}
		return islandStore.SetContent(ctx, id, &islandContent)
	})
}
