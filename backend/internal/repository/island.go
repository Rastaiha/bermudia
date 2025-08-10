package repository

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/models"
	"path"
)

var (
	ErrIslandNotFound = errors.New("island not found")
)

type Island interface {
	GetByID(ctx context.Context, id string) (*models.IslandContent, error)
}

//go:embed data/islands
var islandFiles embed.FS

type jsonIslandRepository struct {
	fs embed.FS
}

func NewJSONIslandRepository() Island {
	return &jsonIslandRepository{
		fs: islandFiles,
	}
}

func (j *jsonIslandRepository) GetByID(_ context.Context, id string) (*models.IslandContent, error) {
	filePath := path.Join("data/islands", fmt.Sprintf("%s.json", id))

	data, err := j.fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIslandNotFound, id)
	}

	var islandContent models.IslandContent
	if err := json.Unmarshal(data, &islandContent); err != nil {
		return nil, fmt.Errorf("failed to parse island JSON: %w", err)
	}

	return &islandContent, nil
}
