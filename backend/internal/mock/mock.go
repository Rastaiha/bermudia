package mock

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/service"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FsFromURL(url string) (fs.FS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Join(os.TempDir(), "content.zip"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	err = ExtractZip(f, stat.Size(), "./content")
	if err != nil {
		return nil, err
	}
	root := os.DirFS("./content")
	return root, nil
}

//go:embed data
var DataFiles embed.FS

func CreateMockData(adminService *service.Admin, mockUsersPassword string, root fs.FS) error {
	slog.Info("Creating mock data...")
	territoryFiles, err1 := fs.Sub(root, "data/territories")
	booksFiles, err2 := fs.Sub(root, "data/books")
	poolSettingsFiles, err3 := fs.Sub(root, "data/pool_settings")
	usersFile, err4 := fs.ReadFile(root, "data/users.json")
	err := errors.Join(err1, err2, err3, err4)
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
	if err := createMockUsers(adminService, usersFile, mockUsersPassword); err != nil {
		return fmt.Errorf("failed to create mock users: %w", err)
	}
	return nil
}

type mockUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createMockUsers(adminService *service.Admin, usersJson []byte, defaultPass string) error {
	var users []mockUser
	err := json.Unmarshal(usersJson, &users)
	if err != nil {
		return err
	}
	ctx := context.Background()
	var errs []error
	for i, u := range users {
		password := u.Password
		if password == "" {
			password = defaultPass
		}
		errs = append(errs, adminService.CreateUser(ctx, int32(1001*(i+1)), u.Username, password))
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
