package repository

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/Rastaiha/bermudia/internal/domain"
)

//go:embed data/territories
var territoryFiles embed.FS

// jsonTerritoryRepository implements Territory using embedded JSON files
type jsonTerritoryRepository struct {
	fs embed.FS
}

// NewJSONTerritoryRepository creates a new JSON-based territory repository
func NewJSONTerritoryRepository() domain.TerritoryStore {
	return &jsonTerritoryRepository{
		fs: territoryFiles,
	}
}

// GetTerritoryByID retrieves a territory by its ID from embedded JSON files
func (r *jsonTerritoryRepository) GetTerritoryByID(ctx context.Context, territoryID string) (*domain.Territory, error) {
	filePath := path.Join("data/territories", fmt.Sprintf("%s.json", territoryID))

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Read embedded file
	data, err := r.fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrTerritoryNotFound, territoryID)
	}

	// Parse JSON
	var territory domain.Territory
	if err := json.Unmarshal(data, &territory); err != nil {
		return nil, fmt.Errorf("failed to parse territory JSON: %w", err)
	}

	return &territory, nil
}

// ListTerritories returns all available territories from embedded files
func (r *jsonTerritoryRepository) ListTerritories(ctx context.Context) ([]domain.Territory, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	dirEntries, err := r.fs.ReadDir("data/territories")
	if err != nil {
		return nil, fmt.Errorf("failed to list territory files: %w", err)
	}

	var territories []domain.Territory
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			filePath := path.Join("data/territories", entry.Name())

			data, err := r.fs.ReadFile(filePath)
			if err != nil {
				continue // Skip files that can't be read
			}

			var territory domain.Territory
			if err := json.Unmarshal(data, &territory); err != nil {
				continue // Skip files that can't be parsed
			}

			territories = append(territories, territory)
		}
	}

	return territories, nil
}

// SQLTerritoryRepository is a placeholder for future SQL implementation
// This shows how easy it will be to swap implementations
type SQLTerritoryRepository struct {
	// db *sql.DB // Will be added when implementing SQL version
}

// NewSQLTerritoryRepository creates a new SQL-based territory repository
// func NewSQLTerritoryRepository(db *sql.DB) *SQLTerritoryRepository {
// 	return &SQLTerritoryRepository{db: db}
// }

// Implement Territory interface methods for SQL version
// func (r *SQLTerritoryRepository) GetTerritoryByID(territoryID string) (*models.Territory, error) {
// 	// SQL implementation will go here
// 	panic("not implemented")
// }

// func (r *SQLTerritoryRepository) ListTerritories() ([]models.Territory, error) {
// 	// SQL implementation will go here
// 	panic("not implemented")
// }
