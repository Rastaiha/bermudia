package mock

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/service"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

//go:embed territories
var territoryFiles embed.FS

//go:embed islands
var islandFiles embed.FS

func CreateMockData(adminService *service.Admin, mockUsersPassword string) error {
	slog.Info("Creating mock data...")
	if mockUsersPassword == "" {
		return errors.New("mock users password is empty")
	}
	if err := createMockTerritories(adminService); err != nil {
		return fmt.Errorf("failed to create mock territories: %w", err)
	}
	if err := createMockIslands(adminService); err != nil {
		return fmt.Errorf("failed to create mock islands: %w", err)
	}
	if err := createMockUsers(adminService, mockUsersPassword); err != nil {
		return fmt.Errorf("failed to create mock users: %w", err)
	}
	return nil
}

func createMockUsers(adminService *service.Admin, password string) error {
	ctx := context.Background()
	errs := []error{adminService.CreateUser(ctx, 100, "alice", password)}
	for i := range 100 {
		i = i + 1
		errs = append(errs, adminService.CreateUser(ctx, int32(100+i), fmt.Sprintf("test%d", i), password))
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("failed to create mock users: %w", err)
	}
	return nil
}

func createMockTerritories(adminService *service.Admin) error {
	ctx := context.Background()
	return fs.WalkDir(territoryFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		content, err := territoryFiles.ReadFile(path)
		if err != nil {
			return err
		}

		var territory domain.Territory
		if err := json.Unmarshal(content, &territory); err != nil {
			return err
		}

		return adminService.SetTerritory(ctx, territory)
	})
}

func createMockIslands(adminService *service.Admin) error {
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
		var islandContent domain.IslandInputContent
		if err := json.Unmarshal(content, &islandContent); err != nil {
			return err
		}
		_, err = adminService.SetIsland(ctx, id, islandContent)
		return err
	})
}
