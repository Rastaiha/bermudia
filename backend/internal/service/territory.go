package service

import (
	"context"

	"github.com/Rastaiha/bermudia/internal/models"
	"github.com/Rastaiha/bermudia/internal/repository"
)

// Territory handles business logic for territories
type Territory struct {
	repo repository.Territory
}

// NewTerritory creates a new territory service
func NewTerritory(repo repository.Territory) *Territory {
	return &Territory{
		repo: repo,
	}
}

// GetTerritory retrieves a territory by ID with any business logic applied
func (s *Territory) GetTerritory(ctx context.Context, territoryID string) (*models.Territory, error) {
	return s.repo.GetTerritoryByID(ctx, territoryID)
}

// ListTerritories retrieves all territories with business logic applied
func (s *Territory) ListTerritories(ctx context.Context) ([]models.Territory, error) {
	territories, err := s.repo.ListTerritories(ctx)
	if err != nil {
		return nil, err
	}

	// Apply business logic here if needed
	// For example:
	// - Sort territories by some criteria
	// - Filter based on user permissions
	// - Add metadata

	return territories, nil
}
