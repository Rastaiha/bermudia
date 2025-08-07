package service

import (
	"context"

	"github.com/Rastaiha/rasta-1404-contest/internal/models"
	"github.com/Rastaiha/rasta-1404-contest/internal/repository"
)

// TerritoryService handles business logic for territories
type TerritoryService struct {
	repo repository.TerritoryRepository
}

// NewTerritoryService creates a new territory service
func NewTerritoryService(repo repository.TerritoryRepository) *TerritoryService {
	return &TerritoryService{
		repo: repo,
	}
}

// GetTerritory retrieves a territory by ID with any business logic applied
func (s *TerritoryService) GetTerritory(ctx context.Context, territoryID string) (*models.Territory, error) {
	territory, err := s.repo.GetTerritoryByID(ctx, territoryID)
	if err != nil {
		return nil, err
	}

	// Apply any business logic here if needed
	// For example, you might want to:
	// - Validate territory data
	// - Apply user-specific filtering
	// - Add computed fields
	// - Log access patterns

	return territory, nil
}

// ListTerritories retrieves all territories with business logic applied
func (s *TerritoryService) ListTerritories(ctx context.Context) ([]models.Territory, error) {
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
