package service

import (
	"context"
	"github.com/Rastaiha/bermudia/internal/models"
	"github.com/Rastaiha/bermudia/internal/repository"
)

type Island struct {
	repo repository.Island
}

func NewIsland(repo repository.Island) *Island {
	return &Island{repo: repo}
}

func (i *Island) GetIsland(ctx context.Context, id string) (*models.IslandContent, error) {
	return i.repo.GetByID(ctx, id)
}
