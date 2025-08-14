package repository

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"path"
)

//go:embed data/islands
var islandFiles embed.FS

type jsonIslandRepository struct {
	fs embed.FS
}

func NewJSONIslandRepository() domain.IslandStore {
	return &jsonIslandRepository{
		fs: islandFiles,
	}
}

func (j *jsonIslandRepository) GetByID(_ context.Context, id string) (*domain.IslandContent, error) {
	filePath := path.Join("data/islands", fmt.Sprintf("%s.json", id))

	data, err := j.fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrIslandNotFound, id)
	}

	var islandContent domain.IslandContent
	if err := json.Unmarshal(data, &islandContent); err != nil {
		return nil, fmt.Errorf("failed to parse island JSON: %w", err)
	}

	return &islandContent, nil
}
