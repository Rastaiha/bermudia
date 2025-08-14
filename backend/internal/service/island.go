package service

import (
	"context"
	"github.com/Rastaiha/bermudia/internal/domain"
)

type Island struct {
	repo domain.IslandStore
}

func NewIsland(repo domain.IslandStore) *Island {
	return &Island{repo: repo}
}

func (i *Island) GetIsland(ctx context.Context, id string) (*domain.IslandContent, error) {
	return i.repo.GetByID(ctx, id)
}
