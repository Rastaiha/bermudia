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
	"math/rand"
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
	f, err := os.Create(filepath.Join(os.TempDir(), fmt.Sprintf("content_%d.zip", rand.Int31())))
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

	dir := filepath.Join(os.TempDir(), "content_"+fmt.Sprint(rand.Int31())+"/")
	err = ExtractZip(f, stat.Size(), dir)
	if err != nil {
		return nil, err
	}
	root := os.DirFS(dir)
	return root, nil
}

//go:embed data
var DataFiles embed.FS

func SetGameContent(adminService *service.Admin, files fs.FS, writeBackPath string, defaultPass string) error {
	slog.Info("Setting game content...")
	if writeBackPath != "" {
		err := os.CopyFS(writeBackPath, DataFiles)
		if err != nil {
			return fmt.Errorf("could not copy fs: %w", err)
		}
	}
	if err := createMockTerritories(adminService, files); err != nil {
		return fmt.Errorf("failed to create mock territories: %w", err)
	}
	if err := createMockBooks(adminService, files, writeBackPath); err != nil {
		return fmt.Errorf("failed to create mock islands: %w", err)
	}
	if err := createPoolSettings(adminService, files, writeBackPath); err != nil {
		return fmt.Errorf("failed to create mock pool settings: %w", err)
	}
	if err := createMockUsers(adminService, files, writeBackPath, defaultPass); err != nil {
		return fmt.Errorf("failed to create mock users: %w", err)
	}
	return nil
}

func createMockUsers(adminService *service.Admin, files fs.FS, writeBack string, defaultPass string) error {
	path := "data/users.json"
	usersJson, err := fs.ReadFile(files, path)
	if err != nil {
		return err
	}
	var users []service.User
	err = json.Unmarshal(usersJson, &users)
	if err != nil {
		return err
	}
	ctx := context.Background()
	var errs []error
	result := make([]service.User, 0, len(users))
	for _, u := range users {
		if u.Password == "" && defaultPass != "" {
			u.Password = defaultPass
		}
		u, err := adminService.CreateUser(ctx, u)
		errs = append(errs, err)
		result = append(result, u)
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("failed to create mock users: %w", err)
	}
	return writeBackData(writeBack, path, result)
}

func createMockTerritories(adminService *service.Admin, territoryFiles fs.FS) error {
	ctx := context.Background()
	return fs.WalkDir(territoryFiles, "data/territories", func(path string, d fs.DirEntry, err error) error {
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

func createMockBooks(adminService *service.Admin, booksFiles fs.FS, writeBack string) error {
	root := "data/books"
	ctx := context.Background()
	islandsDir := filepath.Join(root, "islands/")
	poolDir := filepath.Join(root, "pool/")
	return fs.WalkDir(booksFiles, root, func(path string, d fs.DirEntry, err error) error {
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
		if filepath.Clean(dir) == islandsDir {
			islandId := strings.TrimSuffix(file, filepath.Ext(file))
			book, err = adminService.SetBookAndBindToIsland(ctx, islandId, book)
			if err != nil {
				return err
			}
		} else if pool, ok := strings.CutPrefix(dir, poolDir); ok {
			pool = strings.Trim(pool, "/")
			book, err = adminService.SetBookAndBindToPool(ctx, pool, book)
			if err != nil {
				return err
			}
		} else {
			return errors.New("unknown book with path " + path)
		}

		return writeBackData(writeBack, path, book)
	})
}

func createPoolSettings(adminService *service.Admin, poolSettingsFiles fs.FS, writeBack string) error {
	ctx := context.Background()
	return fs.WalkDir(poolSettingsFiles, "data/pool_settings", func(path string, d fs.DirEntry, err error) error {
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
		bindings, err = adminService.SetTerritoryIslandBindings(ctx, bindings)
		if err != nil {
			return err
		}
		return writeBackData(writeBack, path, bindings)
	})
}

func writeBackData(writeBack string, path string, data any) error {
	if writeBack == "" {
		return nil
	}
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(writeBack, path), j, os.FileMode(os.O_TRUNC|os.O_WRONLY))
	if err != nil {
		return fmt.Errorf("failed to write back to %q: %w", path, err)
	}
	return nil
}
