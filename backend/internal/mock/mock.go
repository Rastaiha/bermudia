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

//go:embed data
var DataFiles embed.FS

func CreateMockData(adminService *service.Admin, mockUsersPassword string, root fs.FS) error {
	slog.Info("Creating mock data...")
	territoryFiles, err1 := fs.Sub(root, "data/territories")
	booksFiles, err2 := fs.Sub(root, "data/books")
	poolSettingsFiles, err3 := fs.Sub(root, "data/pool_settings")
	err := errors.Join(err1, err2, err3)
	if err != nil {
		return fmt.Errorf("bad file structure: %w", err)
	}
	if mockUsersPassword == "" {
		return errors.New("mock users password is empty")
	}
	if err := createMockTerritories(adminService, territoryFiles); err != nil {
		return fmt.Errorf("failed to create mock territories: %w", err)
	}
	if err := createMockBooks(adminService, booksFiles); err != nil {
		return fmt.Errorf("failed to create mock islands: %w", err)
	}
	if err := createPoolSettings(adminService, poolSettingsFiles); err != nil {
		return fmt.Errorf("failed to create mock pool settings: %w", err)
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

func createMockTerritories(adminService *service.Admin, territoryFiles fs.FS) error {
	ctx := context.Background()
	return fs.WalkDir(territoryFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(territoryFiles, path)
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

func createMockBooks(adminService *service.Admin, booksFiles fs.FS) error {
	ctx := context.Background()
	return fs.WalkDir(booksFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(booksFiles, path)
		if err != nil {
			return err
		}
		var book service.BookInput
		if err := json.Unmarshal(content, &book); err != nil {
			return err
		}
		dir, file := filepath.Split(path)
		if dir == "islands/" {
			islandId := strings.TrimSuffix(file, filepath.Ext(file))
			book, err = adminService.SetBookAndBindToIsland(ctx, islandId, book)
			return err
		}
		if pool, ok := strings.CutPrefix(dir, "pool/"); ok {
			pool = strings.Trim(pool, "/")
			book, err = adminService.SetBookAndBindToPool(ctx, pool, book)
			return err
		}
		return errors.New("unknown book with path " + path)
	})
}

func createPoolSettings(adminService *service.Admin, poolSettingsFiles fs.FS) error {
	ctx := context.Background()
	return fs.WalkDir(poolSettingsFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(poolSettingsFiles, path)
		if err != nil {
			return err
		}
		var bindings service.TerritoryIslandBindings
		if err := json.Unmarshal(content, &bindings); err != nil {
			return err
		}
		_, err = adminService.SetTerritoryIslandBindings(ctx, bindings)
		return err
	})
}
